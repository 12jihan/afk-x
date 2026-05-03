package engine_test

import (
	"testing"

	"github.com/12jihan/afk-x/internal/content"
	"github.com/12jihan/afk-x/internal/engine"
	"github.com/12jihan/afk-x/internal/game"
)

// TestComputeRatesDefaultState verifies AC3: ComputeRates with a fresh RunState
// (no upgrades, no active perks) returns values equal to content.BaseRates.
// Also checks the result contains no extra keys beyond content.BaseRates.
func TestComputeRatesDefaultState(t *testing.T) {
	state := game.NewGameState()
	rates := engine.ComputeRates(state.Run)

	if len(rates) != len(content.BaseRates) {
		t.Errorf("ComputeRates result has %d keys, want %d (content.BaseRates length)",
			len(rates), len(content.BaseRates))
	}
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
// Baseline is captured before ComputeRates is called so the assertion is non-vacuous:
// if ComputeRates itself mutated BaseRates, the snapshot would catch it.
func TestComputeRatesReturnsCopy(t *testing.T) {
	// Snapshot BEFORE the call under test.
	originalCPU := content.BaseRates[content.CPUCycles]

	state := game.NewGameState()
	rates := engine.ComputeRates(state.Run)

	// Mutate the returned map.
	rates[content.CPUCycles] = 9999.0

	// content.BaseRates must be unchanged.
	if content.BaseRates[content.CPUCycles] != originalCPU {
		t.Errorf("mutating ComputeRates result corrupted content.BaseRates: got %v, want %v",
			content.BaseRates[content.CPUCycles], originalCPU)
	}
}
