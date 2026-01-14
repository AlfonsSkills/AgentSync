package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"

	"github.com/AlfonsSkills/SkillSync/internal/project"
	"github.com/AlfonsSkills/SkillSync/internal/skill"
	"github.com/AlfonsSkills/SkillSync/internal/target"
)

// InteractiveContext 存储交互式选择的结果
type InteractiveContext struct {
	Providers     []target.ToolProvider // 选中的目标工具
	InstallGlobal bool                  // 是否安装到全局目录
	InstallLocal  bool                  // 是否安装到项目目录
	ProjectRoot   string                // 项目根目录（如果 InstallLocal 为 true）
}

// resolveTargetProviders 解析或交互选择目标工具
// 如果 targetFlags 为空且未显式指定，显示多选框让用户选择
// explicitlySet: 用户是否通过 --target 显式指定了值
func resolveTargetProviders(targetFlags []string) ([]target.ToolProvider, bool, error) {
	// 如果显式指定了 target，直接解析
	if len(targetFlags) > 0 {
		providers, err := target.ParseProviders(targetFlags)
		if err != nil {
			return nil, true, err
		}
		// 显示已选择的工具
		color.Cyan("🎯 Target tools:\n")
		for _, p := range providers {
			color.White("   • %s\n", p.DisplayName())
		}
		fmt.Println()
		return providers, true, nil
	}

	// 未指定，显示交互式多选
	allProviders := target.AllProviders()
	var options []string
	for _, p := range allProviders {
		options = append(options, p.DisplayName())
	}

	var selectedIndices []int
	prompt := &survey.MultiSelect{
		Message:  "Select target tools:",
		Options:  options,
		PageSize: 5,
	}
	if err := survey.AskOne(prompt, &selectedIndices); err != nil {
		return nil, false, fmt.Errorf("selection cancelled: %w", err)
	}

	if len(selectedIndices) == 0 {
		return nil, false, fmt.Errorf("no tools selected")
	}

	var selectedProviders []target.ToolProvider
	for _, idx := range selectedIndices {
		selectedProviders = append(selectedProviders, allProviders[idx])
	}

	return selectedProviders, false, nil
}

// resolveLocalInstall 解析或交互选择是否安装到项目目录
// localFlag: --local 标志的值
// Returns: installGlobal, installLocal, projectRoot, error
func resolveLocalInstall(localFlag bool) (bool, bool, string, error) {
	// 尝试获取项目根目录
	projectRoot, projectErr := project.FindProjectRoot()
	inProject := projectErr == nil

	// 如果显式指定了 --local
	if localFlag {
		if !inProject {
			return false, false, "", fmt.Errorf("not in a git repository, --local requires a project context")
		}
		// 仅安装到项目目录
		color.Cyan("📁 Install scope: Project only\n")
		color.HiBlack("   Project root: %s\n\n", projectRoot)
		return false, true, projectRoot, nil
	}

	// 如果不在项目中，只能安装到全局
	if !inProject {
		color.Cyan("📁 Install scope: Global only\n")
		color.HiBlack("   (Not in a git repository)\n\n")
		return true, false, "", nil
	}

	// 在项目中，询问是否也安装到项目目录
	var alsoLocal bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Also install to project directory?\n   (%s)", projectRoot),
		Default: false,
	}
	if err := survey.AskOne(prompt, &alsoLocal); err != nil {
		return false, false, "", fmt.Errorf("selection cancelled: %w", err)
	}

	if alsoLocal {
		color.Cyan("📁 Install scope: Global + Project\n")
		color.HiBlack("   Project root: %s\n\n", projectRoot)
		return true, true, projectRoot, nil
	}

	color.Cyan("📁 Install scope: Global only\n\n")
	return true, false, "", nil
}

// showInstallPreview 显示安装路径预览
func showInstallPreview(skills []skill.SkillInfo, providers []target.ToolProvider, installGlobal, installLocal bool, projectRoot string) {
	color.Cyan("📍 Installation preview:\n")

	for _, s := range skills {
		color.White("   Skill: %s\n", color.New(color.FgCyan).Sprint(s.Name))

		if installGlobal {
			color.White("   Global:\n")
			for _, p := range providers {
				dir, _ := p.GlobalInstallDir()
				color.White("     📁 %s/%s\n", dir, s.Name)
			}
		}

		if installLocal && projectRoot != "" {
			color.White("   Project:\n")
			for _, p := range providers {
				dir := p.LocalSkillsDir(projectRoot)
				color.White("     📁 %s/%s\n", dir, s.Name)
			}
		}
	}
	fmt.Println()
}

// showRemovePreview 显示删除路径预览
func showRemovePreview(skillName string, providers []target.ToolProvider, removeGlobal, removeLocal bool, projectRoot string) {
	color.Cyan("🗑️  Removal preview:\n")
	color.White("   Skill: %s\n", color.New(color.FgCyan).Sprint(skillName))

	if removeGlobal {
		color.White("   Global:\n")
		for _, p := range providers {
			dir, _ := p.GlobalInstallDir()
			color.White("     📁 %s/%s\n", dir, skillName)
		}
	}

	if removeLocal && projectRoot != "" {
		color.White("   Project:\n")
		for _, p := range providers {
			dir := p.LocalSkillsDir(projectRoot)
			color.White("     📁 %s/%s\n", dir, skillName)
		}
	}
	fmt.Println()
}

// resolveRemoveScope 解析或交互选择删除范围
// 与 resolveLocalInstall 类似，但用于 remove 命令
func resolveRemoveScope(localFlag bool) (bool, bool, string, error) {
	projectRoot, projectErr := project.FindProjectRoot()
	inProject := projectErr == nil

	if localFlag {
		if !inProject {
			return false, false, "", fmt.Errorf("not in a git repository, --local requires a project context")
		}
		color.Cyan("📁 Remove scope: Project only\n")
		color.HiBlack("   Project root: %s\n\n", projectRoot)
		return false, true, projectRoot, nil
	}

	if !inProject {
		color.Cyan("📁 Remove scope: Global only\n")
		color.HiBlack("   (Not in a git repository)\n\n")
		return true, false, "", nil
	}

	// 在项目中，询问是否也从项目目录删除
	var alsoLocal bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Also remove from project directory?\n   (%s)", projectRoot),
		Default: false,
	}
	if err := survey.AskOne(prompt, &alsoLocal); err != nil {
		return false, false, "", fmt.Errorf("selection cancelled: %w", err)
	}

	if alsoLocal {
		color.Cyan("📁 Remove scope: Global + Project\n")
		color.HiBlack("   Project root: %s\n\n", projectRoot)
		return true, true, projectRoot, nil
	}

	color.Cyan("📁 Remove scope: Global only\n\n")
	return true, false, "", nil
}

// checkSkillExistsInProviders 检查 skill 在哪些工具的全局目录中存在
// 返回存在该 skill 的 providers 列表
func checkSkillExistsInProviders(skillName string) []target.ToolProvider {
	allProviders := target.AllProviders()
	var existingProviders []target.ToolProvider

	for _, p := range allProviders {
		globalDir, err := p.GlobalInstallDir()
		if err != nil {
			continue
		}
		skillPath := filepath.Join(globalDir, skillName)
		if _, err := os.Stat(skillPath); err == nil {
			existingProviders = append(existingProviders, p)
		}
	}

	return existingProviders
}

// checkSkillExistsInProject 检查 skill 是否在指定 providers 的项目目录中存在
// 返回是否在任一工具的项目目录存在
func checkSkillExistsInProject(skillName string, providers []target.ToolProvider, projectRoot string) bool {
	if projectRoot == "" {
		return false
	}

	for _, p := range providers {
		localDir := p.LocalSkillsDir(projectRoot)
		skillPath := filepath.Join(localDir, skillName)
		if _, err := os.Stat(skillPath); err == nil {
			return true
		}
	}

	return false
}

// resolveTargetProvidersForRemove 为 remove 命令解析目标工具
// 仅显示存在 skill 的工具选项
func resolveTargetProvidersForRemove(skillName string, targetFlags []string) ([]target.ToolProvider, bool, error) {
	// 如果显式指定了 target，直接解析（不过滤）
	if len(targetFlags) > 0 {
		providers, err := target.ParseProviders(targetFlags)
		if err != nil {
			return nil, true, err
		}
		color.Cyan("🎯 Target tools:\n")
		for _, p := range providers {
			color.White("   • %s\n", p.DisplayName())
		}
		fmt.Println()
		return providers, true, nil
	}

	// 检查 skill 在哪些工具中存在
	existingProviders := checkSkillExistsInProviders(skillName)
	if len(existingProviders) == 0 {
		return nil, false, fmt.Errorf("skill '%s' not found in any tool's global directory", skillName)
	}

	// 构建选项列表（仅包含存在 skill 的工具）
	var options []string
	for _, p := range existingProviders {
		options = append(options, p.DisplayName())
	}

	var selectedIndices []int
	prompt := &survey.MultiSelect{
		Message:  "Select target tools:",
		Options:  options,
		PageSize: 5,
	}
	if err := survey.AskOne(prompt, &selectedIndices); err != nil {
		return nil, false, fmt.Errorf("selection cancelled: %w", err)
	}

	if len(selectedIndices) == 0 {
		return nil, false, fmt.Errorf("no tools selected")
	}

	var selectedProviders []target.ToolProvider
	for _, idx := range selectedIndices {
		selectedProviders = append(selectedProviders, existingProviders[idx])
	}

	return selectedProviders, false, nil
}

// resolveRemoveScopeWithCheck 解析或交互选择删除范围（带存在性检查）
// 仅当项目目录中存在 skill 时才提示是否删除
func resolveRemoveScopeWithCheck(skillName string, providers []target.ToolProvider, localFlag bool) (bool, bool, string, error) {
	projectRoot, projectErr := project.FindProjectRoot()
	inProject := projectErr == nil

	if localFlag {
		if !inProject {
			return false, false, "", fmt.Errorf("not in a git repository, --local requires a project context")
		}
		// 检查项目目录中是否存在 skill
		if !checkSkillExistsInProject(skillName, providers, projectRoot) {
			return false, false, "", fmt.Errorf("skill '%s' not found in project directory", skillName)
		}
		color.Cyan("📁 Remove scope: Project only\n")
		color.HiBlack("   Project root: %s\n\n", projectRoot)
		return false, true, projectRoot, nil
	}

	if !inProject {
		color.Cyan("📁 Remove scope: Global only\n")
		color.HiBlack("   (Not in a git repository)\n\n")
		return true, false, "", nil
	}

	// 检查项目目录中是否存在 skill
	existsInProject := checkSkillExistsInProject(skillName, providers, projectRoot)
	if !existsInProject {
		// 项目目录中不存在，直接返回仅全局删除
		color.Cyan("📁 Remove scope: Global only\n")
		color.HiBlack("   (Skill not found in project directory)\n\n")
		return true, false, "", nil
	}

	// 在项目中且存在 skill，询问是否也从项目目录删除
	var alsoLocal bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Also remove from project directory?\n   (%s)", projectRoot),
		Default: false,
	}
	if err := survey.AskOne(prompt, &alsoLocal); err != nil {
		return false, false, "", fmt.Errorf("selection cancelled: %w", err)
	}

	if alsoLocal {
		color.Cyan("📁 Remove scope: Global + Project\n")
		color.HiBlack("   Project root: %s\n\n", projectRoot)
		return true, true, projectRoot, nil
	}

	color.Cyan("📁 Remove scope: Global only\n\n")
	return true, false, "", nil
}
