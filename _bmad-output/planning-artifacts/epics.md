---
stepsCompleted: ['step-01-validate-prerequisites', 'step-02-design-epics', 'step-03-create-stories', 'step-04-final-validation']
status: 'complete'
completedAt: '2026-05-02'
inputDocuments: ['_bmad-output/planning-artifacts/prd.md', '_bmad-output/planning-artifacts/architecture.md']
---

# afk-x - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for afk-x, decomposing the requirements from the PRD and Architecture into implementable stories.

## Requirements Inventory

### Functional Requirements

FR1: Player can observe resources generating automatically in real time
FR2: Player can accumulate resources while the game is not running (AFK/offline accumulation)
FR3: System calculates and applies offline resource gains on resume
FR4: Player can observe multiple distinct resource types simultaneously (CPU cycles, memory shards, process threads)
FR5: New resource types unlock as the player advances across runs
FR6: Player can purchase upgrades using accumulated resources
FR7: Player can view available upgrades and their costs at any time
FR8: Player can progress through numbered tower floors by meeting floor completion conditions
FR9: Player encounters daemon enemies on each floor that must be overcome to progress
FR10: Player can observe current floor number and progress toward the next floor clear
FR11: Player is presented with a draft of 3 randomly selected perks upon clearing a floor
FR12: Player can select one perk from the draft to apply to the current run
FR13: Each drafted perk has a distinct effect that influences run strategy
FR14: The available perk pool expands as meta-progression unlocks are earned across runs
FR15: Higher-tier perks become accessible in the draft pool as the player advances
FR16: Player can enter command sequences via an in-game text input panel
FR17: System validates entered command sequences against known combos
FR18: Player can execute valid command combos to receive efficiency bonuses in the current zone
FR19: Player can view a reference panel listing available commands and valid combos
FR20: System rejects invalid command input without altering game state
FR21: Player can queue command combos to execute while away from the game
FR22: Player receives a run summary on run end (floors cleared, resources generated, perks drafted, time elapsed)
FR23: Permanent unlocks earned during a run carry over to all subsequent runs
FR24: Player can start a new run immediately after run end
FR25: Run 2+ begins with visibly faster early progression reflecting meta-progression unlocks
FR26: Player's best command combos from prior runs are available as presets in subsequent runs
FR27: Meta-progression unlock screen displays newly earned unlocks before the new run starts
FR28: Game state is automatically persisted to a portable JSON file
FR29: Game state is fully restored on relaunch from the save file
FR30: Player can specify a custom save file location at launch
FR31: Player can wipe save data and start fresh via a launch command, with a confirmation step
FR32: Pending perk drafts persist across game close and reopen until the player makes a selection
FR33: Staged command queues persist while the game is not running
FR34: Player is presented with a boot sequence on first launch that includes narrative flavor text
FR35: Player receives flavor text at designated milestone floors revealing tower lore
FR36: All game elements use terminal/computing terminology and aesthetic
FR37: Player can read the current floor's narrative context within the game UI
FR38: Player can launch the game with no arguments to start or resume a session
FR39: Player can query the installed application version from the command line
FR40: Player can disable ANSI color output via the NO_COLOR environment variable
FR41: System detects terminal dimensions on launch and displays an error if too small to render
FR42: System saves game state on receiving SIGINT or SIGTERM before exiting
FR43: Save files are stored in an XDG-compliant directory by default (~/.local/share/afk-x/)
FR44: Game renders correctly across 5+ mainstream terminal emulators

### NonFunctional Requirements

NFR1: Resource tick calculations complete within 16ms per cycle (≤60fps render budget)
NFR2: AFK accumulation calculation on resume completes within 500ms regardless of offline duration
NFR3: Game launch (binary start to interactive TUI) completes within 2 seconds on a standard developer machine
NFR4: Perk draft screen renders within 100ms of floor clear trigger
NFR5: Save file write completes within 200ms
NFR6: Save file is written atomically — a crash during write must not corrupt the existing save
NFR7: Game state on resume matches game state at last save — no drift or loss from clean or unclean exits
NFR8: AFK accumulation is deterministic — identical offline duration and game state always produce identical resource gain
NFR9: Invalid or corrupted save files produce a clear error message and reset prompt, not a silent crash
NFR10: Game handles SIGINT and SIGTERM without data loss under normal operating conditions
NFR11: Distributed as a single self-contained binary — no runtime dependencies beyond a POSIX-compatible terminal
NFR12: Binary supports Linux (amd64, arm64) and macOS (amd64, arm64) at minimum
NFR13: Game renders correctly in iTerm2, Terminal.app, GNOME Terminal, Windows Terminal, and tmux
NFR14: Save files are human-readable JSON, loadable on any supported platform without modification
NFR15: Binary size ≤ 20MB
NFR16: No game-critical information is encoded in color alone
NFR17: NO_COLOR fully disables ANSI color output without affecting gameplay or readability
NFR18: All game interactions are operable via keyboard only — no mouse dependency
NFR19: Text elements remain legible at standard terminal font sizes

### Additional Requirements

From Architecture document:

- **Project initialization:** `go mod init`, add bubbletea v1.3.10, bubbles v2.0.0, lipgloss v1.1.0
- **Package structure:** 7 packages — `game`, `content`, `engine`, `combat`, `save`, `ui`, `cmd/afk-x`
- **Package dependency rules:** enforced — `engine`/`combat` must be pure functions with no I/O; `save` is the only file I/O package; `ui` is the only terminal I/O package
- **Save schema Version field:** must be present from Story 1 to enable future migration
- **ComboDefinition struct:** `{ID, Sequence []string, BonusType string, BonusValue float64, FlavorText string}` — must be agreed on between `content` and `combat` packages
- **AFK pattern:** always use `engine.ApplyOffline(state, time.Now())` wrapper — never inline the delta calc
- **Atomic save:** always `os.Rename` after temp write — never `os.WriteFile` directly
- **Boot sequence:** `bootDuration = 2 * time.Second` constant in `ui/screens/boot.go`, skippable with any key
- **CI/CD:** GitHub Actions for `go test ./...` + `go vet ./...` on push; GoReleaser on tag for cross-platform binaries
- **Makefile:** `build`, `test`, `lint`, `release` targets
- **GoReleaser:** `-ldflags "-s -w"` for binary size, Linux/macOS amd64+arm64 matrix

### UX Design Requirements

No UX design document exists. Screen layout decisions are deferred to implementation, with the following minimum constraints from the PRD:
- Resource display panel (FR1, FR4)
- Upgrade panel (FR7)
- Floor progress indicator (FR10)
- Command input panel + reference panel (FR16, FR19, `?` toggle)
- Perk draft screen — 3 options, pick 1 (FR11–FR12)
- Boot sequence screen (FR34)
- Run summary screen (FR22)
- Meta-progression unlock screen (FR27)

### FR Coverage Map

FR1: Epic 1 — Resources tick in real time (game shell)
FR2: Epic 5 — AFK accumulation
FR3: Epic 5 — Apply offline gains on resume
FR4: Epic 1 — Multiple resource types visible
FR5: Epic 6 — New resource types unlock across runs
FR6: Epic 2 — Purchase upgrades
FR7: Epic 2 — View upgrades and costs
FR8: Epic 2 — Progress through floors
FR9: Epic 2 — Daemon enemies per floor
FR10: Epic 2 — Floor number and progress indicator
FR11: Epic 3 — 3 random perks on floor clear
FR12: Epic 3 — Select one perk
FR13: Epic 3 — Perk effects influence run strategy
FR14: Epic 3 — Perk pool expansion (stubbed; full expansion in Epic 6)
FR15: Epic 3 — Higher-tier perks accessible (stubbed; unlocked in Epic 6)
FR16: Epic 4 — Command input panel
FR17: Epic 4 — Combo validation
FR18: Epic 4 — Efficiency bonuses on valid combos
FR19: Epic 4 — Reference panel (?  toggle)
FR20: Epic 4 — Invalid input rejection
FR21: Epic 4 — Queue combos while AFK
FR22: Epic 6 — Run summary screen
FR23: Epic 6 — Permanent unlocks carry over
FR24: Epic 6 — Start new run immediately
FR25: Epic 6 — Run 2+ visibly faster
FR26: Epic 6 — Best combos as presets in run 2+
FR27: Epic 6 — Meta-progression unlock screen
FR28: Epic 5 — Persist game state to JSON
FR29: Epic 5 — Restore state on relaunch
FR30: Epic 5 — Custom save path via --save flag
FR31: Epic 5 — Reset via --reset with confirmation
FR32: Epic 5 — Persist pending perk draft
FR33: Epic 5 — Persist staged combo queue
FR34: Epic 1 — Boot sequence with narrative
FR35: Epic 2 — Milestone flavor text on floor clears
FR36: Epic 1 — Terminal/computing aesthetic throughout
FR37: Epic 1 — Current floor narrative context
FR38: Epic 1 — Launch with no arguments
FR39: Epic 1 — --version flag
FR40: Epic 1 — NO_COLOR env var support
FR41: Epic 1 — Terminal size check on launch
FR42: Epic 5 — SIGINT/SIGTERM save before exit
FR43: Epic 5 — XDG-compliant save path
FR44: Epic 1 — Cross-terminal compatibility

## Epic List

### Epic 1: Runnable Game Shell
Player can launch the game, see the boot sequence, and watch resources tick in a styled terminal UI. The project scaffold, CI, and distribution pipeline are also established.
**FRs covered:** FR1, FR4, FR34, FR36, FR37, FR38, FR39, FR40, FR41, FR44
**Also covers:** Project scaffold (go mod init + deps), GitHub Actions CI, Makefile, GoReleaser config, `game` + `content` + `ui` base packages

### Epic 2: Core Idle Loop
Player can purchase upgrades, fight through daemon floors, and see milestone flavor text — the idle loop is playable end-to-end.
**FRs covered:** FR6, FR7, FR8, FR9, FR10, FR35

### Epic 3: Perk Draft System
Player receives 3 random perk choices when clearing a floor and selects one that shapes their run strategy.
**FRs covered:** FR11, FR12, FR13, FR14 (stubbed for Epic 6 expansion), FR15 (stubbed for Epic 6 expansion)

### Epic 4: Command & Combo System
Player can type command sequences in an input panel to execute efficiency combos, queue them while AFK, and reference valid combos.
**FRs covered:** FR16, FR17, FR18, FR19, FR20, FR21

### Epic 5: Save, Persistence & AFK
The game saves and restores automatically; resources accumulate while away; all transient state (drafts, queues) survives process death.
**FRs covered:** FR2, FR3, FR28, FR29, FR30, FR31, FR32, FR33, FR42, FR43

### Epic 6: Meta-Progression & Run Loop
Player completes a full run, earns permanent unlocks, and starts run 2 with visible acceleration — the complete product loop.
**FRs covered:** FR5, FR22, FR23, FR24, FR25, FR26, FR27

---

## Epic 1: Runnable Game Shell

Player can launch the game, see the boot sequence, and watch resources tick in a styled terminal UI. The project scaffold, CI, and distribution pipeline are established.
**FRs covered:** FR1, FR4, FR34, FR36, FR37, FR38, FR39, FR40, FR41, FR44
**Also covers:** Project scaffold (go mod init + deps), GitHub Actions CI, Makefile, GoReleaser config, `game` + `content` + `ui` base packages

### Story 1.1: Project Scaffold & Build Pipeline

As a developer,
I want the Go project initialized with all dependencies, build tooling, and CI pipeline,
So that all subsequent development has a clean, reproducible foundation.

**Acceptance Criteria:**

**Given** the repository is cloned
**When** `make build` is run
**Then** a binary `bin/afk-x` is produced that exits cleanly with code 0

**Given** a commit is pushed to main
**When** GitHub Actions CI runs
**Then** `go test ./...` and `go vet ./...` both pass

**Given** the binary exists
**When** run with `--version`
**Then** it prints a version string and exits with code 0

**Given** `.goreleaser.yml` exists
**When** `goreleaser check` is run
**Then** no configuration errors are reported

---

### Story 1.2: Game State Schema

As a developer,
I want the core game state structs defined with JSON serialization,
So that all packages share a consistent, versioned data contract.

**Acceptance Criteria:**

**Given** a `GameState` is created
**When** marshaled to JSON
**Then** all fields use snake_case keys and `version` is set to 1

**Given** a `GameState` is marshaled and then unmarshaled
**When** the result is compared to the original
**Then** all fields are identical (round-trip test passes)

**Given** `RunState.Resources`
**When** populated with resource values
**Then** it uses `map[string]float64` with resource key constants from the `content` package

**Given** `RunState.PendingDraft`
**When** no draft is pending
**Then** it serializes as `null`, not `[]`

---

### Story 1.3: Content Package — Resource & Narrative Definitions

As a developer,
I want resource key constants, initial rates, and narrative text defined in the content package,
So that all other packages reference game content without magic strings.

**Acceptance Criteria:**

**Given** `content/resources.go` is imported
**When** the constants are accessed
**Then** `CPUCycles`, `MemoryShards`, and `ProcessThreads` key constants are defined and non-empty

**Given** `content/narrative.go`
**When** `BootText` is accessed
**Then** it contains the tower lore flavor text (non-empty string)

**Given** `content/narrative.go`
**When** `MilestoneText(floor int)` is called with a milestone floor number
**Then** it returns the corresponding flavor text string

**Given** `go test ./internal/content/...`
**When** run
**Then** all tests pass with no compilation errors

---

### Story 1.4: Resource Tick Engine

As a player,
I want resources to generate automatically over time,
So that the game progresses without requiring my attention.

**Acceptance Criteria:**

**Given** a `RunState` and `ResourceRates`
**When** `engine.Tick()` is called with a 16ms delta
**Then** each resource value increases by `rate * 0.016`

**Given** identical `RunState` and delta inputs
**When** `engine.Tick()` is called twice
**Then** both calls produce identical output (determinism confirmed)

**Given** `engine.ComputeRates()` called with a default `RunState`
**When** no upgrades or perks are active
**Then** it returns the base rates defined in `content/resources.go`

**Given** `engine.tickCmd()` is scheduled via `tea.Every(16ms, ...)`
**When** the Bubbletea program is running
**Then** a `TickMsg` is delivered approximately every 16 milliseconds

---

### Story 1.5: TUI Shell — Boot Sequence & Resource Display

As a player,
I want to launch the game, see a boot sequence, and watch my resources tick,
So that I know the game is running and my progress is accumulating.

**Acceptance Criteria:**

**Given** the binary is run with no arguments
**When** launched in a supported terminal
**Then** a boot sequence with narrative flavor text displays, auto-advancing after 2 seconds or on any keypress

**Given** the boot sequence completes
**When** the game screen appears
**Then** CPU cycles, memory shards, and process threads are visible and incrementing in real time

**Given** a terminal that is too small to render the UI
**When** the game launches
**Then** an error message is displayed and the game exits cleanly without crashing

**Given** `NO_COLOR=1` is set in the environment
**When** the game runs
**Then** all output is plain text with no ANSI color sequences

**Given** the game is running in iTerm2, Terminal.app, GNOME Terminal, Windows Terminal, and tmux
**When** the layout is displayed
**Then** panels render without visual corruption or broken box-drawing characters

---

## Epic 2: Core Idle Loop

Player can purchase upgrades, fight through daemon floors, and see milestone flavor text — the idle loop is playable end-to-end.
**FRs covered:** FR6, FR7, FR8, FR9, FR10, FR35

### Story 2.1: Upgrade Definitions & Purchase System

As a player,
I want to view available upgrades and spend my resources to purchase them,
So that I can improve my resource generation and progress faster.

**Acceptance Criteria:**

**Given** the game screen is showing
**When** the player views the upgrade panel
**Then** available upgrades are listed with their names, costs, and current level

**Given** a player has sufficient resources for an upgrade
**When** they select and confirm the upgrade
**Then** the resources are deducted and the upgrade level increments

**Given** a player has insufficient resources for an upgrade
**When** they attempt to purchase it
**Then** the purchase is rejected with visual feedback and no resources are deducted

**Given** an upgrade is purchased
**When** `engine.ComputeRates()` is called
**Then** the returned rates reflect the upgrade's effect on resource generation

---

### Story 2.2: Floor Progression & Daemon Encounters

As a player,
I want to advance through tower floors by defeating daemon enemies,
So that I can climb the tower and unlock new content.

**Acceptance Criteria:**

**Given** the player's resources and upgrades meet a floor's completion threshold
**When** `engine.CheckFloorClear()` is evaluated on a tick
**Then** a `FloorClearMsg` is emitted

**Given** a `FloorClearMsg` is received
**When** the floor advances
**Then** the floor number increments and the new floor's daemon enemies are displayed

**Given** the game screen is showing
**When** the current floor has active daemon enemies
**Then** enemy names, types, and progress indicators are visible in the UI

**Given** the player is on any floor
**When** the game screen is viewed
**Then** the current floor number and progress toward floor clear are clearly displayed

---

### Story 2.3: Milestone Flavor Text

As a player,
I want to receive story flavor text when I reach milestone floors,
So that I learn more about the tower and feel rewarded for progressing.

**Acceptance Criteria:**

**Given** the player clears a designated milestone floor (e.g., floor 5, 10, 20)
**When** the floor clear is processed
**Then** the milestone flavor text from `content/narrative.go` is displayed prominently

**Given** a milestone message is displayed
**When** the player continues playing
**Then** the message auto-dismisses after a short delay or on keypress

**Given** a non-milestone floor is cleared
**When** the floor advances
**Then** no milestone message is shown

---

## Epic 3: Perk Draft System

Player receives 3 random perk choices when clearing a floor and selects one that shapes their run strategy.
**FRs covered:** FR11, FR12, FR13, FR14 (stubbed for Epic 6 expansion), FR15 (stubbed for Epic 6 expansion)

### Story 3.1: Perk Definitions & Draft Pool

As a player,
I want a diverse pool of perks available for drafting,
So that each run offers meaningful strategic variety.

**Acceptance Criteria:**

**Given** `content/perks.go`
**When** the perk pool is accessed
**Then** at least 12 distinct `PerkDefinition` entries exist, each with a unique ID, name, effect description, and `BonusType`/`BonusValue`

**Given** `engine.DrawPerkDraft()` is called with a run seed
**When** called with the same seed and state
**Then** it returns the same 3 perks (deterministic draft)

**Given** `engine.DrawPerkDraft()` is called with different seeds
**When** called
**Then** it returns different perk selections (random variety confirmed)

**Given** `MetaState.UnlockedPerks` is empty (run 1)
**When** `DrawPerkDraft()` is called
**Then** it draws from the base pool of 12 perks only

---

### Story 3.2: Perk Draft Screen & Selection

As a player,
I want to see 3 perk options when I clear a floor and choose one to keep,
So that I can shape my build and make each run feel distinct.

**Acceptance Criteria:**

**Given** a floor is cleared
**When** `FloorClearMsg` is received
**Then** the screen transitions to `PerkDraftScreen` showing 3 perk options with names and descriptions

**Given** the perk draft screen is showing
**When** the player presses `1`, `2`, or `3`
**Then** the corresponding perk is added to `RunState.ActivePerks` and the screen returns to `GameScreen`

**Given** the perk draft screen is showing
**When** the player views the options
**Then** each perk card shows its name, effect description, and bonus value

**Given** a perk draft is pending and the game is closed and reopened (after Epic 5)
**When** the game loads
**Then** the perk draft screen is shown with the same 3 options (FR32)

---

### Story 3.3: Active Perk Effects

As a player,
I want my selected perks to meaningfully affect my resource generation,
So that my draft choices have real strategic impact on each run.

**Acceptance Criteria:**

**Given** a player has selected a resource-multiplying perk (e.g., "Overclock" for CPU cycles)
**When** `engine.ComputeRates()` is called
**Then** the affected resource rate is higher than it was before the perk was selected

**Given** a player has no active perks vs. one active perk
**When** comparing `ComputeRates()` output
**Then** the rates differ in the direction described by the perk's `BonusType` and `BonusValue`

**Given** `engine.ApplyPerk()` is called with a valid perk ID and `RunState`
**When** the new `RunState` is inspected
**Then** the perk ID appears in `ActivePerks` and no other fields are modified

---

## Epic 4: Command & Combo System

Player can type command sequences for efficiency bonuses, queue them while AFK, and reference valid combos.
**FRs covered:** FR16, FR17, FR18, FR19, FR20, FR21

### Story 4.1: Combo Definitions, Parser & Validator

As a player,
I want the game to recognize command sequences I type,
So that I can execute combos for efficiency bonuses.

**Acceptance Criteria:**

**Given** `content/combos.go`
**When** accessed
**Then** at least 5 `ComboDefinition` entries exist with unique IDs, sequences, bonus types, and bonus values

**Given** `combat.ParseInput("scan exploit loot")`
**When** called
**Then** it returns `["scan", "exploit", "loot"]`

**Given** `combat.ValidateCombo(["scan", "exploit", "loot"])` where this is a valid combo
**When** called
**Then** it returns the matching `ComboDefinition`

**Given** `combat.ValidateCombo(["scan", "invalid", "loot"])` where no such combo exists
**When** called
**Then** it returns `nil, false` without modifying any state

---

### Story 4.2: Command Input Panel & Reference Panel

As a player,
I want a text input panel to enter commands and a reference panel showing valid combos,
So that I can discover and use combos efficiently.

**Acceptance Criteria:**

**Given** the game screen is showing
**When** the player activates the command input (e.g., presses `/`)
**Then** a text input field is active and accepts keyboard input

**Given** the game screen is showing
**When** the player presses `?`
**Then** a reference panel toggles showing available command tokens and valid combo sequences

**Given** the reference panel is visible
**When** the player presses `?` again
**Then** the reference panel is hidden

**Given** an invalid command is entered and submitted
**When** the input is processed
**Then** a brief inline error message appears and the game state is not modified (FR20)

---

### Story 4.3: Combo Execution & Efficiency Bonus

As a player,
I want valid command combos to grant efficiency bonuses to my current zone,
So that active play is rewarded with faster floor progression.

**Acceptance Criteria:**

**Given** the player types a valid combo and submits it
**When** `combat.ComputeBonus()` is called
**Then** the floor progress rate is increased by the combo's `BonusValue` for the current zone

**Given** a combo is executed
**When** the bonus is applied
**Then** the combo's `FlavorText` is briefly displayed in the game screen status area

**Given** an invalid combo sequence is submitted
**When** the input is processed
**Then** no bonus is applied and inline rejection feedback is shown

---

### Story 4.4: Combo Queue (AFK Execution)

As a player,
I want to queue a command combo to execute automatically while I'm away,
So that I can set up active bonuses before closing the terminal.

**Acceptance Criteria:**

**Given** the player enters a valid combo and queues it
**When** confirmed
**Then** it is added to `RunState.ComboQueue` and a confirmation message is shown

**Given** combos are in `RunState.ComboQueue`
**When** the game loop ticks
**Then** the first queued combo is dequeued, its bonus applied, and removed from the queue

**Given** the game is closed with combos in the queue and then reopened (after Epic 5)
**When** the game loads
**Then** the queued combos are still present and begin executing on resume (FR33)

---

## Epic 5: Save, Persistence & AFK

The game saves and restores automatically; resources accumulate while away; all transient state survives process death.
**FRs covered:** FR2, FR3, FR28, FR29, FR30, FR31, FR32, FR33, FR42, FR43

### Story 5.1: Save Package — Atomic Write & XDG Paths

As a player,
I want the game to save my progress automatically to a portable file,
So that my game state is never lost even if the process crashes.

**Acceptance Criteria:**

**Given** `save.Write(state, path)` is called
**When** the write succeeds
**Then** a valid JSON file exists at the given path containing the full `GameState`

**Given** a crash occurs during the write
**When** the save file is inspected
**Then** the existing save file is intact (atomic write via temp → rename)

**Given** no custom path is provided
**When** `save.XDGSavePath()` is called
**Then** it returns `~/.local/share/afk-x/save.json` (or `$XDG_DATA_HOME/afk-x/save.json` if set)

**Given** `go test ./internal/save/...` (writing to `os.TempDir()`)
**When** run
**Then** all tests pass

---

### Story 5.2: Save Package — Load & Corruption Recovery

As a player,
I want the game to restore my saved state on launch,
So that I can continue exactly where I left off.

**Acceptance Criteria:**

**Given** a valid save file exists
**When** `save.Load(path)` is called
**Then** it returns the `GameState` with `isNew = false` and `err = nil`

**Given** no save file exists
**When** `save.Load(path)` is called
**Then** it returns a default `GameState` with `isNew = true` and `err = nil`

**Given** a corrupted save file (invalid JSON) exists
**When** `save.Load(path)` is called
**Then** it returns `err != nil` and the UI displays a clear error with a reset prompt

**Given** a `GameState` is written and reloaded
**When** compared to the original
**Then** all fields are identical (no drift — NFR7)

---

### Story 5.3: Launch Flags — `--save` & `--reset`

As a player,
I want to specify a custom save file location and reset my game from the command line,
So that my save file is portable and I can start fresh when I want to.

**Acceptance Criteria:**

**Given** the game is launched with `--save ~/my-save.json`
**When** the game saves
**Then** the save file is written to `~/my-save.json` instead of the default XDG path

**Given** the game is launched with `--reset`
**When** the terminal shows the confirmation prompt
**Then** the game waits for `y` or `n` input

**Given** the player confirms reset with `y`
**When** the reset completes
**Then** the save file is deleted and the game starts fresh

**Given** the player cancels reset with `n`
**When** the cancellation is processed
**Then** the save file is unchanged and the game exits cleanly

---

### Story 5.4: AFK Accumulation Engine

As a player,
I want resources to have accumulated when I return after being away,
So that the game rewards time spent away from the terminal.

**Acceptance Criteria:**

**Given** a save file with `saved_at` 4 hours in the past
**When** `engine.ApplyOffline()` is called on load
**Then** resources are increased by `rate * 14400` seconds

**Given** the same `saved_at` and game state
**When** `engine.ApplyOffline()` is called twice
**Then** both calls produce identical results (determinism — NFR8)

**Given** an offline duration exceeding `maxOfflineDuration` (e.g., 30 days)
**When** `engine.ApplyOffline()` is called
**Then** the delta is capped at `maxOfflineDuration` to prevent overflow

**Given** `engine.ApplyOffline()` completes
**When** the updated state is inspected
**Then** `SavedAt` is set to `time.Now()`

---

### Story 5.5: Signal Handling & Persistent Transient State

As a player,
I want my pending perk draft and queued combos to survive closing the terminal,
So that I never lose my progress or queued actions.

**Acceptance Criteria:**

**Given** the game is running with a pending perk draft
**When** `SIGINT` or `SIGTERM` is received
**Then** the game saves state (including `PendingDraft`) and exits cleanly

**Given** the game is reopened after a forced close with a pending draft
**When** the save file loads
**Then** the perk draft screen is shown with the same options (FR32)

**Given** combos are in `RunState.ComboQueue` and the game is closed then reopened
**When** the game loads
**Then** the queued combos are still present and begin executing on resume (FR33)

**Given** the game is running normally
**When** approximately 1 second elapses (60 ticks)
**Then** the game state is auto-saved without player action

---

## Epic 6: Meta-Progression & Run Loop

Player completes a full run, earns permanent unlocks, and starts run 2 with visible acceleration — the complete product loop.
**FRs covered:** FR5, FR22, FR23, FR24, FR25, FR26, FR27

### Story 6.1: Run End Detection & Summary Screen

As a player,
I want to see a summary of my completed run when I've reached the plateau,
So that I can reflect on my performance before starting a new run.

**Acceptance Criteria:**

**Given** all upgrades are maxed and no further floor progression is possible
**When** the run-end condition is evaluated
**Then** a `RunEndMsg` is emitted and the screen transitions to `RunSummaryScreen`

**Given** the run summary screen is showing
**When** the player views it
**Then** it displays: floors cleared, total resources generated, perks drafted, and time elapsed (FR22)

**Given** the run summary screen is showing
**When** the player presses the confirm key
**Then** the screen advances to the meta-progression unlock screen

---

### Story 6.2: Meta-Progression Engine

As a player,
I want permanent unlocks earned during my run to carry forward into future runs,
So that each run makes subsequent runs faster and more varied.

**Acceptance Criteria:**

**Given** a completed `RunState`
**When** `engine.ApplyRunToMeta()` is called
**Then** the returned `MetaState` contains new permanent unlocks based on the run's achievements (FR23)

**Given** a run with recorded effective combos
**When** `ApplyRunToMeta()` processes the run
**Then** the best-performing combo is stored in `MetaState.BestCombos` (FR26)

**Given** a `MetaState` with unlocked perks
**When** `engine.DrawPerkDraft()` is called in run 2+
**Then** the draft pool includes perks from `MetaState.UnlockedPerks` (FR14)

---

### Story 6.3: Meta-Progression Screen & New Run Start

As a player,
I want to see what I've unlocked before starting my next run,
So that I feel the reward of progression before diving back in.

**Acceptance Criteria:**

**Given** the meta-progression screen is showing
**When** the player views it
**Then** newly earned permanent unlocks are listed clearly (FR27)

**Given** the meta-progression screen is showing
**When** the player presses the confirm key
**Then** `engine.NewRunState()` is called and the game transitions to `GameScreen` with a fresh run (FR24)

**Given** a new run starts
**When** `RunState.RunNumber` is inspected
**Then** it is incremented by 1 from the previous run

---

### Story 6.4: Run 2+ Acceleration & Resource Type Unlock

As a player,
I want run 2 and beyond to feel noticeably faster from the start,
So that the meta-progression reward is immediately satisfying and motivates continued play.

**Acceptance Criteria:**

**Given** run 2 starts with meta-progression unlocks applied
**When** the time to clear floor 1 is compared to run 1
**Then** floor 1 clears measurably faster (FR25)

**Given** `MetaState` contains a resource type unlock
**When** `engine.NewRunState()` is called
**Then** the new `RunState.Resources` includes the unlocked resource type as an active resource (FR5)

**Given** run 2 begins
**When** the player views the upgrade panel
**Then** upgrades unlocked via meta-progression are available that were not present in run 1
