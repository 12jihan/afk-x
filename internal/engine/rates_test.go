package engine_test

import (
	"testing"

	"github.com/12jihan/afk-x/internal/content"
	"github.com/12jihan/afk-x/internal/engine"
	"github.com/12jihan/afk-x/internal/game"
)

// TestComputeRatesDefaultState verifies AC3: ComputeRates with a fresh RunState
// (no upgrades, no active perks) returns values equal to content.BaseRates.
func TestComputeRatesDefaultState(t *testing.T) {
	state := game.NewGameState()
	rates := engine.ComputeRates(state.Run)

	for key, want := range content.BaseRates {
		got, ok := rates[key]
		if !ok {
			t.Errorf("ComputeRates result missing key %q", key)
			continue
		}
		if got != want {
			t.Errorf("ComputeRates[%q] = %v, want %v", key, got, want)
		}
	}
}

// TestComputeRatesAllKeysPresent verifies the result contains all three resource keys.
func TestComputeRatesAllKeysPresent(t *testing.T) {
	state := game.NewGameState()
	rates := engine.ComputeRates(state.Run)

	required := []string{
		content.CPUCycles,
		content.MemoryShards,
		content.ProcessThreads,
	}
	for _, key := range required {
		if _, ok := rates[key]; !ok {
			t.Errorf("ComputeRates result missing required key %q", key)
		}
	}
}

// TestComputeRatesReturnsCopy verifies that mutating the returned ResourceRates
// does not affect content.BaseRates (defensive copy — resolves Story 1.3 deferred item).
func TestComputeRatesReturnsCopy(t *testing.T) {
	state := game.NewGameState()
	rates := engine.ComputeRates(state.Run)

	originalCPU := content.BaseRates[content.CPUCycles]

	// Mutate the returned map
	rates[content.CPUCycles] = 9999.0

	// content.BaseRates must be unchanged
	if content.BaseRates[content.CPUCycles] != originalCPU {
		t.Errorf("mutating ComputeRates result corrupted content.BaseRates: got %v, want %v",
			content.BaseRates[content.CPUCycles], originalCPU)
	}
}
