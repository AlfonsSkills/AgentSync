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

var (
	// 全局 flags
	targetFlags []string
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "agentsync",
	Short: "Sync skills from Git repositories to local AI coding tools",
	Long: `AgentSync - Git Skill Sync Tool

Sync skills from Git repositories (default: GitHub) to local AI coding tool directories.
Supports Gemini CLI, Claude Code, and Codex CLI.

Examples:
  # Install skill to all tools
  agentsync install user/repo

  # Install to specific tool
  agentsync install user/repo --target gemini
  agentsync install user/repo -t claude,codex

  # List installed skills
  agentsync list

  # Remove skill
  agentsync remove skill-name`,
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
	rootCmd.SetVersionTemplate(fmt.Sprintf(`AgentSync %s
Git Commit: %s
Build Time: %s
`, Version, GitCommit, BuildTime))

	// 添加全局 flags
	rootCmd.PersistentFlags().StringSliceVarP(&targetFlags, "target", "t", []string{},
		"Target tools (gemini, claude, codex, antigravity), comma-separated, default: all")
}
