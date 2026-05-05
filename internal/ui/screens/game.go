package screens

import (
	"fmt"
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

var resourceShortLabels = map[string]string{
	content.CPUCycles:      "CPU",
	content.MemoryShards:   "MEM",
	content.ProcessThreads: "THR",
}

// GameView renders the main game screen: resource panel, upgrade panel, optional
// status message, and a floor/run footer with key hints.
func GameView(resources map[string]float64, rates engine.ResourceRates, upgrades map[string]int, statusMsg string, width, height int) string {
	if width == 0 || height == 0 {
		return ""
	}

	var resourceRows []string
	for _, key := range resourceOrder {
		resourceRows = append(resourceRows, styles.FormatResource(resourceLabels[key], resources[key]))
	}
	resourcePanel := styles.ResourcePanel(width).Render(strings.Join(resourceRows, "\n"))

	upgradePanel := upgradeView(upgrades, resources, width)

	footer := lipgloss.NewStyle().Foreground(styles.Muted).
		Render("Floor: 1  Run: 1   [1-3] buy upgrade  [q] quit")

	parts := []string{resourcePanel, upgradePanel}
	if statusMsg != "" {
		parts = append(parts, lipgloss.NewStyle().Foreground(styles.Normal).Render(statusMsg))
	}
	parts = append(parts, footer)

	return strings.Join(parts, "\n")
}

// upgradeView renders the upgrade panel listing all upgrades with costs and affordability.
func upgradeView(upgrades map[string]int, resources map[string]float64, width int) string {
	var rows []string
	for i, def := range content.Upgrades {
		level := upgrades[def.ID]
		var row string
		if level >= def.MaxLevel {
			row = fmt.Sprintf("[%d] %-20s MAXED", i+1, def.Name)
			row = lipgloss.NewStyle().Foreground(styles.Accent).Render(row)
		} else {
			cost := content.ScaledCost(def, level)
			costStr := formatCost(cost)
			base := fmt.Sprintf("[%d] %-20s Lv.%-2d  Cost: %s", i+1, def.Name, level, costStr)
			if canAfford(resources, cost) {
				row = lipgloss.NewStyle().Foreground(styles.Accent).Render(base)
			} else {
				row = lipgloss.NewStyle().Foreground(styles.Muted).Render(base)
			}
		}
		rows = append(rows, row)
	}
	return styles.UpgradePanel(width).Render(strings.Join(rows, "\n"))
}

func formatCost(cost map[string]float64) string {
	var parts []string
	for _, key := range resourceOrder {
		if v, ok := cost[key]; ok {
			parts = append(parts, fmt.Sprintf("%.0f %s", v, resourceShortLabels[key]))
		}
	}
	return strings.Join(parts, ", ")
}

func canAfford(resources, cost map[string]float64) bool {
	for k, v := range cost {
		if resources[k] < v {
			return false
		}
	}
	return true
}
