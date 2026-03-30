package agent
import (
	"strings"

	"github.com/niuxh/sts2-go-agent/pkg/models"
)

// EvaluateDeck analyzes the current deck and relics to provide a strong bias for the LLM.
func EvaluateDeck(runState *models.RunState) string {
	if runState == nil || len(runState.Deck) == 0 {
		return ""
	}

	exhaustCount := 0
	discardCount := 0
	orbCount := 0
	starCount := 0
	zeroCostCount := 0

	for _, c := range runState.Deck {
		// Basic heuristic matching based on card names/types (could be expanded)
		name := strings.ToLower(c.Name)
		if strings.Contains(name, "exhaust") || strings.Contains(name, "burn") {
			exhaustCount++
		}
		if strings.Contains(name, "discard") || strings.Contains(name, "sly") {
			discardCount++
		}
		if strings.Contains(name, "channel") || strings.Contains(name, "orb") {
			orbCount++
		}
		if strings.Contains(name, "star") || strings.Contains(name, "light") || strings.Contains(name, "radiate") {
			starCount++
		}
		if c.EnergyCost == 0 {
			zeroCostCount++
		}
	}

	var insights []string

	if exhaustCount >= 3 {
		insights = append(insights, "Exhaust Synergy Active. Value Exhaust/Status cards highly.")
	}
	if discardCount >= 3 {
		insights = append(insights, "Discard/Sly Synergy Active. Prioritize cards that discard or trigger when discarded.")
	}
	if orbCount >= 3 {
		insights = append(insights, "Orb/Focus Synergy Active. Prioritize Focus and Frost/Lightning generation.")
	}
	if starCount >= 3 {
		insights = append(insights, "Star Economy Active. Prioritize Star generation and burst damage over Forge.")
	}
	if zeroCostCount >= 5 {
		insights = append(insights, "Low-Cost Engine Active. Card draw is extremely valuable.")
	}

	if len(insights) == 0 {
		return "DECKMASTER: You are in the early game. Focus on raw damage and basic block. Do not force synergies yet."
	}

	return "DECKMASTER ALGORITHM DETECTED:\n- " + strings.Join(insights, "\n- ")
}
