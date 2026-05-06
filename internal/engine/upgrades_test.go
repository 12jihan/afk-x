package engine_test

import (
	"testing"

	"github.com/12jihan/afk-x/internal/content"
	"github.com/12jihan/afk-x/internal/engine"
	"github.com/12jihan/afk-x/internal/game"
)

func newRunWith(cpu, mem, threads float64) game.RunState {
	s := game.NewGameState()
	s.Run.Resources[content.CPUCycles] = cpu
	s.Run.Resources[content.MemoryShards] = mem
	s.Run.Resources[content.ProcessThreads] = threads
	return s.Run
}

func TestCanPurchase_Affordable(t *testing.T) {
	run := newRunWith(100, 0, 0)
	if !engine.CanPurchase(run, "overclock") {
		t.Error("expected CanPurchase true when resources are sufficient")
	}
}

func TestCanPurchase_Insufficient(t *testing.T) {
	run := newRunWith(10, 0, 0)
	if engine.CanPurchase(run, "overclock") {
		t.Error("expected CanPurchase false when resources are insufficient")
	}
}

func TestCanPurchase_UnknownUpgrade(t *testing.T) {
	run := newRunWith(99999, 99999, 99999)
	if engine.CanPurchase(run, "does_not_exist") {
		t.Error("expected CanPurchase false for unknown upgrade ID")
	}
}

func TestCanPurchase_AtMaxLevel(t *testing.T) {
	run := newRunWith(99999, 99999, 99999)
	def, _ := content.UpgradeByID("overclock")
	run.Upgrades["overclock"] = def.MaxLevel
	if engine.CanPurchase(run, "overclock") {
		t.Error("expected CanPurchase false when upgrade is at MaxLevel")
	}
}

func TestPurchaseUpgrade_IncrementsLevel(t *testing.T) {
	run := newRunWith(100, 0, 0)
	updated, ok := engine.PurchaseUpgrade(run, "overclock")
	if !ok {
		t.Fatal("expected PurchaseUpgrade to succeed")
	}
	if updated.Upgrades["overclock"] != 1 {
		t.Errorf("expected overclock level 1, got %d", updated.Upgrades["overclock"])
	}
}

func TestPurchaseUpgrade_DeductsResources(t *testing.T) {
	run := newRunWith(100, 0, 0)
	updated, ok := engine.PurchaseUpgrade(run, "overclock")
	if !ok {
		t.Fatal("expected PurchaseUpgrade to succeed")
	}
	if updated.Resources[content.CPUCycles] >= 100 {
		t.Error("expected CPU cycles to be deducted after purchase")
	}
}

func TestPurchaseUpgrade_FailsWhenInsufficient(t *testing.T) {
	run := newRunWith(10, 0, 0)
	_, ok := engine.PurchaseUpgrade(run, "overclock")
	if ok {
		t.Error("expected PurchaseUpgrade to fail when resources are insufficient")
	}
}

func TestPurchaseUpgrade_DoesNotMutateOriginal(t *testing.T) {
	run := newRunWith(100, 0, 0)
	originalCPU := run.Resources[content.CPUCycles]
	originalLevel := run.Upgrades["overclock"]

	engine.PurchaseUpgrade(run, "overclock")

	if run.Resources[content.CPUCycles] != originalCPU {
		t.Error("PurchaseUpgrade must not mutate the caller's Resources map")
	}
	if run.Upgrades["overclock"] != originalLevel {
		t.Error("PurchaseUpgrade must not mutate the caller's Upgrades map")
	}
}

func TestComputeRates_WithUpgrade(t *testing.T) {
	run := newRunWith(0, 0, 0)
	run.Upgrades["overclock"] = 2
	rates := engine.ComputeRates(run)

	def, _ := content.UpgradeByID("overclock")
	want := content.BaseRates[content.CPUCycles] + def.BonusValue*2
	if rates[content.CPUCycles] != want {
		t.Errorf("ComputeRates CPU with overclock Lv.2 = %v, want %v", rates[content.CPUCycles], want)
	}
}

func TestComputeRates_ZeroLevelUpgradeIgnored(t *testing.T) {
	run := newRunWith(0, 0, 0)
	run.Upgrades["overclock"] = 0
	rates := engine.ComputeRates(run)

	want := content.BaseRates[content.CPUCycles]
	if rates[content.CPUCycles] != want {
		t.Errorf("ComputeRates CPU with overclock Lv.0 = %v, want base %v", rates[content.CPUCycles], want)
	}
}
