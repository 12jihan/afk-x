package engine_test

import (
	"math"
	"testing"

	"github.com/12jihan/afk-x/internal/content"
	"github.com/12jihan/afk-x/internal/engine"
	"github.com/12jihan/afk-x/internal/game"
)

const floatTol = 1e-9

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) < floatTol
}

// TestTickAppliesRates verifies AC1: each resource increases by rate * delta.Seconds().
// Hardcoded oracle values prevent the test from masking formula bugs:
//   rate=1.0, delta=16ms → 1.0 * 0.016 = 0.016
//   rate=0.5, delta=16ms → 0.5 * 0.016 = 0.008
//   rate=0.25, delta=16ms → 0.25 * 0.016 = 0.004
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
	delta := engine.TickInterval // 16ms

	result := engine.Tick(state.Run, delta, rates)

	cases := []struct {
		key  string
		want float64 // hardcoded oracle: rate * 0.016
	}{
		{content.CPUCycles, 0.016},
		{content.MemoryShards, 0.008},
		{content.ProcessThreads, 0.004},
	}
	for _, c := range cases {
		got := result.Resources[c.key]
		if !approxEqual(got, c.want) {
			t.Errorf("Tick Resources[%q] = %v, want %v (diff %v)", c.key, got, c.want, math.Abs(got-c.want))
		}
	}
}

// TestTickDeterministic verifies AC2: identical inputs always produce identical outputs.
// Also verifies that keys not in rates (ProcessThreads) are preserved unchanged
// in both results — Tick must not drop or alter unrated resources.
func TestTickDeterministic(t *testing.T) {
	state := game.NewGameState()
	state.Run.Resources[content.CPUCycles] = 100.0
	state.Run.Resources[content.MemoryShards] = 50.0
	state.Run.Resources[content.ProcessThreads] = 7.0 // not in rates — must be preserved

	rates := engine.ResourceRates{
		content.CPUCycles:    1.0,
		content.MemoryShards: 0.5,
	}
	delta := engine.TickInterval // 16ms

	result1 := engine.Tick(state.Run, delta, rates)
	result2 := engine.Tick(state.Run, delta, rates)

	// Rate-tracked keys must be identical across calls.
	for key := range rates {
		if result1.Resources[key] != result2.Resources[key] {
			t.Errorf("Tick not deterministic: Resources[%q] first=%v second=%v",
				key, result1.Resources[key], result2.Resources[key])
		}
	}
	// Non-rate key must be preserved unchanged in both results.
	if result1.Resources[content.ProcessThreads] != 7.0 {
		t.Errorf("Tick modified non-rate key in result1: ProcessThreads = %v, want 7.0",
			result1.Resources[content.ProcessThreads])
	}
	if result2.Resources[content.ProcessThreads] != 7.0 {
		t.Errorf("Tick modified non-rate key in result2: ProcessThreads = %v, want 7.0",
			result2.Resources[content.ProcessThreads])
	}
}

// TestTickPure verifies that Tick does not mutate the original RunState's Resources map.
func TestTickPure(t *testing.T) {
	state := game.NewGameState()
	state.Run.Resources[content.CPUCycles] = 100.0

	originalValue := state.Run.Resources[content.CPUCycles]

	rates := engine.ResourceRates{content.CPUCycles: 1.0}
	_ = engine.Tick(state.Run, engine.TickInterval, rates)

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

// TestTickPreservesUnratedKeys verifies that resource keys present in state.Resources
// but absent from rates are copied unchanged into the returned state.
// This is Tick's documented contract for forward-compatibility with future resource types.
func TestTickPreservesUnratedKeys(t *testing.T) {
	state := game.NewGameState()
	state.Run.Resources[content.CPUCycles] = 10.0
	state.Run.Resources[content.MemoryShards] = 20.0
	state.Run.Resources[content.ProcessThreads] = 99.0 // not in rates below

	// Only rate cpu_cycles — memory_shards and process_threads must survive unchanged.
	rates := engine.ResourceRates{
		content.CPUCycles: 1.0,
	}

	result := engine.Tick(state.Run, engine.TickInterval, rates)

	if result.Resources[content.MemoryShards] != 20.0 {
		t.Errorf("Tick dropped/changed MemoryShards (not in rates): got %v, want 20.0",
			result.Resources[content.MemoryShards])
	}
	if result.Resources[content.ProcessThreads] != 99.0 {
		t.Errorf("Tick dropped/changed ProcessThreads (not in rates): got %v, want 99.0",
			result.Resources[content.ProcessThreads])
	}
}
