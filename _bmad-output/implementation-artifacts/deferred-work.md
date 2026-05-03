# Deferred Work

## Deferred from: code review of 1-2-game-state-schema (2026-05-03)

- Use `content` package resource key constants in `TestResourcesMapType` and `TestJSONRoundTrip` instead of string literals — blocked until Story 1.3 defines the constants in `content/resources.go`
- `PendingDraft` nil invariant has no type-level enforcement (no custom `MarshalJSON`) — current test coverage sufficient for MVP; revisit if schema evolves
- Resources zero-value key ambiguity (`"cpu_cycles": 0.0` vs absent key) not tested or documented — game design question, defer to engine/consumer layer
- Version schema mismatch detection (loading `"version": 99`) has no tested enforcement path — belongs in `save` package (Story 5.x), not `game`
- Float precision edge cases (`1.1`, `0.3`, near `math.MaxFloat64`) not exercised in resource round-trip — theoretical for this game's value ranges
