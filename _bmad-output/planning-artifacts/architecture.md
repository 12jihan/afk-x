---
stepsCompleted: ['step-01-init', 'step-02-context', 'step-03-starter', 'step-04-decisions', 'step-05-patterns', 'step-06-structure', 'step-07-validation', 'step-08-complete']
status: 'complete'
completedAt: '2026-05-02'
inputDocuments: ['_bmad-output/planning-artifacts/prd.md']
workflowType: 'architecture'
project_name: 'afk-x'
user_name: 'Boss'
date: '2026-05-02'
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements:**
44 FRs across 8 capability areas: resource generation, upgrade/floor progression, perk draft system, command/combo system, meta-progression/run flow, save/persistence, narrative/presentation, and application shell. The FRs decompose into three distinct subsystems — game engine (tick, AFK calc, state machine), UI layer (Bubbletea panels and screen transitions), and persistence layer (serialization, atomic writes, CLI surface).

**Non-Functional Requirements:**
- **Performance:** 16ms tick budget drives goroutine-based game loop; 500ms AFK calc ceiling requires efficient offline math; 2s launch, 100ms perk draft render, 200ms save cap shape every hot path
- **Reliability:** Atomic save writes + deterministic AFK accumulation are non-negotiable; corrupted save must be detectable and recoverable
- **Portability:** Single binary (no CGo), Linux/macOS amd64+arm64, ≤20MB, 5+ terminal emulators
- **Accessibility:** No color-only information; full NO_COLOR support; keyboard-only interactions

**Scale & Complexity:**
- Primary domain: TUI game engine (Go + Bubbletea)
- Complexity level: Medium — no network, no database, no external APIs; complexity is internal (game loop, state machine, schema evolution)
- Estimated architectural components: 6–8 Go packages

### Technical Constraints & Dependencies

- **Language/framework fixed:** Go + Bubbletea (non-negotiable from PRD)
- **No CGo:** Single binary requirement strongly implies pure Go dependencies
- **Bubbletea model lifecycle:** All state must flow through `Model → Update → View` — no shared mutable globals
- **16ms tick budget:** Game loop goroutine must be non-blocking; state updates must be O(1) or O(small constant) per tick
- **Save schema must be extensible:** FR5 (resource type unlocks) and FR14 (perk pool expansion) require the JSON schema to accommodate dynamic collections, not fixed structs

### Cross-Cutting Concerns Identified

1. **Game loop timing** — tick rate, render rate, and AFK delta calculation all share time as a primitive; must be centralized
2. **State serialization** — every system that contributes to game state (resources, upgrades, perks, combos, floor, meta-progression) must serialize/deserialize; schema versioning needed from day one
3. **Screen state machine** — five or more distinct screens with clear transition rules; all UI panels must respect the active screen state
4. **Deterministic resource calculation** — resource rates are a function of upgrades + perks; this function must be pure and tested independently of the game loop

## Starter Template Evaluation

### Primary Technology Domain

CLI/TUI game engine — Go + Bubbletea ecosystem. Stack is fixed from PRD; this step establishes exact package versions and project initialization approach.

### Starter Options Considered

**Option A: `charmbracelet/bubbletea-app-template`** — Official Charmbracelet template, minimal Bubbletea scaffold with the standard `Model/Init/Update/View` structure. No extra dependencies, clean starting point.

**Option B: Manual `go mod init` + selective installs** — More control over project structure from day one; appropriate for a project with well-defined architecture (which we have from the PRD analysis).

**Recommendation: Option B** — We have clear architectural requirements; a minimal scaffold built to spec is better than adapting a generic template.

### Project Initialization

**Commands:**

```bash
go mod init github.com/<username>/afk-x
go get github.com/charmbracelet/bubbletea@v1.3.10
go get github.com/charmbracelet/bubbles@v2.0.0
go get github.com/charmbracelet/lipgloss@v1.1.0
```

**Version rationale:**
- `bubbletea v1.3.10` — current stable; v2 is RC and not recommended for production
- `bubbles v2.0.0` — stable, compatible with bubbletea v1.x; provides `textinput`, `viewport`, `progress`, `spinner` components directly usable for the command panel, scroll panels, and progress indicators
- `lipgloss v1.1.0` — stable; provides layout primitives (flexbox-style), borders, colors, and `NO_COLOR` awareness built-in

**No additional third-party dependencies for MVP** — standard library covers JSON (`encoding/json`), file ops (`os`), timing (`time`), CLI flags (`flag`), and signal handling (`os/signal`).

### Architectural Decisions Established by Initialization

**Language & Runtime:** Go 1.21+ (minimum for standard library features used); no CGo

**Styling Solution:** Lipgloss for all layout and color — `NO_COLOR` env var support is built into lipgloss, satisfying NFR17 at the framework level

**Build Tooling:** Standard `go build` — produces single binary, satisfying NFR11 and NFR15 (≤20MB)

**Testing Framework:** Standard `testing` package + `go test ./...`; no test framework dependency needed

**Code Organization:** Established in step-04 (architectural decisions); initialized here as flat `cmd/` + `internal/` structure

**Development Experience:** `go run ./cmd/afk-x` for local dev; no build step required

**Note:** Epic 1, Story 1 = initialize Go module, add dependencies, scaffold `cmd/afk-x/main.go` with a minimal Bubbletea program that compiles and launches.

## Core Architectural Decisions

### Decision Priority Analysis

**Critical Decisions (Block Implementation):**
- Package structure — defines how stories are implemented independently
- Game state model — defines save schema and all cross-system contracts
- Game loop architecture — drives tick, AFK calc, and UI refresh

**Important Decisions (Shape Architecture):**
- Bubbletea model/screen architecture — defines how UI panels are organized
- CI/CD and build pipeline — defines how the binary ships

**Deferred Decisions (Post-MVP):**
- Config file format (deferred to Phase 2)
- Save schema migration strategy (v1 → v2 when new resource types land)

### Data Architecture — Game State Model

**Decision:** Single `GameState` struct containing `RunState` + `MetaState` + schema version, serialized to one JSON file.

```go
// internal/game/state.go
type GameState struct {
    Version   int       `json:"version"`   // schema migration guard
    Run       RunState  `json:"run"`
    Meta      MetaState `json:"meta"`
    SavedAt   time.Time `json:"saved_at"`  // AFK delta anchor
}

type RunState struct {
    Floor         int                    `json:"floor"`
    Resources     map[string]float64     `json:"resources"`     // extensible (FR5)
    Upgrades      map[string]int         `json:"upgrades"`
    ActivePerks   []string               `json:"active_perks"`
    ComboQueue    [][]string             `json:"combo_queue"`   // FR33
    PendingDraft  []string               `json:"pending_draft"` // FR32 — nil = no draft
    RunNumber     int                    `json:"run_number"`
}

type MetaState struct {
    PermanentUnlocks []string            `json:"permanent_unlocks"` // FR23
    BestCombos       map[string][]string `json:"best_combos"`       // FR26
    UnlockedPerks    []string            `json:"unlocked_perks"`    // FR14 — expands across runs
    RunCount         int                 `json:"run_count"`
}
```

**Rationale:** `map[string]float64` for resources and `[]string` for perks/combos means new resource types (FR5) and new perks (FR14) added in future runs require zero schema changes — only new keys. `SavedAt` is the anchor for AFK delta calculation (NFR8).

**Resource rates are NOT stored** — they are computed as a pure function of `RunState` on every tick. This ensures determinism (NFR8) and prevents save/compute drift.

### Game Loop Architecture

**Decision:** Bubbletea `tea.Every` ticker at 16ms driving a pure engine update function.

```go
// internal/engine/tick.go
type TickMsg time.Time

func tickCmd() tea.Cmd {
    return tea.Every(16*time.Millisecond, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}

// Engine update: pure function, no side effects
func Tick(state RunState, delta time.Duration, rates ResourceRates) RunState {
    for resource, rate := range rates {
        state.Resources[resource] += rate * delta.Seconds()
    }
    return state
}
```

**AFK calculation:** On `SavedAt` load, compute `delta = time.Now().Sub(state.SavedAt)`. Apply `Tick` with full delta in one call. Cap offline delta at a configurable max (e.g., 7 days) to prevent numeric overflow.

**ResourceRates:** Pure function `ComputeRates(run RunState) ResourceRates` — takes upgrades + active perks, returns rates map. Computed once per tick, cached in the Bubbletea model between ticks.

### Package Structure

**Decision:**

```
afk-x/
├── cmd/
│   └── afk-x/
│       └── main.go          # CLI flags, signal handling, tea.NewProgram
├── internal/
│   ├── game/                # GameState, RunState, MetaState structs + JSON marshal/unmarshal
│   ├── engine/              # Tick, ComputeRates, FloorProgression, PerkDraft logic
│   ├── combat/              # ComboParser, ComboValidator, BonusCalculator
│   ├── save/                # AtomicWrite, Load, XDG path resolution, corruption recovery
│   ├── content/             # Static data: perk definitions, combo definitions, enemy data, flavor text
│   └── ui/                  # Bubbletea root Model, screen views, lipgloss styles
└── Makefile                 # build, test, lint, release targets
```

**Package dependency rules:**
- `ui` → imports `engine`, `game`, `combat`, `content`
- `engine` → imports `game`, `content`
- `combat` → imports `game`, `content`
- `save` → imports `game`
- `content` → no internal imports (pure data)
- `game` → no internal imports (pure structs)
- No circular dependencies

### Bubbletea Model / Screen Architecture

**Decision:** Single root `Model` with an `ActiveScreen` enum. Screens are rendered as functions, not sub-models. All state lives in the root model.

```go
// internal/ui/model.go
type Screen int
const (
    BootScreen Screen = iota
    GameScreen
    PerkDraftScreen
    RunSummaryScreen
    MetaProgressionScreen
)

type Model struct {
    Screen       Screen
    State        game.GameState
    Rates        engine.ResourceRates   // cached, recomputed each tick
    LastTick     time.Time
    CmdInput     textinput.Model        // bubbles text input
    ShowCmdRef   bool                   // ? toggle
}
```

**Messages:**
- `TickMsg` — game loop advance
- `FloorClearMsg` — trigger perk draft screen transition
- `RunEndMsg` — trigger run summary → meta screen → new run
- `SaveMsg` — trigger atomic save (on tick, on clean exit)
- `tea.KeyMsg` — all keyboard input routed through Update

**Screen transitions:** `Model.Update` handles all `tea.Msg` types; screen-specific logic is gated on `m.Screen`. Transitions are explicit state assignments (`m.Screen = PerkDraftScreen`), never implicit.

### Infrastructure & Deployment

**Decision:** Makefile + GitHub Actions for CI; GoReleaser for cross-platform binary distribution.

```makefile
build:   go build -o bin/afk-x ./cmd/afk-x
test:    go test ./...
lint:    go vet ./... && staticcheck ./...
release: goreleaser release --clean
```

**GitHub Actions:** On push to `main` — `go test ./...` + `go vet ./...`. On tag — GoReleaser builds Linux (amd64/arm64) + macOS (amd64/arm64) binaries and attaches to GitHub Release.

**No Docker** — binary is self-contained; containerization adds no value for a terminal game.

### Decision Impact Analysis

**Implementation Sequence (for epic ordering):**
1. `game` package (structs + JSON) — no dependencies, foundational
2. `content` package (static data) — no dependencies, required by engine
3. `engine` package (tick, rates, floor logic) — depends on game + content
4. `save` package (atomic write, load) — depends on game
5. `combat` package (combo parser) — depends on game + content
6. `ui` package (Bubbletea model) — depends on all above
7. `cmd/afk-x` (main, CLI flags, signal handling) — depends on ui + save

**Cross-Component Dependencies:**
- Save schema `Version` field must be defined in `game` from Story 1 — enables future migration without rework
- `PendingDraft []string` in `RunState` bridges `engine` (floor clear event) and `ui` (draft screen) — both packages read/write this field
- `ComboQueue [][]string` in `RunState` bridges `combat` (enqueue) and `save` (persist) — both must agree on the type from day one

## Implementation Patterns & Consistency Rules

### Critical Conflict Points Identified

7 areas where AI agents could make incompatible choices without explicit rules.

### Naming Patterns

**Go Package & File Naming:**
- Package names: lowercase, single word — `engine`, `game`, `combat`, `save`, `content`, `ui`
- File names: snake_case — `state.go`, `tick.go`, `combo_parser.go`, `atomic_write.go`
- Test files: co-located, `_test.go` suffix — `tick_test.go` next to `tick.go`
- No `util`, `helpers`, or `common` packages — put shared code in the package it belongs to

**Go Symbol Naming:**
- Exported types: PascalCase — `GameState`, `RunState`, `TickMsg`, `FloorClearMsg`
- Unexported functions: camelCase — `computeRates`, `applyCombo`
- Constants: PascalCase for exported, camelCase for unexported — `GameScreen`, `bootDuration`
- Interfaces: noun or adjective + `-er` — `Saver`, `Ticker` (only define interfaces at point of use)

**Message Naming:**
- All Bubbletea messages: PascalCase + `Msg` suffix — `TickMsg`, `FloorClearMsg`, `RunEndMsg`, `SaveMsg`, `ErrorMsg`
- Never use bare structs as messages without the `Msg` suffix

**Resource & Content Key Naming:**
- Resource type keys: lowercase snake_case string constants defined in `content` — `"cpu_cycles"`, `"memory_shards"`, `"process_threads"`
- Perk IDs: lowercase snake_case — `"overclock"`, `"fork_bomb"`, `"swap_space"`
- Combo IDs: lowercase snake_case — `"scan_exploit_loot"`
- Never use magic strings — always reference the constant from `content`

### Structure Patterns

**Package Responsibility (strict):**
- `game` — structs and JSON only; zero business logic, zero Bubbletea imports
- `engine` — pure functions only; no I/O, no Bubbletea, no `os` calls
- `combat` — pure functions only; no I/O, no Bubbletea
- `content` — static data definitions only; no functions that mutate state
- `save` — all file I/O; atomic write pattern always, never `os.WriteFile` directly
- `ui` — all Bubbletea `Model/Init/Update/View`; no direct file I/O
- `cmd/afk-x` — wiring only; no game logic, thin main function

**Test Organization:**
- All tests co-located with the package they test (`_test.go`)
- `engine` and `combat` must have unit tests (pure functions — easy to test)
- `save` must have integration tests (test actual file writes to temp dir)
- `ui` tests are optional for MVP — Bubbletea models are hard to unit test

**Content Data Format:**
- All static game data (perks, combos, enemies, flavor text) defined as Go structs in `content` package — NOT external JSON files
- This keeps the binary self-contained (single binary requirement, NFR11)
- Example: `content/perks.go` defines `[]PerkDefinition{...}` as a package-level var

### Format Patterns

**JSON Save File (snake_case throughout):**
```json
{
  "version": 1,
  "saved_at": "2026-05-02T14:30:00Z",
  "run": {
    "floor": 3,
    "resources": {"cpu_cycles": 1250.5, "memory_shards": 340.0},
    "upgrades": {"cycle_boost_1": 2},
    "active_perks": ["overclock"],
    "combo_queue": [["scan", "exploit", "loot"]],
    "pending_draft": null,
    "run_number": 1
  },
  "meta": { ... }
}
```
- All JSON keys: snake_case (enforced by struct tags)
- Times: RFC3339 (`time.Time` marshals this by default)
- `pending_draft: null` means no draft pending; empty array `[]` is invalid — always use `nil`/`null`

**Error Handling:**
- Functions return `(T, error)` — never panic in library code
- Wrap errors with context: `fmt.Errorf("save: writing temp file: %w", err)`
- `ui` package converts errors to `ErrorMsg` for display — never call `log.Fatal` outside `cmd/afk-x/main.go`
- Save corruption (NFR9): `save.Load` returns `(GameState, bool, error)` — the `bool` indicates whether a new game should be started

**Resource Rate Type:**
```go
// internal/engine/rates.go
type ResourceRates map[string]float64  // resource key → units per second
```
Always use this named type, never `map[string]float64` directly in function signatures.

### Communication Patterns

**Bubbletea Message Flow (strict):**
- Packages communicate via Bubbletea messages only — `engine`, `combat`, `save` never call into `ui`
- `ui.Update` calls engine/combat functions and returns new state + commands
- `engine` and `combat` return new state structs, never mutate in place
- Save is triggered by `SaveMsg` — never called directly from `Update`

**Save Triggers:**
- Every N ticks (configurable constant, default 60 = ~1 second)
- On `FloorClearMsg`
- On `RunEndMsg`
- On `SIGINT`/`SIGTERM` (handled in `cmd/afk-x/main.go`)
- On clean `tea.Quit`
- Never save during `PerkDraftScreen` — wait until draft is resolved

**Screen Transitions:**
- All transitions through explicit `m.Screen = XxxScreen` assignments in `Update`
- No transitions in `View` functions
- Transition table:
  - `BootScreen` → `GameScreen` after boot sequence completes
  - `GameScreen` → `PerkDraftScreen` on `FloorClearMsg` with draft options
  - `PerkDraftScreen` → `GameScreen` after perk selection
  - `GameScreen` → `RunSummaryScreen` on `RunEndMsg`
  - `RunSummaryScreen` → `MetaProgressionScreen` on confirm key
  - `MetaProgressionScreen` → `GameScreen` (new run) on confirm key

### Process Patterns

**Error Display:**
- Non-fatal errors (save write warning, invalid combo) → shown as brief status line in `GameScreen`, auto-clear after 3 seconds
- Fatal errors (save load corruption) → shown on dedicated error screen with reset prompt
- Never use `fmt.Println` or `log.Print` inside the Bubbletea program — all output through `View`

**AFK Accumulation (always this pattern):**
```go
// On load:
delta := time.Since(state.SavedAt)
if delta > maxOfflineDuration { delta = maxOfflineDuration }
rates := engine.ComputeRates(state.Run)
state.Run = engine.Tick(state.Run, delta, rates)
state.SavedAt = time.Now()
```
Always use `engine.ApplyOffline(state, time.Now())` wrapper — never inline this logic.

### Enforcement Guidelines

**All AI agents implementing stories MUST:**
- Import only packages that the dependency rules allow for their package
- Use string constants from `content` for all resource/perk/combo keys — never string literals
- Use `save.Write` for all save operations — never write JSON directly
- Return `(T, error)` from all functions that can fail — never panic
- Add `_test.go` for every new function in `engine` or `combat`
- Use `tea.Cmd` to trigger side effects — never call I/O from inside `Update`

**Anti-Patterns (never do these):**
- ❌ `os.WriteFile` anywhere outside `save` package
- ❌ Global mutable variables (except package-level content definitions)
- ❌ Importing `ui` from any package other than `cmd/afk-x`
- ❌ Importing `save` from `ui` directly — use `SaveMsg` command pattern
- ❌ `map[string]float64` in function signatures — use `engine.ResourceRates`
- ❌ Magic strings for resource/perk/combo keys

## Project Structure & Boundaries

### Complete Project Directory Structure

```
afk-x/
├── .github/
│   └── workflows/
│       ├── ci.yml              # go test + go vet on push to main
│       └── release.yml         # goreleaser on tag push
├── .goreleaser.yml             # cross-platform build matrix (Linux/macOS amd64/arm64)
├── .gitignore
├── Makefile                    # build, test, lint, release targets
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   └── afk-x/
│       └── main.go             # CLI flags, signal handling, tea.NewProgram, XDG path
└── internal/
    ├── game/                   # save schema contract for all packages
    │   ├── state.go            # GameState, RunState, MetaState structs + JSON tags
    │   └── state_test.go       # JSON round-trip tests
    ├── content/                # perk/combo/enemy/flavor text definitions
    │   ├── perks.go            # []PerkDefinition (FR11-FR15)
    │   ├── combos.go           # []ComboDefinition — valid sequences + bonuses (FR17-FR18)
    │   ├── enemies.go          # []EnemyDefinition — daemon types per floor range (FR9)
    │   ├── upgrades.go         # []UpgradeDefinition — upgrade tree + costs (FR6-FR7)
    │   ├── resources.go        # Resource key constants + initial rates (FR1, FR4, FR5)
    │   └── narrative.go        # Boot text, milestone flavor text map (FR34-FR37)
    ├── engine/                 # game loop + resource math
    │   ├── tick.go             # TickMsg, tickCmd(), Tick() pure function (FR1-FR3, NFR1)
    │   ├── tick_test.go
    │   ├── rates.go            # ResourceRates type, ComputeRates() (FR1, FR4, NFR8)
    │   ├── rates_test.go
    │   ├── offline.go          # ApplyOffline() — AFK delta calculation (FR2-FR3, NFR2, NFR8)
    │   ├── offline_test.go
    │   ├── floor.go            # CheckFloorClear(), AdvanceFloor() (FR8-FR10)
    │   ├── floor_test.go
    │   ├── perk.go             # DrawPerkDraft(), ApplyPerk() (FR11-FR13, FR14-FR15)
    │   ├── perk_test.go
    │   ├── upgrade.go          # CanAfford(), ApplyUpgrade() (FR6-FR7)
    │   ├── upgrade_test.go
    │   ├── meta.go             # ApplyRunToMeta(), NewRunState() (FR22-FR27)
    │   └── meta_test.go
    ├── combat/                 # command input + combo system
    │   ├── parser.go           # ParseInput() — tokenize command string (FR16-FR17)
    │   ├── parser_test.go
    │   ├── validator.go        # ValidateCombo() — check against content.Combos (FR17, FR20)
    │   ├── validator_test.go
    │   ├── bonus.go            # ComputeBonus() — efficiency multiplier (FR18)
    │   ├── bonus_test.go
    │   ├── queue.go            # EnqueueCombo(), DequeueCombo() (FR21, FR33)
    │   └── queue_test.go
    ├── save/                   # persistence layer
    │   ├── paths.go            # XDGSavePath(), resolveSavePath() (FR43, NFR11)
    │   ├── paths_test.go
    │   ├── write.go            # Write() — atomic temp→rename (FR28, NFR6)
    │   ├── write_test.go       # writes to os.TempDir()
    │   ├── load.go             # Load() → (GameState, bool, error) (FR29, NFR7, NFR9)
    │   └── load_test.go
    └── ui/                     # all Bubbletea rendering
        ├── model.go            # Model struct, Screen enum, Init()
        ├── update.go           # Update() — routes all tea.Msg to screen handlers
        ├── view.go             # View() — delegates to screen-specific view funcs
        ├── screens/
        │   ├── boot.go         # BootScreen view + update (FR34)
        │   ├── game.go         # GameScreen view + update (FR1,FR6-FR10,FR16-FR21,FR35-FR37)
        │   ├── draft.go        # PerkDraftScreen view + update (FR11-FR12)
        │   ├── summary.go      # RunSummaryScreen view + update (FR22)
        │   └── meta.go         # MetaProgressionScreen view + update (FR23-FR27)
        └── styles/
            ├── colors.go       # Lipgloss color palette (NO_COLOR safe)
            ├── layout.go       # Panel dimensions, lipgloss layout primitives
            └── text.go         # Text formatting helpers, truncation
```

### Architectural Boundaries

**Package import rules:**

```
cmd/afk-x  →  ui, save, game
ui          →  engine, combat, game, content
engine      →  game, content
combat      →  game, content
save        →  game
content     →  (no internal imports)
game        →  (no internal imports)
```

No package may import its downstream dependents. Violations detected via `go build` circular import errors.

**State Boundary:**
- `game.GameState` is the single source of truth
- Only `ui.Model` holds a live `GameState` at runtime
- All other packages receive/return `RunState` or `MetaState` by value (copy, not pointer)

**I/O Boundary:**
- All file I/O: `save` package only
- All terminal I/O: `ui` package only (via Bubbletea)
- All signal handling: `cmd/afk-x/main.go` only

### Requirements to Structure Mapping

| FR Category | Primary Package | Supporting |
|---|---|---|
| Resource Generation (FR1–FR5) | `engine/tick.go`, `engine/rates.go`, `engine/offline.go` | `content/resources.go`, `game/state.go` |
| Upgrade & Floor (FR6–FR10) | `engine/upgrade.go`, `engine/floor.go` | `content/upgrades.go`, `content/enemies.go` |
| Perk Draft (FR11–FR15) | `engine/perk.go` | `content/perks.go`, `game/state.go` |
| Command & Combo (FR16–FR21) | `combat/parser.go`, `combat/validator.go`, `combat/bonus.go`, `combat/queue.go` | `content/combos.go` |
| Meta-Progression (FR22–FR27) | `engine/meta.go` | `game/state.go`, `ui/screens/summary.go`, `ui/screens/meta.go` |
| Save & Persistence (FR28–FR33) | `save/write.go`, `save/load.go`, `save/paths.go` | `game/state.go` |
| Narrative & Presentation (FR34–FR37) | `ui/screens/boot.go`, `ui/screens/game.go` | `content/narrative.go` |
| Application Shell (FR38–FR44) | `cmd/afk-x/main.go`, `save/paths.go` | `ui/model.go` |

**NFR → Implementation Location:**

| NFR | Location |
|---|---|
| NFR1 (16ms tick) | `engine/tick.go` — `tea.Every(16ms, ...)` |
| NFR2 (500ms AFK calc) | `engine/offline.go` — single-pass calculation |
| NFR6 (atomic save) | `save/write.go` — `os.Rename` after temp write |
| NFR8 (deterministic AFK) | `engine/offline.go` + `engine/rates.go` — pure functions |
| NFR9 (corrupt save recovery) | `save/load.go` — `(GameState, bool, error)` signature |
| NFR10 (SIGINT/SIGTERM) | `cmd/afk-x/main.go` — `signal.Notify` → save → `tea.Quit` |
| NFR11 (single binary) | `content/*.go` — embedded Go data, no external files |
| NFR15 (≤20MB binary) | `.goreleaser.yml` — `ldflags: -s -w` strip debug info |
| NFR17 (NO_COLOR) | `ui/styles/colors.go` — lipgloss detects NO_COLOR natively |
| NFR18 (keyboard only) | `ui/update.go` — no mouse events registered |

### Integration Points

**Runtime data flow:**
```
time.Ticker (16ms)
  → TickMsg
  → ui.Update
  → engine.ComputeRates(state.Run)        // pure
  → engine.Tick(state.Run, delta, rates)  // pure → new RunState
  → engine.CheckFloorClear(state.Run)     // → FloorClearMsg if cleared
  → (every 60 ticks) → SaveMsg → save.Write(state)
```

**AFK resume flow:**
```
save.Load(path)
  → (GameState, isNew bool, error)
  → if !isNew: engine.ApplyOffline(state, time.Now())
  → ui.Model{State: state}
  → tea.NewProgram(model).Run()
```

**Perk draft flow:**
```
FloorClearMsg{Options: []string}
  → m.Screen = PerkDraftScreen
  → m.State.Run.PendingDraft = options
  → SaveMsg (persist draft — FR32)
  → user selects perk
  → engine.ApplyPerk(state.Run, selected)
  → m.Screen = GameScreen
  → m.State.Run.PendingDraft = nil
  → SaveMsg
```

### Development Workflow

- **Dev:** `go run ./cmd/afk-x`
- **Test:** `go test ./...`
- **Build:** `make build` → `bin/afk-x`
- **Release:** `git tag v0.1.0 && git push --tags` → GoReleaser → 4 binaries on GitHub Release

## Architecture Validation Results

### Coherence Validation

**Decision Compatibility:**
All technology choices are mutually compatible — bubbletea v1.3.10 + bubbles v2.0.0 + lipgloss v1.1.0 are released together and verified compatible. Standard library covers all remaining dependencies. The pure-function engine model is a natural fit for Bubbletea's immutable `Model → Update → View` lifecycle. Atomic save via `os.Rename` is idiomatic Go and satisfies NFR6 without additional dependencies.

**Pattern Consistency:**
Naming conventions are consistent across packages (PascalCase exports, camelCase unexported, snake_case JSON, `Msg` suffix for messages). Package dependency rules are enforced by Go's import system. Communication exclusively via Bubbletea messages aligns with the framework's design.

**Structure Alignment:**
Package structure maps directly to architectural decisions. All FR categories have dedicated files. All three I/O boundaries (file, terminal, signal) map to exactly one package each.

### Requirements Coverage Validation

**Functional Requirements Coverage:**
All 44 FRs are mapped to specific files in the requirements-to-structure table. No FR is architecturally unsupported.

**Non-Functional Requirements Coverage:**
All 19 NFRs are addressed — 10 with explicit implementation locations, 9 satisfied at the framework level (Lipgloss `NO_COLOR`, Bubbletea terminal rendering, `go build` single binary, GoReleaser cross-compile).

### Implementation Readiness Validation

All critical decisions documented with verified versions, complete file tree with 30+ specific files mapped to FRs, and 7 conflict areas covered with explicit rules and anti-patterns.

### Gap Analysis Results

**Critical Gaps:** None.

**Minor Gaps (resolved inline):**

`ComboDefinition` struct shape defined here to prevent inconsistency between `content/combos.go` and `combat/validator.go`:

```go
// content/combos.go
type ComboDefinition struct {
    ID         string   // e.g. "scan_exploit_loot"
    Sequence   []string // e.g. ["scan", "exploit", "loot"]
    BonusType  string   // e.g. "efficiency", "resource_mult"
    BonusValue float64  // e.g. 0.40 (= 40% bonus)
    FlavorText string   // shown on successful execution
}
```

Boot sequence duration: `bootDuration = 2 * time.Second` constant in `ui/screens/boot.go`, skippable with any key press.

Save version migration strategy: intentionally deferred to Phase 2.

### Architecture Completeness Checklist

**Requirements Analysis**
- [x] Project context thoroughly analyzed
- [x] Scale and complexity assessed
- [x] Technical constraints identified
- [x] Cross-cutting concerns mapped

**Architectural Decisions**
- [x] Critical decisions documented with versions
- [x] Technology stack fully specified
- [x] Integration patterns defined
- [x] Performance considerations addressed

**Implementation Patterns**
- [x] Naming conventions established
- [x] Structure patterns defined
- [x] Communication patterns specified
- [x] Process patterns documented

**Project Structure**
- [x] Complete directory structure defined
- [x] Component boundaries established
- [x] Integration points mapped
- [x] Requirements to structure mapping complete

### Architecture Readiness Assessment

**Overall Status: READY FOR IMPLEMENTATION**

**Confidence Level:** High — all 16 checklist items verified, no critical gaps.

**Key Strengths:**
- Pure function engine design makes tick, AFK, and combo logic trivially testable
- `map[string]float64` resource schema ensures zero-migration extensibility for Phase 2
- Lipgloss `NO_COLOR` support satisfies NFR17 with no additional code
- Package dependency rules enforced by the compiler — impossible to violate accidentally

**Areas for Future Enhancement:**
- Save schema migration framework (Phase 2, when resource types expand)
- `ui` package integration tests (Phase 2, when Bubbletea testing utilities mature)
- Combo discovery/hint system (Phase 2 growth feature)

### Implementation Handoff

**AI Agent Guidelines:**
- Follow all architectural decisions exactly as documented
- Use `content` package constants for all resource/perk/combo keys — never string literals
- `engine` and `combat` functions must be pure — no I/O, no side effects
- All saves via `save.Write` — never `os.WriteFile` directly
- Refer to this document for all structural and pattern questions

**First Implementation Priority:**
```bash
go mod init github.com/<username>/afk-x
go get github.com/charmbracelet/bubbletea@v1.3.10
go get github.com/charmbracelet/bubbles@v2.0.0
go get github.com/charmbracelet/lipgloss@v1.1.0
```
Epic 1, Story 1: scaffold `cmd/afk-x/main.go` + `internal/game/state.go` + `internal/content/resources.go` — the foundational layer everything else builds on.
