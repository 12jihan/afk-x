package game_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/12jihan/afk-x/internal/game"
)

// Alias package-level names for brevity in tests.
type GameState = game.GameState
type RunState = game.RunState

func NewGameState() game.GameState { return game.NewGameState() }

// TestVersionSetToOne verifies AC1: Version field is 1 and serializes as "version":1
func TestVersionSetToOne(t *testing.T) {
	state := NewGameState()
	if state.Version != 1 {
		t.Errorf("NewGameState().Version = %d, want 1", state.Version)
	}

	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"version":1`) {
		t.Errorf("JSON does not contain \"version\":1, got: %s", jsonStr)
	}
}

// TestAllFieldsSnakeCase verifies AC1: all JSON keys are snake_case
func TestAllFieldsSnakeCase(t *testing.T) {
	state := NewGameState()
	state.Run.ActivePerks = []string{"overclock"}
	state.Run.ComboQueue = [][]string{{"scan", "exploit"}}
	state.Run.PendingDraft = []string{"perk_a"}
	state.Meta.PermanentUnlocks = []string{"unlock_1"}
	state.Meta.UnlockedPerks = []string{"fork_bomb"}

	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	jsonStr := string(data)

	// Check required snake_case keys exist
	snakeCaseKeys := []string{
		`"version"`,
		`"saved_at"`,
		`"run_number"`,
		`"pending_draft"`,
		`"active_perks"`,
		`"combo_queue"`,
		`"permanent_unlocks"`,
		`"best_combos"`,
		`"unlocked_perks"`,
		`"run_count"`,
	}
	for _, key := range snakeCaseKeys {
		if !strings.Contains(jsonStr, key) {
			t.Errorf("JSON missing expected snake_case key %s; got: %s", key, jsonStr)
		}
	}
}

// TestJSONRoundTrip verifies AC2: marshal → unmarshal produces identical struct
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

	// Version
	if restored.Version != original.Version {
		t.Errorf("Version: got %d, want %d", restored.Version, original.Version)
	}
	// SavedAt — use .Equal() because JSON normalizes to UTC
	if !restored.SavedAt.Equal(original.SavedAt) {
		t.Errorf("SavedAt mismatch: got %v, want %v", restored.SavedAt, original.SavedAt)
	}
	// Floor
	if restored.Run.Floor != original.Run.Floor {
		t.Errorf("Run.Floor: got %d, want %d", restored.Run.Floor, original.Run.Floor)
	}
	// Resources
	for k, v := range original.Run.Resources {
		if restored.Run.Resources[k] != v {
			t.Errorf("Run.Resources[%q]: got %v, want %v", k, restored.Run.Resources[k], v)
		}
	}
	// Upgrades
	for k, v := range original.Run.Upgrades {
		if restored.Run.Upgrades[k] != v {
			t.Errorf("Run.Upgrades[%q]: got %d, want %d", k, restored.Run.Upgrades[k], v)
		}
	}
	// ActivePerks
	if len(restored.Run.ActivePerks) != len(original.Run.ActivePerks) {
		t.Errorf("Run.ActivePerks len: got %d, want %d", len(restored.Run.ActivePerks), len(original.Run.ActivePerks))
	}
	// ComboQueue
	if len(restored.Run.ComboQueue) != len(original.Run.ComboQueue) {
		t.Errorf("Run.ComboQueue len: got %d, want %d", len(restored.Run.ComboQueue), len(original.Run.ComboQueue))
	}
	// PendingDraft
	if len(restored.Run.PendingDraft) != len(original.Run.PendingDraft) {
		t.Errorf("Run.PendingDraft len: got %d, want %d", len(restored.Run.PendingDraft), len(original.Run.PendingDraft))
	}
	// RunNumber
	if restored.Run.RunNumber != original.Run.RunNumber {
		t.Errorf("Run.RunNumber: got %d, want %d", restored.Run.RunNumber, original.Run.RunNumber)
	}
	// Meta.PermanentUnlocks
	if len(restored.Meta.PermanentUnlocks) != len(original.Meta.PermanentUnlocks) {
		t.Errorf("Meta.PermanentUnlocks len: got %d, want %d", len(restored.Meta.PermanentUnlocks), len(original.Meta.PermanentUnlocks))
	}
	// Meta.RunCount
	if restored.Meta.RunCount != original.Meta.RunCount {
		t.Errorf("Meta.RunCount: got %d, want %d", restored.Meta.RunCount, original.Meta.RunCount)
	}
	// Meta.UnlockedPerks
	if len(restored.Meta.UnlockedPerks) != len(original.Meta.UnlockedPerks) {
		t.Errorf("Meta.UnlockedPerks len: got %d, want %d", len(restored.Meta.UnlockedPerks), len(original.Meta.UnlockedPerks))
	}
}

// TestPendingDraftNullWhenNil verifies AC4: nil PendingDraft → JSON null, NOT []
func TestPendingDraftNullWhenNil(t *testing.T) {
	state := NewGameState()
	// PendingDraft must be nil (zero value) from constructor

	if state.Run.PendingDraft != nil {
		t.Errorf("NewGameState().Run.PendingDraft should be nil, got %v", state.Run.PendingDraft)
	}

	data, err := json.Marshal(state.Run)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	jsonStr := string(data)

	// Must contain null, not []
	if strings.Contains(jsonStr, `"pending_draft":[]`) {
		t.Errorf("PendingDraft serialized as [] instead of null: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"pending_draft":null`) {
		t.Errorf("PendingDraft not serialized as null: %s", jsonStr)
	}
}

// TestResourcesMapType verifies AC3: Resources is map[string]float64, round-trips correctly
func TestResourcesMapType(t *testing.T) {
	state := NewGameState()

	// Populate with planned content key strings
	state.Run.Resources["cpu_cycles"] = 100.5
	state.Run.Resources["memory_shards"] = 50.0
	state.Run.Resources["process_threads"] = 10.25

	data, err := json.Marshal(state.Run)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var restored RunState
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	cases := map[string]float64{
		"cpu_cycles":      100.5,
		"memory_shards":   50.0,
		"process_threads": 10.25,
	}
	for key, want := range cases {
		if got := restored.Resources[key]; got != want {
			t.Errorf("Resources[%q]: got %v, want %v", key, got, want)
		}
	}
}

// TestSavedAtRFC3339 verifies SavedAt serializes as RFC3339 (time.Time default)
func TestSavedAtRFC3339(t *testing.T) {
	now := time.Date(2026, 5, 3, 12, 0, 0, 0, time.UTC)
	state := NewGameState()
	state.SavedAt = now

	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	jsonStr := string(data)
	// time.Time marshals to RFC3339Nano by default
	if !strings.Contains(jsonStr, "2026-05-03") {
		t.Errorf("SavedAt not in RFC3339 format in JSON: %s", jsonStr)
	}
}
