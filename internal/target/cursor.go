package target

import (
	"os"
	"path/filepath"
)

// cursorProvider 实现 Cursor IDE 的 ToolProvider 接口
// Cursor 使用以下目录结构：
// - 全局 Skills: ~/.cursor/skills/
// - 项目级 Skills: .cursor/skills/
type cursorProvider struct {
	homeDir string
}

// NewCursorProvider 创建 Cursor IDE Provider 实例
func NewCursorProvider() ToolProvider {
	homeDir, _ := os.UserHomeDir()
	return &cursorProvider{homeDir: homeDir}
}

// Type 返回工具类型枚举
func (c *cursorProvider) Type() ToolType {
	return ToolCursor
}

// DisplayName 返回用户可见名称
func (c *cursorProvider) DisplayName() string {
	return "Cursor IDE"
}

// GlobalSkillsDir 返回全局 skills 扫描目录
// Cursor 使用 ~/.cursor/skills/
func (c *cursorProvider) GlobalSkillsDir() (string, error) {
	return filepath.Join(c.homeDir, ".cursor", "skills"), nil
}

// GlobalInstallDir 返回全局安装目录（与扫描目录相同）
func (c *cursorProvider) GlobalInstallDir() (string, error) {
	return c.GlobalSkillsDir()
}

// LocalSkillsDir 返回项目级 skills 目录
// Cursor 使用 .cursor/skills/
func (c *cursorProvider) LocalSkillsDir(projectRoot string) string {
	return filepath.Join(projectRoot, ".cursor", "skills")
}

// Categories 返回分类子目录列表（Cursor 无分类）
func (c *cursorProvider) Categories() []string {
	return nil
}

// EnsureInstallDir 确保全局安装目录存在
func (c *cursorProvider) EnsureInstallDir() (string, error) {
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
func (c *cursorProvider) EnsureLocalInstallDir(projectRoot string) (string, error) {
	dir := c.LocalSkillsDir(projectRoot)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}
