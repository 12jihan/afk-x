package engine

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/12jihan/afk-x/internal/game"
)

// TickMsg is the Bubbletea message delivered by TickCmd every 16ms.
// ui.Update handles TickMsg by calling Tick() and updating the model state.
type TickMsg time.Time

// TickCmd returns a Bubbletea command that fires a TickMsg every 16ms.
// Called from ui.Init() and returned from ui.Update() after each tick to
// reschedule the next tick. This is the only Bubbletea dependency in the engine.
func TickCmd() tea.Cmd {
	return tea.Every(16*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Tick advances resources by delta time at the given rates.
// Pure function — returns a new RunState with updated Resources.
// The original RunState's Resources map is never mutated.
//
// Resource update formula: newValue = oldValue + (rate * delta.Seconds())
func Tick(state game.RunState, delta time.Duration, rates ResourceRates) game.RunState {
	// Copy the Resources map to avoid mutating the caller's map through the shared reference.
	// game.RunState is passed by value but its Resources field is map[string]float64 — a reference type.
	newResources := make(map[string]float64, len(state.Resources))
	for k, v := range state.Resources {
		newResources[k] = v
	}
	for resource, rate := range rates {
		newResources[resource] += rate * delta.Seconds()
	}
	state.Resources = newResources
	return state
}
