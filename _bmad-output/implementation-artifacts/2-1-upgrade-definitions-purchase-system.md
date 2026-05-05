# Story 2.1: Upgrade Definitions & Purchase System

Status: done

## Story

As a player,
I want to view available upgrades and spend my resources to purchase them,
so that I can improve my resource generation and progress faster.

## Acceptance Criteria

1. **Given** the game screen is showing, **When** the player views the upgrade panel, **Then** available upgrades are listed with their names, costs, and current level
2. **Given** a player has sufficient resources for an upgrade, **When** they press the upgrade's number key (1-3), **Then** the resources are deducted and the upgrade level increments
3. **Given** a player has insufficient resources for an upgrade, **When** they attempt to purchase it, **Then** the purchase is rejected with a status message and no resources are deducted
4. **Given** an upgrade is purchased, **When** `engine.ComputeRates()` is called, **Then** the returned rates reflect the upgrade's effect on resource generation

## Tasks / Subtasks

- [x] Task 1: Create `internal/content/upgrades.go` — UpgradeDefinition struct, Upgrades pool, UpgradeByID, ScaledCost (AC: 1, 4)
- [x] Task 2: Create `internal/engine/upgrades.go` — CanPurchase, PurchaseUpgrade with defensive map copy (AC: 2, 3)
- [x] Task 3: Update `internal/engine/rates.go` — apply rate_add upgrade bonuses in ComputeRates (AC: 4)
- [x] Task 4: Create `internal/engine/upgrades_test.go` — 10 tests covering all ACs (AC: 1-4)
- [x] Task 5: Add UpgradePanel style to `internal/ui/styles/layout.go` (AC: 1)
- [x] Task 6: Update `internal/ui/screens/game.go` — add upgradeView, formatCost, canAfford; update GameView signature (AC: 1, 2, 3)
- [x] Task 7: Update `internal/ui/model.go` — add StatusMsg field, clearStatusMsg, clearStatusAfterCmd (AC: 2, 3)
- [x] Task 8: Update `internal/ui/update.go` — handle "1"/"2"/"3" key presses and clearStatusMsg (AC: 2, 3)
- [x] Task 9: Update `internal/ui/view.go` — pass upgrades and statusMsg to GameView (AC: 1)
- [x] Task 10: Verify build and test suite (all ACs)

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-6

### Completion Notes List

- 3 upgrades defined: Overclock (+1.0 CPU/s), Cache Expander (+0.5 MEM/s), Thread Pool (+0.25 THR/s) — all cost CPU Cycles only with 1.5x scaling per level, MaxLevel 10
- PurchaseUpgrade deep-copies both Resources and Upgrades maps before mutation — same pattern as Tick() — prevents caller aliasing
- ComputeRates now applies "rate_add" bonuses from run.Upgrades; zero-level entries are skipped
- Affordable upgrades render in Accent (green), unaffordable in Muted (gray) — NO_COLOR handled transparently by lipgloss
- Number keys 1-3 trigger purchase on GameScreen; status message auto-clears after 2 seconds via clearStatusAfterCmd
- 19 engine tests passing (10 new upgrade tests + 9 pre-existing), build and vet clean

### File List

- `internal/content/upgrades.go` — new
- `internal/engine/upgrades.go` — new
- `internal/engine/upgrades_test.go` — new
- `internal/engine/rates.go` — updated (upgrade bonus application)
- `internal/ui/model.go` — updated (StatusMsg, clearStatusMsg, clearStatusAfterCmd)
- `internal/ui/update.go` — updated (1/2/3 key handling, clearStatusMsg)
- `internal/ui/styles/layout.go` — updated (UpgradePanel style)
- `internal/ui/screens/game.go` — updated (upgradeView, updated GameView signature)
- `internal/ui/view.go` — updated (pass upgrades + statusMsg)

### Change Log

- 2026-05-04: Story 2.1 complete — upgrade definitions, purchase system, and TUI panel
