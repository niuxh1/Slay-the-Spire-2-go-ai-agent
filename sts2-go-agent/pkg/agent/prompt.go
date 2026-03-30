package agent

import (
	"fmt"
	"strings"

	"github.com/niuxh/sts2-go-agent/pkg/models"
)

func StateToPrompt(state *models.GameState, lastError string) string {
	var b strings.Builder

	if lastError != "" {
		b.WriteString(fmt.Sprintf("### !! PREVIOUS ACTION FAILED !!\nError: %s\nPlease fix your logic.\n\n", lastError))
	}

	// 注入启发式外脑
	deckAdvice := EvaluateDeck(state.Run)
	if deckAdvice != "" {
		b.WriteString(fmt.Sprintf("### !! HEURISTIC ADVICE !!\n%s\n\n", deckAdvice))
	}

	b.WriteString(fmt.Sprintf("### Current Screen: %s\n", state.Screen))
	
	if state.Run != nil {
		b.WriteString(fmt.Sprintf("**Player Stats**: HP %d/%d, Gold: %d, Floor: %d, Energy: %d/%d\n", 
			state.Run.CurrentHP, state.Run.MaxHP, state.Run.Gold, state.Run.Floor, 
			state.Run.MaxEnergy, state.Run.MaxEnergy))
		
		if len(state.Run.Relics) > 0 {
			relics := []string{}
			for _, r := range state.Run.Relics { relics = append(relics, r.Name) }
			b.WriteString(fmt.Sprintf("**Relics**: %s\n", strings.Join(relics, ", ")))
		}

		if len(state.Run.Potions) > 0 {
			pots := []string{}
			for _, p := range state.Run.Potions {
				if p.Occupied {
					pots = append(pots, fmt.Sprintf("[%d] %s", p.Index, p.Name))
				}
			}
			if len(pots) > 0 {
				b.WriteString(fmt.Sprintf("**Potions**: %s\n", strings.Join(pots, ", ")))
			}
		}

		// 全知视野：让 LLM 看到自己牌库里的所有牌
		if len(state.Run.Deck) > 0 {
			b.WriteString("**Current Deck List**:\n")
			cardCounts := make(map[string]int)
			for _, c := range state.Run.Deck {
				cardCounts[c.Name]++
			}
			for name, count := range cardCounts {
				b.WriteString(fmt.Sprintf("- %dx %s\n", count, name))
			}
		}
	}

	switch state.Screen {
	case "COMBAT":
		if state.Combat != nil {
			c := state.Combat
			b.WriteString(fmt.Sprintf("\n**Combat State**: Current Energy: %d, Stars: %d, Current Block: %d\n", 
				c.Player.Energy, c.Player.Stars, c.Player.Block))
			
			b.WriteString("#### Enemies (Indexes for target_index):\n")
			for _, e := range c.Enemies {
				if !e.IsAlive { continue }
				intentStr := ""
				for _, i := range e.Intents {
					intentStr += fmt.Sprintf("[%s: %s, Damage: %v] ", i.IntentType, i.Label, i.TotalDamage)
				}
				b.WriteString(fmt.Sprintf("- [%d] %s: HP %d/%d, Block %d, Intents: %s\n", 
					e.Index, e.Name, e.CurrentHP, e.MaxHP, e.Block, intentStr))
			}

			b.WriteString("#### Hand Cards (Indexes for card_index):\n")
			for _, card := range c.Hand {
				playable := "NO"
				if card.Playable { playable = "YES" }
				b.WriteString(fmt.Sprintf("- [%d] %s: Cost %d, Playable: %s, Needs Target: %v\n", 
					card.Index, card.Name, card.EnergyCost, playable, card.RequiresTarget))
			}
		}

	case "MAP":
		if state.Map != nil {
			navAdvice := CalculateOptimalPath(state.Map)
			if navAdvice != "" {
				b.WriteString(fmt.Sprintf("\n#### Map Analysis:\n%s\n", navAdvice))
			}

			b.WriteString("\n#### Available Map Nodes (Indexes for option_index):\n")
			for _, node := range state.Map.AvailableNodes {
				b.WriteString(fmt.Sprintf("- [%d] %s (Row: %d, Col: %d)\n", 
					node.Index, node.NodeType, node.Row, node.Col))
			}
		}

	case "REWARD":
		if state.Reward != nil {
			if state.Reward.PendingCardChoice {
				b.WriteString("\n#### Card Reward Selection (Indexes for option_index):\n")
				for _, card := range state.Reward.CardOptions {
					b.WriteString(fmt.Sprintf("- [%d] %s (Cost: %d)\n", card.Index, card.Name, card.EnergyCost))
				}
			} else {
				b.WriteString("\n#### Rewards to Claim (Indexes for option_index):\n")
				for _, r := range state.Reward.Rewards {
					if r.Claimable {
						b.WriteString(fmt.Sprintf("- [%d] %s (%s)\n", r.Index, r.Description, r.RewardType))
					}
				}
			}
		}
	case "SHOP":
		if state.Shop != nil && state.Shop.IsOpen {
			b.WriteString("\n#### Shop Items (Indexes for option_index):\n")
			for _, item := range state.Shop.Cards {
				if item.Available {
					b.WriteString(fmt.Sprintf("- [%d] Card: %s (Price: %d)\n", item.Index, item.Name, item.Price))
				}
			}
		}
	case "REST":
		if state.Rest != nil {
			b.WriteString("\n#### Rest Options (Indexes for option_index):\n")
			for _, opt := range state.Rest.Options {
				if opt.IsEnabled {
					b.WriteString(fmt.Sprintf("- [%d] %s: %s\n", opt.Index, opt.Title, opt.OptionID))
				}
			}
		}
	case "EVENT":
		if state.Event != nil {
			b.WriteString(fmt.Sprintf("\n#### Event: %s\n%s\n\nOptions (Indexes for option_index):\n", state.Event.Title, state.Event.Description))
			for _, opt := range state.Event.Options {
				b.WriteString(fmt.Sprintf("- [%d] %s: %s\n", opt.Index, opt.Title, opt.Description))
			}
		}
	case "CARD_SELECTION":
		if state.Selection != nil {
			b.WriteString(fmt.Sprintf("\n#### Selection Prompt: %s\nOptions (Indexes for option_index):\n", state.Selection.Prompt))
			for _, c := range state.Selection.Cards {
				b.WriteString(fmt.Sprintf("- [%d] %s\n", c.Index, c.Name))
			}
		}
	}

	b.WriteString(fmt.Sprintf("\n**Available Action Names**: %s\n", strings.Join(state.AvailableActions, ", ")))
	return b.String()
}

const SystemPrompt = `You are a supreme Slay the Spire 2 AI Agent. Token limits are not an issue; you must utilize extreme analytical depth to achieve the highest possible win rate.

### CORE OPERATING PRINCIPLES:
1. **DECK ARCHETYPE & SYNERGY**: You have full visibility of your "Current Deck List". Before picking ANY card at a REWARD or SHOP screen, deeply analyze your current deck. DO NOT pick a card unless it actively supports your existing archetype (e.g., Exhaust, Discard, Frost). If none fit, use the "skip_reward_cards" or "proceed" action. Bloated decks lead to death.
2. **COMBAT MATHEMATICS**: You must calculate EXACT lethal thresholds and incoming damage. Do not blindly attack if you can't kill them and they are hitting you hard. Block incoming damage optimally.
3. **EXTREME PACIFIST ROUTING**: Avoid unnecessary combats. Calculate the path to the Boss with the absolute minimum number of Elite and Monster nodes. Prioritize Events (?), Shops, and Rest sites.
4. **POTION ECONOMY**: Use a potion if it saves you >= 12 HP in combat. Do not die hoarding potions.
5. **HARD GATE**: If you have energy and playable cards that can deal damage or block safely, play them before calling 'end_turn'.

### JSON OUTPUT FORMAT:
You must output ONLY a valid JSON object. No markdown tags around it.
{
  "deck_analysis": "Examine current deck synergies. What does the deck need? What should be avoided?",
  "situation_analysis": "Calculate exact incoming damage vs current block. Or evaluate the map/event choices.",
  "reasoning": "Explain step-by-step why the chosen action is mathematically and strategically optimal.",
  "action": "Select EXACTLY ONE from Available Action Names",
  "card_index": integer or null,
  "target_index": integer or null,
  "option_index": integer or null
}`
