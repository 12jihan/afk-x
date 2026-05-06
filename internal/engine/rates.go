package engine

import (
	"github.com/12jihan/afk-x/internal/content"
	"github.com/12jihan/afk-x/internal/game"
)

// ResourceRates maps resource keys to their generation rate (units per second).
// Always use this named type — never map[string]float64 in engine function signatures.
type ResourceRates map[string]float64

// ComputeRates returns the current resource generation rates for a run.
// Starts from content.BaseRates, then adds flat bonuses from active upgrades.
// Returns a defensive copy — callers may mutate the result safely without
// affecting content.BaseRates (resolves Story 1.3 deferred: BaseRates mutation risk).
func ComputeRates(run game.RunState) ResourceRates {
	rates := make(ResourceRates, len(content.BaseRates))
	for k, v := range content.BaseRates {
		rates[k] = v
	}
	for id, level := range run.Upgrades {
		if level == 0 {
			continue
		}
		def, ok := content.UpgradeByID(id)
		if !ok {
			continue
		}
		if def.BonusType == "rate_add" {
			rates[def.BonusTarget] += def.BonusValue * float64(level)
		}
	}
	// TODO: Story 3.x — apply active perk bonuses from run.ActivePerks
	return rates
}
