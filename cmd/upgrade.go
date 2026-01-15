// Package cmd 实现 CLI 命令
package cmd

import (
	"fmt"
	"os"

	"github.com/AlfonsSkills/SkillSync/internal/updater"
	"github.com/spf13/cobra"
)

var checkOnly bool

// upgradeCmd upgrade 命令
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade SkillSync to the latest version",
	Long: `Check for updates and upgrade SkillSync to the latest version.

Examples:
  # Check for updates only
  skillsync upgrade --check

  # Upgrade to latest version
  skillsync upgrade`,
	Run: func(cmd *cobra.Command, args []string) {
		if checkOnly {
			// 仅检查更新
			result, err := updater.CheckLatestVersion(Version)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ %v\n", err)
				os.Exit(1)
			}

			if result.IsLatest {
				fmt.Println("✓ You are using the latest version")
			} else {
				fmt.Printf("⚠ A new version is available: %s\n", result.LatestVersion)
				fmt.Printf("  Current: %s\n", result.CurrentVersion)
				fmt.Printf("  Run: skillsync upgrade\n")
			}
			return
		}

		// 执行升级
		if err := updater.Upgrade(Version); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	upgradeCmd.Flags().BoolVarP(&checkOnly, "check", "c", false, "Only check for updates, don't upgrade")
	rootCmd.AddCommand(upgradeCmd)
}
