package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
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
  skillsync remove my-skill
  skillsync remove my-skill --target gemini
  skillsync remove my-skill --local`,
	Args: cobra.ExactArgs(1),
	RunE: runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&localRemove, "local", "l", false, "Remove from project-local skills directories only")
}

func runRemove(cmd *cobra.Command, args []string) error {
	skillName := args[0]

	color.Cyan("🗑️  Preparing to remove: %s\n\n", skillName)

	// Step 1: Resolve target providers (仅显示存在该 skill 的工具)
	providers, _, err := resolveTargetProvidersForRemove(skillName, targetFlags)
	if err != nil {
		return err
	}

	// Step 2: Resolve remove scope (仅当项目目录存在 skill 时才提示)
	removeGlobal, removeLocal, projectRoot, err := resolveRemoveScopeWithCheck(skillName, providers, localRemove)
	if err != nil {
		return err
	}

	// Step 3: Show removal preview
	showRemovePreview(skillName, providers, removeGlobal, removeLocal, projectRoot)

	// Step 4: Confirm removal
	var confirmRemove bool
	confirmPrompt := &survey.Confirm{
		Message: "Proceed with removal?",
		Default: false,
	}
	if err := survey.AskOne(confirmPrompt, &confirmRemove); err != nil {
		return fmt.Errorf("cancelled: %w", err)
	}
	if !confirmRemove {
		color.Yellow("Removal cancelled\n")
		return nil
	}

	// Step 5: Execute removal
	color.Cyan("\n🗑️  Removing skill: %s\n", skillName)
	removedCount := 0

	for _, p := range providers {
		// Remove from global directory
		if removeGlobal {
			globalDir, err := p.GlobalInstallDir()
			if err == nil {
				skillPath := filepath.Join(globalDir, skillName)
				if _, err := os.Stat(skillPath); os.IsNotExist(err) {
					color.Yellow("   ⚠ %s: not found\n", p.DisplayName())
				} else if err := os.RemoveAll(skillPath); err != nil {
					color.Red("   ❌ %s: failed to remove - %v\n", p.DisplayName(), err)
				} else {
					color.Green("   ✓ Removed from %s\n", p.DisplayName())
					removedCount++
				}
			}
		}

		// Remove from project directory
		if removeLocal && projectRoot != "" {
			localDir := p.LocalSkillsDir(projectRoot)
			skillPath := filepath.Join(localDir, skillName)
			if _, err := os.Stat(skillPath); os.IsNotExist(err) {
				color.Yellow("   ⚠ .%s/skills: not found\n", p.Type())
			} else if err := os.RemoveAll(skillPath); err != nil {
				color.Red("   ❌ .%s/skills: failed to remove - %v\n", p.Type(), err)
			} else {
				color.Green("   ✓ Removed from .%s/skills\n", p.Type())
				removedCount++
			}
		}
	}

	if removedCount > 0 {
		color.Green("\n✅ Skill '%s' removed successfully!\n", skillName)
	} else {
		color.Yellow("\n⚠ No files were actually removed\n")
	}

	return nil
}
