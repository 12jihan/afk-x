package engine_test

import (
	"testing"

	"github.com/12jihan/afk-x/internal/content"
	"github.com/12jihan/afk-x/internal/engine"
	"github.com/12jihan/afk-x/internal/game"
)

func newFloorRun(floor int, cpu float64) game.RunState {
	s := game.NewGameState()
	s.Run.Floor = floor
	s.Run.Resources[content.CPUCycles] = cpu
	return s.Run
}

func TestFloorThreshold_Floor1(t *testing.T) {
	if engine.FloorThreshold(1) != 100.0 {
		t.Errorf("FloorThreshold(1) = %v, want 100.0", engine.FloorThreshold(1))
	}
}

func TestFloorThreshold_ScalesExponentially(t *testing.T) {
	cases := []struct{ floor int; want float64 }{
		{2, 200.0},
		{3, 400.0},
		{4, 800.0},
	}
	for _, c := range cases {
		got := engine.FloorThreshold(c.floor)
		if got != c.want {
			t.Errorf("FloorThreshold(%d) = %v, want %v", c.floor, got, c.want)
		}
	}
}

func TestFloorProgress_BelowThreshold(t *testing.T) {
	run := newFloorRun(1, 50)
	got := engine.FloorProgress(run)
	if got != 0.5 {
		t.Errorf("FloorProgress with 50 CPU on floor 1 = %v, want 0.5", got)
	}
}

func TestFloorProgress_ClampsAtOne(t *testing.T) {
	run := newFloorRun(1, 9999)
	got := engine.FloorProgress(run)
	if got != 1.0 {
		t.Errorf("FloorProgress with 9999 CPU on floor 1 = %v, want 1.0", got)
	}
}

func TestCheckFloorClear_NotCleared(t *testing.T) {
	run := newFloorRun(1, 99)
	if engine.CheckFloorClear(run) {
		t.Error("CheckFloorClear with 99 CPU on floor 1 should return false")
	}
}

func TestCheckFloorClear_ExactlyCleared(t *testing.T) {
	run := newFloorRun(1, 100)
	if !engine.CheckFloorClear(run) {
		t.Error("CheckFloorClear with exactly 100 CPU on floor 1 should return true")
	}
}

func TestAdvanceFloor_IncrementsFloor(t *testing.T) {
	run := newFloorRun(1, 100)
	updated := engine.AdvanceFloor(run)
	if updated.Floor != 2 {
		t.Errorf("AdvanceFloor: floor = %d, want 2", updated.Floor)
	}
}

func TestAdvanceFloor_DeductsCost(t *testing.T) {
	run := newFloorRun(1, 250)
	updated := engine.AdvanceFloor(run)
	want := 150.0 // 250 - 100
	if updated.Resources[content.CPUCycles] != want {
		t.Errorf("AdvanceFloor: CPU = %v, want %v", updated.Resources[content.CPUCycles], want)
	}
}

func TestAdvanceFloor_ClampsToZero(t *testing.T) {
	run := newFloorRun(1, 100)
	updated := engine.AdvanceFloor(run)
	if updated.Resources[content.CPUCycles] != 0 {
		t.Errorf("AdvanceFloor: CPU = %v, want 0", updated.Resources[content.CPUCycles])
	}
}

func TestAdvanceFloor_DoesNotMutateOriginal(t *testing.T) {
	run := newFloorRun(1, 250)
	originalCPU := run.Resources[content.CPUCycles]
	originalFloor := run.Floor

	engine.AdvanceFloor(run)

	if run.Resources[content.CPUCycles] != originalCPU {
		t.Error("AdvanceFloor must not mutate the caller's Resources map")
	}
	if run.Floor != originalFloor {
		t.Error("AdvanceFloor must not mutate the caller's Floor field")
	}
}
