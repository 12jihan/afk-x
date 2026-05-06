package ui

import (
	"time"

	"github.com/12jihan/afk-x/internal/content"
	"github.com/12jihan/afk-x/internal/engine"
	"github.com/12jihan/afk-x/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// Update routes all Bubbletea messages and returns the updated model plus any
// follow-up commands. All state mutations happen here — never in View or Init.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width, m.Height = msg.Width, msg.Height
		if m.Width < styles.MinWidth || m.Height < styles.MinHeight {
			m.tooSmall = true
			return m, tea.Quit
		}

	case bootDoneMsg:
		m.LastTick = time.Now()
		m.Screen = GameScreen
		return m, engine.TickCmd()

	case engine.TickMsg:
		if m.Screen != GameScreen {
			return m, nil
		}
		now := time.Time(msg)
		delta := now.Sub(m.LastTick)
		m.LastTick = now
		m.Rates = engine.ComputeRates(m.State.Run)
		m.State.Run = engine.Tick(m.State.Run, delta, m.Rates)
		return m, engine.TickCmd()

	case clearStatusMsg:
		m.StatusMsg = ""

	case tea.KeyMsg:
		if m.Screen == BootScreen {
			m.LastTick = time.Now()
			m.Screen = GameScreen
			return m, engine.TickCmd()
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1", "2", "3":
			if m.Screen == GameScreen {
				idx := int(msg.String()[0] - '1')
				if idx >= 0 && idx < len(content.Upgrades) {
					def := content.Upgrades[idx]
					updated, ok := engine.PurchaseUpgrade(m.State.Run, def.ID)
					if ok {
						m.State.Run = updated
						m.StatusMsg = "Purchased: " + def.Name
					} else {
						m.StatusMsg = "Insufficient resources for " + def.Name
					}
					return m, clearStatusAfterCmd(2 * time.Second)
				}
			}
		}
	}

	return m, nil
}
