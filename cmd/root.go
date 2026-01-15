// Package cmd 实现 CLI 命令
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/AlfonsSkills/SkillSync/internal/updater"
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
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// 跳过 upgrade、help 和 version 命令的更新检查
		if cmd.Name() == "upgrade" || cmd.Name() == "help" || cmd.Name() == "skillsync" {
			return
		}
		// 检查更新（带超时）
		checkUpdateInBackground()
	},
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

// checkUpdateInBackground 检查更新（带超时）
func checkUpdateInBackground() {
	// dev 版本不检查更新
	if Version == "dev" {
		return
	}

	// 设置超时，避免阻塞太久
	done := make(chan struct{})
	go func() {
		defer close(done)

		result, err := updater.CheckLatestVersion(Version)
		if err != nil {
			return // 静默失败，不影响用户体验
		}

		if !result.IsLatest {
			fmt.Printf("\n⚠ A new version is available: %s (current: %s)\n", result.LatestVersion, result.CurrentVersion)
			fmt.Println("  Run: skillsync upgrade")
		}
	}()

	// 最多等待 3 秒
	select {
	case <-done:
	case <-time.After(3 * time.Second):
	}
}
