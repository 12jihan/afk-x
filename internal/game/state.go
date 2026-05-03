// Package game defines the core game state types and JSON serialization contract
// shared by all packages in afk-x. This package contains structs and a constructor
// only — zero business logic, zero Bubbletea imports, zero internal package imports.
package game

import "time"

// GameState is the top-level save schema. It is the single source of truth
// for all persisted game data. Version must be checked on load for future
// schema migrations.
type GameState struct {
	Version int       `json:"version"`  // schema migration guard — always 1 for MVP
	Run     RunState  `json:"run"`
	Meta    MetaState `json:"meta"`
	SavedAt time.Time `json:"saved_at"` // AFK delta anchor — serializes as RFC3339
}

// RunState holds all mutable state for the current run. Resources uses
// map[string]float64 so new resource types (FR5) can be added in future runs
// without a schema change. All slice fields are nil by default — callers must
// not initialize them as empty slices, which would serialize as [] instead of null.
type RunState struct {
	Floor        int                `json:"floor"`
	Resources    map[string]float64 `json:"resources"`      // keys from content.Resource* constants
	Upgrades     map[string]int     `json:"upgrades"`
	ActivePerks  []string           `json:"active_perks"`
	ComboQueue   [][]string         `json:"combo_queue"`    // FR33: persisted AFK combo queue
	PendingDraft []string           `json:"pending_draft"`  // FR32: nil = no draft pending; null in JSON
	RunNumber    int                `json:"run_number"`
}

// MetaState holds permanent cross-run progression. Fields accumulate across
// all runs and are never reset by starting a new run.
type MetaState struct {
	PermanentUnlocks []string            `json:"permanent_unlocks"` // FR23
	BestCombos       map[string][]string `json:"best_combos"`       // FR26: zone → best combo sequence
	UnlockedPerks    []string            `json:"unlocked_perks"`    // FR14: expands perk draft pool
	RunCount         int                 `json:"run_count"`
}

// NewGameState returns a freshly initialized GameState ready for run 1.
// Maps are initialized to prevent nil-map write panics. Slice fields are
// intentionally left nil so they serialize as JSON null rather than [].
func NewGameState() GameState {
	return GameState{
		Version: 1,
		Run: RunState{
			Floor:     1,
			Resources: make(map[string]float64),
			Upgrades:  make(map[string]int),
			// ActivePerks, ComboQueue, PendingDraft: nil (zero value) — serializes as null
			RunNumber: 1,
		},
		Meta: MetaState{
			BestCombos: make(map[string][]string),
			// PermanentUnlocks, UnlockedPerks: nil (zero value) — serializes as null
		},
		SavedAt: time.Now(),
	}
}
