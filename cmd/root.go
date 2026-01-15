// Package cmd 实现 CLI 命令
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// 版本信息（通过 ldflags 注入）
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

// 项目元数据
const (
	Author     = "Alfons <alfonsxh@gmail.com>"
	ProjectURL = "https://github.com/AlfonsSkills/SkillSync"
)

var (
	// 全局 flags
	targetFlags []string
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "skillsync",
	Short: "Sync skills from Git repositories to local AI coding tools",
	Long: `SkillSync - Git Skill Sync Tool

Sync skills from Git repositories (default: GitHub) to local AI coding tool directories.
Supports Gemini CLI, Claude Code, and Codex CLI.

Examples:
  # Install skill to all tools
  skillsync install user/repo

  # Install to specific tool
  skillsync install user/repo --target gemini
  skillsync install user/repo -t claude,codex

  # List installed skills
  skillsync list

  # Remove skill
  skillsync remove skill-name`,
	Version: Version,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// 设置版本模板
	rootCmd.SetVersionTemplate(fmt.Sprintf(`SkillSync %s
Git Commit: %s
Build Time: %s

Author:  %s
Project: %s
`, Version, GitCommit, BuildTime, Author, ProjectURL))

	// 添加全局 flags
	rootCmd.PersistentFlags().StringSliceVarP(&targetFlags, "target", "t", []string{},
		"Target tools (gemini, claude, codex, opencode, goose, crush, antigravity, copilot, cursor, cline, droid, kilocode, roocode, vscode), comma-separated, default: all")
}
