# AgentSync

[![CI](https://github.com/AlfonsSkills/AgentSync/actions/workflows/ci.yml/badge.svg)](https://github.com/AlfonsSkills/AgentSync/actions/workflows/ci.yml)
[![Release](https://github.com/AlfonsSkills/AgentSync/actions/workflows/release.yml/badge.svg)](https://github.com/AlfonsSkills/AgentSync/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/AlfonsSkills/AgentSync)](https://goreportcard.com/report/github.com/AlfonsSkills/AgentSync)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[中文文档](README_CN.md)

Sync skills from Git repositories to local AI coding tools (Gemini CLI / Claude Code / Codex CLI).

## Features

- 📦 **Install skills** from any Git repository (GitHub by default)
- 📋 **List skills** installed locally across all tools
- 🗑️ **Remove skills** from specific or all tools
- 🎯 **Target selection** - choose which tools to sync to

## Installation

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/AlfonsSkills/AgentSync/main/install.sh | bash
```

### From Release

Download the latest binary from [Releases](https://github.com/AlfonsSkills/AgentSync/releases).

### From Source

```bash
git clone https://github.com/AlfonsSkills/AgentSync.git
cd AgentSync
make build
# Binary will be at ./build/agentsync
```

## Usage

```bash
# Install skills from monorepo (interactive selection)
agentsync install anthropics/skills

# Install to specific tool
agentsync install anthropics/skills --target gemini
agentsync install AlfonsSkills/skills -t claude,codex

# Install from other Git platforms
agentsync install https://gitlab.com/user/skill-repo.git

# List installed skills
agentsync list
agentsync list --target gemini

# Remove skill
agentsync remove skill-name
agentsync remove skill-name --target claude
```

## Supported Tools

| Tool | Skills Directory | Flag |
|------|-----------------|------|
| Gemini CLI | `~/.gemini/skills/` | `-t gemini` |
| Claude Code | `~/.claude/skills/` | `-t claude` |
| Codex CLI | `~/.codex/skills/public/` | `-t codex` |

## Skill Repository Requirements

A valid skill repository must contain a `SKILL.md` file:

```
my-skill/
├── SKILL.md          # Required: Skill definition
├── references/       # Optional: Reference docs
│   └── *.md
└── scripts/          # Optional: Scripts
    └── *.sh
```

## License

MIT License - see [LICENSE](LICENSE) for details.
