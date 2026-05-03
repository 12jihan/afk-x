# Story 1.3: Content Package — Resource & Narrative Definitions

Status: done

## Story

As a developer,
I want resource key constants, initial rates, and narrative text defined in the content package,
so that all other packages reference game content without magic strings.

## Acceptance Criteria

1. **Given** `content/resources.go` is imported, **When** the constants are accessed, **Then** `CPUCycles`, `MemoryShards`, and `ProcessThreads` key constants are defined and non-empty
2. **Given** `content/narrative.go`, **When** `BootText` is accessed, **Then** it contains the tower lore flavor text (non-empty string)
3. **Given** `content/narrative.go`, **When** `MilestoneText(floor int)` is called with a milestone floor number, **Then** it returns the corresponding flavor text string
4. **Given** `go test ./internal/content/...`, **When** run, **Then** all tests pass with no compilation errors

## Tasks / Subtasks

- [x] Task 1: Implement `internal/content/resources.go` (AC: 1)
  - [x] Define exported string constants: `CPUCycles = "cpu_cycles"`, `MemoryShards = "memory_shards"`, `ProcessThreads = "process_threads"`
  - [x] Define `BaseRates map[string]float64` package-level var with an entry for each constant (all values > 0)
  - [x] No functions, no internal imports — pure constant/var declarations only

- [x] Task 2: Implement `internal/content/narrative.go` (AC: 2, 3)
  - [x] Define exported `BootText string` constant with multi-line terminal-aesthetic tower lore (non-empty)
  - [x] Define `MilestoneText(floor int) string` function that returns flavor text for milestone floors (5, 10, 20, 25) and `""` for all others
  - [x] Milestone texts must be non-empty strings with distinct lore per floor

- [x] Task 3: Replace `internal/content/content.go` stub (AC: 4)
  - [x] Replace the stub with a proper package doc comment — or delete it entirely (individual files each carry `package content`)
  - [x] Ensure `go build ./...` still passes after the change

- [x] Task 4: Create `internal/content/content_test.go` with all AC tests (AC: 1, 2, 3, 4)
  - [x] `TestResourceKeyConstants` — `CPUCycles`, `MemoryShards`, `ProcessThreads` are all non-empty strings and distinct from each other
  - [x] `TestBaseRatesDefined` — `BaseRates` has an entry for each of the three constants, all values `> 0`
  - [x] `TestBootTextNonEmpty` — `BootText` is a non-empty string
  - [x] `TestMilestoneTextForMilestoneFloors` — floors 5, 10, 20, 25 each return non-empty strings
  - [x] `TestMilestoneTextEmptyForNonMilestone` — floors 1, 2, 3, 4, 6, 7 return `""`
  - [x] `TestResourceKeysMatchBaseRates` — every constant value appears as a key in `BaseRates`

- [x] Task 5: Resolve deferred item from Story 1.2 (AC: 4, regression check)
  - [x] Update `internal/game/state_test.go` to import `content` and replace string literals `"cpu_cycles"`, `"memory_shards"`, `"process_threads"` with `content.CPUCycles`, `content.MemoryShards`, `content.ProcessThreads`
  - [x] Run `go test ./...` to confirm no regressions

### Senior Developer Review (AI)

**Review Date:** 2026-05-03
**Outcome:** Changes Requested
**Layers:** Blind Hunter, Edge Case Hunter, Acceptance Auditor

#### Action Items

- [x] [Review][Patch] `TestMilestoneTextsAreDistinct` overwrites `seen[text]` after collision — third duplicate would report wrong origin floor [internal/content/content_test.go:106]
- [x] [Review][Patch] `milestoneTexts` var lacks "treat as read-only" comment — inconsistent with `BaseRates` documentation pattern [internal/content/narrative.go:19]
- [x] [Review][Defer] `BaseRates` external mutation risk — any importer can corrupt the baseline; mitigation belongs in Story 1.4 (`ComputeRates` must copy the map, not reference directly) [internal/content/resources.go:14] — deferred, pre-existing
- [x] [Review][Defer] `TestMilestoneTextForMilestoneFloors` hard-codes milestone floor list — must be manually updated if milestones change; by design for now [internal/content/content_test.go:87] — deferred, pre-existing
- [x] [Review][Defer] `TestMilestoneTextEmptyForNonMilestone` doesn't cover floor 0 or negative values — behavior is correct but contract is undocumented [internal/content/content_test.go:97] — deferred, pre-existing
- [x] [Review][Defer] `TestBaseRatesDefined`/`TestResourceKeysMatchBaseRates` don't assert `len(BaseRates) == 3` — no extra-key invariant enforced [internal/content/content_test.go:49–74] — deferred, pre-existing
- [x] [Review][Defer] `BaseRates`-to-engine integration test missing — belongs in Story 1.4 when `engine/rates.go` is implemented — deferred, pre-existing

## Dev Notes

### What This Story Does

Story 1.2 used string literals (`"cpu_cycles"` etc.) as placeholders in tests because this package didn't exist. This story creates the two files that define the ground truth for all resource keys and narrative text. Every other package that touches resources or game narrative must import these constants — never string literals.

**Files changed:**
- `internal/content/content.go` — replace stub (or delete; see Task 3)
- `internal/content/resources.go` — new
- `internal/content/narrative.go` — new
- `internal/content/content_test.go` — new
- `internal/game/state_test.go` — minor update (deferred item from 1.2, Task 5)

### Architecture Compliance (MUST FOLLOW)

**Package responsibility — `content` is static data ONLY:**
- ✅ Allowed: exported constants, exported package-level vars, one pure query function (`MilestoneText`)
- ❌ FORBIDDEN: any internal package imports (`game`, `engine`, `combat`, `save`, `ui`)
- ❌ FORBIDDEN: mutable state (no `init()` that modifies vars, no global state)
- ❌ FORBIDDEN: any I/O, file reads, or external data sources — all content is embedded Go literals
- ❌ FORBIDDEN: Bubbletea, Lipgloss, or any Charm imports

**Import rule for `content` package:**
```go
// No imports needed for resources.go or narrative.go — pure Go literals
package content
```

**Why embedded Go literals, not external files:**
NFR11 requires a single self-contained binary. External JSON/YAML data files would require embedding or bundling. Pure Go literals are compiled in — zero runtime dependency.

**Dependency graph:**
```
content  →  (no internal imports)
game     →  (no internal imports)
engine   →  game, content   ← Story 1.4 will import content constants
```

[Source: `architecture.md` → Package Structure, Package Responsibility, Content Data Format]

### Exact Resource Constants

```go
// internal/content/resources.go
package content

// Resource key constants — used as map keys in RunState.Resources (map[string]float64).
// All keys are lowercase snake_case to match the JSON save format.
const (
    CPUCycles      = "cpu_cycles"
    MemoryShards   = "memory_shards"
    ProcessThreads = "process_threads"
)

// BaseRates defines the initial resource generation rate (units per second)
// for a new run with no upgrades or perks active.
// ComputeRates() in engine/rates.go returns these as the starting point.
var BaseRates = map[string]float64{
    CPUCycles:      1.0,   // 1 CPU cycle per second baseline
    MemoryShards:   0.5,   // 0.5 memory shards per second baseline
    ProcessThreads: 0.25,  // 0.25 process threads per second baseline
}
```

**Why `var`, not `const` for `BaseRates`:** Go does not support const maps. Use a package-level `var` — treat it as read-only (never mutate it; engine takes a copy via `ComputeRates`).

### Narrative Content — BootText & MilestoneText

The game uses a terminal/computing aesthetic (FR36). All narrative text should feel like system logs from an ancient, self-aware computing tower. Think: old Unix system messages + haunted mainframe.

**BootText pattern** (FR34 — shown during boot sequence):

```go
const BootText = `SYSTEM INITIALIZE — THE TOWER
> Detecting hardware...          [OK]
> Mounting sector 0x00...        [OK]
> Daemon census: 1024 active processes
> WARNING: Unauthorized access detected on FLOOR 01
> Administrator terminal connected.
>
> Your resources are limited.
> The tower is not.
>
> Ascend.`
```

**MilestoneText pattern** (FR35 — shown on designated floor clears):

Milestone floors for MVP: **5, 10, 20, 25**

```go
// internal/content/narrative.go
package content

var milestoneTexts = map[int]string{
    5: `SECTOR 0x05 CLEARED
The daemon cluster on this floor was running a single process,
unchanged since the tower was seeded. Its last instruction: WAIT.
It has been waiting for 40 years.
You interrupted it.`,

    10: `SECTOR 0x0A CLEARED
You have cleared the first tenth of the lower stack.
The tower grows quieter here. Fewer daemons. More entropy.
The architects built this place to last.
They did not build it to be finished.`,

    20: `SECTOR 0x14 CLEARED
A recursive function. No base case.
It has been calling itself since initialization.
You broke the loop.
Somewhere in memory, a value that was never freed
has finally been collected.`,

    25: `SECTOR 0x19 CLEARED
ADMIN NOTE — ARCHIVED LOG:
"We built the tower to solve a problem we no longer remember.
The solution is still running."`,
}

// MilestoneText returns the flavor text for a milestone floor.
// Returns "" (empty string) if floor is not a milestone — callers should
// check for non-empty before displaying.
func MilestoneText(floor int) string {
    return milestoneTexts[floor]
}
```

**Unexported `milestoneTexts` map, exported `MilestoneText` function:** The map is implementation detail; the function is the contract. Story 2.3 will call `MilestoneText(floor)` and check `if text != ""`.

### Task 3 — Handling `content.go` Stub

The stub at `internal/content/content.go` currently contains:
```go
// Package content provides static game content definitions.
// TODO: Implemented in Story 1.3
package content
```

**Recommended approach:** Replace the file with a proper package doc:
```go
// Package content provides static game content definitions for afk-x.
// All content is embedded as Go literals — no external files, no I/O.
// This package has no internal imports and no mutable state.
package content
```

This preserves the package documentation entry point while removing the TODO. Alternatively, delete the file entirely — Go doesn't require a "main" file per package. Either approach works as long as `go build ./...` passes.

### Task 5 — Resolving the Story 1.2 Deferred Item

The `internal/game/state_test.go` file currently uses string literals for resource keys:
```go
original.Run.Resources["cpu_cycles"] = 1250.5   // ← string literal
state.Run.Resources["cpu_cycles"] = 100.5        // ← string literal
```

After Task 1 is complete, update those lines to use constants:
```go
import "github.com/12jihan/afk-x/internal/content"

original.Run.Resources[content.CPUCycles] = 1250.5    // ← constant
state.Run.Resources[content.CPUCycles] = 100.5         // ← constant
```

**Important:** `game/state_test.go` is in `package game_test`, which CAN import `content` (no import cycle — `game` itself still has zero internal imports; only the *test binary* imports both). The production `game` package remains unchanged.

**Locations to update in `state_test.go`:**
- `TestJSONRoundTrip`: `Resources["cpu_cycles"]`, `Resources["memory_shards"]`
- `TestResourcesMapType`: the `cases` map and `state.Run.Resources[...]` writes
- `TestNilSlicesSerializeAsNull`: no resource keys used — no change needed

### Testing Standards

- Tests live at `internal/content/content_test.go` — use `package content_test` (external/black-box)
- Import only `testing` and `github.com/12jihan/afk-x/internal/content`
- No test helpers needed — all assertions are simple equality checks
- `go test ./internal/content/...` must pass (AC4)
- `go test ./...` must still pass (CI requirement)

### Previous Story Intelligence (from Story 1.2)

- **External test package pattern:** Use `package content_test` + explicit import (same as Story 1.2 used `package game_test`)
- **Deferred item explicitly tracked:** `deferred-work.md` entry "Update TestResourcesMapType and TestJSONRoundTrip to use content package resource key constants" — Task 5 resolves this
- **No type aliases needed:** Unlike Story 1.2 which aliased `GameState` and `RunState`, this story's test only calls functions and accesses constants — no type aliases required
- **Module path confirmed:** `github.com/12jihan/afk-x`

### Project Structure Notes

- `content` package: `internal/content/` — no subdirectories
- Do NOT create `perks.go`, `combos.go`, `enemies.go`, or `upgrades.go` in this story — those belong to Stories 3.1 and 4.1 respectively
- Go version: 1.26.2 (do not change `go.mod`)
- No new dependencies — `content` package has zero imports

### References

- Resource key names: `architecture.md` → "Resource & Content Key Naming"
- `content` package responsibility: `architecture.md` → "Package Responsibility (strict)"
- `BaseRates` usage: `architecture.md` → "Game Loop Architecture" (`ComputeRates` starts from base rates)
- FR34 (boot sequence): `epics.md` → Story 1.5 AC1 (boot sequence displays flavor text)
- FR35 (milestone text): `epics.md` → Story 2.3 AC1 (milestone flavor text shown on floor clear)
- FR36 (terminal aesthetic): `prd.md` → Functional Requirements
- Deferred item: `_bmad-output/implementation-artifacts/deferred-work.md` → Story 1.2 entry

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-6

### Debug Log References

### Completion Notes List

- `content` package has zero internal imports — pure Go literals, fully architecture-compliant
- TDD: test file written first (RED — build failures on undefined symbols), then implementation (GREEN — 8/8 pass)
- `resources.go`: CPUCycles/MemoryShards/ProcessThreads constants + BaseRates var (map[string]float64, all > 0)
- `narrative.go`: BootText const (multi-line terminal lore) + unexported milestoneTexts map + exported MilestoneText(floor int) string
- `content.go` stub replaced with proper package doc comment
- Task 5: `game/state_test.go` updated — added `content` import, replaced all string literals with content constants in TestJSONRoundTrip and TestResourcesMapType; resolves Story 1.2 deferred item
- `go test ./...` — all packages pass, zero regressions

### File List

- `internal/content/resources.go` — new
- `internal/content/narrative.go` — new
- `internal/content/content.go` — stub replaced with package doc
- `internal/content/content_test.go` — new (8 tests)
- `internal/game/state_test.go` — updated (content import + constants replacing string literals)

### Change Log

- 2026-05-02: Story 1.3 complete — content package implemented with resource constants, base rates, boot/milestone narrative text, 8-test suite; Story 1.2 deferred item resolved
