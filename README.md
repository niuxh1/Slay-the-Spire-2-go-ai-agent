# STS2-Agent-Pro: LLM-Driven High-Precision AI for Slay the Spire 2 🃏🤖

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue.svg?style=for-the-badge&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)
[![Protocol](https://img.shields.io/badge/MCP-Compatible-blueviolet.svg?style=for-the-badge)](https://modelcontextprotocol.io)
[![Status](https://img.shields.io/badge/Status-Beta-orange.svg?style=for-the-badge)](https://github.com/niuxh/sts2-go-agent)

[**English**](#english) | [**中文说明**](#chinese)

---

<a name="english"></a>
## 🇬🇧 English Description

**STS2-Agent-Pro** is a high-performance, fully autonomous AI decision-making system built for *Slay the Spire 2*. By leveraging native C# memory injection and the **Model Context Protocol (MCP)**, the project brings the reasoning power of Large Language Models (LLMs) into the complex strategic space of Roguelike deck-builders.

### 🧠 Technical Highlights
- **Heuristic SOP Engine**: Built-in Standard Operating Procedures with HP-Equivalent Value Models ($\Delta HP \ge 12$) and Act 1 greedy path algorithms.
- **Physical Layer Hard Gate**: A strict validation layer designed to eliminate LLM hallucinations. Automatically aligns game memory before `end_turn` to prevent unforced errors.
- **Semantic State Reduction**: Transforms massive raw Godot engine data into compact Markdown streams, significantly reducing token costs.
- **Class-Specific Logic**: Formula-based decision-making for Ironclad, Necrobinder, Regent, and more.

---

<a name="chinese"></a>
## 🇨🇳 中文说明

**STS2-Agent-Pro** 是一款专为《杀戮尖塔 2》(Slay the Spire 2) 构建的高性能、全自动 AI 决策系统。通过原生 C# 内存注入与 **模型上下文协议 (MCP)**，该项目将大语言模型 (LLM) 的逻辑推理能力引入了复杂的 Roguelike 塔防博弈中。

### 🧠 核心技术亮点
- **启发式 SOP 决策引擎**：内置标准作业程序，支持血量等效价值模型 ($\Delta HP \ge 12$) 与 Act 1 贪婪路径演算法。
- **物理层硬性门控 (Hard Gate)**：针对 LLM 幻觉设计的强制状态校验机制。在 `end_turn` 前自动对齐游戏内存，杜绝非受迫性空过。
- **语义化状态降维**：将 Godot 引擎的海量原始数据转化为紧凑的 Markdown 信息流，显著降低 Token 成本。
- **多职业专属逻辑**：针对铁甲战士、死灵缚者、摄政王等职业的资源分配进行了数学推演与公式化决策。

---

## 📚 Technical Documentation / 技术文档

| Document / 文档 | Description / 描述 |
| :--- | :--- |
| [**System Architecture**](./docs/ARCHITECTURE.md) | Decoupling C# Mod (Senses) and Go Engine (Brain). <br> 解析 C# Mod (感官) 与 Go Engine (大脑) 的解耦通讯模型。 |
| [**Decision Logic & SOP**](./docs/DECISION_LOGIC.md) | Math formulas, class heuristics, and gating mechanisms. <br> 深入探讨背后的数学期望公式、职业算法调优及门控机制。 |
| [**Dev & Contribution**](./docs/DEVELOPMENT.md) | How to add new classes, models, or optimize prompts. <br> 如何为智能体添加新职业支持、接入新模型或优化 Prompt 模板。 |

## 🚀 Quick Start / 快速启动

### 1. Backend C# Mod (Prerequisite) / 安装后端 C# Mod
This project requires the **STS2-Agent** mod to communicate with the game.
本项目需要 **STS2-Agent** 模组作为与游戏通信的桥梁。

1. **Clone & Install / 克隆与安装**:
   ```bash
   git clone https://github.com/CharTyr/STS2-Agent.git
   ```
2. **Setup / 环境**: Install [.NET 9 SDK](https://dotnet.microsoft.com/download/dotnet/9.0).
3. **Build & Deploy / 构建与部署**:
   ```powershell
   cd STS2-Agent
   powershell -File scripts/build-mod.ps1 -GameRoot "Your_Game_Path"
   ```

### 2. Start Go Brain / 启动 Go 智能体
```powershell
cd sts2-go-agent
go build -o STS2-Agent-Pro.exe cmd/agent/main.go
.\STS2-Agent-Pro.exe
```

### 3. Dashboard / 连接控制中心
Open `http://localhost:8090` in your browser.

---

## 🙏 Acknowledgements / 特别致谢
Special thanks to the developers of **STS2-Agent**. Your contributions to game protocol reversing and MCP standardization made this project possible.
特别感谢 **STS2-Agent 项目的研发者**。没有你们提供的稳定 API，本项目将无法实现对游戏环境的实时感知。
