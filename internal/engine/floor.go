package engine

import (
	"math"

	"github.com/12jihan/afk-x/internal/content"
	"github.com/12jihan/afk-x/internal/game"
)

// FloorThreshold returns the CPU Cycles required to clear the given floor.
// Uses exponential scaling: 100 * 2^(floor-1).
// Floor 1 = 100, Floor 2 = 200, Floor 3 = 400, Floor 4 = 800, ...
func FloorThreshold(floor int) float64 {
	return 100.0 * math.Pow(2.0, float64(floor-1))
}

// FloorProgress returns the player's progress toward clearing the current floor
// as a value in [0.0, 1.0]. Returns 1.0 if already at or past the threshold.
func FloorProgress(run game.RunState) float64 {
	threshold := FloorThreshold(run.Floor)
	if threshold == 0 {
		return 1.0
	}
	p := run.Resources[content.CPUCycles] / threshold
	if p > 1.0 {
		p = 1.0
	}
	return p
}

// CheckFloorClear reports whether the player has accumulated enough CPU Cycles
// to clear the current floor.
func CheckFloorClear(run game.RunState) bool {
	return run.Resources[content.CPUCycles] >= FloorThreshold(run.Floor)
}

// AdvanceFloor deducts the current floor's CPU cost, increments the floor counter,
// and returns the updated RunState. Deep-copies the Resources map before mutation
// to avoid aliasing the caller's state.
func AdvanceFloor(run game.RunState) game.RunState {
	threshold := FloorThreshold(run.Floor)
	newResources := make(map[string]float64, len(run.Resources))
	for k, v := range run.Resources {
		newResources[k] = v
	}
	cpu := newResources[content.CPUCycles] - threshold
	if cpu < 0 {
		cpu = 0
	}
	newResources[content.CPUCycles] = cpu
	run.Resources = newResources
	run.Floor++
	return run
}
