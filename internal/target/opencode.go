package target

import (
	"os"
	"path/filepath"
)

// opencodeProvider 实现 OpenCode 的 ToolProvider 接口
// OpenCode 使用以下目录结构（注意：使用 skill 单数而非 skills）：
// - 全局 Skills: ~/.config/opencode/skill/
// - 项目级 Skills: .opencode/skill/
type opencodeProvider struct {
	homeDir string
}

// NewOpencodeProvider 创建 OpenCode Provider 实例
func NewOpencodeProvider() ToolProvider {
	homeDir, _ := os.UserHomeDir()
	return &opencodeProvider{homeDir: homeDir}
}

// Type 返回工具类型枚举
func (o *opencodeProvider) Type() ToolType {
	return ToolOpencode
}

// DisplayName 返回用户可见名称
func (o *opencodeProvider) DisplayName() string {
	return "OpenCode"
}

// GlobalSkillsDir 返回全局 skills 扫描目录
// OpenCode 使用 ~/.config/opencode/skill/（单数形式）
func (o *opencodeProvider) GlobalSkillsDir() (string, error) {
	return filepath.Join(o.homeDir, ".config", "opencode", "skill"), nil
}

// GlobalInstallDir 返回全局安装目录（与扫描目录相同）
func (o *opencodeProvider) GlobalInstallDir() (string, error) {
	return o.GlobalSkillsDir()
}

// LocalSkillsDir 返回项目级 skills 目录
// OpenCode 使用 .opencode/skill/（单数形式）
func (o *opencodeProvider) LocalSkillsDir(projectRoot string) string {
	return filepath.Join(projectRoot, ".opencode", "skill")
}

// Categories 返回分类子目录列表（OpenCode 无分类）
func (o *opencodeProvider) Categories() []string {
	return nil
}

// EnsureInstallDir 确保全局安装目录存在
func (o *opencodeProvider) EnsureInstallDir() (string, error) {
	dir, err := o.GlobalInstallDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

// EnsureLocalInstallDir 确保项目级安装目录存在
func (o *opencodeProvider) EnsureLocalInstallDir(projectRoot string) (string, error) {
	dir := o.LocalSkillsDir(projectRoot)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}
