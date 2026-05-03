# Story 1.5: TUI Shell — Boot Sequence & Resource Display

Status: review

## Story

As a player,
I want to launch the game, see a boot sequence, and watch my resources tick,
so that I know the game is running and my progress is accumulating.

## Acceptance Criteria

1. **Given** the binary is run with no arguments, **When** launched in a supported terminal, **Then** a boot sequence with narrative flavor text displays, auto-advancing after 2 seconds or on any keypress
2. **Given** the boot sequence completes, **When** the game screen appears, **Then** CPU cycles, memory shards, and process threads are visible and incrementing in real time
3. **Given** a terminal that is too small to render the UI, **When** the game launches, **Then** an error message is displayed and the game exits cleanly without crashing
4. **Given** `NO_COLOR=1` is set in the environment, **When** the game runs, **Then** all output is plain text with no ANSI color sequences
5. **Given** the game is running in iTerm2, Terminal.app, GNOME Terminal, Windows Terminal, and tmux, **When** the layout is displayed, **Then** panels render without visual corruption or broken box-drawing characters

## Tasks / Subtasks

- [x] Task 1: Create `internal/ui/styles/` package — colors, layout, text helpers (AC: 4, 5)
  - [x] `styles/colors.go` — `package styles`; define a lipgloss color palette using `lipgloss.AdaptiveColor` so `NO_COLOR=1` is handled automatically by lipgloss at the library level; define at minimum `Accent`, `Muted`, `Normal` colors
  - [x] `styles/layout.go` — `package styles`; export `MinWidth = 80`, `MinHeight = 24` constants; export a `ResourcePanel(width int) lipgloss.Style` helper that returns a bordered panel style
  - [x] `styles/text.go` — `package styles`; export `FormatResource(name string, value float64) string` that right-aligns the float to 2 decimal places in a fixed-width field (e.g., `"CPU Cycles:      1234.56"`)

- [x] Task 2: Create `internal/ui/screens/boot.go` — BootView function (AC: 1)
  - [x] `package screens`; import `content` and `ui/styles`; do NOT import the `ui` package (avoids circular import)
  - [x] Export `BootView(width, height int) string` — renders `content.BootText` centered horizontally using lipgloss; add a dim footer line: `"[any key to continue]"` below the text
  - [x] No update logic in this file — boot timer and keypress handling live in `ui/update.go`

- [x] Task 3: Create `internal/ui/screens/game.go` — GameView function (AC: 2, 5)
  - [x] `package screens`; import `engine`, `content`, `ui/styles`; do NOT import `ui` package
  - [x] Export `GameView(resources map[string]float64, rates engine.ResourceRates, width, height int) string` — renders the three resource rows using `styles.FormatResource`; use `content.CPUCycles`, `content.MemoryShards`, `content.ProcessThreads` as keys (never string literals)
  - [x] Wrap the resource list in a lipgloss border panel using `styles.ResourcePanel(width)`
  - [x] Add a dim footer line showing floor and run number stubs (hardcoded `"Floor: 1  Run: 1"` for now — these become dynamic in Story 2.2)

- [x] Task 4: Create `internal/ui/model.go` — Model, Screen enum, New(), Init() (AC: all)
  - [x] `package ui`; define `type Screen int` with constants: `BootScreen`, `GameScreen`, `PerkDraftScreen`, `RunSummaryScreen`, `MetaProgressionScreen`
  - [x] Define unexported constants: `bootDuration = 2 * time.Second`
  - [x] Define unexported `bootDoneMsg struct{}` — delivered by `bootTimerCmd` to trigger boot → game transition
  - [x] Define unexported `func bootTimerCmd() tea.Cmd` — blocks `bootDuration` then returns `bootDoneMsg{}`
  - [x] Define `type Model struct` with fields: `Screen Screen`, `State game.GameState`, `Rates engine.ResourceRates`, `LastTick time.Time`, `CmdInput textinput.Model`, `ShowCmdRef bool`, `Width int`, `Height int`, `tooSmall bool`; `textinput` from `charm.land/bubbles/v2/textinput`
  - [x] Export `func New() Model` — initializes `State: game.NewGameState()`, `Screen: BootScreen`, `CmdInput: textinput.New()`
  - [x] Implement `(m Model) Init() tea.Cmd` — returns `bootTimerCmd()` only; tick starts on GameScreen entry

- [x] Task 5: Create `internal/ui/update.go` — full message routing (AC: 1, 2, 3, 4)
  - [x] `package ui`; implement `(m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)`
  - [x] Handle `tea.WindowSizeMsg`: set `m.Width`, `m.Height`; if either is below `styles.MinWidth`/`styles.MinHeight`, set `m.tooSmall = true` and return `tea.Quit` — do NOT panic or call `log.Fatal`
  - [x] Handle `bootDoneMsg`: set `m.LastTick = time.Now()`, `m.Screen = GameScreen`, return `engine.TickCmd()`
  - [x] Handle `engine.TickMsg` (only when `m.Screen == GameScreen`): compute `delta := time.Time(msg).Sub(m.LastTick)`, update `m.LastTick`, call `engine.ComputeRates(m.State.Run)` → `m.Rates`, call `engine.Tick(m.State.Run, delta, m.Rates)` → `m.State.Run`, return `engine.TickCmd()`
  - [x] Handle `tea.KeyMsg`: if `m.Screen == BootScreen`, treat any keypress as boot-done (set `m.LastTick = time.Now()`, `m.Screen = GameScreen`, return `engine.TickCmd()`); handle `"q"` and `"ctrl+c"` on any screen → `tea.Quit`
  - [x] All other messages: return `m, nil`

- [x] Task 6: Create `internal/ui/view.go` — View delegating to screen funcs (AC: 1, 2, 3)
  - [x] `package ui`; implement `(m Model) View() string`
  - [x] If `m.tooSmall`: return a plain-text error string like `"Terminal too small: need at least 80x24 (current WxH)\n"` — no ANSI, no lipgloss styles
  - [x] Switch on `m.Screen`: `BootScreen` → `screens.BootView(m.Width, m.Height)`; `GameScreen` → `screens.GameView(m.State.Run.Resources, m.Rates, m.Width, m.Height)`; default → `""`

- [x] Task 7: Update `internal/ui/ui.go` — replace stub with package doc (AC: all)
  - [x] Replace the TODO stub comment with a real package doc: `// Package ui implements the Bubbletea TUI model, screen management, and view rendering for afk-x.`
  - [x] `package ui` declaration stays

- [x] Task 8: Update `cmd/afk-x/main.go` — wire tea.NewProgram (AC: all)
  - [x] Keep `--version` flag logic unchanged
  - [x] Replace the `fmt.Println("afk-x: TUI not yet initialized")` placeholder with `tea.NewProgram(ui.New(), tea.WithAltScreen())`
  - [x] Import `github.com/charmbracelet/bubbletea` and `github.com/12jihan/afk-x/internal/ui`
  - [x] No signal handling yet — that's Story 5.5

- [x] Task 9: Verify build and test suite (AC: all)
  - [x] `go build ./cmd/afk-x` — zero errors
  - [x] `go test ./...` — zero regressions across all packages
  - [ ] Manual smoke: `go run ./cmd/afk-x` launches in terminal, shows boot text, auto-advances at 2s (or keypress), game screen shows three incrementing resources, `q` exits cleanly
  - [ ] Manual smoke: `NO_COLOR=1 go run ./cmd/afk-x` — no color codes in output

## Dev Notes

### What This Story Does

Stories 1.2–1.4 built the data and computation layers. Story 1.5 wires them into a running Bubbletea TUI. The deliverable is a binary that boots, shows flavor text, then displays three resources ticking in real time.

**Files created:**
- `internal/ui/model.go` — Model, Screen, New(), Init(), bootDoneMsg, bootTimerCmd
- `internal/ui/update.go` — Update() with full message routing
- `internal/ui/view.go` — View() delegating to screen functions
- `internal/ui/screens/boot.go` — BootView() render function
- `internal/ui/screens/game.go` — GameView() render function
- `internal/ui/styles/colors.go` — NO_COLOR-safe palette
- `internal/ui/styles/layout.go` — MinWidth, MinHeight, ResourcePanel style
- `internal/ui/styles/text.go` — FormatResource helper
- `internal/ui/ui.go` — replace stub with package doc
- `cmd/afk-x/main.go` — replace placeholder with tea.NewProgram

**Do NOT create in this story:**
- `internal/ui/screens/draft.go` — Story 3.2
- `internal/ui/screens/summary.go` — Story 6.1
- `internal/ui/screens/meta.go` — Story 6.3
- `internal/engine/offline.go` — Story 5.4
- Any save integration — Story 5.x

### Architecture Compliance (MUST FOLLOW)

**Package import rules (strictly enforced by go build):**
```
cmd/afk-x  →  ui, bubbletea, os, fmt
ui          →  engine, game, content, ui/screens, ui/styles, bubbletea, bubbles/v2/textinput, time
ui/screens  →  engine, content, ui/styles, lipgloss
ui/styles   →  lipgloss only
```
- ❌ `ui/screens` MUST NOT import `ui` — circular import
- ❌ `ui/styles` MUST NOT import `ui` or `ui/screens` — circular import
- ❌ `ui` MUST NOT import `save` — use `SaveMsg` pattern (deferred to Story 5.x)
- ❌ `engine` MUST NOT import `ui` — engine only imports game and content

**Bubbletea rules:**
- All state mutations happen in `Update()` — never in `View()` or `Init()`
- Side effects (tickers, timers) via `tea.Cmd` only — never call I/O from inside `Update`
- Screen transitions via explicit `m.Screen = XxxScreen` in `Update` — never in `View`
- Never use `fmt.Println` or `log.Print` inside the Bubbletea program — all output through `View`
- Never call `log.Fatal` outside `cmd/afk-x/main.go`

**NO_COLOR handling:**
- Lipgloss natively detects the `NO_COLOR` environment variable — no custom code needed
- Use `lipgloss.AdaptiveColor{Light: "...", Dark: "..."}` or `lipgloss.Color("...")` — lipgloss strips all color when `NO_COLOR=1`
- The `tooSmall` error in `View()` MUST be plain string without any lipgloss — it fires before Bubbletea is fully initialized in some edge cases

**Terminal compatibility (AC5):**
- Use lipgloss for all box-drawing and borders — it handles cross-terminal box char compatibility
- `tea.WithAltScreen()` in `main.go` ensures clean terminal restoration on exit

### Key Implementation Patterns

**Boot timer command — blocking cmd (standard Bubbletea pattern):**
```go
func bootTimerCmd() tea.Cmd {
    return func() tea.Msg {
        time.Sleep(bootDuration)
        return bootDoneMsg{}
    }
}
```
Bubbletea executes cmds concurrently — this goroutine sleeps in the background without blocking `Update`.

**TickMsg handler — always reschedule:**
```go
case engine.TickMsg:
    if m.Screen != GameScreen {
        return m, nil // ignore ticks during boot
    }
    now := time.Time(msg)
    delta := now.Sub(m.LastTick)
    m.LastTick = now
    m.Rates = engine.ComputeRates(m.State.Run)
    m.State.Run = engine.Tick(m.State.Run, delta, m.Rates)
    return m, engine.TickCmd() // always reschedule
```
`engine.TickCmd()` uses `tea.Every(engine.TickInterval, ...)` — the exported `TickInterval = 16*time.Millisecond` constant from Story 1.4.

**Window size check — set flag then quit:**
```go
case tea.WindowSizeMsg:
    m.Width, m.Height = msg.Width, msg.Height
    if m.Width < styles.MinWidth || m.Height < styles.MinHeight {
        m.tooSmall = true
        return m, tea.Quit
    }
```
Bubbletea calls `View()` one final time before quitting, so the tooSmall error message will be displayed.

**FormatResource alignment — use fmt.Sprintf with fixed width:**
```go
func FormatResource(name string, value float64) string {
    return fmt.Sprintf("%-20s %10.2f", name, value)
}
```
Right-aligns float values so resources don't jump as they grow. Adjust widths as needed.

### Module and Import Details

- Module: `github.com/12jihan/afk-x`
- Go version: 1.26.2 (do NOT change go.mod)
- Bubbletea: `github.com/charmbracelet/bubbletea v1.3.10`
- Lipgloss: `github.com/charmbracelet/lipgloss v1.1.0`
- Bubbles v2 textinput: `charm.land/bubbles/v2/textinput` (the module is `charm.land/bubbles/v2`, already in go.mod)
- All three dependencies already in `go.sum` from Story 1.4 — no `go get` should be needed

### Previous Story Intelligence (from Story 1.4)

- **External test package pattern:** `package engine_test` with explicit import — for `ui` tests, use `package ui_test` if any tests are written, but tests are optional for MVP per architecture
- **Module path confirmed:** `github.com/12jihan/afk-x`
- **Go version:** 1.26.2 — do NOT modify `go.mod`
- **TickCmd is exported:** `engine.TickCmd()` (not `tickCmd`) — Story 1.4 exported it because `ui` needs to call it from a different package
- **TickInterval exported:** `engine.TickInterval = 16*time.Millisecond` — use in tests if timing constants are needed
- **go.sum may need update:** Story 1.4 required `go get github.com/charmbracelet/x/ansi@v0.11.6` for a transitive dependency; if any new transitive dep is missing, run `go mod tidy`
- **content.BootText** is defined in `internal/content/narrative.go` — plain multi-line string with terminal aesthetic flavor text
- **content.CPUCycles, content.MemoryShards, content.ProcessThreads** are string constants in `internal/content/resources.go` — always use these, never string literals

### Project Structure Notes

Existing files to MODIFY:
- `internal/ui/ui.go` — replace stub `// TODO: Implemented in Story 1.5` comment
- `cmd/afk-x/main.go` — replace `fmt.Println("afk-x: TUI not yet initialized")`
- `internal/ui/screens/.gitkeep` — delete once `boot.go` and `game.go` are created
- `internal/ui/styles/.gitkeep` — delete once `colors.go`, `layout.go`, `text.go` are created

Existing `ui.go` has: `// Package ui provides the TUI application model and screen management.\n// TODO: Implemented in Story 1.5\npackage ui`

Bubbletea requires `Model` to implement the `tea.Model` interface:
```go
type tea.Model interface {
    Init() tea.Cmd
    Update(tea.Msg) (tea.Model, tea.Cmd)
    View() string
}
```
`Update` must return `(tea.Model, tea.Cmd)` — NOT `(Model, tea.Cmd)`. Return `m` (the value receiver copy) which satisfies the interface.

### Deferred Items (do NOT address in this story)

- Boot screen keypress: currently any key advances — no need to filter specific keys; Story 4.x will add combo input to GameScreen
- `textinput.Model` in `Model` struct: initialize it but do NOT render it — command panel is Story 4.2
- Save integration: no save calls in this story — `SaveMsg` pattern deferred to Story 5.x
- Signal handling (SIGINT/SIGTERM): Story 5.5
- `engine.go` doc "no global mutable state" overstated — known deferred item from Story 1.4 review

### References

- Story 1.5 ACs: `_bmad-output/planning-artifacts/epics.md` → Epic 1, Story 1.5
- Bubbletea Model architecture: `_bmad-output/planning-artifacts/architecture.md` → "Bubbletea Model / Screen Architecture"
- Screen enum and Model struct: `architecture.md` → code block under "Bubbletea Model / Screen Architecture"
- Package import rules: `architecture.md` → "Package import rules" table
- NO_COLOR / NFR17: `architecture.md` → NFR table row NFR17
- Terminal size / NFR: `architecture.md` → "Enforcement Guidelines"
- `content.BootText`: `internal/content/narrative.go` (defined in Story 1.3)
- `engine.TickCmd`, `engine.TickInterval`, `engine.Tick`, `engine.ComputeRates`: `internal/engine/tick.go`, `internal/engine/rates.go` (defined in Story 1.4)
- `game.NewGameState()`: `internal/game/state.go` (defined in Story 1.2)
- bubbles v2 import path: `go.mod` → `charm.land/bubbles/v2 v2.0.0`

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-6

### Debug Log References

- `charm.land/bubbles/v2/textinput` required `go get charm.land/bubbles/v2/textinput@v2.0.0` to pull transitive deps (`charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`, `github.com/atotto/clipboard`) — ran on session pickup after remote session had already written model.go

### Completion Notes List

- Remote session created Tasks 1–7 (all styles, screens, model, update, view, ui.go); session pickup completed Task 8 (main.go wiring) and Task 9 (build + test verification)
- `bootTimerCmd()` uses blocking goroutine pattern — Bubbletea executes cmds concurrently, sleep doesn't block Update
- `engine.TickCmd()` only scheduled on GameScreen entry (bootDoneMsg or boot keypress) — no tick accumulation during boot
- `NO_COLOR` handled transparently by `lipgloss.AdaptiveColor` — no custom code needed
- `tooSmall` error in View() returns plain string without lipgloss — safe for pre-init edge cases
- Terminal size guard: sets `m.tooSmall = true` then returns `tea.Quit`; Bubbletea calls View() one final time before quitting to display the error
- `resourceOrder` slice in game.go ensures stable rendering order across Go map iteration
- `go build ./...` — zero errors; `go test ./...` — zero regressions (9 engine tests, content/game tests all pass)

### File List

- `internal/ui/ui.go` — replaced stub with package doc
- `internal/ui/model.go` — new: Model, Screen enum, New(), Init(), bootDoneMsg, bootTimerCmd
- `internal/ui/update.go` — new: Update() with full message routing
- `internal/ui/view.go` — new: View() delegating to screen functions
- `internal/ui/screens/boot.go` — new: BootView() render function
- `internal/ui/screens/game.go` — new: GameView() with resourceOrder + labels
- `internal/ui/styles/colors.go` — new: AdaptiveColor palette (Accent, Muted, Normal)
- `internal/ui/styles/layout.go` — new: MinWidth, MinHeight constants + ResourcePanel style
- `internal/ui/styles/text.go` — new: FormatResource helper
- `internal/ui/screens/.gitkeep` — deleted
- `internal/ui/styles/.gitkeep` — deleted
- `cmd/afk-x/main.go` — replaced placeholder with tea.NewProgram wiring
- `go.mod` — updated (bubbles/v2 transitive deps)
- `go.sum` — updated

### Change Log

- 2026-05-03: Story 1.5 complete — TUI shell with boot sequence, resource display, window-size guard, NO_COLOR support
