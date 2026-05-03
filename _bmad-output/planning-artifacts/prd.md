---
stepsCompleted: ['step-01-init', 'step-02-discovery', 'step-02b-vision', 'step-02c-executive-summary', 'step-03-success', 'step-04-journeys', 'step-05-domain', 'step-06-innovation', 'step-07-project-type', 'step-08-scoping', 'step-09-functional', 'step-10-nonfunctional', 'step-11-polish', 'step-12-complete']
releaseMode: phased
inputDocuments: []
workflowType: 'prd'
briefCount: 0
researchCount: 0
brainstormingCount: 0
projectDocsCount: 0
classification:
  projectType: cli_tool
  domain: gaming-idle-roguelite
  complexity: medium
  projectContext: greenfield
  platform: terminal (Go + Bubbletea)
  theme: fantasy-dungeon x terminal/computer hybrid
  activeInputModel: semi-idle
  runStructure: roguelite-draft
  metaProgression: permanent unlocks expand draft pool across runs
  social: solo only
---

# Product Requirements Document — afk-x

**Author:** Boss
**Date:** 2026-05-02

## Executive Summary

afk-x is a semi-idle incremental dungeon-crawler for the terminal, built in Go with Bubbletea. It targets developers and terminal-native users with work downtime — the game advances while AFK but rewards skill when the player is present. Core loop: generate resources, spend them on upgrades, draft perks at floor clears, and climb a corrupted tower. Each completed run feeds permanent meta-progression that makes the next run faster and more varied. MVP success criterion: player finishes run 1 and immediately starts run 2.

### What Makes This Special

afk-x is almost certainly the only AFK idle game native to the terminal. The terminal aesthetic is load-bearing: resources are CPU cycles and memory shards, enemies are rogue daemons, abilities are shell commands, dungeon floors are corrupted OS processes. Players aren't hiding a browser tab — they're running something that belongs in their environment. The roguelite draft layer (3 random perks per floor clear, pick 1) ensures no two runs feel identical while keeping active input near-zero. The game earns retention structurally through the first loop, not through notifications or FOMO.

## Project Classification

| Field | Value |
|---|---|
| Project Type | CLI/TUI Application (Go + Bubbletea) |
| Domain | Gaming — idle/incremental + roguelite |
| Complexity | Medium |
| Context | Greenfield |
| Platform | Terminal (any shell) |
| Social | Solo only |
| Active Input | Semi-idle — AFK progression with optional command-chaining and resource management |

## Success Criteria

### User Success

- Player reaches first floor clear within ≤20 minutes of a fresh start
- First run plateaus at 5–8 hours — loop feels complete, not abandoned
- Player starts run 2 immediately after run 1; meta-progression visibly accelerates early floors within the first 5 minutes
- Sessions are drop-in/drop-out — AFK accumulation means any window of attention is rewarding

### Business Success

- Project ships publicly on GitHub and accumulates stars over a 6-month horizon
- Codebase is clean enough to attract contributors or forks
- Distributed via `go install` or binary release — no infrastructure required

### Technical Success

- Runs correctly in iTerm2, Terminal.app, GNOME Terminal, Windows Terminal, and tmux
- Save file is portable human-readable JSON, transferable between machines
- Single binary, no runtime dependencies

### Measurable Outcomes

| Outcome | Target |
|---|---|
| Time to first floor clear | ≤ 20 minutes |
| First run duration before plateau | 5–8 hours |
| Run 2 early acceleration | Noticeable within first 5 min |
| Terminal compatibility | 5+ mainstream terminals |
| Save file format | Portable JSON |

## Product Scope & Roadmap

### MVP Strategy

Experience MVP — the game loop is the product. Validation: does the player start run 2 immediately after run 1? Solo developer constraint: every deferred feature is a ship-sooner feature.

### Phase 1 — MVP

All four user journeys supported (First Run, AFK Check-in, Active Floor Clear, End of Run / New Run Start).

**Must-Have Capabilities:**
- Resource tick engine — CPU cycles, memory shards, process threads generating in real time
- AFK accumulation — offline time calculated and applied on resume
- Floor-based tower progression with corrupted daemon enemies
- Upgrade purchase system
- Roguelite perk draft — 3 random perks per floor clear, pick 1
- Command input panel with combo execution and efficiency bonuses
- In-game command reference panel (`?` toggle)
- Meta-progression — permanent unlocks carry across runs, expand draft pool
- Run summary screen + new run start flow
- Portable JSON save file (atomic write, XDG-compliant path)
- `--version`, `--save <path>`, `--reset` launch flags
- `NO_COLOR` env var support
- Cross-terminal compatibility (5+ mainstream terminals)
- Boot sequence with light narrative flavor text
- Milestone flavor text on floor clears
- Graceful `SIGINT`/`SIGTERM` with save-before-exit
- Minimum terminal size check on launch

### Phase 2 — Growth

- Expanded story and tower lore
- Achievement system
- Broader enemy/daemon variety with unique mechanics
- Boss encounters with dedicated active resource management
- Additional prestige tiers beyond first reset
- Config file support (keybinds, color themes)
- Tab completion in command input panel

### Phase 3 — Vision

- Rich narrative layer (full lore, named characters, branching flavor text)
- Community run-sharing (seed codes, build exports)
- Mod/plugin support for custom resource types or floor packs
- Optional web dashboard for passive run tracking

### Risk Mitigation

| Risk | Mitigation |
|---|---|
| Terminal rendering fragmentation | Bubbletea abstraction; test on 5+ terminals before release |
| Save corruption on crash/interrupt | Atomic write (temp → rename); `SIGINT`/`SIGTERM` hooks |
| First run pacing off (20-min milestone) | Progression constants in config struct for easy tuning; playtest before release |
| Command combo balancing | Start conservative; additive bonuses only until playtested |
| Command-chaining too complex for idle players | Combos are optional — AFK path always works; active play is upside, not required |
| Novel mechanic confusing at first launch | Boot sequence + first floor acts as embedded tutorial |

## User Journeys

### Journey 1: New Player — First Run

**Alex** is a backend developer three hours into a slow Friday afternoon. The deploy is done, Slack is quiet. He spots afk-x on GitHub and gives it five minutes.

He runs `go install`, types `afk-x`. A boot sequence — corrupted ASCII, a fragment of lore about a tower no process has ever escaped — then a clean TUI. CPU cycles ticking at 1/sec. A single upgrade available. He buys it. Cycles tick faster.

Fifteen minutes in he's cleared floor 1, drafted his first perk (_Overclock_ over _Memory Leak_ and _Fork Bomb_), and has three resource types generating. At the 20-minute mark a milestone fires — flavor text about a daemon that corrupted the kernel years ago. He wants to know what's up there.

Five hours later the progression slows. Upgrades are maxed. The tower sends: _Signal fading. This run ends here._ He presses `R` without thinking. Early floors dissolve in minutes. His second run already feels different.

**Capabilities revealed:** Boot/intro flow, resource tick engine, upgrade purchase, floor progression, perk draft UI, milestone flavor text, AFK accumulation, run-end prompt, meta-progression carry-over.

---

### Journey 2: Returning Player — AFK Check-in

Alex has been in back-to-back meetings since 9am. Four minutes before his next call, he opens the terminal and types `afk-x`. Save loads silently. The resource counter shows four hours of accumulated cycles and shards. A floor clear he was 20 minutes away from last night is auto-completed. Two upgrade slots available.

He buys upgrades, queues a command combo for the active zone, and closes the terminal. When he returns at 3pm, the combo has executed, the floor is cleared, and a perk draft is waiting. He picks in 30 seconds and gets back to work.

**Capabilities revealed:** Persistent save/resume, offline accumulation on load, staged command queue, perk draft persistence.

---

### Journey 3: Active Player — Floor Clear

Quiet Tuesday. Alex has a real lunch break. He opens the command input panel and types `scan → exploit → loot`. A 40% efficiency bonus fires. He tries `isolate → drain → corrupt` — different bonus. He clears the zone 3x faster than idle.

Floor clear triggers the draft. _Parallel Threads_ (passive tick boost), _Root Access_ (combo damage multiplier), _Swap Space_ (memory shard overflow converts to cycles). He's building memory-heavy — takes _Swap Space_. Synergy is immediate. He finishes lunch having cleared two floors.

**Capabilities revealed:** Command input panel, combo parser and validator, combo efficiency bonus system, perk synergy legibility in draft UI.

---

### Journey 4: End of Run / New Run Start

Run 1, hour 7. Upgrades maxed, resources generating but nowhere to spend them. The tower sends: _The kernel is stable. No further corruption detected. This signal ends here._

Run summary: floors cleared, resources generated, perks drafted, time elapsed. Then meta-progression: draft pool expands from 12 to 18 perks, process threads unlock for run 2, best combo saved as a preset.

He presses `R`. Floor 1 clears in 45 seconds. He's seeing perks he's never seen before. Ten minutes in, he's further than he was at hour 1 of run 1. The loop is unbroken.

**Capabilities revealed:** Run summary screen, meta-progression unlock display, persistent combo presets, expanded draft pool, new resource type unlock, run 2 speed differential feedback.

---

### Journey Requirements Summary

| Capability Area | Revealed By |
|---|---|
| Resource tick engine + AFK accumulation | Journeys 1, 2 |
| Floor progression + enemy encounters | Journeys 1, 3 |
| Perk draft UI (3 options, pick 1) | Journeys 1, 3 |
| Command input + combo system | Journeys 2, 3 |
| Staged command queue (AFK-safe) | Journey 2 |
| Persistent save/resume | Journeys 1, 2 |
| Milestone flavor text | Journey 1 |
| Run summary + meta-progression screen | Journey 4 |
| Expanded draft pool across runs | Journey 4 |
| Run 2 speed differential feedback | Journey 4 |

## Innovation & Novel Patterns

### Detected Innovation Areas

**1. Terminal as Native Game Medium**
afk-x is designed *for* the terminal, not merely compatible with it. The TUI rendering, command-input mechanic, and tech-themed resource layer are all load-bearing. The game looks like work, runs where developers live, and uses the terminal's affordances (text, commands, ASCII) as first-class gameplay elements. No existing idle game occupies this space.

**2. Command-Chaining as Skill Expression**
Typed command sequences (`scan → exploit → loot`) are the skill mechanic — novel in idle games, which typically use clicking or passive timers. Command-chaining borrows from real-world terminal fluency: experienced users feel at home; new players discover a transferable skill loop. The mechanic is thematically coherent (hacking a dungeon OS) and mechanically distinct from any existing idle game.

**3. Semi-Idle + Roguelite Draft in Terminal**
AFK accumulation, roguelite perk drafting, and command-chaining haven't been assembled in this form before — let alone in a terminal context. Each element is individually proven; the combination and medium are the innovation.

### Market Context & Competitive Landscape

- **Idle/incremental games** (Cookie Clicker, Magic Research, NGU Idle): Web or mobile, passive-first, no terminal presence
- **Terminal toys/tools** (pipes, cmatrix, nyancat): Aesthetic novelty only, no game loop
- **Roguelite idle hybrids** (Idle Slayer, Tap Wizard 2): Mobile/web, no command mechanic, no terminal identity
- **afk-x gap:** No meaningful competitor in "terminal-native idle game." The closest analogue is a browser idle game in a tmux pane — a workaround, not a product.

### Validation Approach

- **Obsession metric:** First loop → immediate re-run (defined in success criteria)
- **Terminal identity:** GitHub stars from developer-adjacent audience validates the identity hypothesis
- **Command-chaining engagement:** Voluntary use of command input over AFK confirms the mechanic has pull

## CLI/TUI Specific Requirements

### Launch Interface

| Command | Behavior |
|---|---|
| `afk-x` | Launch game; load save if present, otherwise start fresh |
| `afk-x --version` | Print version string and exit |
| `afk-x --save <path>` | Override default save file location |
| `afk-x --reset` | Wipe save file after in-terminal confirmation prompt |

No subcommands. No piped input/output. Launching without a TTY exits gracefully with an error message.

### Configuration

No config file in MVP. Environment variables only:

| Variable | Behavior |
|---|---|
| `XDG_DATA_HOME` | Base path for save file (falls back to `~/.local/share/afk-x/`) |
| `NO_COLOR` | Disable all ANSI color output |

### In-Game Command System

- **Input panel:** Text field accepting command sequences (`scan`, `exploit`, `loot`, etc.)
- **Combo execution:** Sequences validated on submit; efficiency bonus applied for correct chains
- **Reference panel:** Available commands and valid combos, toggled with `?`
- **Invalid input:** Rejected with inline feedback; no state mutation
- **Queued combos:** Staged for execution while AFK; persist across game close/reopen

### Output

- **TUI:** All game output rendered via Bubbletea
- **Save file:** JSON, written atomically on state change and clean exit
- **Run summary:** Displayed in-TUI at run end; not written to disk
- No external logs for MVP

## Functional Requirements

### Resource Generation

- **FR1:** Player can observe resources generating automatically in real time
- **FR2:** Player can accumulate resources while the game is not running (AFK/offline accumulation)
- **FR3:** System calculates and applies offline resource gains on resume
- **FR4:** Player can observe multiple distinct resource types simultaneously (CPU cycles, memory shards, process threads)
- **FR5:** New resource types unlock as the player advances across runs

### Upgrade & Floor Progression

- **FR6:** Player can purchase upgrades using accumulated resources
- **FR7:** Player can view available upgrades and their costs at any time
- **FR8:** Player can progress through numbered tower floors by meeting floor completion conditions
- **FR9:** Player encounters daemon enemies on each floor that must be overcome to progress
- **FR10:** Player can observe current floor number and progress toward the next floor clear

### Perk Draft System

- **FR11:** Player is presented with a draft of 3 randomly selected perks upon clearing a floor
- **FR12:** Player can select one perk from the draft to apply to the current run
- **FR13:** Each drafted perk has a distinct effect that influences run strategy
- **FR14:** The available perk pool expands as meta-progression unlocks are earned across runs
- **FR15:** Higher-tier perks become accessible in the draft pool as the player advances

### Command & Combo System

- **FR16:** Player can enter command sequences via an in-game text input panel
- **FR17:** System validates entered command sequences against known combos
- **FR18:** Player can execute valid command combos to receive efficiency bonuses in the current zone
- **FR19:** Player can view a reference panel listing available commands and valid combos
- **FR20:** System rejects invalid command input without altering game state
- **FR21:** Player can queue command combos to execute while away from the game

### Meta-Progression & Run Flow

- **FR22:** Player receives a run summary on run end (floors cleared, resources generated, perks drafted, time elapsed)
- **FR23:** Permanent unlocks earned during a run carry over to all subsequent runs
- **FR24:** Player can start a new run immediately after run end
- **FR25:** Run 2+ begins with visibly faster early progression reflecting meta-progression unlocks
- **FR26:** Player's best command combos from prior runs are available as presets in subsequent runs
- **FR27:** Meta-progression unlock screen displays newly earned unlocks before the new run starts

### Save & Persistence

- **FR28:** Game state is automatically persisted to a portable JSON file
- **FR29:** Game state is fully restored on relaunch from the save file
- **FR30:** Player can specify a custom save file location at launch
- **FR31:** Player can wipe save data and start fresh via a launch command, with a confirmation step
- **FR32:** Pending perk drafts persist across game close and reopen until the player makes a selection
- **FR33:** Staged command queues persist while the game is not running

### Narrative & Presentation

- **FR34:** Player is presented with a boot sequence on first launch that includes narrative flavor text
- **FR35:** Player receives flavor text at designated milestone floors revealing tower lore
- **FR36:** All game elements — resources, enemies, abilities, upgrades — use terminal/computing terminology and aesthetic
- **FR37:** Player can read the current floor's narrative context within the game UI

### Application Shell

- **FR38:** Player can launch the game with no arguments to start or resume a session
- **FR39:** Player can query the installed application version from the command line
- **FR40:** Player can disable ANSI color output via the `NO_COLOR` environment variable
- **FR41:** System detects terminal dimensions on launch and displays an error if too small to render
- **FR42:** System saves game state on receiving `SIGINT` or `SIGTERM` before exiting
- **FR43:** Save files are stored in an XDG-compliant directory by default (`~/.local/share/afk-x/`)
- **FR44:** Game renders correctly across 5+ mainstream terminal emulators

## Non-Functional Requirements

### Performance

- **NFR1:** Resource tick calculations complete within 16ms per cycle (≤60fps render budget)
- **NFR2:** AFK accumulation calculation on resume completes within 500ms regardless of offline duration
- **NFR3:** Game launch (binary start to interactive TUI) completes within 2 seconds on a standard developer machine
- **NFR4:** Perk draft screen renders within 100ms of floor clear trigger
- **NFR5:** Save file write completes within 200ms

### Reliability

- **NFR6:** Save file is written atomically — a crash during write must not corrupt the existing save
- **NFR7:** Game state on resume matches game state at last save — no drift or loss from clean or unclean exits
- **NFR8:** AFK accumulation is deterministic — identical offline duration and game state always produce identical resource gain
- **NFR9:** Invalid or corrupted save files produce a clear error message and reset prompt, not a silent crash
- **NFR10:** Game handles `SIGINT` and `SIGTERM` without data loss under normal operating conditions

### Portability

- **NFR11:** Distributed as a single self-contained binary — no runtime dependencies beyond a POSIX-compatible terminal
- **NFR12:** Binary supports Linux (amd64, arm64) and macOS (amd64, arm64) at minimum
- **NFR13:** Game renders correctly in iTerm2, Terminal.app, GNOME Terminal, Windows Terminal, and tmux
- **NFR14:** Save files are human-readable JSON, loadable on any supported platform without modification
- **NFR15:** Binary size ≤ 20MB

### Accessibility

- **NFR16:** No game-critical information is encoded in color alone — all information conveyed through text and symbols
- **NFR17:** `NO_COLOR` fully disables ANSI color output without affecting gameplay or readability
- **NFR18:** All game interactions are operable via keyboard only — no mouse dependency
- **NFR19:** Text elements remain legible at standard terminal font sizes — no reliance on specific font rendering
