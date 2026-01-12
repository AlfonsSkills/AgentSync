// Package target 管理目标工具的路径配置
package target

import (
	"fmt"
	"os"
	"path/filepath"
)

// Target 表示支持的 AI 编码工具
type Target string

const (
	Gemini Target = "gemini"
	Claude Target = "claude"
	Codex  Target = "codex"
)

// AllTargets 返回所有支持的目标工具
func AllTargets() []Target {
	return []Target{Gemini, Claude, Codex}
}

// ParseTargets 解析目标字符串列表，返回有效的 Target 切片
// 如果输入为空，返回所有目标
func ParseTargets(targets []string) ([]Target, error) {
	if len(targets) == 0 {
		return AllTargets(), nil
	}

	var result []Target
	for _, t := range targets {
		target := Target(t)
		if !target.IsValid() {
			return nil, fmt.Errorf("invalid target: %s, valid targets are: gemini, claude, codex", t)
		}
		result = append(result, target)
	}
	return result, nil
}

// IsValid 检查目标是否有效
func (t Target) IsValid() bool {
	switch t {
	case Gemini, Claude, Codex:
		return true
	default:
		return false
	}
}

// SkillsDir 返回目标工具的 skills 目录路径（用于扫描）
func (t Target) SkillsDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	var dir string
	switch t {
	case Gemini:
		dir = filepath.Join(homeDir, ".gemini", "skills")
	case Claude:
		dir = filepath.Join(homeDir, ".claude", "skills")
	case Codex:
		dir = filepath.Join(homeDir, ".codex", "skills")
	default:
		return "", fmt.Errorf("unknown target: %s", t)
	}

	return dir, nil
}

// InstallDir 返回目标工具安装 skill 的目录路径
// Codex 使用 public/ 子目录，其他工具直接使用 skills/
func (t Target) InstallDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	var dir string
	switch t {
	case Gemini:
		dir = filepath.Join(homeDir, ".gemini", "skills")
	case Claude:
		dir = filepath.Join(homeDir, ".claude", "skills")
	case Codex:
		// Codex 用户 skills 安装到 public/ 子目录
		dir = filepath.Join(homeDir, ".codex", "skills", "public")
	default:
		return "", fmt.Errorf("unknown target: %s", t)
	}

	return dir, nil
}

// EnsureInstallDir 确保安装目录存在
func (t Target) EnsureInstallDir() (string, error) {
	dir, err := t.InstallDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create install directory %s: %w", dir, err)
	}

	return dir, nil
}

// EnsureSkillsDir 确保 skills 目录存在（向后兼容）
func (t Target) EnsureSkillsDir() (string, error) {
	return t.EnsureInstallDir()
}

// String 返回目标的字符串表示
func (t Target) String() string {
	return string(t)
}

// DisplayName 返回目标的显示名称
func (t Target) DisplayName() string {
	switch t {
	case Gemini:
		return "Gemini CLI"
	case Claude:
		return "Claude Code"
	case Codex:
		return "Codex CLI"
	default:
		return string(t)
	}
}
