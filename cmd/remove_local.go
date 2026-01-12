package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"

	"github.com/AlfonsSkills/AgentSync/internal/project"
)

// runLocalRemove handles --local removal from project directories
func runLocalRemove(skillName string) error {
	// Find project root
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		color.Red("❌ %v\n", err)
		color.Yellow("💡 The --local flag requires a git repository\n")
		return err
	}

	color.Cyan("📁 Project root: %s\n", projectRoot)

	// Remove from all three project directories
	targets := []string{"gemini", "claude", "codex"}
	removedCount := 0

	// Build target names for confirmation
	var targetNames []string
	for _, targetName := range targets {
		targetNames = append(targetNames, fmt.Sprintf(".%s/skills", targetName))
	}

	// Confirmation prompt
	color.Yellow("⚠ This will remove '%s' from: %s\n", skillName, fmt.Sprintf("%v", targetNames))

	var confirm bool
	prompt := &survey.Confirm{
		Message: "Are you sure you want to continue?",
		Default: false,
	}
	if err := survey.AskOne(prompt, &confirm); err != nil {
		return fmt.Errorf("cancelled: %w", err)
	}
	if !confirm {
		color.Yellow("Cancelled\n")
		return nil
	}

	color.Cyan("🗑️  Removing skill: %s\n", skillName)

	for _, targetName := range targets {
		// Get project-local skills directory
		skillsDir, err := project.GetLocalSkillsDir(targetName)
		if err != nil {
			continue
		}

		skillPath := filepath.Join(skillsDir, skillName)

		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			color.Yellow("   ⚠ .%s/skills: not found\n", targetName)
			continue
		}

		if err := os.RemoveAll(skillPath); err != nil {
			color.Red("   ❌ .%s/skills: failed to remove - %v\n", targetName, err)
			continue
		}

		color.Green("   ✓ Removed from .%s/skills\n", targetName)
		removedCount++
	}

	if removedCount > 0 {
		color.Green("\n✅ Skill '%s' removed from project directories!\n", skillName)
	} else {
		color.Yellow("\n⚠ No files were actually removed\n")
	}

	return nil
}
