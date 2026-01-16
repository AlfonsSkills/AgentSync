// Package git 提供 Git 仓库拉取功能
package git

import "fmt"

// TreeURLInfo 表示仓库子目录 URL 的通用解析结果
// 用于存储从各平台（GitHub, GitLab 等）解析的 Tree URL 信息
type TreeURLInfo struct {
	Platform string // 平台标识 (e.g., "github", "gitlab")
	Owner    string // 仓库所有者
	Repo     string // 仓库名
	Branch   string // 分支名
	Path     string // 子目录路径，支持多级 (e.g., "category/skill-name")
}

// CloneURL 返回用于 git clone 的标准 URL
func (t *TreeURLInfo) CloneURL() string {
	switch t.Platform {
	case "github":
		return fmt.Sprintf("https://github.com/%s/%s.git", t.Owner, t.Repo)
	case "gitlab":
		return fmt.Sprintf("https://gitlab.com/%s/%s.git", t.Owner, t.Repo)
	default:
		return fmt.Sprintf("https://%s/%s/%s.git", t.Platform, t.Owner, t.Repo)
	}
}

// RepoSlug 返回 owner/repo 格式的仓库标识
func (t *TreeURLInfo) RepoSlug() string {
	return fmt.Sprintf("%s/%s", t.Owner, t.Repo)
}

// TreeURLParser 定义仓库 Tree URL 解析器接口
// 每个平台（GitHub, GitLab 等）实现此接口以支持其特定的 URL 格式
type TreeURLParser interface {
	// Match 检测 URL 是否匹配该解析器的格式
	Match(url string) bool
	// Parse 解析 URL 返回 TreeURLInfo
	Parse(url string) (*TreeURLInfo, error)
	// Platform 返回平台名称标识
	Platform() string
}

// 注册所有支持的解析器
var parsers = []TreeURLParser{
	&GitHubTreeURLParser{}, // GitHub 解析器
	// 未来可添加更多解析器:
	// &GitLabTreeURLParser{},
	// &BitbucketTreeURLParser{},
}

// IsTreeURL 检测 URL 是否为任意支持平台的 Tree URL
func IsTreeURL(url string) bool {
	for _, p := range parsers {
		if p.Match(url) {
			return true
		}
	}
	return false
}

// ParseTreeURL 工厂方法，自动选择匹配的解析器解析 URL
func ParseTreeURL(url string) (*TreeURLInfo, error) {
	for _, p := range parsers {
		if p.Match(url) {
			return p.Parse(url)
		}
	}
	return nil, fmt.Errorf("unsupported tree URL format: %s", url)
}
