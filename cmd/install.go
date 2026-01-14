package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/AlfonsSkills/SkillSync/internal/git"
	"github.com/AlfonsSkills/SkillSync/internal/skill"
)

var (
	localInstall bool
)

// installCmd install command
var installCmd = &cobra.Command{
	Use:   "install <repository>",
	Short: "Install skills to target tools",
	Long: `Install skills from a Git repository to local AI coding tool directories.

Repository formats:
  user/repo              Use GitHub (default)
  https://github.com/... Full URL

Examples:
  skillsync install AlfonsSkills/skills
  skillsync install AlfonsSkills/skills --target gemini
  skillsync install AlfonsSkills/skills --local
  skillsync install https://github.com/AlfonsSkills/skills.git -t claude,codex`,
	Args: cobra.ExactArgs(1),
	RunE: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVarP(&localInstall, "local", "l", false, "Install to project-local skills directories only")
}

func runInstall(cmd *cobra.Command, args []string) error {
	source := args[0]

	// Create Git fetcher and clone
	fetcher := git.NewFetcher()
	color.Cyan("📦 Cloning repository...\n")
	color.White("   Source: %s\n\n", fetcher.NormalizeURL(source))

	tempDir, err := fetcher.CloneToTemp(source)
	if err != nil {
		color.Red("❌ Clone failed: %v\n", err)
		return err
	}
	defer os.RemoveAll(tempDir)

	// Scan skills in repository
	skills, err := skill.ScanSkills(tempDir)
	if err != nil {
		color.Red("❌ Scan failed: %v\n", err)
		return err
	}

	// Handle single-skill repo (root is the skill)
	if len(skills) == 0 {
		if err := skill.ValidateSkillDir(tempDir); err != nil {
			color.Red("❌ No valid skills found in repository\n")
			return fmt.Errorf("no skills found in repository")
		}
		repoName := skill.ExtractSkillName(source)
		skills = []skill.SkillInfo{{
			Name: repoName,
			Path: tempDir,
		}}
	}

	// Step 1: Select skills to install
	color.Green("✓ Found %d skill(s)\n\n", len(skills))

	var options []string
	cyan := color.New(color.FgCyan).SprintFunc()
	for _, s := range skills {
		if s.Desc != "" {
			options = append(options, fmt.Sprintf("%s - %s", cyan(s.Name), s.Desc))
		} else {
			options = append(options, cyan(s.Name))
		}
	}

	var selectedIndices []int
	skillPrompt := &survey.MultiSelect{
		Message:  "Select skills to install:",
		Options:  options,
		PageSize: 10,
	}
	if err := survey.AskOne(skillPrompt, &selectedIndices); err != nil {
		return fmt.Errorf("selection cancelled: %w", err)
	}

	if len(selectedIndices) == 0 {
		color.Yellow("⚠ No skills selected\n")
		return nil
	}

	var selectedSkills []skill.SkillInfo
	for _, idx := range selectedIndices {
		selectedSkills = append(selectedSkills, skills[idx])
	}
	fmt.Println()

	// Step 2: Resolve target providers (interactive if not specified)
	providers, _, err := resolveTargetProviders(targetFlags)
	if err != nil {
		return err
	}

	// Step 3: Resolve install scope (global/local)
	installGlobal, installLocal, projectRoot, err := resolveLocalInstall(localInstall)
	if err != nil {
		return err
	}

	// Step 4: Show installation preview
	showInstallPreview(selectedSkills, providers, installGlobal, installLocal, projectRoot)

	// Step 5: Confirm and execute installation
	var confirmInstall bool
	confirmPrompt := &survey.Confirm{
		Message: "Proceed with installation?",
		Default: true,
	}
	if err := survey.AskOne(confirmPrompt, &confirmInstall); err != nil {
		return fmt.Errorf("cancelled: %w", err)
	}
	if !confirmInstall {
		color.Yellow("Installation cancelled\n")
		return nil
	}

	// Execute installation
	copyOpts := skill.DefaultCopyOptions()
	totalInstalled := 0

	for _, s := range selectedSkills {
		color.Cyan("\n📦 Installing: %s\n", s.Name)
		installedCount := 0

		for _, p := range providers {
			// Install to global directory
			if installGlobal {
				globalDir, err := p.EnsureInstallDir()
				if err != nil {
					color.Yellow("   ⚠ Skipping %s (global): %v\n", p.DisplayName(), err)
				} else {
					destDir := filepath.Join(globalDir, s.Name)
					if _, err := os.Stat(destDir); !os.IsNotExist(err) {
						os.RemoveAll(destDir)
					}
					if err := skill.CopyDir(s.Path, destDir, copyOpts); err != nil {
						color.Yellow("   ⚠ Copy to %s failed: %v\n", p.DisplayName(), err)
					} else {
						color.Green("   ✓ %s: %s\n", p.DisplayName(), destDir)
						installedCount++
					}
				}
			}

			// Install to project directory
			if installLocal && projectRoot != "" {
				localDir, err := p.EnsureLocalInstallDir(projectRoot)
				if err != nil {
					color.Yellow("   ⚠ Skipping %s (project): %v\n", p.DisplayName(), err)
				} else {
					destDir := filepath.Join(localDir, s.Name)
					if _, err := os.Stat(destDir); !os.IsNotExist(err) {
						os.RemoveAll(destDir)
					}
					if err := skill.CopyDir(s.Path, destDir, copyOpts); err != nil {
						color.Yellow("   ⚠ Copy to .%s/skills failed: %v\n", p.Type(), err)
					} else {
						color.Green("   ✓ .%s/skills: %s\n", p.Type(), destDir)
						installedCount++
					}
				}
			}
		}

		if installedCount > 0 {
			totalInstalled++
		}
	}

	if totalInstalled == 0 {
		color.Red("\n❌ No skills installed successfully\n")
		return fmt.Errorf("installation failed")
	}

	color.Green("\n✅ Installation complete! %d skill(s) installed\n", totalInstalled)
	return nil
}
