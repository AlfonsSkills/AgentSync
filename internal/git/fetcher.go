// Package git 提供 Git 仓库拉取功能
package git

import (
	"fmt"
	"os"
	"os/exec"
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

	// 克隆到临时目录
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

	// 克隆到临时目录
	if err := f.CloneWithBranch(source, tempDir, branch); err != nil {
		os.RemoveAll(tempDir) // 清理临时目录
		return "", err
	}

	return tempDir, nil
}
