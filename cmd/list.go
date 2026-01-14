package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/AlfonsSkills/SkillSync/internal/skill"
	"github.com/AlfonsSkills/SkillSync/internal/target"
)

// listCmd list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed skills",
	Long: `List locally installed skills (scans actual directories).

Examples:
  skillsync list
  skillsync list --target gemini`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// LocalSkill represents a locally discovered skill
type LocalSkill struct {
	Name     string
	Path     string
	Provider target.ToolProvider // 使用 Provider 替代 Target
	Valid    bool                // Contains SKILL.md
	Category string              // Category (e.g., public, .system, or empty for root)
}

func runList(cmd *cobra.Command, args []string) error {
	// Parse target filter using Provider interface
	providers, err := target.ParseProviders(targetFlags)
	if err != nil {
		return err
	}

	// Scan skills in target directories (global)
	var allSkills []LocalSkill
	for _, p := range providers {
		skills, err := scanLocalSkillsWithProvider(p)
		if err != nil {
			color.Yellow("⚠ Failed to scan %s: %v\n", p.DisplayName(), err)
			continue
		}
		allSkills = append(allSkills, skills...)
	}

	// Also scan project-local skills if in a git repository
	projectSkills := scanProjectSkillsWithProviders(providers)
	allSkills = append(allSkills, projectSkills...)

	if len(allSkills) == 0 {
		color.Yellow("📭 No installed skills found\n")
		return nil
	}

	// Group by provider for display
	skillsByProvider := make(map[target.ToolType][]LocalSkill)
	for _, s := range allSkills {
		skillsByProvider[s.Provider.Type()] = append(skillsByProvider[s.Provider.Type()], s)
	}

	color.Cyan("📦 Installed Skills:\n\n")

	for _, p := range providers {
		skills := skillsByProvider[p.Type()]
		if len(skills) == 0 {
			continue
		}

		// Target header
		color.White("  %s (%d):\n", color.New(color.Bold).Sprint(p.DisplayName()), len(skills))

		skillsDir, _ := p.GlobalSkillsDir()
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

		// Then show categorized skills (from Provider.Categories())
		categories := p.Categories()
		for _, cat := range categories {
			if catSkills, ok := byCategory[cat]; ok {
				color.HiBlack("    [%s]\n", cat)
				for _, s := range catSkills {
					printSkill(s)
				}
			}
		}

		// Show any other categories not in Provider.Categories()
		for cat, catSkills := range byCategory {
			if cat == "" {
				continue
			}
			// Skip if already shown via Provider.Categories()
			found := false
			for _, c := range categories {
				if c == cat {
					found = true
					break
				}
			}
			if found {
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

// scanLocalSkillsWithProvider scans skills in the specified provider's directory
func scanLocalSkillsWithProvider(p target.ToolProvider) ([]LocalSkill, error) {
	skillsDir, err := p.GlobalSkillsDir()
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

	// Get known category subdirectories from Provider
	categories := p.Categories()
	knownCategories := make(map[string]bool)
	for _, cat := range categories {
		knownCategories[cat] = true
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
			catSkills, err := scanCategoryDirWithProvider(entryPath, name, p)
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
				Provider: p,
				Valid:    valid,
				Category: "",
			})
		}
	}

	return skills, nil
}

// scanCategoryDirWithProvider scans a category subdirectory
func scanCategoryDirWithProvider(dir, category string, p target.ToolProvider) ([]LocalSkill, error) {
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
			Provider: p,
			Valid:    valid,
			Category: category,
		})
	}

	return skills, nil
}

// scanProjectSkillsWithProviders scans project-local skills for the given providers
func scanProjectSkillsWithProviders(providers []target.ToolProvider) []LocalSkill {
	var projectSkills []LocalSkill

	// Try to find project root (silently fail if not in a git repo)
	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil
	}

	for _, p := range providers {
		// Get project-local skills directory using Provider interface
		skillsDir := p.LocalSkillsDir(projectRoot)

		// Check if directory exists
		if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
			continue
		}

		// Read directory contents
		entries, err := os.ReadDir(skillsDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			name := entry.Name()
			if strings.HasPrefix(name, ".") {
				continue
			}

			entryPath := filepath.Join(skillsDir, name)
			valid := skill.ValidateSkillDir(entryPath) == nil

			projectSkills = append(projectSkills, LocalSkill{
				Name:     name,
				Path:     entryPath,
				Provider: p,
				Valid:    valid,
				Category: fmt.Sprintf("project:%s", filepath.Base(projectRoot)),
			})
		}
	}

	return projectSkills
}

// findProjectRoot searches for project root by looking for .git directory
func findProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := currentDir
	for {
		gitPath := filepath.Join(dir, ".git")
		if info, err := os.Stat(gitPath); err == nil && info.IsDir() {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not in a git repository")
		}
		dir = parent
	}
}
