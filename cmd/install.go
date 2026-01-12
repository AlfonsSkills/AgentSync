package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/AlfonsSkills/AgentSync/internal/git"
	"github.com/AlfonsSkills/AgentSync/internal/skill"
	"github.com/AlfonsSkills/AgentSync/internal/target"
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
  agentsync install AlfonsSkills/skills
  agentsync install AlfonsSkills/skills --target gemini
  agentsync install https://github.com/AlfonsSkills/skills.git -t claude,codex`,
	Args: cobra.ExactArgs(1),
	RunE: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	source := args[0]

	// Parse target tools
	targets, err := target.ParseTargets(targetFlags)
	if err != nil {
		return err
	}

	// Create Git fetcher
	fetcher := git.NewFetcher()

	color.Cyan("📦 Cloning repository...\n")
	color.White("   Source: %s\n", fetcher.NormalizeURL(source))

	// Clone to temp directory
	tempDir, err := fetcher.CloneToTemp(source)
	if err != nil {
		color.Red("❌ Clone failed: %v\n", err)
		return err
	}
	defer os.RemoveAll(tempDir) // Cleanup temp directory

	// Scan skills in repository
	skills, err := skill.ScanSkills(tempDir)
	if err != nil {
		color.Red("❌ Scan failed: %v\n", err)
		return err
	}

	// Check if any skills found
	if len(skills) == 0 {
		// Try to validate root directory as a skill
		if err := skill.ValidateSkillDir(tempDir); err != nil {
			color.Red("❌ No valid skills found in repository\n")
			return fmt.Errorf("no skills found in repository")
		}
		// Root directory is a skill
		repoName := skill.ExtractSkillName(source)
		skills = []skill.SkillInfo{{
			Name: repoName,
			Path: tempDir,
		}}
	}

	// Select skills to install
	color.Green("✓ Found %d skill(s)\n\n", len(skills))

	// Build options list
	var options []string
	for _, s := range skills {
		if s.Desc != "" {
			options = append(options, fmt.Sprintf("%s - %s", s.Name, s.Desc))
		} else {
			options = append(options, s.Name)
		}
	}

	// Interactive multi-select
	var selectedIndices []int
	prompt := &survey.MultiSelect{
		Message: "Select skills to install:",
		Options: options,
	}
	if err := survey.AskOne(prompt, &selectedIndices); err != nil {
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

	// Install selected skills
	copyOpts := skill.DefaultCopyOptions()
	totalInstalled := 0

	for _, s := range selectedSkills {
		color.Cyan("\n📦 Installing: %s\n", s.Name)
		installedCount := 0

		for _, t := range targets {
			skillsDir, err := t.EnsureSkillsDir()
			if err != nil {
				color.Yellow("   ⚠ Skipping %s: %v\n", t.DisplayName(), err)
				continue
			}

			destDir := skillsDir + "/" + s.Name

			// Remove existing if exists
			if _, err := os.Stat(destDir); !os.IsNotExist(err) {
				os.RemoveAll(destDir)
			}

			if err := skill.CopyDir(s.Path, destDir, copyOpts); err != nil {
				color.Yellow("   ⚠ Copy to %s failed: %v\n", t.DisplayName(), err)
				continue
			}

			color.Green("   ✓ Installed to %s: %s\n", t.DisplayName(), destDir)
			installedCount++
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
