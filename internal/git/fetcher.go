// Package git 提供 Git 仓库拉取功能
package git

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Fetcher 用于拉取 Git 仓库
type Fetcher struct {
	// DefaultHost 默认的 Git 托管平台
	DefaultHost string
}

// NewFetcher 创建一个新的 Fetcher
func NewFetcher() *Fetcher {
	return &Fetcher{
		DefaultHost: "github.com",
	}
}

// RepoKey 生成仓库的规范化标识，用于缓存定位。
// 入参：source（支持短格式、HTTPS/SSH URL、Tree URL）。
// 返回：host/owner/repo 格式的规范化 key；若无法解析则返回 error。
// 规则：SSH 与 HTTPS 的同仓库会映射到相同 key。
func (f *Fetcher) RepoKey(source string) (string, error) {
	normalized := strings.TrimSpace(source)
	normalized = strings.TrimRight(normalized, "/")
	if normalized == "" {
		return "", fmt.Errorf("empty repository source")
	}

	// 关键步骤：若为 Tree URL，先转换为 Clone URL
	if IsTreeURL(normalized) {
		treeURL, err := ParseTreeURL(normalized)
		if err != nil {
			return "", err
		}
		normalized = treeURL.CloneURL()
	}

	// 处理 scp-like SSH：git@host:owner/repo(.git)
	if strings.HasPrefix(normalized, "git@") && strings.Contains(normalized, ":") && !strings.Contains(normalized, "://") {
		atIndex := strings.Index(normalized, "@")
		colonIndex := strings.Index(normalized, ":")
		if atIndex >= 0 && colonIndex > atIndex+1 {
			host := normalized[atIndex+1 : colonIndex]
			pathPart := normalized[colonIndex+1:]
			return buildRepoKey(host, pathPart)
		}
		return "", fmt.Errorf("invalid ssh repository format: %s", source)
	}

	// 处理带协议的 URL（https/http/ssh 等）
	if strings.Contains(normalized, "://") {
		parsed, err := url.Parse(normalized)
		if err != nil {
			return "", fmt.Errorf("invalid repository url: %w", err)
		}
		host := parsed.Host
		pathPart := strings.TrimPrefix(parsed.Path, "/")
		return buildRepoKey(host, pathPart)
	}

	// 处理无协议格式：github.com/owner/repo 或 owner/repo
	parts := strings.Split(normalized, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid repository source: %s", source)
	}

	// 若包含域名，则第一段为 host
	if strings.Contains(parts[0], ".") {
		host := parts[0]
		pathPart := strings.Join(parts[1:], "/")
		return buildRepoKey(host, pathPart)
	}

	// 短格式：owner/repo
	host := f.DefaultHost
	pathPart := strings.Join(parts, "/")
	return buildRepoKey(host, pathPart)
}

// buildRepoKey 解析 owner/repo，并生成 host/owner/repo 规范化标识。
// 入参：host（可带端口），pathPart（owner/repo 或 owner/repo/...）。
// 返回：规范化 key 或 error。
func buildRepoKey(host, pathPart string) (string, error) {
	host = strings.TrimSpace(host)
	if host == "" {
		return "", fmt.Errorf("invalid repository host")
	}

	pathPart = strings.Trim(pathPart, "/")
	pathPart = strings.TrimSuffix(pathPart, ".git")
	if pathPart == "" {
		return "", fmt.Errorf("invalid repository path")
	}

	segments := strings.Split(pathPart, "/")
	if len(segments) < 2 {
		return "", fmt.Errorf("invalid repository path: %s", pathPart)
	}

	owner := strings.TrimSpace(segments[0])
	repo := strings.TrimSpace(segments[1])
	if owner == "" || repo == "" {
		return "", fmt.Errorf("invalid repository path: %s", pathPart)
	}

	// 关键步骤：统一大小写，确保 SSH/HTTPS/短格式的 key 一致
	return fmt.Sprintf("%s/%s/%s",
		strings.ToLower(host),
		strings.ToLower(owner),
		strings.ToLower(repo),
	), nil
}

// cacheRepoPath 返回缓存仓库的固定路径（基于 RepoKey 的哈希）。
// 入参：source（仓库输入）。
// 返回：缓存路径或 error。
func (f *Fetcher) cacheRepoPath(source string) (string, error) {
	repoKey, err := f.RepoKey(source)
	if err != nil {
		return "", err
	}

	cacheRoot := filepath.Join(os.TempDir(), "skillsync-cache", "git")
	hash := sha256.Sum256([]byte(repoKey))
	cacheDir := hex.EncodeToString(hash[:]) + ".git"
	return filepath.Join(cacheRoot, cacheDir), nil
}

// ensureCacheRepo 确保缓存仓库存在并更新到最新。
// 入参：source（仓库输入）。
// 返回：缓存仓库路径或 error。
func (f *Fetcher) ensureCacheRepo(source string) (string, error) {
	cachePath, err := f.cacheRepoPath(source)
	if err != nil {
		return "", err
	}

	// 关键步骤：创建缓存根目录
	if err := os.MkdirAll(filepath.Dir(cachePath), 0o755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	if _, statErr := os.Stat(cachePath); os.IsNotExist(statErr) {
		// 首次：使用 --mirror 创建缓存仓库
		url := f.NormalizeURL(source)
		cmd := exec.Command("git", "clone", "--mirror", url, cachePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to create cache repository: %w", err)
		}
		return cachePath, nil
	}

	// 已存在：确保 remote 使用当前输入的 URL（便于 SSH/HTTPS 统一）
	url := f.NormalizeURL(source)
	setURLCmd := exec.Command("git", "-C", cachePath, "remote", "set-url", "origin", url)
	setURLCmd.Stdout = os.Stdout
	setURLCmd.Stderr = os.Stderr
	if err := setURLCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to update cache remote url: %w", err)
	}

	// 已存在：拉取更新
	cmd := exec.Command("git", "-C", cachePath, "fetch", "--prune", "--tags")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to update cache repository: %w", err)
	}

	return cachePath, nil
}

// cloneFromCache 使用本地缓存仓库创建临时工作区。
// 入参：cachePath（缓存仓库路径）、destDir（目标目录）、branch（可选分支）。
// 返回：error。
func (f *Fetcher) cloneFromCache(cachePath, destDir, branch string) error {
	args := []string{"clone", "--shared"}
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	args = append(args, cachePath, destDir)

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone from cache: %w", err)
	}
	return nil
}

// NormalizeURL 将简短格式的仓库地址转换为完整的 Git URL
// 支持格式：
//   - user/repo -> https://github.com/user/repo.git
//   - https://github.com/user/repo -> https://github.com/user/repo.git
//   - git@github.com:user/repo.git -> git@github.com:user/repo.git
func (f *Fetcher) NormalizeURL(source string) string {
	// 去除末尾的斜杠
	source = strings.TrimRight(source, "/")

	// 如果已经是完整的 URL，直接返回
	if strings.HasPrefix(source, "https://") || strings.HasPrefix(source, "git@") {
		// 确保有 .git 后缀
		if !strings.HasSuffix(source, ".git") {
			source += ".git"
		}
		return source
	}

	// 简短格式：user/repo
	return fmt.Sprintf("https://%s/%s.git", f.DefaultHost, source)
}

// Clone 使用系统 git 命令克隆仓库到指定目录
// 使用系统 git 以继承用户的代理配置 (http.proxy / https.proxy)
func (f *Fetcher) Clone(source, destDir string) error {
	url := f.NormalizeURL(source)

	// 确保目标目录不存在
	if _, err := os.Stat(destDir); !os.IsNotExist(err) {
		if err := os.RemoveAll(destDir); err != nil {
			return fmt.Errorf("failed to clean destination directory: %w", err)
		}
	}

	// 使用系统 git 命令进行浅克隆
	cmd := exec.Command("git", "clone", "--depth", "1", url, destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository %s: %w", url, err)
	}

	return nil
}

// CloneToTemp 克隆仓库到临时目录，返回临时目录路径
func (f *Fetcher) CloneToTemp(source string) (string, error) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "skillsync-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// 关键步骤：优先使用缓存仓库，失败则回退为直接 clone
	if cachePath, cacheErr := f.ensureCacheRepo(source); cacheErr == nil {
		if err := f.cloneFromCache(cachePath, tempDir, ""); err == nil {
			return tempDir, nil
		}
	}

	// 回退：克隆到临时目录
	if err := f.Clone(source, tempDir); err != nil {
		os.RemoveAll(tempDir) // 清理临时目录
		return "", err
	}

	return tempDir, nil
}

// CloneWithBranch 克隆指定分支的仓库到指定目录
// 使用 --branch 参数确保克隆正确的分支
func (f *Fetcher) CloneWithBranch(source, destDir, branch string) error {
	url := f.NormalizeURL(source)

	// 确保目标目录不存在
	if _, err := os.Stat(destDir); !os.IsNotExist(err) {
		if err := os.RemoveAll(destDir); err != nil {
			return fmt.Errorf("failed to clean destination directory: %w", err)
		}
	}

	// 使用系统 git 命令进行浅克隆，指定分支
	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", branch, url, destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository %s (branch: %s): %w", url, branch, err)
	}

	return nil
}

// CloneToTempWithBranch 克隆指定分支到临时目录，返回临时目录路径
func (f *Fetcher) CloneToTempWithBranch(source, branch string) (string, error) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "skillsync-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// 关键步骤：优先使用缓存仓库，失败则回退为直接 clone
	if cachePath, cacheErr := f.ensureCacheRepo(source); cacheErr == nil {
		if err := f.cloneFromCache(cachePath, tempDir, branch); err == nil {
			return tempDir, nil
		}
	}

	// 回退：克隆到临时目录
	if err := f.CloneWithBranch(source, tempDir, branch); err != nil {
		os.RemoveAll(tempDir) // 清理临时目录
		return "", err
	}

	return tempDir, nil
}
