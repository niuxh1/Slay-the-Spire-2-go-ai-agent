# Decision Logic & Heuristic SOP
# 决策逻辑与启发式 SOP 详解

[English](#english) | [中文](#chinese)

---

<a name="english"></a>
## 🇬🇧 English: Strategic Foundations

The intelligence of **STS2-Agent-Pro** is rooted in a combination of **Probabilistic Expectations** and **Class-Specific Heuristics**.

### 1. The Potion-Value Model ($\Delta HP$)
To avoid "hoarding syndrome," we treat potions as health-preservation equivalents:
$$\Delta HP = \text{Damage}_{\text{raw}} - \text{Damage}_{\text{mitigated}}$$
- **SOP Action**: If $\Delta HP \ge 12$, use the potion immediately. This logic ensures that AI uses resources when facing elite spikes rather than waiting for a Boss that may never come.

### 2. Hard-Gate Validation Layer
Before committing any `end_turn` action, the Go Engine performs a **Memory Alignment Check**:
1. **Fetch**: `get_screen_state` to verify actual Energy/Stars.
2. **Scan**: Verify if any playable cards remain in hand.
3. **Interrupt**: If `Energy > 0` and `Card.Playable == true`, the `end_turn` command is blocked, and the Reasoning loop restarts.

### 3. Necrobinder: Doom Survival Logic
Doom is a处决 (execution) mechanic, but its结算 (settlement) is delayed:
$$\text{Safety} = (\text{Block} + \text{Osty.HP}) > \text{Incoming.Damage}$$
If `Safety` is false, the agent MUST prioritize defense even if the enemy is marked for Doom-death next turn.

---

<a name="chinese"></a>
## 🇨🇳 中文：策略基石

**STS2-Agent-Pro** 的智能核心源自**概率期望 (Probabilistic Expectations)** 与**职业启发式算法 (Class Heuristics)** 的有机结合。

### 1. 药水价值模型 ($\Delta HP$)
为避免“资源囤积综合征”，智能体将药水视为生命值保护等效物：
$$\Delta HP = \text{原始预估伤害} - \text{使用药水后的伤害}$$
- **SOP 动作**: 若 $\Delta HP \ge 12$，立即使用药水。这确保了智能体在面对精英怪的爆发伤害时敢于消耗资源，而非在首领战前意外死亡。

### 2. 硬性指令门控 (Hard-Gate)
在执行 `end_turn` 指令前，Go 引擎会进行物理层的**内存对齐校验**：
1. **抓取**: 调用 `get_screen_state` 确认当前实际能量/星辰。
2. **扫描**: 确认手牌中是否仍有可打出的卡牌。
3. **中断**: 若 `能量 > 0` 且手牌可出，则强行拦截结束回合指令，重新触发推理循环。这杜绝了 LLM 因为计算错误导致的提前空过。

### 3. 各职业核心数学推演
- **死灵缚者 (Necrobinder)**: 将随从 Osty 的生命值视为“永久护甲”。在处理**末日 (Doom)** 处决时，植入存活判定方程：若当回合无法抵挡敌方攻击，即使敌方下回合必死，智能体也会优先叠甲而非叠末日。
- **摄政王 (The Regent)**: 采用星辰 CEC (Converted Energy Cost) 模型：
  $$CEC = \text{卡牌能量消耗} + \frac{\text{星辰消耗}}{2}$$
  这统一了两种资源的价值标尺，使智能体能够量化对比不同卡牌的收益。

### 4. Act 1 贪婪路由演算法 (Greedy Routing)
智能体在 Act 1 初期被赋予了极高的风险偏好：
- **前 3 层**: 强制寻找普通战斗节点，以累积卡牌奖励。
- **中期**: 优先规划“营火 -> 精英 -> 营火”的路径结构，最大化遗物收益。
