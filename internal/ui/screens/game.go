package screens

import (
	"strings"

	"github.com/12jihan/afk-x/internal/content"
	"github.com/12jihan/afk-x/internal/engine"
	"github.com/12jihan/afk-x/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

var resourceOrder = []string{
	content.CPUCycles,
	content.MemoryShards,
	content.ProcessThreads,
}

var resourceLabels = map[string]string{
	content.CPUCycles:      "CPU Cycles",
	content.MemoryShards:   "Memory Shards",
	content.ProcessThreads: "Process Threads",
}

// GameView renders the main game screen: a resource panel showing current values
// for all three base resources, plus a stub floor/run footer.
func GameView(resources map[string]float64, rates engine.ResourceRates, width, height int) string {
	if width == 0 || height == 0 {
		return ""
	}

	var rows []string
	for _, key := range resourceOrder {
		label := resourceLabels[key]
		rows = append(rows, styles.FormatResource(label, resources[key]))
	}

	panel := styles.ResourcePanel(width).Render(strings.Join(rows, "\n"))

	footer := lipgloss.NewStyle().
		Foreground(styles.Muted).
		Render("Floor: 1  Run: 1")

	return panel + "\n" + footer
}
