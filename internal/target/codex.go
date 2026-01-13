package target

import (
	"os"
	"path/filepath"
)

// codexProvider 实现 Codex CLI 的 ToolProvider 接口
// Codex 有特殊的目录结构：安装到 public/ 子目录，同时支持 .system/ 分类
type codexProvider struct {
	homeDir string
}

// NewCodexProvider 创建 Codex CLI Provider 实例
func NewCodexProvider() ToolProvider {
	homeDir, _ := os.UserHomeDir()
	return &codexProvider{homeDir: homeDir}
}

// Type 返回工具类型枚举
func (c *codexProvider) Type() ToolType {
	return ToolCodex
}

// DisplayName 返回用户可见名称
func (c *codexProvider) DisplayName() string {
	return "Codex CLI"
}

// GlobalSkillsDir 返回全局 skills 扫描目录
// 返回根目录，扫描时会递归查找 public/ 和 .system/ 子目录
func (c *codexProvider) GlobalSkillsDir() (string, error) {
	return filepath.Join(c.homeDir, ".codex", "skills"), nil
}

// GlobalInstallDir 返回全局安装目录
// Codex 用户 skills 默认安装到 public/ 子目录
func (c *codexProvider) GlobalInstallDir() (string, error) {
	return filepath.Join(c.homeDir, ".codex", "skills", "public"), nil
}

// LocalSkillsDir 返回项目级 skills 目录
func (c *codexProvider) LocalSkillsDir(projectRoot string) string {
	return filepath.Join(projectRoot, ".codex", "skills")
}

// Categories 返回 Codex 特有的分类子目录
func (c *codexProvider) Categories() []string {
	return []string{"public", ".system"}
}

// EnsureInstallDir 确保全局安装目录存在
func (c *codexProvider) EnsureInstallDir() (string, error) {
	dir, err := c.GlobalInstallDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

// EnsureLocalInstallDir 确保项目级安装目录存在
func (c *codexProvider) EnsureLocalInstallDir(projectRoot string) (string, error) {
	dir := c.LocalSkillsDir(projectRoot)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}
