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
  user/repo                                   Use GitHub (default)
  https://github.com/user/repo                Full URL
  https://github.com/user/repo/tree/br/path   Specific skill path (with default selection)

Examples:
  skillsync install AlfonsSkills/skills
  skillsync install AlfonsSkills/skills --target gemini
  skillsync install AlfonsSkills/skills --local
  skillsync install https://github.com/AlfonsSkills/skills.git -t claude,codex
  skillsync install https://github.com/AlfonsSkills/skills/tree/main/all-money-back-my-home`,
	Args: cobra.ExactArgs(1),
	RunE: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVarP(&localInstall, "local", "l", false, "Install to project-local skills directories only")
}

func runInstall(cmd *cobra.Command, args []string) error {
	source := args[0]

	// Create Git fetcher
	fetcher := git.NewFetcher()

	var tempDir string
	var targetPath string // 指定的子目录路径
	var targetFullPath string // 解析后的目标目录完整路径
	var err error

	// 检测是否为 Tree URL (支持 GitHub, 未来可扩展 GitLab 等)
	if git.IsTreeURL(source) {
		treeURL, parseErr := git.ParseTreeURL(source)
		if parseErr != nil {
			color.Red("❌ Invalid tree URL: %v\n", parseErr)
			return parseErr
		}

		color.Cyan("📦 Cloning repository (branch: %s)...\n", treeURL.Branch)
		color.White("   Source: %s\n", treeURL.CloneURL())
		color.White("   Target Path: %s\n\n", treeURL.Path)

		tempDir, err = fetcher.CloneToTempWithBranch(treeURL.CloneURL(), treeURL.Branch)
		targetPath = treeURL.Path
	} else {
		// 原有逻辑
		color.Cyan("📦 Cloning repository...\n")
		color.White("   Source: %s\n\n", fetcher.NormalizeURL(source))
		tempDir, err = fetcher.CloneToTemp(source)
	}

	if err != nil {
		color.Red("❌ Clone failed: %v\n", err)
		return err
	}
	defer os.RemoveAll(tempDir)

	// 如果指定了 targetPath，验证路径是否存在
	if targetPath != "" {
		// 关键步骤：根据 Tree URL 计算 skill 目标目录
		targetFullPath = filepath.Join(tempDir, targetPath)
		if _, statErr := os.Stat(targetFullPath); os.IsNotExist(statErr) {
			color.Red("❌ Target path not found: %s\n", targetPath)
			color.Yellow("   The specified path does not exist in the repository.\n")
			color.Yellow("   Please check the URL and try again.\n")
			return fmt.Errorf("target path not found: %s", targetPath)
		}
		// 验证目标是否为有效的 skill 目录
		if err := skill.ValidateSkillDir(targetFullPath); err != nil {
			color.Red("❌ Target path is not a valid skill: %s\n", targetPath)
			color.Yellow("   The directory must contain a SKILL.md file.\n")
			return fmt.Errorf("target path is not a valid skill: %s", targetPath)
		}
	}

	// Step 1: Build skill list (Tree URL 指定时仅选择该 skill)
	var skills []skill.SkillInfo
	if targetPath != "" {
		// 关键步骤：tree URL 已明确 skill，补充读取描述并复用展示格式
		desc := skill.ReadSkillDescription(targetFullPath)
		skills = []skill.SkillInfo{{
			Name: filepath.Base(targetFullPath),
			Path: targetFullPath,
			Desc: desc,
		}}
	} else {
		// Scan skills in repository
		skills, err = skill.ScanSkills(tempDir)
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
	}

	// Step 2: Select skills to install
	color.Green("✓ Found %d skill(s)\n\n", len(skills))

	var selectedSkills []skill.SkillInfo
	if targetPath != "" {
		// Tree URL 已明确 skill 路径，直接使用该 skill，并复用名称/描述展示格式
		cyan := color.New(color.FgCyan).SprintFunc()
		s := skills[0]
		if s.Desc != "" {
			color.White("   %s - %s\n\n", cyan(s.Name), s.Desc)
		} else {
			color.White("   %s\n\n", cyan(s.Name))
		}
		selectedSkills = skills
	} else {
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

		for _, idx := range selectedIndices {
			selectedSkills = append(selectedSkills, skills[idx])
		}
		fmt.Println()
	}

	// Step 3: Resolve target providers (interactive if not specified)
	providers, _, err := resolveTargetProviders(targetFlags)
	if err != nil {
		return err
	}

	// Step 4: Resolve install scope (global/local)
	installGlobal, installLocal, projectRoot, err := resolveLocalInstall(localInstall)
	if err != nil {
		return err
	}

	// Step 5: Show installation preview
	showInstallPreview(selectedSkills, providers, installGlobal, installLocal, projectRoot)

	// Step 6: Confirm and execute installation
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

	// 检查更新（利用已有网络连接）
	checkUpdateInBackground()

	return nil
}
