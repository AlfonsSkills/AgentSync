// Package git 提供 Git 仓库拉取功能
package git

import (
	"fmt"
	"regexp"
	"strings"
)

// GitHubTreeURLParser 实现 TreeURLParser 接口，用于解析 GitHub Tree URL
type GitHubTreeURLParser struct{}

// githubTreeURLRegex 匹配 GitHub tree URL 格式
// 格式: https://github.com/{owner}/{repo}/tree/{branch}/{path...}
var githubTreeURLRegex = regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+)/tree/([^/]+)/(.+)$`)

// Platform 返回平台名称
func (p *GitHubTreeURLParser) Platform() string {
	return "github"
}

// Match 检测 URL 是否为 GitHub tree 格式
func (p *GitHubTreeURLParser) Match(url string) bool {
	normalizedURL, ok := normalizeGitHubTreeURL(url)
	if !ok {
		return false
	}
	return githubTreeURLRegex.MatchString(normalizedURL)
}

// Parse 解析 GitHub tree URL
// 输入: https://github.com/AlfonsSkills/skills/tree/main/all-money-back-my-home
// 输出: &TreeURLInfo{Platform: "github", Owner: "AlfonsSkills", ...}
func (p *GitHubTreeURLParser) Parse(url string) (*TreeURLInfo, error) {
	normalizedURL, ok := normalizeGitHubTreeURL(url)
	if !ok {
		return nil, fmt.Errorf("invalid GitHub tree URL format: %s", url)
	}

	matches := githubTreeURLRegex.FindStringSubmatch(normalizedURL)
	if matches == nil {
		return nil, fmt.Errorf("invalid GitHub tree URL format: %s", url)
	}

	return &TreeURLInfo{
		Platform: p.Platform(),
		Owner:    matches[1],
		Repo:     matches[2],
		Branch:   matches[3],
		Path:     matches[4],
	}, nil
}

// normalizeGitHubTreeURL 将输入转换为标准 GitHub Tree URL
// 入参: rawURL 支持 https://github.com/..., github.com/..., owner/repo/...
// 返回: 标准化后的 https://github.com/... 或 ok=false 表示无法识别
func normalizeGitHubTreeURL(rawURL string) (string, bool) {
	normalized := strings.TrimRight(strings.TrimSpace(rawURL), "/")
	if normalized == "" {
		return "", false
	}

	// 直接处理完整 URL
	if strings.HasPrefix(normalized, "https://github.com/") {
		return normalized, true
	}
	if strings.HasPrefix(normalized, "http://github.com/") {
		return "https://" + strings.TrimPrefix(normalized, "http://"), true
	}

	// 处理无协议但含 github.com 的情况
	if strings.HasPrefix(normalized, "github.com/") {
		return "https://" + normalized, true
	}
	if strings.HasPrefix(normalized, "www.github.com/") {
		return "https://" + strings.TrimPrefix(normalized, "www."), true
	}

	// 若已包含协议但非 GitHub，则不处理
	if strings.HasPrefix(normalized, "https://") || strings.HasPrefix(normalized, "http://") || strings.HasPrefix(normalized, "git@") {
		return "", false
	}

	// 处理 owner/repo/tree/... 形式，默认 GitHub
	firstSegment := strings.SplitN(normalized, "/", 2)[0]
	if strings.Contains(firstSegment, ".") {
		// 有域名但非 github.com，视为非 GitHub
		return "", false
	}

	return "https://github.com/" + normalized, true
}

// --- 以下为向后兼容的函数（已废弃，将在未来版本移除）---

// GitHubTreeURL 已废弃，请使用 TreeURLInfo
// Deprecated: Use TreeURLInfo instead
type GitHubTreeURL = TreeURLInfo

// IsGitHubTreeURL 已废弃，请使用 IsTreeURL
// Deprecated: Use IsTreeURL instead
func IsGitHubTreeURL(url string) bool {
	return (&GitHubTreeURLParser{}).Match(url)
}

// ParseGitHubTreeURL 已废弃，请使用 ParseTreeURL
// Deprecated: Use ParseTreeURL instead
func ParseGitHubTreeURL(url string) (*TreeURLInfo, error) {
	return (&GitHubTreeURLParser{}).Parse(url)
}
