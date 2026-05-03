package ui

import (
	"fmt"

	"github.com/12jihan/afk-x/internal/ui/screens"
	"github.com/12jihan/afk-x/internal/ui/styles"
)

// View returns the string rendered to the terminal. It delegates to screen-specific
// render functions. No state mutations happen here.
func (m Model) View() string {
	if m.tooSmall {
		return fmt.Sprintf("Terminal too small: need at least %dx%d (current %dx%d)\n",
			styles.MinWidth, styles.MinHeight, m.Width, m.Height)
	}

	switch m.Screen {
	case BootScreen:
		return screens.BootView(m.Width, m.Height)
	case GameScreen:
		return screens.GameView(m.State.Run.Resources, m.Rates, m.Width, m.Height)
	default:
		return ""
	}
}
