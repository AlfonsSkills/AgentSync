// Package updater 提供版本检查和自动更新功能
package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	// GitHub API 地址
	releaseAPI = "https://api.github.com/repos/AlfonsSkills/SkillSync/releases/latest"
	// 下载地址模板
	downloadURL = "https://github.com/AlfonsSkills/SkillSync/releases/download/%s/skillsync-%s-%s%s"
)

// Release GitHub release 信息
type Release struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

// CheckResult 版本检查结果
type CheckResult struct {
	CurrentVersion string
	LatestVersion  string
	IsLatest       bool
	ReleaseURL     string
}

// CheckLatestVersion 检查是否有新版本
func CheckLatestVersion(currentVersion string) (*CheckResult, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(releaseAPI)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to check for updates: HTTP %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}

	// 规范化版本号比较
	current := normalizeVersion(currentVersion)
	latest := normalizeVersion(release.TagName)

	return &CheckResult{
		CurrentVersion: currentVersion,
		LatestVersion:  release.TagName,
		IsLatest:       current == latest || currentVersion == "dev",
		ReleaseURL:     release.HTMLURL,
	}, nil
}

// Upgrade 执行升级
func Upgrade(currentVersion string) error {
	result, err := CheckLatestVersion(currentVersion)
	if err != nil {
		return err
	}

	if result.IsLatest {
		fmt.Println("✓ You are already using the latest version")
		return nil
	}

	fmt.Printf("📦 Upgrading from %s to %s...\n", result.CurrentVersion, result.LatestVersion)

	// 构建下载 URL
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	ext := ""
	if goos == "windows" {
		ext = ".exe"
	}

	url := fmt.Sprintf(downloadURL, result.LatestVersion, goos, goarch, ext)

	// 下载新版本
	fmt.Printf("⬇️  Downloading from %s...\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: HTTP %d", resp.StatusCode)
	}

	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp(filepath.Dir(execPath), "skillsync-upgrade-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// 写入下载内容
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to save download: %w", err)
	}
	tmpFile.Close()

	// 设置执行权限
	if err := os.Chmod(tmpPath, 0755); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// 备份当前版本（可选）
	backupPath := execPath + ".bak"
	_ = os.Remove(backupPath) // 忽略错误，可能不存在
	if err := os.Rename(execPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current version: %w", err)
	}

	// 替换为新版本
	if err := os.Rename(tmpPath, execPath); err != nil {
		// 尝试恢复备份
		_ = os.Rename(backupPath, execPath)
		return fmt.Errorf("failed to install new version: %w", err)
	}

	// 删除备份
	_ = os.Remove(backupPath)

	fmt.Printf("✅ Successfully upgraded to %s\n", result.LatestVersion)
	return nil
}

// normalizeVersion 规范化版本号（去除 v 前缀）
func normalizeVersion(version string) string {
	return strings.TrimPrefix(strings.TrimSpace(version), "v")
}
