package screens

import (
	"github.com/12jihan/afk-x/internal/content"
	"github.com/12jihan/afk-x/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// BootView renders the boot sequence screen — content.BootText centered in the
// terminal with a footer prompt. No update logic lives here; timer and keypress
// handling are in ui/update.go.
func BootView(width, height int) string {
	if width == 0 || height == 0 {
		return content.BootText
	}

	bootStyle := lipgloss.NewStyle().Foreground(styles.Accent)
	footerStyle := lipgloss.NewStyle().Foreground(styles.Muted)

	block := bootStyle.Render(content.BootText) + "\n\n" + footerStyle.Render("[any key to continue]")
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, block)
}
