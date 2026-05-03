package engine

import (
	"time"

	"github.com/12jihan/afk-x/internal/game"
	tea "github.com/charmbracelet/bubbletea"
)

// TickInterval is the target duration between game ticks (~62.5 Hz).
// Used by TickCmd and referenced in tests to avoid duplicating the magic number.
const TickInterval = 16 * time.Millisecond

// TickMsg is the Bubbletea message delivered by TickCmd every TickInterval.
// ui.Update handles TickMsg by calling Tick() and updating the model state.
type TickMsg time.Time

// TickCmd returns a Bubbletea command that fires a TickMsg every TickInterval.
// Called from ui.Init() and returned from ui.Update() after each tick to
// reschedule the next tick. This is the only Bubbletea dependency in the engine.
func TickCmd() tea.Cmd {
	return tea.Every(TickInterval, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Tick advances resources by delta time at the given rates.
// Pure function — returns a new RunState with updated Resources.
// The original RunState's Resources map is never mutated.
//
// Resource update formula: newValue = oldValue + (rate * delta.Seconds())
//
// Behavioral notes:
//   - Keys present in rates but absent from state.Resources are created with a
//     base of 0 before the increment is applied (Go map zero-value).
//   - Keys present in state.Resources but absent from rates are copied unchanged.
//   - A nil state.Resources is safe: the copy loop is a no-op and rates keys are
//     created from zero as above.
//   - Precondition: delta must be >= 0. A negative delta silently decrements
//     resources; callers are responsible for clamping before passing.
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
