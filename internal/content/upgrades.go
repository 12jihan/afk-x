package content

import "math"

// UpgradeDefinition describes a purchasable improvement to resource generation.
type UpgradeDefinition struct {
	ID          string
	Name        string
	Description string
	BaseCost    map[string]float64 // resource costs at level 0 → 1
	CostScaling float64            // each level costs BaseCost * CostScaling^currentLevel
	BonusType   string             // "rate_add": adds BonusValue units/s per level to BonusTarget
	BonusTarget string             // resource key from content.* constants
	BonusValue  float64            // per-level flat bonus
	MaxLevel    int
}

// Upgrades is the ordered list of purchasable upgrades available each run.
var Upgrades = []UpgradeDefinition{
	{
		ID:          "overclock",
		Name:        "Overclock",
		Description: "+1.0 CPU cycle/s per level",
		BaseCost:    map[string]float64{CPUCycles: 50},
		CostScaling: 1.5,
		BonusType:   "rate_add",
		BonusTarget: CPUCycles,
		BonusValue:  1.0,
		MaxLevel:    10,
	},
	{
		ID:          "cache_expander",
		Name:        "Cache Expander",
		Description: "+0.5 memory shard/s per level",
		BaseCost:    map[string]float64{CPUCycles: 100},
		CostScaling: 1.5,
		BonusType:   "rate_add",
		BonusTarget: MemoryShards,
		BonusValue:  0.5,
		MaxLevel:    10,
	},
	{
		ID:          "thread_pool",
		Name:        "Thread Pool",
		Description: "+0.25 process thread/s per level",
		BaseCost:    map[string]float64{CPUCycles: 200},
		CostScaling: 1.5,
		BonusType:   "rate_add",
		BonusTarget: ProcessThreads,
		BonusValue:  0.25,
		MaxLevel:    10,
	},
}

// UpgradeByID returns the UpgradeDefinition for the given ID and ok=true.
// Returns a zero value and ok=false if not found.
func UpgradeByID(id string) (UpgradeDefinition, bool) {
	for _, def := range Upgrades {
		if def.ID == id {
			return def, true
		}
	}
	return UpgradeDefinition{}, false
}

// ScaledCost returns the resource costs to purchase the next level of def,
// given the current level. Cost grows as BaseCost * CostScaling^currentLevel,
// rounded up to the nearest integer unit.
func ScaledCost(def UpgradeDefinition, currentLevel int) map[string]float64 {
	cost := make(map[string]float64, len(def.BaseCost))
	multiplier := math.Pow(def.CostScaling, float64(currentLevel))
	for k, v := range def.BaseCost {
		cost[k] = math.Ceil(v * multiplier)
	}
	return cost
}
