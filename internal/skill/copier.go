// Package skill 提供文件拷贝功能
package skill

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CopyOptions 拷贝选项
type CopyOptions struct {
	ExcludeDirs  []string // 要排除的目录名
	ExcludeFiles []string // 要排除的文件名
}

// DefaultCopyOptions 返回默认的拷贝选项
func DefaultCopyOptions() CopyOptions {
	return CopyOptions{
		ExcludeDirs:  []string{".git"},
		ExcludeFiles: []string{".gitignore", ".gitattributes"},
	}
}

// CopyDir 将源目录内容拷贝到目标目录
func CopyDir(src, dst string, opts CopyOptions) error {
	// 获取源目录信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source directory: %w", err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory: %s", src)
	}

	// 创建目标目录
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// 遍历源目录
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 检查是否需要排除目录
			if shouldExclude(entry.Name(), opts.ExcludeDirs) {
				continue
			}
			// 递归拷贝子目录
			if err := CopyDir(srcPath, dstPath, opts); err != nil {
				return err
			}
		} else {
			// 检查是否需要排除文件
			if shouldExclude(entry.Name(), opts.ExcludeFiles) {
				continue
			}
			// 拷贝文件
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile 拷贝单个文件
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}

// shouldExclude 检查名称是否在排除列表中
func shouldExclude(name string, excludeList []string) bool {
	for _, exclude := range excludeList {
		if strings.EqualFold(name, exclude) {
			return true
		}
	}
	return false
}

// SkillInfo 表示一个 skill 的信息
type SkillInfo struct {
	Name string // skill 名称（目录名）
	Path string // skill 完整路径
	Desc string // 从 SKILL.md 提取的描述
}

// ScanSkills scans directory recursively for valid skills (directories containing SKILL.md)
// Supports nested structures common in monorepos (e.g., root/skills/skill-name/SKILL.md)
func ScanSkills(dir string) ([]SkillInfo, error) {
	var skills []SkillInfo

	// Recursively scan for SKILL.md files
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git and other hidden directories
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		// Found SKILL.md
		if !info.IsDir() && info.Name() == "SKILL.md" {
			// Get parent directory as skill directory
			skillDir := filepath.Dir(path)

			// Skip root directory SKILL.md (likely a template)
			if skillDir == dir {
				return nil
			}

			// Extract skill name from directory
			skillName := filepath.Base(skillDir)

			// Extract description
			desc := extractSkillDescription(path)

			skills = append(skills, SkillInfo{
				Name: skillName,
				Path: skillDir,
				Desc: desc,
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}

	return skills, nil
}

// extractSkillDescription 从 SKILL.md 提取 description 字段
// 入参: skillFile SKILL.md 文件路径
// 返回: description 字段内容，读取失败或不存在则为空
func extractSkillDescription(skillFile string) string {
	content, err := os.ReadFile(skillFile)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "description:") {
			desc := strings.TrimPrefix(line, "description:")
			return strings.TrimSpace(desc)
		}
	}
	return ""
}

// ReadSkillDescription 读取指定 skill 目录中的描述信息
// 入参: skillDir skill 目录路径
// 返回: description 字段内容，读取失败或不存在则为空
func ReadSkillDescription(skillDir string) string {
	skillFile := filepath.Join(skillDir, "SKILL.md")
	return extractSkillDescription(skillFile)
}

// ValidateSkillDir 验证目录是否是有效的 skill 目录
// 有效的 skill 目录必须包含 SKILL.md 文件
func ValidateSkillDir(dir string) error {
	skillFile := filepath.Join(dir, "SKILL.md")
	if _, err := os.Stat(skillFile); os.IsNotExist(err) {
		return fmt.Errorf("invalid skill directory: SKILL.md not found in %s", dir)
	}
	return nil
}

// ExtractSkillName 从仓库 URL 或路径中提取 skill 名称
func ExtractSkillName(source string) string {
	// 移除 .git 后缀
	source = strings.TrimSuffix(source, ".git")

	// 获取最后一个路径段
	parts := strings.Split(source, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return source
}
