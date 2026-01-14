package target

import (
	"os"
	"path/filepath"
)

// antigravityProvider 实现 Antigravity IDE 的 ToolProvider 接口
// Antigravity 使用独特的目录结构：
// - 全局 Skills: ~/.gemini/antigravity/skills/
// - 项目级 Skills: .agent/skills/
type antigravityProvider struct {
	homeDir string
}

// NewAntigravityProvider 创建 Antigravity IDE Provider 实例
func NewAntigravityProvider() ToolProvider {
	homeDir, _ := os.UserHomeDir()
	return &antigravityProvider{homeDir: homeDir}
}

// Type 返回工具类型枚举
func (a *antigravityProvider) Type() ToolType {
	return ToolAntigravity
}

// DisplayName 返回用户可见名称
func (a *antigravityProvider) DisplayName() string {
	return "Antigravity IDE"
}

// GlobalSkillsDir 返回全局 skills 扫描目录
// Antigravity 使用 ~/.gemini/antigravity/skills/
func (a *antigravityProvider) GlobalSkillsDir() (string, error) {
	return filepath.Join(a.homeDir, ".gemini", "antigravity", "skills"), nil
}

// GlobalInstallDir 返回全局安装目录（与扫描目录相同）
func (a *antigravityProvider) GlobalInstallDir() (string, error) {
	return a.GlobalSkillsDir()
}

// LocalSkillsDir 返回项目级 skills 目录
// Antigravity 使用 .agent/skills/ 而非 .antigravity/skills/
func (a *antigravityProvider) LocalSkillsDir(projectRoot string) string {
	return filepath.Join(projectRoot, ".agent", "skills")
}

// Categories 返回分类子目录列表（Antigravity 无分类）
func (a *antigravityProvider) Categories() []string {
	return nil
}

// EnsureInstallDir 确保全局安装目录存在
func (a *antigravityProvider) EnsureInstallDir() (string, error) {
	dir, err := a.GlobalInstallDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

// EnsureLocalInstallDir 确保项目级安装目录存在
func (a *antigravityProvider) EnsureLocalInstallDir(projectRoot string) (string, error) {
	dir := a.LocalSkillsDir(projectRoot)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}
