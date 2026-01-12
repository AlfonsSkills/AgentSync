# AgentSync

[![CI](https://github.com/AlfonsSkills/AgentSync/actions/workflows/ci.yml/badge.svg)](https://github.com/AlfonsSkills/AgentSync/actions/workflows/ci.yml)
[![Release](https://github.com/AlfonsSkills/AgentSync/actions/workflows/release.yml/badge.svg)](https://github.com/AlfonsSkills/AgentSync/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/AlfonsSkills/AgentSync)](https://goreportcard.com/report/github.com/AlfonsSkills/AgentSync)
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
curl -fsSL https://raw.githubusercontent.com/AlfonsSkills/AgentSync/main/install.sh | bash
```

### 从 Release 下载

从 [Releases](https://github.com/AlfonsSkills/AgentSync/releases) 下载最新版本。

### 从源码构建

```bash
git clone https://github.com/AlfonsSkills/AgentSync.git
cd AgentSync
make build
# 二进制文件位于 ./build/agentsync
```

## 使用方法

```bash
# 安装技能到所有工具
agentsync install user/repo

# 安装到指定工具
agentsync install user/repo --target gemini
agentsync install user/repo -t claude,codex

# 从其他 Git 平台安装
agentsync install https://gitlab.com/user/skill-repo.git

# 列出已安装的技能
agentsync list
agentsync list --target gemini

# 移除技能
agentsync remove skill-name
agentsync remove skill-name --target claude
```

## 支持的工具

| 工具 | Skills 目录 | 参数 |
|------|------------|------|
| Gemini CLI | `~/.gemini/skills/` | `-t gemini` |
| Claude Code | `~/.claude/skills/` | `-t claude` |
| Codex CLI | `~/.codex/skills/public/` | `-t codex` |

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
