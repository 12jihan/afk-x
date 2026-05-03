---
stepsCompleted: ['step-01-document-discovery', 'step-02-prd-analysis', 'step-03-epic-coverage-validation', 'step-04-ux-alignment', 'step-05-epic-quality-review', 'step-06-final-assessment']
documentsFound:
  prd: '_bmad-output/planning-artifacts/prd.md'
  architecture: null
  epics: null
  ux: null
---

# Implementation Readiness Assessment Report

**Date:** 2026-05-02
**Project:** afk-x

## PRD Analysis

### Functional Requirements

**Resource Generation (FR1–FR5)**
- FR1: Player can observe resources generating automatically in real time
- FR2: Player can accumulate resources while the game is not running (AFK/offline accumulation)
- FR3: System calculates and applies offline resource gains on resume
- FR4: Player can observe multiple distinct resource types simultaneously (CPU cycles, memory shards, process threads)
- FR5: New resource types unlock as the player advances across runs

**Upgrade & Floor Progression (FR6–FR10)**
- FR6: Player can purchase upgrades using accumulated resources
- FR7: Player can view available upgrades and their costs at any time
- FR8: Player can progress through numbered tower floors by meeting floor completion conditions
- FR9: Player encounters daemon enemies on each floor that must be overcome to progress
- FR10: Player can observe current floor number and progress toward the next floor clear

**Perk Draft System (FR11–FR15)**
- FR11: Player is presented with a draft of 3 randomly selected perks upon clearing a floor
- FR12: Player can select one perk from the draft to apply to the current run
- FR13: Each drafted perk has a distinct effect that influences run strategy
- FR14: The available perk pool expands as meta-progression unlocks are earned across runs
- FR15: Higher-tier perks become accessible in the draft pool as the player advances

**Command & Combo System (FR16–FR21)**
- FR16: Player can enter command sequences via an in-game text input panel
- FR17: System validates entered command sequences against known combos
- FR18: Player can execute valid command combos to receive efficiency bonuses in the current zone
- FR19: Player can view a reference panel listing available commands and valid combos
- FR20: System rejects invalid command input without altering game state
- FR21: Player can queue command combos to execute while away from the game

**Meta-Progression & Run Flow (FR22–FR27)**
- FR22: Player receives a run summary on run end (floors cleared, resources generated, perks drafted, time elapsed)
- FR23: Permanent unlocks earned during a run carry over to all subsequent runs
- FR24: Player can start a new run immediately after run end
- FR25: Run 2+ begins with visibly faster early progression reflecting meta-progression unlocks
- FR26: Player's best command combos from prior runs are available as presets in subsequent runs
- FR27: Meta-progression unlock screen displays newly earned unlocks before the new run starts

**Save & Persistence (FR28–FR33)**
- FR28: Game state is automatically persisted to a portable JSON file
- FR29: Game state is fully restored on relaunch from the save file
- FR30: Player can specify a custom save file location at launch
- FR31: Player can wipe save data and start fresh via a launch command, with a confirmation step
- FR32: Pending perk drafts persist across game close and reopen until the player makes a selection
- FR33: Staged command queues persist while the game is not running

**Narrative & Presentation (FR34–FR37)**
- FR34: Player is presented with a boot sequence on first launch that includes narrative flavor text
- FR35: Player receives flavor text at designated milestone floors revealing tower lore
- FR36: All game elements use terminal/computing terminology and aesthetic
- FR37: Player can read the current floor's narrative context within the game UI

**Application Shell (FR38–FR44)**
- FR38: Player can launch the game with no arguments to start or resume a session
- FR39: Player can query the installed application version from the command line
- FR40: Player can disable ANSI color output via the NO_COLOR environment variable
- FR41: System detects terminal dimensions on launch and displays an error if too small to render
- FR42: System saves game state on receiving SIGINT or SIGTERM before exiting
- FR43: Save files are stored in an XDG-compliant directory by default (~/.local/share/afk-x/)
- FR44: Game renders correctly across 5+ mainstream terminal emulators

**Total FRs: 44**

### Non-Functional Requirements

**Performance (NFR1–NFR5)**
- NFR1: Resource tick calculations complete within 16ms per cycle (≤60fps render budget)
- NFR2: AFK accumulation calculation on resume completes within 500ms regardless of offline duration
- NFR3: Game launch completes within 2 seconds on a standard developer machine
- NFR4: Perk draft screen renders within 100ms of floor clear trigger
- NFR5: Save file write completes within 200ms

**Reliability (NFR6–NFR10)**
- NFR6: Save file is written atomically — crash during write must not corrupt existing save
- NFR7: Game state on resume matches game state at last save — no drift or loss
- NFR8: AFK accumulation is deterministic — identical inputs always produce identical resource gain
- NFR9: Invalid or corrupted save files produce a clear error message and reset prompt, not a silent crash
- NFR10: Game handles SIGINT and SIGTERM without data loss under normal operating conditions

**Portability (NFR11–NFR15)**
- NFR11: Distributed as a single self-contained binary — no runtime dependencies beyond a POSIX-compatible terminal
- NFR12: Binary supports Linux (amd64, arm64) and macOS (amd64, arm64) at minimum
- NFR13: Game renders correctly in iTerm2, Terminal.app, GNOME Terminal, Windows Terminal, and tmux
- NFR14: Save files are human-readable JSON, loadable on any supported platform without modification
- NFR15: Binary size ≤ 20MB

**Accessibility (NFR16–NFR19)**
- NFR16: No game-critical information encoded in color alone
- NFR17: NO_COLOR fully disables ANSI color output without affecting gameplay or readability
- NFR18: All game interactions operable via keyboard only — no mouse dependency
- NFR19: Text elements remain legible at standard terminal font sizes

**Total NFRs: 19**

### Additional Requirements

**Launch Interface (from CLI/TUI section):**
- `afk-x` — launch/resume
- `afk-x --version` — version query
- `afk-x --save <path>` — custom save path
- `afk-x --reset` — wipe save with confirmation

**Environment Variables:**
- `XDG_DATA_HOME` — base path for save file
- `NO_COLOR` — disable ANSI color

**Implicit constraints:**
- No config file in MVP (deferred to Phase 2)
- No tab completion in MVP (deferred to Phase 2)
- No piped input/output — interactive TUI only
- Solo player only — no multiplayer or leaderboards

### PRD Completeness Assessment

PRD is well-structured with 44 FRs across 8 capability areas and 19 NFRs across 4 quality categories. Vision, success criteria, user journeys, and scope are clearly defined. No architecture, UX, or epics documents exist yet — this assessment will evaluate PRD readiness for handoff to those phases.

## Epic Coverage Validation

### Coverage Matrix

No epics document found. All 44 FRs are currently unimplemented in any epic.

| FR Range | Capability Area | Epic Coverage | Status |
|---|---|---|---|
| FR1–FR5 | Resource Generation | No epics | ❌ Not yet planned |
| FR6–FR10 | Upgrade & Floor Progression | No epics | ❌ Not yet planned |
| FR11–FR15 | Perk Draft System | No epics | ❌ Not yet planned |
| FR16–FR21 | Command & Combo System | No epics | ❌ Not yet planned |
| FR22–FR27 | Meta-Progression & Run Flow | No epics | ❌ Not yet planned |
| FR28–FR33 | Save & Persistence | No epics | ❌ Not yet planned |
| FR34–FR37 | Narrative & Presentation | No epics | ❌ Not yet planned |
| FR38–FR44 | Application Shell | No epics | ❌ Not yet planned |

### Coverage Statistics

- Total PRD FRs: 44
- FRs covered in epics: 0
- Coverage percentage: 0% — epics not yet created (expected at this stage)

## UX Alignment Assessment

### UX Document Status

Not found. No formal UX design document exists.

### UX Implied Assessment

UX is implied by the PRD. afk-x is a user-facing TUI application. The PRD references multiple distinct UI surfaces:
- Resource display (FR1, FR4)
- Upgrade panel (FR7)
- Floor progress indicator (FR10)
- Perk draft screen — 3 options, pick 1 (FR11–FR12)
- Command input panel + reference panel (FR16, FR19)
- Boot sequence screen (FR34)
- Floor narrative context panel (FR37)
- Run summary + meta-progression unlock screen (FR22, FR27)

### Warnings

⚠️ **No UX document** — screen layouts, component arrangement, and keyboard navigation flows are undefined. For a Bubbletea TUI, this means the following are unspecified:
- Panel layout and proportions (how the screen is divided between resource display, floor status, command input, etc.)
- Keyboard navigation map (which keys do what outside of command input)
- Visual hierarchy of information (what the player sees at a glance vs. on demand)
- State transitions (what the screen looks like during perk draft vs. active floor vs. run summary)

**Severity:** Medium — not a blocker for architecture or epic planning, but screen layout decisions will need to be made during implementation. Recommend defining at minimum a rough ASCII wireframe of the main game screen before epic breakdown to avoid rework.

### Alignment Issues

No misalignment issues (no UX document to compare against). All UX decisions are deferred to implementation.

## Epic Quality Review

### Status

No epics document exists. Quality review cannot be performed against actual epics.

### Pre-Creation Standards (Guidance for When Epics Are Written)

The following criteria must be met when `/bmad-create-epics-and-stories` is run:

#### Epic Structure Requirements
- Each epic must deliver **player-visible value** — no "Setup Game Engine" or "Create Data Models" as standalone epics
- Epics must be independently shippable — a player should be able to do something meaningful after each epic
- Suggested epic boundaries based on PRD capability areas:
  - **Epic 1:** Core resource tick + upgrade purchase (FR1–FR7) — player can generate and spend resources
  - **Epic 2:** Floor progression + enemy encounters (FR8–FR10, FR34–FR37) — player can climb the tower
  - **Epic 3:** Perk draft system (FR11–FR15) — player gets roguelite choices at floor clears
  - **Epic 4:** Command & combo system (FR16–FR21) — player can actively engage
  - **Epic 5:** AFK accumulation + save/persistence (FR2–FR3, FR28–FR33) — game works while away
  - **Epic 6:** Meta-progression & run flow (FR22–FR27, FR4–FR5) — run end and run 2 loop
  - **Epic 7:** Application shell + distribution (FR38–FR44) — shippable binary

#### Story Requirements
- Each story must be completable without depending on future stories in the same epic
- Database/state schema should be created story-by-story, not all upfront
- Acceptance criteria must follow Given/When/Then format with measurable outcomes
- Greenfield setup story required: Epic 1 Story 1 = "Initialize Go module, Bubbletea project scaffold, CI build"

#### 🟡 Pre-emptive Concerns
- **Save schema evolution:** If save file JSON schema is defined in Epic 5 but resource types are added in Epic 6 (FR5), the save schema will need migration — plan this dependency explicitly
- **Perk pool expansion (FR14):** Depends on meta-progression unlocks (FR23) which are in Epic 6 — Epic 3 must stub this correctly so Epic 6 can extend it without rework
- **Command queue persistence (FR33):** Staged combos persisting while AFK touches both the combo system (Epic 4) and the save system (Epic 5) — story ordering matters here

## Summary and Recommendations

### Overall Readiness Status

**PRD: READY** — The PRD is solid and ready for handoff to architecture and epic breakdown.
**Epics: NOT CREATED** — No epics exist yet; this is expected and the correct next step.
**UX: ADVISORY** — No UX document; medium-severity gap that should be addressed before or alongside epic breakdown.

### Issues by Severity

#### 🟢 No Blockers Found
The PRD itself has no critical gaps. All 44 FRs are well-defined, testable, and implementation-agnostic. All 19 NFRs are specific and measurable. Vision, success criteria, user journeys, and scope are complete and internally consistent.

#### 🟠 Medium: No UX / Screen Layout Definition
The PRD defines 8+ distinct UI surfaces (resource display, upgrade panel, perk draft screen, command input, reference panel, boot sequence, floor narrative, run summary) but provides no layout or navigation specification. Bubbletea requires explicit component layout decisions. Without at minimum a rough wireframe, each UI surface will be designed ad-hoc during implementation, risking inconsistency and rework.

**Impact:** Medium — affects implementation speed and UX consistency, not correctness.

#### 🟡 Low: Three Cross-Epic Dependencies to Watch
1. **Save schema + resource type unlock** (FR5 + FR28): Process threads unlock in run 2 must be accounted for in the save schema from the start, or a migration story will be needed.
2. **Perk pool expansion stub** (FR14 + FR23): The perk draft system must be built with extension hooks for the meta-progression unlock layer.
3. **Command queue in save file** (FR21 + FR33): Staged combos must be persisted — the save schema and combo system must agree on this from story 1.

### Recommended Next Steps

1. **(Optional but recommended) Create a TUI wireframe** — a rough ASCII layout of the main game screen, perk draft screen, and run summary screen. This takes 30 minutes and prevents layout rework during implementation. Use `/bmad-create-ux-design` or sketch it manually.

2. **Run `/bmad-create-architecture`** — define the Go package structure, Bubbletea model/update/view architecture, game loop tick design, and save file schema. The three cross-epic dependencies above must be resolved here.

3. **Run `/bmad-create-epics-and-stories`** — break the 44 FRs into epics and stories using the suggested groupings in this report. Epic 1 Story 1 must be the project scaffold.

4. **Review this report** before epic breakdown — the pre-creation standards section has specific guidance on story ordering and dependency management for this project.

### Final Note

This assessment found **0 blocking issues** in the PRD and **3 low-severity dependency risks** to watch during epic breakdown. The PRD is ready to proceed. Address the UX gap and cross-epic dependencies during architecture design, and the project will be well-positioned for clean implementation.

**Report saved to:** `_bmad-output/planning-artifacts/implementation-readiness-report-2026-05-02.md`
