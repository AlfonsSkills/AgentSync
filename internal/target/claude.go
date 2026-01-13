package target

import (
	"os"
	"path/filepath"
)

// claudeProvider 实现 Claude Code 的 ToolProvider 接口
type claudeProvider struct {
	homeDir string
}

// NewClaudeProvider 创建 Claude Code Provider 实例
func NewClaudeProvider() ToolProvider {
	homeDir, _ := os.UserHomeDir()
	return &claudeProvider{homeDir: homeDir}
}

// Type 返回工具类型枚举
func (c *claudeProvider) Type() ToolType {
	return ToolClaude
}

// DisplayName 返回用户可见名称
func (c *claudeProvider) DisplayName() string {
	return "Claude Code"
}

// GlobalSkillsDir 返回全局 skills 扫描目录
func (c *claudeProvider) GlobalSkillsDir() (string, error) {
	return filepath.Join(c.homeDir, ".claude", "skills"), nil
}

// GlobalInstallDir 返回全局安装目录（与扫描目录相同）
func (c *claudeProvider) GlobalInstallDir() (string, error) {
	return c.GlobalSkillsDir()
}

// LocalSkillsDir 返回项目级 skills 目录
func (c *claudeProvider) LocalSkillsDir(projectRoot string) string {
	return filepath.Join(projectRoot, ".claude", "skills")
}

// Categories 返回分类子目录列表（Claude 无分类）
func (c *claudeProvider) Categories() []string {
	return nil
}

// EnsureInstallDir 确保全局安装目录存在
func (c *claudeProvider) EnsureInstallDir() (string, error) {
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
func (c *claudeProvider) EnsureLocalInstallDir(projectRoot string) (string, error) {
	dir := c.LocalSkillsDir(projectRoot)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}
