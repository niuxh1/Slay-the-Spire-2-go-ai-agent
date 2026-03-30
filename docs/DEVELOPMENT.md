# Development & Contribution Guide
# 开发环境配置与贡献指南

[English](#english) | [中文](#chinese)

---

<a name="english"></a>
## 🇬🇧 English: Development Overview

We welcome contributions! Whether it's adding a new class heuristic or optimizing the prompt structure, this guide will help you get started.

### 1. Prerequisites
- **Go SDK 1.21+** (for the Brain Layer)
- **.NET 9 SDK** (for the Senses Layer/Mod)
- **Visual Studio 2022** (if you plan to modify the C# Mod)
- **OpenAI/Azure API Key** (for inference)

### 2. Project Layout
- `sts2-go-agent/cmd/agent/main.go`: The main entry point. Orchestrates the control loop.
- `sts2-go-agent/pkg/agent/prompt.go`: Defines how game state is translated into Markdown for LLMs. **The best place for prompt tuning.**
- `sts2-go-agent/pkg/models/`: Go structs representing the game state and action JSON.
- `STS2-Agent/`: (Submodule/External) The C# source for the game mod. [Repository Link](https://github.com/CharTyr/STS2-Agent)

### 3. Adding Support for a New Class
To add a new class (e.g., a Mod-added character):
1. Update `pkg/models/state.go` to include any unique resource trackers (like Stars or Souls).
2. Modify `pkg/agent/prompt.go` to inject class-specific keywords and heuristics into the system prompt.
3. Test using a local LLM or the debug dashboard to verify reasoning accuracy.

---

<a name="chinese"></a>
## 🇨🇳 中文：开发概述

我们非常欢迎社区贡献！无论是添加新职业的启发式逻辑，还是优化 Prompt 结构，本指南都将助你快速上手。

### 1. 环境准备
- **Go SDK 1.21+**: 用于大脑层逻辑开发。
- **.NET 9 SDK**: 用于感官层 (C# Mod) 的开发与构建。
- **Visual Studio 2022**: 如果你需要修改底层的内存读取或注入逻辑。
- **API Key**: 用于推理的 LLM 密钥。

### 2. 项目目录结构
- `sts2-go-agent/cmd/agent/main.go`: 程序主入口，协调整个控制循环。
- `sts2-go-agent/pkg/agent/prompt.go`: 定义如何将游戏状态转化为 LLM 易读的 Markdown。**这是进行 Prompt 调优的核心位置。**
- `sts2-go-agent/pkg/models/`: 对应游戏状态与动作的 Go 结构体定义。
- `STS2-Agent/`: (子模块/外部项目) C# Mod 的源代码。

### 3. 核心开发流程
#### 3.1 增加新职业支持
若要支持新加入的角色：
1. 在 `pkg/models/state.go` 中加入该职业的特殊资源（如：某种独特的充能点）。
2. 在 `pkg/agent/prompt.go` 中为该职业编写专属的 Prompt 逻辑，注入该职业的核心关键词（如“消耗”、“狡诈”等）。
3. 使用调试控制台（Dashboard）观察 AI 的推理过程，并根据表现微调启发式权重。

#### 3.2 优化推理精度 (Prompt Tuning)
我们采用语义化 Markdown 表格形式。如果你发现 AI 对某类卡牌的理解有偏差，请尝试在 `prompt.go` 中增加该卡牌的关联遗物上下文。

### 4. 提交规范
在发起 Pull Request 前，请确保：
- 代码已通过 `go fmt ./...`。
- 新增逻辑已在至少一个 Act 1 周期内进行了测试。
- 遵循中英双语的文档注释习惯。
