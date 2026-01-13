package target

import (
	"os"
	"path/filepath"
)

// geminiProvider 实现 Gemini CLI 的 ToolProvider 接口
type geminiProvider struct {
	homeDir string
}

// NewGeminiProvider 创建 Gemini CLI Provider 实例
func NewGeminiProvider() ToolProvider {
	homeDir, _ := os.UserHomeDir()
	return &geminiProvider{homeDir: homeDir}
}

// Type 返回工具类型枚举
func (g *geminiProvider) Type() ToolType {
	return ToolGemini
}

// DisplayName 返回用户可见名称
func (g *geminiProvider) DisplayName() string {
	return "Gemini CLI"
}

// GlobalSkillsDir 返回全局 skills 扫描目录
func (g *geminiProvider) GlobalSkillsDir() (string, error) {
	return filepath.Join(g.homeDir, ".gemini", "skills"), nil
}

// GlobalInstallDir 返回全局安装目录（与扫描目录相同）
func (g *geminiProvider) GlobalInstallDir() (string, error) {
	return g.GlobalSkillsDir()
}

// LocalSkillsDir 返回项目级 skills 目录
func (g *geminiProvider) LocalSkillsDir(projectRoot string) string {
	return filepath.Join(projectRoot, ".gemini", "skills")
}

// Categories 返回分类子目录列表（Gemini 无分类）
func (g *geminiProvider) Categories() []string {
	return nil
}

// EnsureInstallDir 确保全局安装目录存在
func (g *geminiProvider) EnsureInstallDir() (string, error) {
	dir, err := g.GlobalInstallDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

// EnsureLocalInstallDir 确保项目级安装目录存在
func (g *geminiProvider) EnsureLocalInstallDir(projectRoot string) (string, error) {
	dir := g.LocalSkillsDir(projectRoot)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}
