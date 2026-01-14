package target

import (
	"os"
	"path/filepath"
)

// copilotProvider 实现 GitHub Copilot / VSCode 的 ToolProvider 接口
// GitHub Copilot 使用以下目录结构：
// - 全局 Skills: ~/.copilot/skills/
// - 项目级 Skills: .github/skills/
type copilotProvider struct {
	homeDir string
}

// NewCopilotProvider 创建 GitHub Copilot Provider 实例
func NewCopilotProvider() ToolProvider {
	homeDir, _ := os.UserHomeDir()
	return &copilotProvider{homeDir: homeDir}
}

// Type 返回工具类型枚举
func (c *copilotProvider) Type() ToolType {
	return ToolCopilot
}

// DisplayName 返回用户可见名称
func (c *copilotProvider) DisplayName() string {
	return "GitHub Copilot / VSCode"
}

// GlobalSkillsDir 返回全局 skills 扫描目录
// GitHub Copilot 使用 ~/.copilot/skills/
func (c *copilotProvider) GlobalSkillsDir() (string, error) {
	return filepath.Join(c.homeDir, ".copilot", "skills"), nil
}

// GlobalInstallDir 返回全局安装目录（与扫描目录相同）
func (c *copilotProvider) GlobalInstallDir() (string, error) {
	return c.GlobalSkillsDir()
}

// LocalSkillsDir 返回项目级 skills 目录
// GitHub Copilot 使用 .github/skills/
func (c *copilotProvider) LocalSkillsDir(projectRoot string) string {
	return filepath.Join(projectRoot, ".github", "skills")
}

// Categories 返回分类子目录列表（Copilot 无分类）
func (c *copilotProvider) Categories() []string {
	return nil
}

// EnsureInstallDir 确保全局安装目录存在
func (c *copilotProvider) EnsureInstallDir() (string, error) {
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
func (c *copilotProvider) EnsureLocalInstallDir(projectRoot string) (string, error) {
	dir := c.LocalSkillsDir(projectRoot)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}
