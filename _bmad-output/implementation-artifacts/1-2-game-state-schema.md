# Story 1.2: Game State Schema

Status: done

## Story

As a developer,
I want the core game state structs defined with JSON serialization,
so that all packages share a consistent, versioned data contract.

## Acceptance Criteria

1. **Given** a `GameState` is created, **When** marshaled to JSON, **Then** all fields use snake_case keys and `version` is set to 1
2. **Given** a `GameState` is marshaled and then unmarshaled, **When** the result is compared to the original, **Then** all fields are identical (round-trip test passes)
3. **Given** `RunState.Resources`, **When** populated with resource values, **Then** it uses `map[string]float64` with resource key constants from the `content` package
4. **Given** `RunState.PendingDraft`, **When** no draft is pending, **Then** it serializes as `null`, not `[]`

## Tasks / Subtasks

- [x] Task 1: Implement `internal/game/state.go` with all structs and constructor (AC: 1, 2, 3, 4)
  - [x] Define `GameState` struct with `Version int`, `Run RunState`, `Meta MetaState`, `SavedAt time.Time` — all snake_case JSON tags
  - [x] Define `RunState` struct with all fields per architecture spec
  - [x] Define `MetaState` struct with all fields per architecture spec
  - [x] Implement `NewGameState() GameState` constructor that sets `Version: 1`, initializes all maps, and leaves `PendingDraft` as **nil** (not `[]string{}`)
  - [x] Verify package imports only `time` — zero business logic, zero Bubbletea, zero internal imports

- [x] Task 2: Create `internal/game/state_test.go` with all AC tests (AC: 1, 2, 3, 4)
  - [x] `TestVersionSetToOne` — `NewGameState().Version == 1` and JSON contains `"version":1`
  - [x] `TestAllFieldsSnakeCase` — marshal and verify JSON keys are snake_case (spot-check `saved_at`, `run_number`, `pending_draft`, `active_perks`, `combo_queue`, `permanent_unlocks`, `best_combos`, `unlocked_perks`, `run_count`)
  - [x] `TestJSONRoundTrip` — populate every field, marshal to JSON, unmarshal into new struct, deep-compare all fields
  - [x] `TestPendingDraftNullWhenNil` — marshal `RunState{PendingDraft: nil}` → JSON must contain `"pending_draft":null`, NOT `"pending_draft":[]`
  - [x] `TestResourcesMapType` — verify `Resources` is `map[string]float64`; populate with string keys matching planned content constants (`"cpu_cycles"`, `"memory_shards"`, `"process_threads"`) and verify round-trip

### Review Findings (AI)

- [x] [Review][Patch] Round-trip test verifies slice/map lengths only — add element-value assertions for `ActivePerks`, `ComboQueue` inner elements, `BestCombos` map values, `PendingDraft`, and `UnlockedPerks` [internal/game/state_test.go — TestJSONRoundTrip]
- [x] [Review][Patch] Missing null-serialization tests for `ComboQueue`, `ActivePerks`, `PermanentUnlocks`, `UnlockedPerks` — only `PendingDraft` has a `TestPendingDraftNullWhenNil`; add equivalent assertions for the other nil slice fields [internal/game/state_test.go]
- [x] [Review][Patch] `TestSavedAtRFC3339` only checks for a date substring (`"2026-05-03"`) rather than validating RFC3339 format — use `time.Parse(time.RFC3339Nano, ...)` on the extracted value [internal/game/state_test.go — TestSavedAtRFC3339]
- [x] [Review][Defer] Update `TestResourcesMapType` and `TestJSONRoundTrip` to use `content` package resource key constants instead of string literals — blocked until Story 1.3 defines the constants [internal/game/state_test.go] — deferred, pre-existing
- [x] [Review][Defer] `PendingDraft` nil invariant has no type-level enforcement (no custom `MarshalJSON`) — current test coverage is sufficient for MVP; revisit if schema complexity grows — deferred, pre-existing
- [x] [Review][Defer] Resources zero-value key ambiguity (`"cpu_cycles":0.0` vs absent key) not tested or documented — game design question, defer to engine layer — deferred, pre-existing
- [x] [Review][Defer] Version schema mismatch detection (e.g. loading `"version":99`) has no tested enforcement path — belongs in the `save` package (Story 5.x), not `game` — deferred, pre-existing
- [x] [Review][Defer] Float precision edge cases (`1.1`, `0.3`, near `math.MaxFloat64`) not exercised in resource round-trip tests — theoretical for this game's value ranges; defer — deferred, pre-existing

## Dev Notes

### What This Story Does

Story 1.1 left `internal/game/state.go` as a stub with only a package declaration. This story **replaces that stub** with the full struct definitions that every other package in the project depends on.

**The output of this story:** `internal/game/state.go` with 3 structs + 1 constructor, and `internal/game/state_test.go` with 5+ passing tests. That's it — no other files change.

### Architecture Compliance (MUST FOLLOW)

**Package responsibility rule — `game` is structs + JSON ONLY:**
- ✅ Allowed: struct definitions, JSON tags, `time.Time` for timestamps, `NewGameState()` constructor
- ❌ FORBIDDEN: any business logic (calculations, validations, mutations)
- ❌ FORBIDDEN: any Bubbletea imports (`tea`, `lipgloss`, `bubbles`)
- ❌ FORBIDDEN: any internal package imports (`engine`, `content`, `combat`, `save`, `ui`)
- ❌ FORBIDDEN: `os`, `fmt` (except in `_test.go`), file I/O of any kind

**Import rule for `game` package:**
```go
import "time"  // only stdlib import needed
```

**Dependency graph (do not violate):**
```
game  →  (no internal imports)
```

[Source: `_bmad-output/planning-artifacts/architecture.md` → Package Structure, Package Responsibility]

### Exact Struct Definitions

These are canonical — do not deviate from field names, types, or JSON tags:

```go
// internal/game/state.go
package game

import "time"

type GameState struct {
    Version int       `json:"version"`   // schema migration guard — always 1 for MVP
    Run     RunState  `json:"run"`
    Meta    MetaState `json:"meta"`
    SavedAt time.Time `json:"saved_at"`  // AFK delta anchor — RFC3339 in JSON
}

type RunState struct {
    Floor        int                `json:"floor"`
    Resources    map[string]float64 `json:"resources"`      // extensible — string keys from content package
    Upgrades     map[string]int     `json:"upgrades"`
    ActivePerks  []string           `json:"active_perks"`
    ComboQueue   [][]string         `json:"combo_queue"`    // FR33: persisted combo queue
    PendingDraft []string           `json:"pending_draft"`  // FR32: nil = no draft pending
    RunNumber    int                `json:"run_number"`
}

type MetaState struct {
    PermanentUnlocks []string            `json:"permanent_unlocks"` // FR23
    BestCombos       map[string][]string `json:"best_combos"`       // FR26
    UnlockedPerks    []string            `json:"unlocked_perks"`    // FR14
    RunCount         int                 `json:"run_count"`
}
```

[Source: `_bmad-output/planning-artifacts/architecture.md` → Data Architecture — Game State Model]

### Constructor Pattern

```go
// NewGameState returns a freshly initialized GameState ready for a new run.
func NewGameState() GameState {
    return GameState{
        Version: 1,
        Run: RunState{
            Floor:     1,
            Resources: make(map[string]float64),
            Upgrades:  make(map[string]int),
            // ActivePerks, ComboQueue, PendingDraft: nil (zero value) — correct
            RunNumber: 1,
        },
        Meta: MetaState{
            BestCombos: make(map[string][]string),
            // PermanentUnlocks, UnlockedPerks: nil (zero value) — correct
        },
        SavedAt: time.Now(),
    }
}
```

**Why maps are initialized but slices are not:** Maps panic on write if nil; nil slices behave correctly for JSON (`null`) and range iteration. Do NOT `make([]string, 0)` for any slice field — that produces `[]` in JSON, which breaks AC4.

### AC4 Critical Detail — `null` vs `[]`

This is the most common mistake. In Go:
- `nil` slice → JSON `null` ✅
- `[]string{}` or `make([]string, 0)` → JSON `[]` ❌

The architecture explicitly requires `PendingDraft: null` in JSON to signal "no draft pending" (FR32). An empty array `[]` is **invalid** per the data contract.

```go
// ✅ CORRECT — nil slice → "pending_draft": null
state := RunState{}
// state.PendingDraft is nil by default

// ❌ WRONG — empty slice → "pending_draft": []
state.PendingDraft = make([]string, 0)
state.PendingDraft = []string{}
```

The same rule applies to `ActivePerks`, `ComboQueue`, `PermanentUnlocks`, `UnlockedPerks` — prefer nil (zero value) over empty slice initialization. However, `PendingDraft` is the one tested by AC4.

[Source: `_bmad-output/planning-artifacts/architecture.md` → Format Patterns, `pending_draft: null` note]

### AC3 — Resource Keys and the `content` Package

AC3 says `Resources` uses keys from `content` package constants. Important:
- The `game` package **cannot** import `content` (violates dependency rules)
- The `Resources map[string]float64` type is correct — any string can be a key
- In tests, use the planned string literals directly: `"cpu_cycles"`, `"memory_shards"`, `"process_threads"`
- Story 1.3 will define these as proper constants in `content/resources.go`
- Consumers of `game` (engine, combat, ui) will use the content constants — `game` just holds the map

**Test pattern for AC3:**
```go
func TestResourcesMapType(t *testing.T) {
    state := NewGameState()
    state.Run.Resources["cpu_cycles"] = 100.5
    state.Run.Resources["memory_shards"] = 50.0
    state.Run.Resources["process_threads"] = 10.25

    data, err := json.Marshal(state.Run)
    if err != nil {
        t.Fatal(err)
    }

    var restored RunState
    if err := json.Unmarshal(data, &restored); err != nil {
        t.Fatal(err)
    }

    if restored.Resources["cpu_cycles"] != 100.5 {
        t.Errorf("cpu_cycles: got %v, want 100.5", restored.Resources["cpu_cycles"])
    }
}
```

### JSON Round-Trip Test Pattern

```go
func TestJSONRoundTrip(t *testing.T) {
    original := NewGameState()
    original.Run.Floor = 5
    original.Run.Resources["cpu_cycles"] = 1250.5
    original.Run.Resources["memory_shards"] = 340.0
    original.Run.Upgrades["cycle_boost_1"] = 2
    original.Run.ActivePerks = []string{"overclock"}
    original.Run.ComboQueue = [][]string{{"scan", "exploit", "loot"}}
    original.Run.PendingDraft = []string{"perk_a", "perk_b", "perk_c"}
    original.Run.RunNumber = 2
    original.Meta.PermanentUnlocks = []string{"unlock_1"}
    original.Meta.BestCombos = map[string][]string{"zone_1": {"scan", "exploit"}}
    original.Meta.UnlockedPerks = []string{"fork_bomb"}
    original.Meta.RunCount = 3

    data, err := json.Marshal(original)
    if err != nil {
        t.Fatalf("marshal failed: %v", err)
    }

    var restored GameState
    if err := json.Unmarshal(data, &restored); err != nil {
        t.Fatalf("unmarshal failed: %v", err)
    }

    // Compare key fields
    if restored.Version != original.Version {
        t.Errorf("Version: got %d, want %d", restored.Version, original.Version)
    }
    if restored.Run.Floor != original.Run.Floor {
        t.Errorf("Floor: got %d, want %d", restored.Run.Floor, original.Run.Floor)
    }
    if !restored.SavedAt.Equal(original.SavedAt) {
        t.Errorf("SavedAt mismatch")
    }
    // ... continue for all fields
}
```

**Note on `time.Time` comparison:** Use `.Equal()` not `==` — JSON round-trip normalizes timezone to UTC. `time.Time` values with the same instant but different timezone representations are `!=` but `.Equal()` returns true.

### Module Path

The actual module is `github.com/12jihan/afk-x` (set during Story 1.1). Use this exact path if any import is needed in tests.

### Testing Standards

- Tests live at `internal/game/state_test.go` — co-located with the package (not a separate `tests/` dir)
- Use `package game` (white-box) or `package game_test` (black-box) — either is acceptable; black-box is slightly preferred since we're testing the public API
- Use only `encoding/json`, `testing`, and `reflect` from stdlib in tests
- `go test ./internal/game/...` must pass
- `go test ./...` must still pass (CI requirement from Story 1.1)

### Existing File State

`internal/game/state.go` currently contains:
```go
// Package game defines the core game state types.
// TODO: Implemented in Story 1.2
package game
```

**Replace this file entirely.** Do not append to it.

### Project Structure Notes

- Only two files change in this story: `internal/game/state.go` (replace) and `internal/game/state_test.go` (new)
- No other packages change — `content` remains a stub, `engine` remains a stub, etc.
- Module path: `github.com/12jihan/afk-x`
- Go version in `go.mod`: 1.26.2 (do not change)
- `bubbles` import path in `go.mod` is `charm.land/bubbles/v2` — irrelevant to this story, `game` has no Charm imports

### References

- Struct definitions: `_bmad-output/planning-artifacts/architecture.md` → "Data Architecture — Game State Model"
- Package rules: `_bmad-output/planning-artifacts/architecture.md` → "Package Responsibility (strict)" and "Package Structure"
- JSON format spec: `_bmad-output/planning-artifacts/architecture.md` → "Format Patterns"
- null vs []: `_bmad-output/planning-artifacts/architecture.md` → "`pending_draft: null` means no draft pending; empty array `[]` is invalid"
- FR32/FR33: `_bmad-output/planning-artifacts/epics.md` → Story 5.5 (PendingDraft + ComboQueue must survive process death)
- AC source: `_bmad-output/planning-artifacts/epics.md` → Epic 1, Story 1.2

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-6

### Debug Log References

- External test package (`package game_test`) requires explicit import of `github.com/12jihan/afk-x/internal/game` — added type aliases and wrapper function for test ergonomics

### Completion Notes List

- Replaced stub `internal/game/state.go` with full `GameState`, `RunState`, `MetaState` structs + `NewGameState()` constructor
- `game` package imports only `"time"` from stdlib — zero internal imports, zero business logic, zero Bubbletea
- `NewGameState()` initializes `Resources` and `Upgrades` maps (prevent nil-map panics) but leaves all slice fields nil (serializes as JSON `null`, not `[]`)
- 6 tests written and passing: `TestVersionSetToOne`, `TestAllFieldsSnakeCase`, `TestJSONRoundTrip`, `TestPendingDraftNullWhenNil`, `TestResourcesMapType`, `TestSavedAtRFC3339`
- All 4 ACs verified: snake_case keys ✅, round-trip ✅, map[string]float64 resources ✅, null pending_draft ✅
- `go test ./...` and `go vet ./...` pass with zero regressions

### File List

- `internal/game/state.go`
- `internal/game/state_test.go`

### Change Log

- 2026-05-03: Story 1.2 complete — GameState schema implemented with full JSON serialization and 6-test suite
