# Deferred Work

## Deferred from: code review of 1-4-resource-tick-engine (2026-05-03)

- `content.BaseRates` mutable `var` with no write guard — `ComputeRates` defensive copy prevents engine-side corruption but any package can still write `content.BaseRates[key] = X` directly; making BaseRates an accessor function would be a real fix [internal/content/resources.go]
- `ComputeRates` `run game.RunState` parameter currently unused — intentional placeholder for Story 2.x upgrade multipliers and Story 3.x perk bonuses [internal/engine/rates.go:16]
- Float accumulation drift over many ticks — `float64` addition over thousands of 16ms ticks is non-associative; acceptable for idle game session lengths but worth revisiting for save/load equality comparisons [internal/engine/tick.go]
- No test verifies `TickMsg` type delivered by `TickCmd` — AC4 requires `TickMsg` delivery but the cmd callback type is only testable by running a Bubbletea event loop; defer to Story 1.5 integration testing [internal/engine/tick_test.go]
- `engine.go` doc "no global mutable state" is technically overstated — `ComputeRates` reads `content.BaseRates`, a mutable package-level var; doc should be updated to acknowledge this once BaseRates immutability is addressed [internal/engine/engine.go:2]

## Deferred from: code review of 1-3-content-package-resource-narrative-definitions (2026-05-03)

- `BaseRates` external mutation risk — any importer can write `content.BaseRates[key] = X` and corrupt the baseline; Story 1.4 `ComputeRates()` must copy the map rather than reference it directly [internal/content/resources.go:14]
- `TestMilestoneTextForMilestoneFloors` hard-codes the milestone floor list `{5, 10, 20, 25}` — must be manually updated if milestone set changes in a future story [internal/content/content_test.go:87]
- `TestMilestoneTextEmptyForNonMilestone` doesn't exercise floor 0 or negative values — behavior is correct (missing key returns "") but the contract is undocumented in tests [internal/content/content_test.go:97]
- `TestBaseRatesDefined`/`TestResourceKeysMatchBaseRates` don't assert `len(content.BaseRates) == 3` — no enforcement that BaseRates has no extra unknown keys [internal/content/content_test.go:49–74]
- `BaseRates`-to-engine integration test missing — no test asserts that a fresh run's computed rates equal BaseRates values; belongs in Story 1.4 when `engine/rates.go` is written

## Deferred from: code review of 1-2-game-state-schema (2026-05-03)

- Use `content` package resource key constants in `TestResourcesMapType` and `TestJSONRoundTrip` instead of string literals — blocked until Story 1.3 defines the constants in `content/resources.go`
- `PendingDraft` nil invariant has no type-level enforcement (no custom `MarshalJSON`) — current test coverage sufficient for MVP; revisit if schema evolves
- Resources zero-value key ambiguity (`"cpu_cycles": 0.0` vs absent key) not tested or documented — game design question, defer to engine/consumer layer
- Version schema mismatch detection (loading `"version": 99`) has no tested enforcement path — belongs in `save` package (Story 5.x), not `game`
- Float precision edge cases (`1.1`, `0.3`, near `math.MaxFloat64`) not exercised in resource round-trip — theoretical for this game's value ranges
