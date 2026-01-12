// Package project provides utilities for project root detection
package project

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindProjectRoot searches for project root by looking for .git directory
// Returns the absolute path to project root, or error if not found
func FindProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	dir := currentDir
	for {
		// Check if .git exists in current directory
		gitPath := filepath.Join(dir, ".git")
		if info, err := os.Stat(gitPath); err == nil && info.IsDir() {
			return dir, nil
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root without finding .git
			return "", fmt.Errorf("not in a git repository (no .git directory found)")
		}
		dir = parent
	}
}

// GetLocalSkillsDir returns the project-level skills directory for a target
// targetName should be "gemini", "claude", or "codex"
func GetLocalSkillsDir(targetName string) (string, error) {
	root, err := FindProjectRoot()
	if err != nil {
		return "", err
	}

	// Return .{target}/skills/ path
	return filepath.Join(root, fmt.Sprintf(".%s", targetName), "skills"), nil
}
