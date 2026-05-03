package engine_test

import (
	"math"
	"testing"
	"time"

	"github.com/12jihan/afk-x/internal/content"
	"github.com/12jihan/afk-x/internal/engine"
	"github.com/12jihan/afk-x/internal/game"
)

const floatTol = 1e-9

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) < floatTol
}

// TestTickAppliesRates verifies AC1: each resource increases by rate * delta.Seconds().
func TestTickAppliesRates(t *testing.T) {
	state := game.NewGameState()
	state.Run.Resources[content.CPUCycles] = 0
	state.Run.Resources[content.MemoryShards] = 0
	state.Run.Resources[content.ProcessThreads] = 0

	rates := engine.ResourceRates{
		content.CPUCycles:      1.0,
		content.MemoryShards:   0.5,
		content.ProcessThreads: 0.25,
	}
	delta := 16 * time.Millisecond

	result := engine.Tick(state.Run, delta, rates)

	cases := []struct {
		key  string
		rate float64
	}{
		{content.CPUCycles, 1.0},
		{content.MemoryShards, 0.5},
		{content.ProcessThreads, 0.25},
	}
	for _, c := range cases {
		want := c.rate * delta.Seconds()
		got := result.Resources[c.key]
		if !approxEqual(got, want) {
			t.Errorf("Tick Resources[%q] = %v, want %v (diff %v)", c.key, got, want, math.Abs(got-want))
		}
	}
}

// TestTickDeterministic verifies AC2: identical inputs always produce identical outputs.
func TestTickDeterministic(t *testing.T) {
	state := game.NewGameState()
	state.Run.Resources[content.CPUCycles] = 100.0
	state.Run.Resources[content.MemoryShards] = 50.0

	rates := engine.ResourceRates{
		content.CPUCycles:    1.0,
		content.MemoryShards: 0.5,
	}
	delta := 16 * time.Millisecond

	result1 := engine.Tick(state.Run, delta, rates)
	result2 := engine.Tick(state.Run, delta, rates)

	for key := range rates {
		if result1.Resources[key] != result2.Resources[key] {
			t.Errorf("Tick not deterministic: Resources[%q] first=%v second=%v",
				key, result1.Resources[key], result2.Resources[key])
		}
	}
}

// TestTickPure verifies that Tick does not mutate the original RunState's Resources map.
func TestTickPure(t *testing.T) {
	state := game.NewGameState()
	state.Run.Resources[content.CPUCycles] = 100.0

	originalValue := state.Run.Resources[content.CPUCycles]

	rates := engine.ResourceRates{content.CPUCycles: 1.0}
	_ = engine.Tick(state.Run, 16*time.Millisecond, rates)

	if state.Run.Resources[content.CPUCycles] != originalValue {
		t.Errorf("Tick mutated original Resources map: got %v, want %v",
			state.Run.Resources[content.CPUCycles], originalValue)
	}
}

// TestTickZeroDelta verifies that a zero delta leaves resource values unchanged.
func TestTickZeroDelta(t *testing.T) {
	state := game.NewGameState()
	state.Run.Resources[content.CPUCycles] = 42.0

	rates := engine.ResourceRates{content.CPUCycles: 1.0}
	result := engine.Tick(state.Run, 0, rates)

	if result.Resources[content.CPUCycles] != 42.0 {
		t.Errorf("Tick with zero delta changed resource: got %v, want 42.0",
			result.Resources[content.CPUCycles])
	}
}

// TestTickCmdNonNil verifies AC4: TickCmd returns a non-nil tea.Cmd.
func TestTickCmdNonNil(t *testing.T) {
	cmd := engine.TickCmd()
	if cmd == nil {
		t.Error("TickCmd() returned nil, want non-nil tea.Cmd")
	}
}
