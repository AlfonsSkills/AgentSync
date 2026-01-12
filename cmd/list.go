package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/AlfonsSkills/AgentSync/internal/skill"
	"github.com/AlfonsSkills/AgentSync/internal/target"
)

// listCmd list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed skills",
	Long: `List locally installed skills (scans actual directories).

Examples:
  agentsync list
  agentsync list --target gemini`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// LocalSkill represents a locally discovered skill
type LocalSkill struct {
	Name     string
	Path     string
	Target   target.Target
	Valid    bool   // Contains SKILL.md
	Category string // Category (e.g., public, .system, or empty for root)
}

func runList(cmd *cobra.Command, args []string) error {
	// Parse target filter
	targets, err := target.ParseTargets(targetFlags)
	if err != nil {
		return err
	}

	// Scan skills in target directories
	var allSkills []LocalSkill
	for _, t := range targets {
		skills, err := scanLocalSkills(t)
		if err != nil {
			color.Yellow("⚠ Failed to scan %s: %v\n", t.DisplayName(), err)
			continue
		}
		allSkills = append(allSkills, skills...)
	}

	if len(allSkills) == 0 {
		color.Yellow("📭 No installed skills found\n")
		return nil
	}

	// Group by target for display
	skillsByTarget := make(map[target.Target][]LocalSkill)
	for _, s := range allSkills {
		skillsByTarget[s.Target] = append(skillsByTarget[s.Target], s)
	}

	color.Cyan("📦 Installed Skills:\n\n")

	for _, t := range targets {
		skills := skillsByTarget[t]
		if len(skills) == 0 {
			continue
		}

		// Target header
		color.White("  %s (%d):\n", color.New(color.Bold).Sprint(t.DisplayName()), len(skills))

		skillsDir, _ := t.SkillsDir()
		color.HiBlack("  📁 %s\n", skillsDir)

		// Group by category
		byCategory := make(map[string][]LocalSkill)
		for _, s := range skills {
			byCategory[s.Category] = append(byCategory[s.Category], s)
		}

		// Show root skills first
		if rootSkills, ok := byCategory[""]; ok {
			for _, s := range rootSkills {
				printSkill(s)
			}
		}

		// Then show categorized skills
		for cat, catSkills := range byCategory {
			if cat == "" {
				continue
			}
			color.HiBlack("    [%s]\n", cat)
			for _, s := range catSkills {
				printSkill(s)
			}
		}

		fmt.Println()
	}

	return nil
}

// printSkill prints a single skill
func printSkill(s LocalSkill) {
	prefix := "    "
	if s.Category != "" {
		prefix = "      "
	}
	if s.Valid {
		color.Green("%s✓ %s\n", prefix, s.Name)
	} else {
		color.Yellow("%s⚠ %s (missing SKILL.md)\n", prefix, s.Name)
	}
}

// scanLocalSkills scans skills in the specified target directory
func scanLocalSkills(t target.Target) ([]LocalSkill, error) {
	skillsDir, err := t.SkillsDir()
	if err != nil {
		return nil, err
	}

	// Check if directory exists
	if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
		return nil, nil // Directory doesn't exist, return empty list
	}

	// Read directory contents
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var skills []LocalSkill

	// Known category subdirectories (need recursive scanning)
	knownCategories := map[string]bool{
		"public":  true,
		".system": true,
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue // Skip files, only process directories
		}

		name := entry.Name()
		entryPath := filepath.Join(skillsDir, name)

		// Check if it's a known category directory
		if knownCategories[name] {
			// Recursively scan category directory
			catSkills, err := scanCategoryDir(entryPath, name, t)
			if err != nil {
				continue
			}
			skills = append(skills, catSkills...)
		} else {
			// Regular skill directory
			valid := skill.ValidateSkillDir(entryPath) == nil

			// Skip hidden directories (if not a valid skill)
			if strings.HasPrefix(name, ".") && !valid {
				continue
			}

			skills = append(skills, LocalSkill{
				Name:     name,
				Path:     entryPath,
				Target:   t,
				Valid:    valid,
				Category: "",
			})
		}
	}

	return skills, nil
}

// scanCategoryDir scans a category subdirectory
func scanCategoryDir(dir, category string, t target.Target) ([]LocalSkill, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var skills []LocalSkill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Skip hidden files
		if strings.HasPrefix(name, ".") {
			continue
		}

		entryPath := filepath.Join(dir, name)
		valid := skill.ValidateSkillDir(entryPath) == nil

		skills = append(skills, LocalSkill{
			Name:     name,
			Path:     entryPath,
			Target:   t,
			Valid:    valid,
			Category: category,
		})
	}

	return skills, nil
}
