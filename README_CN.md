# SkillSync

[![CI](https://github.com/AlfonsSkills/SkillSync/actions/workflows/ci.yml/badge.svg)](https://github.com/AlfonsSkills/SkillSync/actions/workflows/ci.yml)
[![Release](https://github.com/AlfonsSkills/SkillSync/actions/workflows/release.yml/badge.svg)](https://github.com/AlfonsSkills/SkillSync/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/AlfonsSkills/SkillSync)](https://goreportcard.com/report/github.com/AlfonsSkills/SkillSync)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[English](README.md)

从 Git 仓库同步 Skill 到本地 AI 编码工具（Gemini CLI / Claude Code / Codex CLI）。

## 功能特性

- 📦 **安装技能** - 从任意 Git 仓库安装（默认 GitHub）
- 📋 **列出技能** - 查看所有工具中已安装的技能
- 🗑️ **移除技能** - 从指定或所有工具中移除
- 🎯 **目标选择** - 选择同步到哪些工具

## 安装

### 一键安装（推荐）

```bash
curl -fsSL https://raw.githubusercontent.com/AlfonsSkills/SkillSync/main/install.sh | bash
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

## 使用方法

```bash
# 从 monorepo 安装技能（交互式选择）
skillsync install anthropics/skills

# 安装到指定工具
skillsync install anthropics/skills --target gemini
skillsync install AlfonsSkills/skills -t claude,codex

# 安装到项目本地目录 (.gemini/skills, .claude/skills, .codex/skills)
skillsync install anthropics/skills --local

# 从其他 Git 平台安装
skillsync install https://gitlab.com/user/skill-repo.git

# 列出已安装的技能（全局 + 项目本地）
skillsync list
skillsync list --target gemini

# 输出示例：
#   📦 Installed Skills:
#
#     Gemini CLI (1):
#     📁 ~/.gemini/skills
#       ✓ devops
#
#     Claude Code (2):
#     📁 ~/.claude/skills
#       ✓ devops
#       [project:MyProject]
#         ✓ docx
#
#     Codex CLI (5):
#     📁 ~/.codex/skills
#       ✓ gh-address-comments
#       [public]
#         ✓ devops
#       [.system]
#         ✓ skill-creator

# 移除技能
skillsync remove skill-name
skillsync remove skill-name --target claude
skillsync remove skill-name --local  # 从项目目录移除
```

## 支持的工具

| 工具 | Skills 目录 | 参数 |
|------|------------|------|
| Gemini CLI | `~/.gemini/skills/` | `-t gemini` |
| Claude Code | `~/.claude/skills/` | `-t claude` |
| Codex CLI | `~/.codex/skills/public/` | `-t codex` |
| OpenCode | `~/.config/opencode/skill/` | `-t opencode` |
| Antigravity IDE | `~/.gemini/antigravity/skills/` | `-t antigravity` |
| Copilot / VSCode | `~/.copilot/skills/` | `-t copilot` |
| Cursor | `~/.cursor/skills/` | `-t cursor` |

## Skill 仓库要求

有效的 Skill 仓库必须包含 `SKILL.md` 文件：

```
my-skill/
├── SKILL.md          # 必需：Skill 定义文件
├── references/       # 可选：参考文档
│   └── *.md
└── scripts/          # 可选：脚本
    └── *.sh
```

## 许可证

MIT License - 详见 [LICENSE](LICENSE)。
