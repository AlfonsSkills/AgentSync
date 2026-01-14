// Package target 管理目标工具的路径配置和 Provider 接口
package target

// ToolType 表示支持的 AI 编码工具类型（枚举）
type ToolType string

const (
	ToolGemini      ToolType = "gemini"
	ToolClaude      ToolType = "claude"
	ToolCodex       ToolType = "codex"
	ToolAntigravity ToolType = "antigravity"
	ToolCopilot     ToolType = "copilot"
	ToolCursor      ToolType = "cursor"
	ToolOpencode    ToolType = "opencode"
)

// String 返回工具类型的字符串表示
func (t ToolType) String() string {
	return string(t)
}

// ToolProvider 定义 AI 编码工具的通用行为接口
// 所有支持的工具（Gemini/Claude/Codex）都必须实现此接口
type ToolProvider interface {
	// Type 返回工具的类型枚举
	Type() ToolType

	// DisplayName 返回用户可见的工具名称
	DisplayName() string

	// GlobalSkillsDir 返回全局 skills 扫描目录路径
	// 用于 list 命令扫描已安装的 skills
	GlobalSkillsDir() (string, error)

	// GlobalInstallDir 返回全局安装目录路径
	// 用于 install 命令安装 skills
	// 注意：某些工具（如 Codex）的安装目录可能与扫描目录不同
	GlobalInstallDir() (string, error)

	// LocalSkillsDir 返回项目级 skills 目录路径
	// projectRoot 为项目根目录（包含 .git 的目录）
	LocalSkillsDir(projectRoot string) string

	// Categories 返回该工具支持的分类子目录列表
	// 例如 Codex 支持 ["public", ".system"]，其他工具返回 nil
	Categories() []string

	// EnsureInstallDir 确保安装目录存在，如不存在则创建
	// 返回创建的目录路径
	EnsureInstallDir() (string, error)

	// EnsureLocalInstallDir 确保项目级安装目录存在
	// projectRoot 为项目根目录
	EnsureLocalInstallDir(projectRoot string) (string, error)
}
