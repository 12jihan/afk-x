# Story 1.4: Resource Tick Engine

Status: review

## Story

As a player,
I want resources to generate automatically over time,
so that the game progresses without requiring my attention.

## Acceptance Criteria

1. **Given** a `RunState` and `ResourceRates`, **When** `engine.Tick()` is called with a 16ms delta, **Then** each resource value in the returned state increases by `rate * delta.Seconds()` (i.e. `rate * 0.016` for 16ms)
2. **Given** identical `RunState` and delta inputs, **When** `engine.Tick()` is called twice with the same arguments, **Then** both calls return identical `RunState` values (pure / deterministic)
3. **Given** `engine.ComputeRates()` called with a default `RunState` (no upgrades, no active perks), **When** evaluated, **Then** it returns a `ResourceRates` map whose values equal `content.BaseRates` for every key
4. **Given** `engine.tickCmd()` is called, **When** evaluated, **Then** it returns a non-nil `tea.Cmd` that uses `tea.Every(16ms, ...)` to deliver a `TickMsg` approximately every 16 milliseconds

## Tasks / Subtasks

- [x] Task 1: Replace `internal/engine/engine.go` stub with package doc (AC: all)
  - [x] Replace the TODO stub with a proper package doc comment
  - [x] `go build ./...` must still pass after change

- [x] Task 2: Implement `internal/engine/rates.go` — ResourceRates type + ComputeRates() (AC: 3)
  - [x] Define `type ResourceRates map[string]float64` — always use this named type, never raw `map[string]float64` in signatures
  - [x] Implement `ComputeRates(run game.RunState) ResourceRates` — copies `content.BaseRates` into a new `ResourceRates` map (defensive copy — never return a reference to the package-level var)
  - [x] Stub space for future upgrade/perk multipliers with a `// TODO: Story 2.x` comment
  - [x] Imports: `game`, `content` only — NO Bubbletea in this file

- [x] Task 3: Implement `internal/engine/tick.go` — TickMsg + tickCmd() + Tick() (AC: 1, 2, 4)
  - [x] Define `type TickMsg time.Time`
  - [x] Implement `TickCmd() tea.Cmd` (exported) — returns `tea.Every(16*time.Millisecond, func(t time.Time) tea.Msg { return TickMsg(t) })`
  - [x] Implement `Tick(state game.RunState, delta time.Duration, rates ResourceRates) game.RunState` — pure function: for each key in rates, `state.Resources[key] += rate * delta.Seconds()`; return new state
  - [x] Imports: `time`, `github.com/charmbracelet/bubbletea`, `github.com/12jihan/afk-x/internal/game`
  - [x] `Tick()` creates a copy of Resources map before mutating — original state not affected

- [x] Task 4: Write `internal/engine/rates_test.go` (AC: 3)
  - [x] `TestComputeRatesDefaultState` — `ComputeRates(game.NewGameState().Run)` returns rates equal to `content.BaseRates` for all three resource keys
  - [x] `TestComputeRatesReturnsCopy` — mutate the returned `ResourceRates`; verify `content.BaseRates` is unchanged
  - [x] `TestComputeRatesAllKeysPresent` — result contains the three keys: `content.CPUCycles`, `content.MemoryShards`, `content.ProcessThreads`

- [x] Task 5: Write `internal/engine/tick_test.go` (AC: 1, 2, 4)
  - [x] `TestTickAppliesRates` — call `Tick()` with 16ms delta and known rates; verify each resource increased by `rate * 0.016` (float tolerance `1e-9`)
  - [x] `TestTickDeterministic` — call `Tick()` twice with identical inputs; verify results are equal
  - [x] `TestTickPure` — verify original `RunState`'s Resources map is not mutated after `Tick()` returns
  - [x] `TestTickCmdNonNil` — `TickCmd()` returns non-nil `tea.Cmd`
  - [x] `TestTickZeroDelta` — `Tick()` with 0 delta returns state with unchanged resource values

- [x] Task 6: Verify full test suite passes (AC: all)
  - [x] `go test ./internal/engine/...` — 8 tests pass
  - [x] `go test ./...` — no regressions

## Dev Notes

### What This Story Does

Stories 1.2 and 1.3 built the data layer (`game` structs, `content` constants). This story builds the computation layer: the two pure functions that drive the idle loop.

- `ComputeRates(run RunState) ResourceRates` — answers "how fast are resources growing right now?"
- `Tick(state RunState, delta Duration, rates ResourceRates) RunState` — answers "given this rate, what's the new state after `delta` time?"
- `tickCmd() tea.Cmd` + `TickMsg` — the Bubbletea scheduling bridge that delivers a tick event every 16ms

Story 1.5 (TUI shell) will wire `tickCmd()` into the Bubbletea `Init()` and handle `TickMsg` in `Update()`. Story 1.4 only creates the engine functions — no UI, no wiring.

**Files changed:**
- `internal/engine/engine.go` — stub replaced with package doc
- `internal/engine/rates.go` — new
- `internal/engine/rates_test.go` — new
- `internal/engine/tick.go` — new
- `internal/engine/tick_test.go` — new

**Do NOT create in this story:**
- `internal/engine/offline.go` — AFK accumulation, Story 5.x
- `internal/engine/floor.go` — floor clear logic, Story 2.2
- `internal/engine/perk.go` — perk application, Story 3.x
- `internal/engine/upgrade.go` — upgrade logic, Story 2.1
- Any UI integration — Story 1.5

### Architecture Compliance (MUST FOLLOW)

**`engine` package responsibility:**
- ✅ Pure functions: `Tick()`, `ComputeRates()` — no side effects, no I/O
- ✅ Bubbletea bridge ONLY in `tick.go`: `TickMsg` type + `tickCmd()` function — this is the intentional exception to the "no Bubbletea" rule; only these two items may use `tea`
- ❌ FORBIDDEN: any `os` calls, file I/O, global mutable state
- ❌ FORBIDDEN: importing `ui`, `save`, `combat` — engine only imports `game` and `content`
- ❌ FORBIDDEN: `map[string]float64` in function signatures — always use `engine.ResourceRates`

**Import map for this story:**
```
engine/rates.go  → game, content
engine/tick.go   → time, bubbletea, game
engine/engine.go → (no imports — package doc only)
```

**Dependency graph (do not violate):**
```
content  →  (none)
game     →  (none)
engine   →  game, content      ← tick.go also → bubbletea (intentional bridge only)
```

[Source: `architecture.md` → Package Structure, Package Responsibility, Game Loop Architecture]

### Exact Function Signatures

These are canonical — match exactly:

```go
// internal/engine/rates.go
package engine

import (
    "github.com/12jihan/afk-x/internal/content"
    "github.com/12jihan/afk-x/internal/game"
)

// ResourceRates maps resource keys to their generation rate (units per second).
// Always use this named type — never map[string]float64 in signatures.
type ResourceRates map[string]float64

// ComputeRates returns the current resource generation rates for a run.
// For MVP (no upgrades or perks active) this equals content.BaseRates.
// Returns a defensive copy — callers may mutate the result safely.
func ComputeRates(run game.RunState) ResourceRates {
    rates := make(ResourceRates, len(content.BaseRates))
    for k, v := range content.BaseRates {
        rates[k] = v
    }
    // TODO: Story 2.x — apply upgrade multipliers from run.Upgrades
    // TODO: Story 3.x — apply active perk bonuses from run.ActivePerks
    return rates
}
```

```go
// internal/engine/tick.go
package engine

import (
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/12jihan/afk-x/internal/game"
)

// TickMsg is the Bubbletea message delivered by tickCmd every 16ms.
type TickMsg time.Time

// tickCmd returns a Bubbletea command that fires a TickMsg every 16ms.
// Unexported — called by ui.Init() and returned from ui.Update() after each tick.
func tickCmd() tea.Cmd {
    return tea.Every(16*time.Millisecond, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}

// Tick advances resources by delta time at the given rates.
// Pure function — returns a new RunState with updated Resources;
// the original RunState is not mutated.
func Tick(state game.RunState, delta time.Duration, rates ResourceRates) game.RunState {
    // Copy the Resources map to avoid mutating the original
    newResources := make(map[string]float64, len(state.Resources))
    for k, v := range state.Resources {
        newResources[k] = v
    }
    for resource, rate := range rates {
        newResources[resource] += rate * delta.Seconds()
    }
    state.Resources = newResources
    return state
}
```

### Critical Implementation Detail — Tick() Map Copy

`game.RunState` is passed by value, but its `Resources` field is `map[string]float64` — a reference type. If you write directly to `state.Resources[key]`, you mutate the caller's map through the shared reference. **Always create a new map and copy before writing.**

```go
// ❌ WRONG — mutates the caller's Resources map
func Tick(state game.RunState, delta time.Duration, rates ResourceRates) game.RunState {
    for resource, rate := range rates {
        state.Resources[resource] += rate * delta.Seconds() // mutates original!
    }
    return state
}

// ✅ CORRECT — defensive copy
func Tick(state game.RunState, delta time.Duration, rates ResourceRates) game.RunState {
    newResources := make(map[string]float64, len(state.Resources))
    for k, v := range state.Resources {
        newResources[k] = v
    }
    for resource, rate := range rates {
        newResources[resource] += rate * delta.Seconds()
    }
    state.Resources = newResources
    return state
}
```

This is what `TestTickPure` tests — the original state's Resources map must be unchanged.

### Float Precision — Use Tolerance in Tests

`16 * time.Millisecond` = 16,000,000 nanoseconds. `delta.Seconds()` = 0.016. However, `0.016` is **not exactly representable** in IEEE 754 float64 (it equals 16/1000 = 2/125 — 5 is not a power of 2). Always compare float64 results with a tolerance, not `==`:

```go
const floatTol = 1e-9

func approxEqual(a, b float64) bool {
    diff := a - b
    if diff < 0 { diff = -diff }
    return diff < floatTol
}
```

For `TestTickAppliesRates`: assert `approxEqual(got, want)`, not `got == want`.

### ResourceRates Named Type Rule

The architecture explicitly forbids using raw `map[string]float64` in function signatures:
```
❌ func Tick(state game.RunState, delta time.Duration, rates map[string]float64) game.RunState
✅ func Tick(state game.RunState, delta time.Duration, rates ResourceRates) game.RunState
```

`ResourceRates` is a named type defined in `rates.go`. Every function in `engine` that works with rate data must use this type.

### Resolving Story 1.3 Deferred Item

Story 1.3 code review deferred:
> `BaseRates` external mutation risk — any importer can corrupt the baseline; mitigation belongs in Story 1.4 (`ComputeRates` must copy the map, not reference directly)

**This is resolved by Task 2.** `ComputeRates()` always creates a new `ResourceRates` map and copies values from `content.BaseRates` — it never returns `content.BaseRates` directly or takes a reference to it. `TestComputeRatesReturnsCopy` verifies this.

### Testing Standards

- Test files co-located with source: `internal/engine/rates_test.go`, `internal/engine/tick_test.go`
- Use **`package engine_test`** (external black-box) — same pattern as `game_test` and `content_test`
- Import `github.com/12jihan/afk-x/internal/engine` explicitly
- `go test ./internal/engine/...` must pass (all new tests)
- `go test ./...` must still pass (no regressions)

**Test package imports:**
```go
// rates_test.go
package engine_test

import (
    "testing"
    "github.com/12jihan/afk-x/internal/content"
    "github.com/12jihan/afk-x/internal/engine"
    "github.com/12jihan/afk-x/internal/game"
)

// tick_test.go
package engine_test

import (
    "math"
    "testing"
    "time"
    "github.com/12jihan/afk-x/internal/content"
    "github.com/12jihan/afk-x/internal/engine"
    "github.com/12jihan/afk-x/internal/game"
)
```

**Note on `tickCmd()` test:** `tickCmd()` is unexported. Story 1.4 tests only verify that the exported surface works correctly. If `tickCmd` needs to be tested, either:
1. Export it as `TickCmd()` (recommended for testability — see AC4)
2. Test indirectly through integration in Story 1.5

**Architecture recommendation:** Export as `TickCmd()` so it can be called from `ui.Init()` and tested directly. The architecture shows it lowercase but Story 1.5 needs to call it from the `ui` package — it must be exported.

### tickCmd Export Decision

The architecture shows `tickCmd()` as unexported, but `ui` (a separate package) must call it in `Init()`. **Export it as `TickCmd()`** for this story. AC4 also tests it via `engine.TickCmd()`.

```go
// ✅ Export — called from ui package in Story 1.5
func TickCmd() tea.Cmd {
    return tea.Every(16*time.Millisecond, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}
```

### Previous Story Intelligence (from Story 1.3)

- **External test package pattern:** Use `package engine_test` + explicit import — same as `content_test` and `game_test`
- **No type aliases needed:** Engine tests call functions directly, no type aliases required
- **Module path confirmed:** `github.com/12jihan/afk-x`
- **Go version:** 1.26.2 (do not change `go.mod`)
- **Bubbletea import path:** `github.com/charmbracelet/bubbletea` (already in `go.mod` as `v1.3.10`)
- **Testing pattern from Story 1.2:** `TestTickDeterministic` mirrors the `TestJSONRoundTrip` pattern — populate, run, compare field-by-field

### Project Structure Notes

```
internal/engine/
├── engine.go       ← replace stub with package doc
├── rates.go        ← new: ResourceRates type + ComputeRates()
├── rates_test.go   ← new: 3 tests
├── tick.go         ← new: TickMsg + TickCmd() + Tick()
└── tick_test.go    ← new: 5 tests
```

Do NOT create `offline.go`, `floor.go`, `perk.go`, or `upgrade.go` — those belong to later stories.

### References

- Story 1.4 AC: `epics.md` → Story 1.4 section
- `Tick()` and `ResourceRates`: `architecture.md` → "Game Loop Architecture"
- Engine package rules: `architecture.md` → "Package Responsibility (strict)"
- Anti-pattern `map[string]float64`: `architecture.md` → "Anti-Patterns"
- `tickCmd()` / `TickMsg`: `architecture.md` → "Game Loop Architecture" code snippet
- `ComputeRates` defensive copy: `deferred-work.md` → Story 1.3 deferred item
- AFK accumulation pattern (for future reference): `architecture.md` → "AFK Accumulation (always this pattern)"
- NFR1 (16ms budget): `epics.md` → NFR1 / `architecture.md` → Performance

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-6

### Debug Log References

- `go.sum` required `go get github.com/charmbracelet/x/ansi@v0.11.6` after importing bubbletea in tick.go — transitive dependency not yet in sum file

### Completion Notes List

- TDD: rates_test.go written first (RED — undefined engine.ComputeRates), rates.go implemented (GREEN); same for tick_test.go / tick.go
- `engine.go` stub replaced with full package doc covering dependency rules and Bubbletea exception
- `rates.go`: `ResourceRates` named type + `ComputeRates()` with defensive copy of `content.BaseRates` — resolves Story 1.3 deferred mutation risk
- `tick.go`: `TickMsg` type + `TickCmd()` (exported for ui package) + `Tick()` pure function with Resources map copy
- `Tick()` copies the Resources map before writing — prevents mutation of caller's state through shared map reference
- `tickCmd` exported as `TickCmd()` — architecture showed lowercase but ui package must call it from a different package
- `go test ./internal/engine/...` — 8/8 pass; `go test ./...` — zero regressions

### File List

- `internal/engine/engine.go` — stub replaced with package doc
- `internal/engine/rates.go` — new: ResourceRates type + ComputeRates()
- `internal/engine/rates_test.go` — new: 3 tests
- `internal/engine/tick.go` — new: TickMsg + TickCmd() + Tick()
- `internal/engine/tick_test.go` — new: 5 tests
- `go.sum` — updated (bubbletea transitive dependency)

### Change Log

- 2026-05-03: Story 1.4 complete — resource tick engine with ResourceRates type, ComputeRates(), Tick(), TickCmd(), 8-test suite
