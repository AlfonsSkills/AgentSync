package target

import (
	"fmt"
	"sync"
)

// registry 存储所有已注册的 Provider
var (
	providers     map[ToolType]ToolProvider
	providersOnce sync.Once
)

// initProviders 初始化 Provider 注册表（懒加载）
func initProviders() {
	providersOnce.Do(func() {
		providers = map[ToolType]ToolProvider{
			ToolGemini:      NewGeminiProvider(),
			ToolClaude:      NewClaudeProvider(),
			ToolCodex:       NewCodexProvider(),
			ToolAntigravity: NewAntigravityProvider(),
		}
	})
}

// AllToolTypes 返回所有支持的工具类型（有序）
func AllToolTypes() []ToolType {
	return []ToolType{ToolGemini, ToolClaude, ToolCodex, ToolAntigravity}
}

// AllProviders 返回所有已注册的 Provider 列表
func AllProviders() []ToolProvider {
	initProviders()
	result := make([]ToolProvider, 0, len(providers))
	// 保持固定顺序
	for _, t := range AllToolTypes() {
		if p, ok := providers[t]; ok {
			result = append(result, p)
		}
	}
	return result
}

// GetProvider 根据工具类型获取对应的 Provider
func GetProvider(toolType ToolType) (ToolProvider, error) {
	initProviders()
	if p, ok := providers[toolType]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("unknown provider: %s, valid providers are: gemini, claude, codex, antigravity", toolType)
}

// GetProviderByName 根据名称字符串获取对应的 Provider
func GetProviderByName(name string) (ToolProvider, error) {
	return GetProvider(ToolType(name))
}

// ParseProviders 解析 Provider 名称列表，返回对应的 Provider 切片
// 如果输入为空，返回所有 Provider
func ParseProviders(names []string) ([]ToolProvider, error) {
	if len(names) == 0 {
		return AllProviders(), nil
	}

	result := make([]ToolProvider, 0, len(names))
	for _, name := range names {
		p, err := GetProviderByName(name)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}
