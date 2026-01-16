# SkillSync

[![CI](https://github.com/AlfonsSkills/SkillSync/actions/workflows/ci.yml/badge.svg)](https://github.com/AlfonsSkills/SkillSync/actions/workflows/ci.yml)
[![Release](https://github.com/AlfonsSkills/SkillSync/actions/workflows/release.yml/badge.svg)](https://github.com/AlfonsSkills/SkillSync/releases)
[![GitHub release](https://img.shields.io/github/v/release/AlfonsSkills/SkillSync)](https://github.com/AlfonsSkills/SkillSync/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/AlfonsSkills/SkillSync)](https://goreportcard.com/report/github.com/AlfonsSkills/SkillSync)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[English](README.md)

**一条命令，将 Git 仓库中的 Skill 同步到 14+ 种 AI 编码工具。**

## 预览

```
📦 Installed Skills:

  Gemini CLI (1):
  📁 ~/.gemini/skills
    ✓ devops

  Claude Code (2):
  📁 ~/.claude/skills
    ✓ devops
    [project:MyProject]
      ✓ docx

  Codex CLI (5):
  📁 ~/.codex/skills
    ✓ gh-address-comments
    [public]
      ✓ devops
    [.system]
      ✓ skill-creator
```

## 安装

```bash
curl -fsSL https://raw.githubusercontent.com/AlfonsSkills/SkillSync/main/install.sh | bash
```

**一键安装 Skill 仓库：**

```bash
# 安装 SkillSync 并同时安装指定的 Skill 仓库
curl -fsSL https://raw.githubusercontent.com/AlfonsSkills/SkillSync/main/install.sh | bash -s -- install AlfonsSkills/skills

# 指定目标工具
curl -fsSL https://raw.githubusercontent.com/AlfonsSkills/SkillSync/main/install.sh | bash -s -- install AlfonsSkills/skills -t claude,gemini
```

### 从 Release 下载

从 [Releases](https://github.com/AlfonsSkills/SkillSync/releases) 下载最新版本。

### 从源码构建

```bash
git clone https://github.com/AlfonsSkills/SkillSync.git
cd SkillSync
make build
# 二进制文件位于 ./build/skillsync
```

## 快速开始

```bash
# 从仓库安装技能
skillsync install anthropics/skills

# 安装到指定工具
skillsync install anthropics/skills -t gemini

# 安装到项目本地目录
skillsync install anthropics/skills --local

# 列出已安装的技能
skillsync list

# 移除技能
skillsync remove skill-name

# 安装到多个工具
skillsync install AlfonsSkills/skills -t claude,codex,gemini

# 从 GitLab 或其他平台安装
skillsync install https://gitlab.com/user/skill-repo.git

# 从 tree URL 安装单个 Skill（GitHub）
skillsync install https://github.com/davila7/claude-code-templates/tree/main/cli-tool/components/skills/creative-design/canvas-design

# 不带域名时默认 GitHub
skillsync install davila7/claude-code-templates/tree/main/cli-tool/components/skills/creative-design/canvas-design

# 列出指定工具的技能
skillsync list --target gemini

# 从指定工具移除
skillsync remove skill-name --target claude

# 仅从项目目录移除
skillsync remove skill-name --local
```

## 支持的工具

SkillSync 支持 **14 种 AI 编码工具**，涵盖终端和 IDE 环境。

### 终端工具

| 工具 | Skills 目录 | 参数 |
|------|------------|------|
| Gemini CLI | `~/.gemini/skills/` | `-t gemini` |
| Claude Code | `~/.claude/skills/` | `-t claude` |
| Codex CLI | `~/.codex/skills/public/` | `-t codex` |
| OpenCode | `~/.config/opencode/skill/` | `-t opencode` |
| Goose AI | `~/.config/goose/skills/` | `-t goose` |
| Crush | `~/.config/crush/skills/` | `-t crush` |

### IDE 工具

| 工具 | Skills 目录 | 参数 |
|------|------------|------|
| Antigravity IDE | `~/.gemini/antigravity/skills/` | `-t antigravity` |
| Copilot | `~/.copilot/skills/` | `-t copilot` |
| Cursor | `~/.cursor/skills/` | `-t cursor` |
| Cline IDE | `~/.cline/skills/` | `-t cline` |
| Droid (Factory AI) | `~/.factory/skills/` | `-t droid` |
| Kilo Code | `~/.kilocode/skills/` | `-t kilocode` |
| Roo Code | `~/.roo/skills/` | `-t roocode` |
| VSCode (Copilot) | `~/.copilot/skills/` | `-t vscode` |

## Skill 格式

有效的 Skill 仓库必须包含 `SKILL.md` 文件：

```
my-skill/
├── SKILL.md          # 必需：Skill 定义文件
├── references/       # 可选：参考文档
└── scripts/          # 可选：脚本
```

## 许可证

MIT License - 详见 [LICENSE](LICENSE)。
