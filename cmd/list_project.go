package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlfonsSkills/AgentSync/internal/project"
	"github.com/AlfonsSkills/AgentSync/internal/skill"
	"github.com/AlfonsSkills/AgentSync/internal/target"
)

func scanProjectSkills(targets []target.Target) []LocalSkill {
	var projectSkills []LocalSkill

	// Try to find project root (silently fail if not in a git repo)
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return nil
	}

	targetNames := map[target.Target]string{
		target.Gemini: "gemini",
		target.Claude: "claude",
		target.Codex:  "codex",
	}

	for _, t := range targets {
		targetName, ok := targetNames[t]
		if !ok {
			continue
		}

		// Get project-local skills directory
		skillsDir, err := project.GetLocalSkillsDir(targetName)
		if err != nil {
			continue
		}

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
				Target:   t,
				Valid:    valid,
				Category: fmt.Sprintf("project:%s", filepath.Base(projectRoot)),
			})
		}
	}

	return projectSkills
}
