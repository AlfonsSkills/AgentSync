---
name: add-tool-provider
description: 指导如何为 AgentSync 添加新的 AI 工具 Provider。当需要扩展 AgentSync 以支持新的 IDE，如 Antigravity、Cursor 等时使用。
---

# 添加工具 Provider 技能

本技能指导如何为 AgentSync 添加新的 AI 编码工具支持。

## 使用场景

- 添加新的 AI IDE 或 CLI 工具支持
- 理解 ToolProvider 接口模式
- 扩展 AgentSync 的多工具架构

## 架构概览

```
┌─────────────────────────────────────────────┐
│           ToolProvider 接口                  │
├───────────┬───────────┬───────────┬─────────┤
│  Gemini   │  Claude   │  Codex    │   新    │
│  Provider │  Provider │  Provider │ Provider│
└───────────┴───────────┴───────────┴─────────┘
```

## 实现步骤

### 1. 添加类型常量

编辑 `internal/target/provider.go`：

```go
const (
    ToolGemini      ToolType = "gemini"
    ToolClaude      ToolType = "claude"
    ToolCodex       ToolType = "codex"
    ToolNewTool     ToolType = "newtool"  // 添加新常量
)
```

### 2. 创建 Provider 文件

创建 `internal/target/newtool.go`，实现 `ToolProvider` 接口：

```go
package target

import (
    "os"
    "path/filepath"
)

type newtoolProvider struct {
    homeDir string
}

func NewNewtoolProvider() ToolProvider {
    homeDir, _ := os.UserHomeDir()
    return &newtoolProvider{homeDir: homeDir}
}

func (n *newtoolProvider) Type() ToolType           { return ToolNewTool }
func (n *newtoolProvider) DisplayName() string      { return "NewTool IDE" }
func (n *newtoolProvider) GlobalSkillsDir() (string, error) {
    return filepath.Join(n.homeDir, ".newtool", "skills"), nil
}
func (n *newtoolProvider) GlobalInstallDir() (string, error) {
    return n.GlobalSkillsDir()
}
func (n *newtoolProvider) LocalSkillsDir(projectRoot string) string {
    return filepath.Join(projectRoot, ".newtool", "skills")
}
func (n *newtoolProvider) Categories() []string { return nil }
// ... 实现 EnsureInstallDir 和 EnsureLocalInstallDir
```

### 3. 注册 Provider

编辑 `internal/target/registry.go`：

```go
providers = map[ToolType]ToolProvider{
    ToolGemini:  NewGeminiProvider(),
    ToolClaude:  NewClaudeProvider(),
    ToolCodex:   NewCodexProvider(),
    ToolNewTool: NewNewtoolProvider(),  // 注册
}

func AllToolTypes() []ToolType {
    return []ToolType{ToolGemini, ToolClaude, ToolCodex, ToolNewTool}
}
```

### 4. 更新 CLI 帮助

编辑 `cmd/root.go`，更新 `--target` 参数说明。

### 5. 验证

```bash
make build
./build/agentsync list --target=newtool
```


## 变更文件

| 文件 | 变更 |
|:-----|:-----|
| `internal/target/antigravity.go` | 新增 Provider 实现 |
| `internal/target/provider.go` | 添加 `ToolAntigravity` 常量 |
| `internal/target/registry.go` | 注册 Provider |
| `cmd/root.go` | 更新帮助文本 |
