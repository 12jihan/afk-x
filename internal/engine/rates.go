package engine

import (
	"github.com/12jihan/afk-x/internal/content"
	"github.com/12jihan/afk-x/internal/game"
)

// ResourceRates maps resource keys to their generation rate (units per second).
// Always use this named type — never map[string]float64 in engine function signatures.
type ResourceRates map[string]float64

// ComputeRates returns the current resource generation rates for a run.
// For MVP (no upgrades or perks active) this equals content.BaseRates.
// Returns a defensive copy — callers may mutate the result safely without
// affecting content.BaseRates (resolves Story 1.3 deferred: BaseRates mutation risk).
func ComputeRates(run game.RunState) ResourceRates {
	rates := make(ResourceRates, len(content.BaseRates))
	for k, v := range content.BaseRates {
		rates[k] = v
	}
	// TODO: Story 2.x — apply upgrade multipliers from run.Upgrades
	// TODO: Story 3.x — apply active perk bonuses from run.ActivePerks
	return rates
}
