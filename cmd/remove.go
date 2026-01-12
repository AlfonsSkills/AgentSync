package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/AlfonsSkills/AgentSync/internal/target"
)

var (
	localRemove bool
)

// removeCmd remove command
var removeCmd = &cobra.Command{
	Use:   "remove <skill-name>",
	Short: "Remove installed skill",
	Long: `Remove an installed skill.

Examples:
  agentsync remove my-skill
  agentsync remove my-skill --target gemini
  agentsync remove my-skill --local  # Remove from project directories`,
	Args: cobra.ExactArgs(1),
	RunE: runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&localRemove, "local", "l", false, "Remove from project-local skills directories")
}

func runRemove(cmd *cobra.Command, args []string) error {
	skillName := args[0]

	// Handle local removal
	if localRemove {
		return runLocalRemove(skillName)
	}

	// Parse target tools for global removal
	targets, err := target.ParseTargets(targetFlags)
	if err != nil {
		return err
	}

	// Build target names for display
	var targetNames []string
	for _, t := range targets {
		targetNames = append(targetNames, t.DisplayName())
	}

	// Always require confirmation with target list
	color.Yellow("⚠ This will remove '%s' from: %s\n", skillName, strings.Join(targetNames, ", "))

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

	removedCount := 0

	for _, t := range targets {
		// Use InstallDir instead of SkillsDir for correct path (e.g., Codex uses public/)
		skillsDir, err := t.InstallDir()
		if err != nil {
			continue
		}

		skillPath := skillsDir + "/" + skillName

		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			color.Yellow("   ⚠ %s: not found\n", t.DisplayName())
			continue
		}

		if err := os.RemoveAll(skillPath); err != nil {
			color.Red("   ❌ %s: failed to remove - %v\n", t.DisplayName(), err)
			continue
		}

		color.Green("   ✓ Removed from %s\n", t.DisplayName())
		removedCount++
	}

	if removedCount > 0 {
		color.Green("\n✅ Skill '%s' removed successfully!\n", skillName)
	} else {
		color.Yellow("\n⚠ No files were actually removed\n")
	}

	return nil
}
