# Story 2.2: Floor Progression & Daemon Encounters

Status: review

## Story

As a player,
I want to advance through tower floors by defeating daemon enemies,
so that I can climb the tower and unlock new content.

## Acceptance Criteria

1. **Given** the player's resources and upgrades meet a floor's completion threshold, **When** `engine.CheckFloorClear()` is evaluated on a tick, **Then** a `FloorClearMsg` is emitted
2. **Given** a `FloorClearMsg` is received, **When** the floor advances, **Then** the floor number increments and the new floor's daemon enemies are displayed
3. **Given** the game screen is showing, **When** the current floor has active daemon enemies, **Then** enemy names and types are visible in the UI floor panel
4. **Given** the player is on any floor, **When** the game screen is viewed, **Then** the current floor number and a progress bar toward floor clear are clearly displayed

## Tasks / Subtasks

- [ ] Task 1: Create `internal/content/enemies.go` — EnemyDefinition struct and floor-keyed enemy pool (AC: 2, 3)
  - [ ] Define `EnemyDefinition` struct: `ID`, `Name`, `Type`, `Description`, `FloorMin`, `FloorMax int` (0 = no limit)
  - [ ] Define `var Enemies = []EnemyDefinition{...}` with 5 tiers: floors 1-3, 4-6, 7-10, 11-15, 16+
  - [ ] Export `func EnemiesForFloor(floor int) []EnemyDefinition` — returns all defs where `FloorMin <= floor && (FloorMax == 0 || floor <= FloorMax)`

- [ ] Task 2: Create `internal/engine/floor.go` — floor engine pure functions (AC: 1, 2, 4)
  - [ ] `func FloorThreshold(floor int) float64` — CPU Cycles needed: `100.0 * math.Pow(2.0, float64(floor-1))` (floor 1=100, 2=200, 3=400…)
  - [ ] `func FloorProgress(run game.RunState) float64` — `min(run.Resources[CPUCycles] / FloorThreshold(run.Floor), 1.0)`; clamp to [0,1]
  - [ ] `func CheckFloorClear(run game.RunState) bool` — returns `run.Resources[content.CPUCycles] >= FloorThreshold(run.Floor)`
  - [ ] `func AdvanceFloor(run game.RunState) game.RunState` — deep-copies Resources map, deducts `FloorThreshold(run.Floor)` from `CPUCycles` (clamp to 0), increments `run.Floor`

- [ ] Task 3: Create `internal/engine/floor_test.go` — tests for all four floor functions (AC: 1, 4)
  - [ ] `TestFloorThreshold_Floor1` — asserts threshold == 100.0
  - [ ] `TestFloorThreshold_ScalesExponentially` — floor 2 = 200, floor 3 = 400
  - [ ] `TestFloorProgress_BelowThreshold` — 50 CPU on floor 1 → 0.5
  - [ ] `TestFloorProgress_ClampsAtOne` — 9999 CPU on floor 1 → 1.0
  - [ ] `TestCheckFloorClear_NotCleared` — 99 CPU on floor 1 → false
  - [ ] `TestCheckFloorClear_ExactlyCleared` — exactly 100 CPU on floor 1 → true
  - [ ] `TestAdvanceFloor_IncrementsFloor` — floor 1 → floor 2
  - [ ] `TestAdvanceFloor_DeductsCost` — 250 CPU after clearing floor 1 (cost 100) → 150 CPU
  - [ ] `TestAdvanceFloor_ClampsToZero` — exactly 100 CPU → 0 CPU after clear
  - [ ] `TestAdvanceFloor_DoesNotMutateOriginal` — source RunState Resources map unchanged

- [ ] Task 4: Add `FloorPanel` style to `internal/ui/styles/layout.go` (AC: 4)
  - [ ] `func FloorPanel(width int) lipgloss.Style` — same dimensions as ResourcePanel/UpgradePanel, `BorderForeground(Accent)`

- [ ] Task 5: Add `floorClearMsg` and `floorClearCmd` to `internal/ui/model.go` (AC: 1)
  - [ ] `type floorClearMsg struct{}` (unexported — only ui package emits and handles it)
  - [ ] `func floorClearCmd() tea.Cmd` — returns `func() tea.Msg { return floorClearMsg{} }`

- [ ] Task 6: Update `internal/ui/update.go` — integrate floor clear into tick loop and handle `floorClearMsg` (AC: 1, 2)
  - [ ] In `engine.TickMsg` handler (GameScreen only): after `engine.Tick(...)`, call `engine.CheckFloorClear(m.State.Run)`
  - [ ] If true: return `m, tea.Batch(engine.TickCmd(), floorClearCmd())` — keep ticking AND trigger advance
  - [ ] Add `case floorClearMsg:` handler: call `engine.AdvanceFloor(m.State.Run)` → set `m.State.Run`, set `m.StatusMsg = fmt.Sprintf("Floor %d cleared!", prevFloor)`, return `m, clearStatusAfterCmd(3 * time.Second)`
  - [ ] Do NOT transition to `PerkDraftScreen` here — that is Story 3.2
  - [ ] Import `fmt` if not present

- [ ] Task 7: Update `internal/ui/screens/game.go` — add floor panel rendering, update `GameView` signature (AC: 3, 4)
  - [ ] Add `func FloorView(floor int, progress float64, enemies []content.EnemyDefinition, width int) string`
    - [ ] Render progress bar: 20 chars wide, filled = `█`, empty = `░`; e.g. `[████████░░░░░░░░░░░░]` with percentage
    - [ ] Render enemy list: names joined by ` · ` or one per line if long
    - [ ] Wrap in `styles.FloorPanel(width)`
  - [ ] Update `GameView` signature: add `floor int`, `runNumber int`, `floorProgress float64`, `enemies []content.EnemyDefinition` parameters
  - [ ] Replace hardcoded `"Floor: 1  Run: 1"` footer with `fmt.Sprintf("Floor: %d  Run: %d", floor, runNumber)`
  - [ ] Insert `FloorView(...)` between the resource panel and upgrade panel in `GameView`

- [ ] Task 8: Update `internal/ui/view.go` — pass new args to `GameView` (AC: 3, 4)
  - [ ] Compute `floorProgress := engine.FloorProgress(m.State.Run)` locally in `View()`
  - [ ] Compute `enemies := content.EnemiesForFloor(m.State.Run.Floor)` locally
  - [ ] Pass `m.State.Run.Floor`, `m.State.Run.RunNumber`, `floorProgress`, `enemies` to `screens.GameView`

- [ ] Task 9: Verify build and tests (AC: all)
  - [ ] `make build` — zero errors
  - [ ] `make test` — all tests pass (new floor tests + all regressions)
  - [ ] `make lint` — clean

## Dev Notes

### What This Story Builds

Floor progression is the primary advancement loop of afk-x. This story wires up:
1. **Engine layer**: pure `floor.go` functions that compute thresholds, progress, and advance state
2. **Content layer**: `enemies.go` with 5 tiers of daemon enemies keyed to floor ranges
3. **UI layer**: a floor progress panel between the resource and upgrade panels, showing the progress bar and active enemies
4. **Message loop**: `floorClearMsg` emitted via `tea.Batch` so ticking continues uninterrupted during the floor advance

**Do NOT implement in this story:**
- Perk draft screen transition on `floorClearMsg` — that is Story 3.2. The floor simply advances and stays on `GameScreen`.
- Milestone flavor text on floor clear — that is Story 2.3. Only the bare `StatusMsg` "Floor X cleared!" is shown here.
- Save on floor clear — that is Story 5.x. No `SaveMsg` in this story.

### Architecture Compliance (MUST FOLLOW)

From `architecture.md`:

**Package import rules (strictly enforced by go build):**
```
engine      →  game, content           (no ui, no save, no combat)
content     →  (no internal imports)
ui/screens  →  engine, content, ui/styles   (NOT game directly)
ui          →  engine, game, content
```

- `engine/floor.go` MUST import only `game` and `content` — no Bubbletea, no `os`, no `math` wait, `math` is stdlib so it's fine
- `screens/game.go` receives enemy data as `[]content.EnemyDefinition` passed in — does NOT call `content.EnemiesForFloor` itself; that call is in `ui/view.go` — this preserves the screen-as-pure-renderer pattern
- `view.go` (ui package) calls `content.EnemiesForFloor` and `engine.FloorProgress` and passes results to `GameView`

**Bubbletea rules:**
- `engine.CheckFloorClear` returns `bool` — pure function, no messages
- `floorClearMsg` is created in `ui` package only, never in `engine`
- `tea.Batch(engine.TickCmd(), floorClearCmd())` — both commands fire; engine tick continues while floor advance processes
- All mutations in `Update()` — `View()` is read-only

**Pure function invariants for engine/floor.go:**
- All four functions receive `game.RunState` by value and return new state — no pointer receivers, no global mutation
- `AdvanceFloor` must deep-copy `run.Resources` before mutation (same pattern as `Tick` and `PurchaseUpgrade`)

### Key Implementation Patterns

**Floor threshold formula:**
```go
// FloorThreshold returns CPU cycles needed to clear the given floor.
// Uses exponential scaling: 100 * 2^(floor-1)
// Floor 1 = 100, Floor 2 = 200, Floor 3 = 400, Floor 4 = 800 ...
func FloorThreshold(floor int) float64 {
    return 100.0 * math.Pow(2.0, float64(floor-1))
}
```

**AdvanceFloor deep-copy pattern (MUST follow — same as Tick and PurchaseUpgrade):**
```go
func AdvanceFloor(run game.RunState) game.RunState {
    threshold := FloorThreshold(run.Floor)
    newResources := make(map[string]float64, len(run.Resources))
    for k, v := range run.Resources {
        newResources[k] = v
    }
    cpu := newResources[content.CPUCycles] - threshold
    if cpu < 0 { cpu = 0 }
    newResources[content.CPUCycles] = cpu
    run.Resources = newResources
    run.Floor++
    return run
}
```

**FloorProgress clamp:**
```go
func FloorProgress(run game.RunState) float64 {
    threshold := FloorThreshold(run.Floor)
    if threshold == 0 { return 1.0 }
    p := run.Resources[content.CPUCycles] / threshold
    if p > 1.0 { p = 1.0 }
    return p
}
```

**tea.Batch for floor clear (keeps ticking AND triggers advance):**
```go
case engine.TickMsg:
    if m.Screen != GameScreen { return m, nil }
    // ... tick logic ...
    m.State.Run = engine.Tick(m.State.Run, delta, m.Rates)
    if engine.CheckFloorClear(m.State.Run) {
        return m, tea.Batch(engine.TickCmd(), floorClearCmd())
    }
    return m, engine.TickCmd()

case floorClearMsg:
    prevFloor := m.State.Run.Floor
    m.State.Run = engine.AdvanceFloor(m.State.Run)
    m.StatusMsg = fmt.Sprintf("Floor %d cleared!", prevFloor)
    return m, clearStatusAfterCmd(3 * time.Second)
    // No TickCmd here — already running from tea.Batch above
```

**Progress bar rendering:**
```go
func progressBar(progress float64, width int) string {
    filled := int(progress * float64(width))
    empty := width - filled
    return "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"
}
```

**GameView signature update (add after `upgrades map[string]int`):**
```go
func GameView(
    resources map[string]float64,
    rates engine.ResourceRates,
    upgrades map[string]int,
    floor int,
    runNumber int,
    floorProgress float64,
    enemies []content.EnemyDefinition,
    statusMsg string,
    width, height int,
) string
```

**Enemy definition structure:**
```go
type EnemyDefinition struct {
    ID          string
    Name        string
    Type        string
    Description string
    FloorMin    int
    FloorMax    int // 0 = no upper limit
}
```

**Enemy tiers (5 tiers covering all floors):**
- Floors 1-3: "Ghost Process" — a dangling reference still executing after its parent was killed
- Floors 4-6: "Zombie Daemon" — terminated but unreleased from the process table
- Floors 7-10: "Race Condition" — two threads competing for the same lock; one will survive
- Floors 11-15: "Memory Leak" — slow, growing, consuming resources without bound
- Floors 16+ (FloorMax=0): "Kernel Panic" — unrecoverable error state given malevolent form

### Previous Story Intelligence (from Story 2.1)

- **Map deep-copy pattern**: Always copy `run.Resources` (and `run.Upgrades` if modified) before mutation. See `PurchaseUpgrade` in `engine/upgrades.go` and `Tick` in `engine/tick.go`.
- **`clearStatusMsg` already exists**: `model.go` already has `clearStatusMsg`, `clearStatusAfterCmd`. Just add `floorClearMsg` and `floorClearCmd` alongside it.
- **`StatusMsg` field on Model**: Already added in Story 2.1 — use it directly. Do NOT add another status field.
- **`GameView` signature**: Currently `(resources, rates, upgrades, statusMsg, width, height)`. This story adds `floor`, `runNumber`, `floorProgress`, `enemies` between `upgrades` and `statusMsg`.
- **`content` and `engine` already imported in `update.go`**: No new imports needed there except `fmt` for `Sprintf`.
- **`styles.UpgradePanel` and `styles.ResourcePanel` pattern**: `FloorPanel` should follow the exact same structure in `layout.go`.
- **Test file pattern**: Use `package engine_test` (external test package), helper func `newRunWith(cpu, mem, threads)` already defined in `upgrades_test.go` — note: you are creating a NEW test file (`floor_test.go`) so you must re-define the helper OR factor it out. Since Go test packages are per-file, define it again or use unexported package-level helpers in a shared `engine_test_helpers_test.go`. Simplest: redefine `newRunWith` locally in `floor_test.go` since it's small.
- **Module path**: `github.com/12jihan/afk-x`

### File-by-File Current State

**`internal/ui/update.go` — what changes:**
Current `TickMsg` handler: ticks resources, no floor check. Add `CheckFloorClear` call after tick and return `tea.Batch` if cleared. Add new `case floorClearMsg:` handler.

**`internal/ui/screens/game.go` — what changes:**
Current `GameView` footer hardcodes `"Floor: 1  Run: 1"`. Change to use passed `floor`/`runNumber` params. Insert new `FloorView(...)` call between resource panel and upgrade panel. Current `GameView` has 6 params → becomes 10 params. All callers live in `view.go` only.

**`internal/ui/view.go` — what changes:**
Current: `screens.GameView(m.State.Run.Resources, m.Rates, m.State.Run.Upgrades, m.StatusMsg, m.Width, m.Height)`. Add the four new args: floor, runNumber, floorProgress, enemies.

**`internal/ui/model.go` — what changes:**
Add `floorClearMsg` struct and `floorClearCmd()` function alongside existing `clearStatusMsg`.

**`internal/ui/styles/layout.go` — what changes:**
Add `FloorPanel(width int) lipgloss.Style` — same pattern as `ResourcePanel` and `UpgradePanel`.

### References

- Story 2.2 ACs: `_bmad-output/planning-artifacts/epics.md` → Epic 2, Story 2.2
- Architecture: `_bmad-output/planning-artifacts/architecture.md` → "Bubbletea Message Flow", "Game Loop Architecture", "Package Structure"
- Prev story pattern: `_bmad-output/implementation-artifacts/2-1-upgrade-definitions-purchase-system.md`
- Engine package: `internal/engine/upgrades.go` — deep-copy map pattern
- Engine package: `internal/engine/tick.go` — Tick deep-copy + TickCmd
- UI model: `internal/ui/model.go` — clearStatusMsg/clearStatusAfterCmd pattern to replicate
- UI update: `internal/ui/update.go` — current TickMsg handler to extend

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-6

### Debug Log References

### Completion Notes List

- 5 enemy tiers defined (Ghost Process floors 1-3, Zombie Daemon 4-6, Race Condition 7-10, Memory Leak 11-15, Kernel Panic 16+) with FloorMax=0 for open-ended tiers
- FloorThreshold uses `100 * 2^(floor-1)` — floor 1=100, 2=200, 3=400, etc.
- AdvanceFloor deep-copies Resources map before mutation (same pattern as Tick and PurchaseUpgrade)
- Floor clear uses `tea.Batch(engine.TickCmd(), floorClearCmd())` — keeps ticking while advance processes
- Perk draft transition intentionally deferred to Story 3.2 — `floorClearMsg` handler stays on GameScreen
- FloorView renders 24-char progress bar with █/░ fill, plus enemy names and types
- GameView signature expanded from 6 → 10 params; all callers are in view.go only (single call site)
- 29 engine tests passing (10 new floor tests + 19 prior), build and vet clean

### File List

- `internal/content/enemies.go` — new
- `internal/engine/floor.go` — new
- `internal/engine/floor_test.go` — new
- `internal/ui/model.go` — updated (floorClearMsg, floorClearCmd)
- `internal/ui/update.go` — updated (floor clear in TickMsg, floorClearMsg handler, fmt import)
- `internal/ui/screens/game.go` — updated (FloorView, progressBar, updated GameView signature)
- `internal/ui/view.go` — updated (engine.FloorProgress, content.EnemiesForFloor, new GameView args)
- `internal/ui/styles/layout.go` — updated (FloorPanel style)
- `_bmad-output/implementation-artifacts/2-2-floor-progression-daemon-encounters.md` — this file

### Change Log

- 2026-05-04: Story 2.2 complete — floor progression, daemon encounters, progress bar
