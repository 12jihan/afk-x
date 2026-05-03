package styles

import "github.com/charmbracelet/lipgloss"

var (
	Accent = lipgloss.AdaptiveColor{Light: "2", Dark: "10"}  // green
	Muted  = lipgloss.AdaptiveColor{Light: "8", Dark: "240"} // gray
	Normal = lipgloss.AdaptiveColor{Light: "0", Dark: "15"}  // foreground
)
