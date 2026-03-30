package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/niuxh/sts2-go-agent/pkg/api"
	"github.com/niuxh/sts2-go-agent/pkg/models"
)

type LLMProvider interface {
	Chat(systemPrompt, userPrompt string) (string, error)
}

type Agent struct {
	Client     *api.Client
	LLM        LLMProvider
	LogHandler func(string)
}

func NewAgent(client *api.Client, llm LLMProvider, logHandler func(string)) *Agent {
	return &Agent{
		Client:     client,
		LLM:        llm,
		LogHandler: logHandler,
	}
}

func (a *Agent) logger(msg string) {
	if a.LogHandler != nil {
		a.LogHandler(fmt.Sprintf("[%s] %s\n", time.Now().Format("15:04:05"), msg))
	}
	log.Println(msg)
}

type AgentDecision struct {
	DeckAnalysis      string `json:"deck_analysis"`
	SituationAnalysis string `json:"situation_analysis"`
	Reasoning         string `json:"reasoning"`
	Action            string `json:"action"`
	CardIndex         *int   `json:"card_index"`
	TargetIndex       *int   `json:"target_index"`
	OptionIndex       *int   `json:"option_index"`
}

func (a *Agent) Run(ctx context.Context) {
	a.logger("Super-Heuristic Single-Core Agent Initialized.")
	
	lastError := ""
	lastPrompt := ""

	for {
		select {
		case <-ctx.Done():
			a.logger("Agent Shutdown signal received.")
			return
		default:
			state, err := a.Client.GetState()
			if err != nil {
				a.logger(fmt.Sprintf("State Sync Error: %v", err))
				time.Sleep(2 * time.Second)
				continue
			}

			if state.Screen == "UNKNOWN" || len(state.AvailableActions) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}

			prompt := StateToPrompt(state, lastError)
			
			// 去重，防止重复请求
			if prompt == lastPrompt && lastError == "" {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			lastPrompt = prompt
			lastError = ""
			
			a.logger(fmt.Sprintf("Analyzing Screen: %s | Energy: %d", state.Screen, getEnergy(state)))
			
			respStr, err := a.LLM.Chat(SystemPrompt, prompt)
			if err != nil {
				a.logger(fmt.Sprintf("LLM Error: %v", err))
				time.Sleep(5 * time.Second)
				continue
			}

			var decision AgentDecision
			if err := json.Unmarshal([]byte(respStr), &decision); err != nil {
				lastError = "JSON Parse Error. Output must be strictly valid JSON without Markdown blocks."
				a.logger(fmt.Sprintf("Parser Error: %v", err))
				continue
			}

			a.logger(fmt.Sprintf("[SYNERGY] %s", decision.DeckAnalysis))
			a.logger(fmt.Sprintf("[LOGIC] %s", decision.Reasoning))
			
			// 硬性门控
			if decision.Action == "end_turn" {
				verifyState, _ := a.Client.GetState()
				if getEnergy(verifyState) > 0 && hasPlayableCards(verifyState) {
					a.logger("[SYSTEM] Gate REJECTED: Attempted end_turn with playable cards.")
					lastError = "System Feedback: You attempted to end_turn but you still have playable cards and energy. Play them!"
					continue
				}
			}

			a.logger(fmt.Sprintf(">> Executing: %s", decision.Action))
			actionReq := &models.ActionRequest{
				Action:      decision.Action,
				CardIndex:   decision.CardIndex,
				TargetIndex: decision.TargetIndex,
				OptionIndex: decision.OptionIndex,
			}

			actionResp, err := a.Client.SendAction(actionReq)
			if err != nil {
				a.logger(fmt.Sprintf("Execution Error: %v", err))
				lastError = fmt.Sprintf("Action execution failed: %v", err)
				lastPrompt = ""
				continue
			}

			if actionResp.Status == "pending" {
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}

func getEnergy(s *models.GameState) int {
	if s.Combat != nil {
		return s.Combat.Player.Energy
	}
	return 0
}

func hasPlayableCards(s *models.GameState) bool {
	if s.Combat == nil { return false }
	for _, c := range s.Combat.Hand {
		if c.Playable { return true }
	}
	return false
}
