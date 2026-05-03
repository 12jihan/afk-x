package game_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/12jihan/afk-x/internal/content"
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

// TestJSONRoundTrip verifies AC2: marshal → unmarshal produces identical struct.
// P1 fix: all slice and map fields verified by element values, not just length.
func TestJSONRoundTrip(t *testing.T) {
	original := NewGameState()
	original.Run.Floor = 5
	original.Run.Resources[content.CPUCycles] = 1250.5
	original.Run.Resources[content.MemoryShards] = 340.0
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
	// Resources — value per key
	for k, v := range original.Run.Resources {
		if restored.Run.Resources[k] != v {
			t.Errorf("Run.Resources[%q]: got %v, want %v", k, restored.Run.Resources[k], v)
		}
	}
	// Upgrades — value per key
	for k, v := range original.Run.Upgrades {
		if restored.Run.Upgrades[k] != v {
			t.Errorf("Run.Upgrades[%q]: got %d, want %d", k, restored.Run.Upgrades[k], v)
		}
	}
	// ActivePerks — element values (P1)
	if len(restored.Run.ActivePerks) != len(original.Run.ActivePerks) {
		t.Fatalf("Run.ActivePerks len: got %d, want %d", len(restored.Run.ActivePerks), len(original.Run.ActivePerks))
	}
	for i, want := range original.Run.ActivePerks {
		if restored.Run.ActivePerks[i] != want {
			t.Errorf("Run.ActivePerks[%d]: got %q, want %q", i, restored.Run.ActivePerks[i], want)
		}
	}
	// ComboQueue — inner slice element values (P1)
	if len(restored.Run.ComboQueue) != len(original.Run.ComboQueue) {
		t.Fatalf("Run.ComboQueue len: got %d, want %d", len(restored.Run.ComboQueue), len(original.Run.ComboQueue))
	}
	for i, origInner := range original.Run.ComboQueue {
		restoredInner := restored.Run.ComboQueue[i]
		if len(restoredInner) != len(origInner) {
			t.Fatalf("Run.ComboQueue[%d] len: got %d, want %d", i, len(restoredInner), len(origInner))
		}
		for j, want := range origInner {
			if restoredInner[j] != want {
				t.Errorf("Run.ComboQueue[%d][%d]: got %q, want %q", i, j, restoredInner[j], want)
			}
		}
	}
	// PendingDraft — element values (P1)
	if len(restored.Run.PendingDraft) != len(original.Run.PendingDraft) {
		t.Fatalf("Run.PendingDraft len: got %d, want %d", len(restored.Run.PendingDraft), len(original.Run.PendingDraft))
	}
	for i, want := range original.Run.PendingDraft {
		if restored.Run.PendingDraft[i] != want {
			t.Errorf("Run.PendingDraft[%d]: got %q, want %q", i, restored.Run.PendingDraft[i], want)
		}
	}
	// RunNumber
	if restored.Run.RunNumber != original.Run.RunNumber {
		t.Errorf("Run.RunNumber: got %d, want %d", restored.Run.RunNumber, original.Run.RunNumber)
	}
	// Meta.PermanentUnlocks — element values (P1)
	if len(restored.Meta.PermanentUnlocks) != len(original.Meta.PermanentUnlocks) {
		t.Fatalf("Meta.PermanentUnlocks len: got %d, want %d", len(restored.Meta.PermanentUnlocks), len(original.Meta.PermanentUnlocks))
	}
	for i, want := range original.Meta.PermanentUnlocks {
		if restored.Meta.PermanentUnlocks[i] != want {
			t.Errorf("Meta.PermanentUnlocks[%d]: got %q, want %q", i, restored.Meta.PermanentUnlocks[i], want)
		}
	}
	// Meta.BestCombos — key presence and inner slice values (P1)
	for zone, origCombos := range original.Meta.BestCombos {
		restoredCombos, ok := restored.Meta.BestCombos[zone]
		if !ok {
			t.Errorf("Meta.BestCombos missing key %q after round-trip", zone)
			continue
		}
		if len(restoredCombos) != len(origCombos) {
			t.Fatalf("Meta.BestCombos[%q] len: got %d, want %d", zone, len(restoredCombos), len(origCombos))
		}
		for i, want := range origCombos {
			if restoredCombos[i] != want {
				t.Errorf("Meta.BestCombos[%q][%d]: got %q, want %q", zone, i, restoredCombos[i], want)
			}
		}
	}
	// Meta.RunCount
	if restored.Meta.RunCount != original.Meta.RunCount {
		t.Errorf("Meta.RunCount: got %d, want %d", restored.Meta.RunCount, original.Meta.RunCount)
	}
	// Meta.UnlockedPerks — element values (P1)
	if len(restored.Meta.UnlockedPerks) != len(original.Meta.UnlockedPerks) {
		t.Fatalf("Meta.UnlockedPerks len: got %d, want %d", len(restored.Meta.UnlockedPerks), len(original.Meta.UnlockedPerks))
	}
	for i, want := range original.Meta.UnlockedPerks {
		if restored.Meta.UnlockedPerks[i] != want {
			t.Errorf("Meta.UnlockedPerks[%d]: got %q, want %q", i, restored.Meta.UnlockedPerks[i], want)
		}
	}
}

// TestPendingDraftNullWhenNil verifies AC4: nil PendingDraft → JSON null, NOT []
func TestPendingDraftNullWhenNil(t *testing.T) {
	state := NewGameState()

	if state.Run.PendingDraft != nil {
		t.Errorf("NewGameState().Run.PendingDraft should be nil, got %v", state.Run.PendingDraft)
	}

	data, err := json.Marshal(state.Run)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	jsonStr := string(data)
	if strings.Contains(jsonStr, `"pending_draft":[]`) {
		t.Errorf("PendingDraft serialized as [] instead of null: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"pending_draft":null`) {
		t.Errorf("PendingDraft not serialized as null: %s", jsonStr)
	}
}

// TestNilSlicesSerializeAsNull verifies that all nil slice fields produce JSON null,
// not []. Covers ComboQueue, ActivePerks, PermanentUnlocks, UnlockedPerks (P2).
func TestNilSlicesSerializeAsNull(t *testing.T) {
	state := NewGameState()

	// Verify all slice fields are nil on a fresh state
	if state.Run.ComboQueue != nil {
		t.Errorf("NewGameState().Run.ComboQueue should be nil, got %v", state.Run.ComboQueue)
	}
	if state.Run.ActivePerks != nil {
		t.Errorf("NewGameState().Run.ActivePerks should be nil, got %v", state.Run.ActivePerks)
	}
	if state.Meta.PermanentUnlocks != nil {
		t.Errorf("NewGameState().Meta.PermanentUnlocks should be nil, got %v", state.Meta.PermanentUnlocks)
	}
	if state.Meta.UnlockedPerks != nil {
		t.Errorf("NewGameState().Meta.UnlockedPerks should be nil, got %v", state.Meta.UnlockedPerks)
	}

	runData, err := json.Marshal(state.Run)
	if err != nil {
		t.Fatalf("json.Marshal(Run) failed: %v", err)
	}
	metaData, err := json.Marshal(state.Meta)
	if err != nil {
		t.Fatalf("json.Marshal(Meta) failed: %v", err)
	}

	runJSON := string(runData)
	metaJSON := string(metaData)

	nullChecks := []struct {
		label  string
		json   string
		field  string
	}{
		{"Run.ComboQueue", runJSON, `"combo_queue":null`},
		{"Run.ActivePerks", runJSON, `"active_perks":null`},
		{"Meta.PermanentUnlocks", metaJSON, `"permanent_unlocks":null`},
		{"Meta.UnlockedPerks", metaJSON, `"unlocked_perks":null`},
	}
	for _, tc := range nullChecks {
		if !strings.Contains(tc.json, tc.field) {
			t.Errorf("%s: want %s in JSON, got: %s", tc.label, tc.field, tc.json)
		}
		// Also assert no empty-array serialization
		badField := strings.Replace(tc.field, "null", "[]", 1)
		if strings.Contains(tc.json, badField) {
			t.Errorf("%s: serialized as [] instead of null: %s", tc.label, tc.json)
		}
	}
}

// TestResourcesMapType verifies AC3: Resources is map[string]float64, round-trips correctly
func TestResourcesMapType(t *testing.T) {
	state := NewGameState()

	state.Run.Resources[content.CPUCycles] = 100.5
	state.Run.Resources[content.MemoryShards] = 50.0
	state.Run.Resources[content.ProcessThreads] = 10.25

	data, err := json.Marshal(state.Run)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var restored RunState
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	cases := map[string]float64{
		content.CPUCycles:      100.5,
		content.MemoryShards:   50.0,
		content.ProcessThreads: 10.25,
	}
	for key, want := range cases {
		if got := restored.Resources[key]; got != want {
			t.Errorf("Resources[%q]: got %v, want %v", key, got, want)
		}
	}
}

// TestSavedAtRFC3339 verifies SavedAt serializes as a valid RFC3339Nano timestamp (P3).
func TestSavedAtRFC3339(t *testing.T) {
	now := time.Date(2026, 5, 3, 12, 0, 0, 0, time.UTC)
	state := NewGameState()
	state.SavedAt = now

	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	// Extract the saved_at value from JSON and parse it as RFC3339Nano
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal into raw map failed: %v", err)
	}
	savedAtRaw, ok := raw["saved_at"]
	if !ok {
		t.Fatal("JSON missing 'saved_at' key")
	}
	var savedAtStr string
	if err := json.Unmarshal(savedAtRaw, &savedAtStr); err != nil {
		t.Fatalf("saved_at is not a JSON string: %v", err)
	}
	parsed, err := time.Parse(time.RFC3339Nano, savedAtStr)
	if err != nil {
		t.Errorf("saved_at %q is not valid RFC3339Nano: %v", savedAtStr, err)
	}
	if !parsed.Equal(now) {
		t.Errorf("saved_at round-trip: got %v, want %v", parsed, now)
	}
}
