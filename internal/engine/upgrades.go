package engine

import (
	"github.com/12jihan/afk-x/internal/content"
	"github.com/12jihan/afk-x/internal/game"
)

// CanPurchase reports whether run has sufficient resources to buy the next level
// of upgradeID. Returns false if the upgrade is unknown or already at MaxLevel.
func CanPurchase(run game.RunState, upgradeID string) bool {
	def, ok := content.UpgradeByID(upgradeID)
	if !ok {
		return false
	}
	currentLevel := run.Upgrades[upgradeID]
	if currentLevel >= def.MaxLevel {
		return false
	}
	cost := content.ScaledCost(def, currentLevel)
	for resource, amount := range cost {
		if run.Resources[resource] < amount {
			return false
		}
	}
	return true
}

// PurchaseUpgrade deducts costs and increments the upgrade level.
// Returns the updated RunState and true on success.
// Returns the original state and false if the upgrade is unknown, at MaxLevel, or unaffordable.
// Deep-copies the Resources and Upgrades maps before mutation to avoid aliasing the caller's state.
func PurchaseUpgrade(run game.RunState, upgradeID string) (game.RunState, bool) {
	if !CanPurchase(run, upgradeID) {
		return run, false
	}

	def, _ := content.UpgradeByID(upgradeID)
	cost := content.ScaledCost(def, run.Upgrades[upgradeID])

	newResources := make(map[string]float64, len(run.Resources))
	for k, v := range run.Resources {
		newResources[k] = v
	}
	newUpgrades := make(map[string]int, len(run.Upgrades)+1)
	for k, v := range run.Upgrades {
		newUpgrades[k] = v
	}
	run.Resources = newResources
	run.Upgrades = newUpgrades

	for resource, amount := range cost {
		run.Resources[resource] -= amount
	}
	run.Upgrades[upgradeID]++
	return run, true
}
