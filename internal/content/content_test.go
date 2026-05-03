package content_test

import (
	"testing"

	"github.com/12jihan/afk-x/internal/content"
)

// TestResourceKeyConstants verifies AC1: constants are defined, non-empty, and distinct.
func TestResourceKeyConstants(t *testing.T) {
	constants := map[string]string{
		"CPUCycles":      content.CPUCycles,
		"MemoryShards":   content.MemoryShards,
		"ProcessThreads": content.ProcessThreads,
	}

	for name, val := range constants {
		if val == "" {
			t.Errorf("%s is empty string", name)
		}
	}

	// Must be distinct
	if content.CPUCycles == content.MemoryShards {
		t.Errorf("CPUCycles == MemoryShards (%q)", content.CPUCycles)
	}
	if content.CPUCycles == content.ProcessThreads {
		t.Errorf("CPUCycles == ProcessThreads (%q)", content.CPUCycles)
	}
	if content.MemoryShards == content.ProcessThreads {
		t.Errorf("MemoryShards == ProcessThreads (%q)", content.MemoryShards)
	}
}

// TestResourceKeyValues verifies the exact key strings match the JSON save format.
func TestResourceKeyValues(t *testing.T) {
	if content.CPUCycles != "cpu_cycles" {
		t.Errorf("CPUCycles = %q, want %q", content.CPUCycles, "cpu_cycles")
	}
	if content.MemoryShards != "memory_shards" {
		t.Errorf("MemoryShards = %q, want %q", content.MemoryShards, "memory_shards")
	}
	if content.ProcessThreads != "process_threads" {
		t.Errorf("ProcessThreads = %q, want %q", content.ProcessThreads, "process_threads")
	}
}

// TestBaseRatesDefined verifies AC1: BaseRates has an entry for each constant, all > 0.
func TestBaseRatesDefined(t *testing.T) {
	requiredKeys := []string{
		content.CPUCycles,
		content.MemoryShards,
		content.ProcessThreads,
	}
	for _, key := range requiredKeys {
		rate, ok := content.BaseRates[key]
		if !ok {
			t.Errorf("BaseRates missing key %q", key)
			continue
		}
		if rate <= 0 {
			t.Errorf("BaseRates[%q] = %v, want > 0", key, rate)
		}
	}
}

// TestResourceKeysMatchBaseRates verifies every constant appears as a key in BaseRates.
func TestResourceKeysMatchBaseRates(t *testing.T) {
	keys := []string{content.CPUCycles, content.MemoryShards, content.ProcessThreads}
	for _, k := range keys {
		if _, ok := content.BaseRates[k]; !ok {
			t.Errorf("constant %q not found as key in BaseRates", k)
		}
	}
}

// TestBootTextNonEmpty verifies AC2: BootText is a non-empty string.
func TestBootTextNonEmpty(t *testing.T) {
	if content.BootText == "" {
		t.Error("BootText is empty")
	}
}

// TestMilestoneTextForMilestoneFloors verifies AC3: milestone floors return non-empty text.
func TestMilestoneTextForMilestoneFloors(t *testing.T) {
	milestones := []int{5, 10, 20, 25}
	for _, floor := range milestones {
		text := content.MilestoneText(floor)
		if text == "" {
			t.Errorf("MilestoneText(%d) returned empty string, want non-empty", floor)
		}
	}
}

// TestMilestoneTextEmptyForNonMilestone verifies AC3: non-milestone floors return "".
func TestMilestoneTextEmptyForNonMilestone(t *testing.T) {
	nonMilestones := []int{1, 2, 3, 4, 6, 7, 8, 9, 11, 15, 99}
	for _, floor := range nonMilestones {
		text := content.MilestoneText(floor)
		if text != "" {
			t.Errorf("MilestoneText(%d) = %q, want empty string", floor, text)
		}
	}
}

// TestMilestoneTextsAreDistinct verifies each milestone floor has unique flavor text.
func TestMilestoneTextsAreDistinct(t *testing.T) {
	milestones := []int{5, 10, 20, 25}
	seen := make(map[string]int)
	for _, floor := range milestones {
		text := content.MilestoneText(floor)
		if prev, dup := seen[text]; dup {
			t.Errorf("MilestoneText(%d) == MilestoneText(%d) — texts are not distinct", floor, prev)
		}
		seen[text] = floor
	}
}
