package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"

	"github.com/AlfonsSkills/AgentSync/internal/git"
	"github.com/AlfonsSkills/AgentSync/internal/project"
	"github.com/AlfonsSkills/AgentSync/internal/skill"
)

// runLocalInstall handles --local installation to project directories
func runLocalInstall(source string) error {
	// Find project root
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		color.Red("❌ %v\n", err)
		color.Yellow("💡 The --local flag requires a git repository\n")
		return err
	}

	color.Cyan("📁 Project root: %s\n", projectRoot)

	// Clone repository
	fetcher := git.NewFetcher()
	color.Cyan("📦 Cloning repository...\n")
	color.White("   Source: %s\n", fetcher.NormalizeURL(source))

	tempDir, err := fetcher.CloneToTemp(source)
	if err != nil {
		color.Red("❌ Clone failed: %v\n", err)
		return err
	}
	defer os.RemoveAll(tempDir)

	// Scan skills
	skills, err := skill.ScanSkills(tempDir)
	if err != nil {
		color.Red("❌ Scan failed: %v\n", err)
		return err
	}

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

	// Interactive selection
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
	prompt := &survey.MultiSelect{
		Message:  "Select skills to install:",
		Options:  options,
		PageSize: 10,
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

	// Install to all three project directories
	targets := []string{"gemini", "claude", "codex"}
	copyOpts := skill.DefaultCopyOptions()
	totalInstalled := 0

	for _, s := range selectedSkills {
		color.Cyan("\n📦 Installing: %s\n", s.Name)
		installedCount := 0

		for _, targetName := range targets {
			// Get project-local skills directory
			skillsDir, err := project.GetLocalSkillsDir(targetName)
			if err != nil {
				color.Yellow("   ⚠ Failed to get %s directory: %v\n", targetName, err)
				continue
			}

			// Ensure directory exists
			if err := os.MkdirAll(skillsDir, 0755); err != nil {
				color.Yellow("   ⚠ Failed to create %s directory: %v\n", targetName, err)
				continue
			}

			destDir := filepath.Join(skillsDir, s.Name)

			// Remove existing if exists
			if _, err := os.Stat(destDir); !os.IsNotExist(err) {
				os.RemoveAll(destDir)
			}

			if err := skill.CopyDir(s.Path, destDir, copyOpts); err != nil {
				color.Yellow("   ⚠ Copy to .%s/skills failed: %v\n", targetName, err)
				continue
			}

			color.Green("   ✓ Installed to .%s/skills: %s\n", targetName, destDir)
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

	color.Green("\n✅ Installation complete! %d skill(s) installed to project directories\n", totalInstalled)
	return nil
}
