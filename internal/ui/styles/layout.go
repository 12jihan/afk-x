package styles

import "github.com/charmbracelet/lipgloss"

const (
	MinWidth  = 80
	MinHeight = 24
)

// ResourcePanel returns a bordered panel style sized to fit within the given terminal width.
// Total horizontal footprint: Width() + 4 padding + 2 border = Width() + 6.
func ResourcePanel(width int) lipgloss.Style {
	inner := width - 6
	if inner < 20 {
		inner = 20
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Accent).
		Padding(1, 2).
		Width(inner)
}

// UpgradePanel returns a bordered panel style for the upgrade list.
func UpgradePanel(width int) lipgloss.Style {
	inner := width - 6
	if inner < 20 {
		inner = 20
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Muted).
		Padding(1, 2).
		Width(inner)
}
