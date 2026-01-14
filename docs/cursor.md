# Agent Skills

> **注意**：Agent Skills 仅在 Nightly 更新渠道中可用。
>
> 要切换更新渠道，打开 Cursor 设置（Mac: `Cmd+Shift+J`，Windows/Linux: `Ctrl+Shift+J`），选择 Beta，然后将更新渠道设置为 Nightly。更新完成后，你可能需要重新启动 Cursor。

Agent Skills 是一种用于为 AI Agent 扩展专门能力的开放标准。Skills 将特定领域的知识和工作流打包，Agent 可以利用这些内容来执行特定任务。

## 什么是技能？

技能是一种可移植、受版本控制的包，用于教会 Agent 如何执行特定领域的任务。

### 可移植

技能可以在任何支持 Agent Skills 标准的 Agent 中使用。

### 版本控制

技能以文件形式存储，可以在你的代码仓库中进行管理和追踪，或通过 GitHub 仓库链接安装。

## 技能如何工作

Cursor 启动时会自动从技能目录中发现技能，并将它们提供给 Agent。Agent 会查看可用的技能，并根据上下文决定何时使用。

也可以在 Agent 对话中输入 `/`，通过搜索技能名称手动调用技能。

## 技能目录

技能会自动从以下位置加载：

| 位置 | 范围 |
|------|------|
| `.cursor/skills/` | 项目级 |
| `.claude/skills/` | 项目级（Claude 兼容） |
| `~/.cursor/skills/` | 用户级（全局） |
| `~/.claude/skills/` | 用户级（全局） |

每个技能都应是一个包含 `SKILL.md` 文件的文件夹：

```
.cursor/
└── skills/
    └── my-skill/
        └── SKILL.md
```

## SKILL.md 文件格式

每个技能都在一个带有 YAML frontmatter 的 `SKILL.md` 文件中定义：

```markdown
---
name: my-skill
description: Short description of what this skill does and when to use it.
---

# My Skill

Detailed instructions for the agent.

## When to Use

- Use this skill when...
- This skill is helpful for...

## Instructions

- Step-by-step guidance for the agent
- Domain-specific conventions
- Best practices and patterns
```

### Frontmatter 字段

| 字段 | 是否必填 | 说明 |
|------|---------|------|
| `description` | 是 | 在菜单中显示的简短描述。Agent 会用它来判断何时使用该技能。 |
| `name` | 否 | 便于阅读的名称。如果省略，将使用父文件夹的名称。 |

## 查看技能

要查看已发现的技能：

1. 打开 Cursor Settings（Mac: `Cmd+Shift+J`，Windows/Linux: `Ctrl+Shift+J`）
2. 导航到 **Rules**
3. 技能会显示在 **Agent Decides** 部分

## 从 GitHub 安装技能

你可以从 GitHub 仓库导入技能：

1. 打开 **Cursor Settings** → **Rules**
2. 在 **Project Rules** 部分，点击 **Add Rule**
3. 选择 **Remote Rule (Github)**
4. 输入 GitHub 仓库 URL

## 了解更多

Agent Skills 是一个开放标准。了解详情请访问 [agentskills.io](https://agentskills.io)。
