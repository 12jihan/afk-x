package ui

import (
	"time"

	"charm.land/bubbles/v2/textinput"
	"github.com/12jihan/afk-x/internal/engine"
	"github.com/12jihan/afk-x/internal/game"
	tea "github.com/charmbracelet/bubbletea"
)

// Screen identifies which screen is currently active.
type Screen int

const (
	BootScreen Screen = iota
	GameScreen
	PerkDraftScreen
	RunSummaryScreen
	MetaProgressionScreen
)

const bootDuration = 2 * time.Second

// bootDoneMsg is delivered by bootTimerCmd after the boot sequence timer elapses.
type bootDoneMsg struct{}

// floorClearMsg is delivered by floorClearCmd when a floor is cleared during a tick.
type floorClearMsg struct{}

// floorClearCmd fires a floorClearMsg to trigger floor advancement in Update.
func floorClearCmd() tea.Cmd {
	return func() tea.Msg { return floorClearMsg{} }
}

// clearStatusMsg is delivered by clearStatusAfterCmd to dismiss the status bar message.
type clearStatusMsg struct{}

// clearStatusAfterCmd schedules a clearStatusMsg after duration d.
func clearStatusAfterCmd(d time.Duration) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(d)
		return clearStatusMsg{}
	}
}

func bootTimerCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(bootDuration)
		return bootDoneMsg{}
	}
}

// Model is the root Bubbletea model. All game state lives here.
type Model struct {
	Screen     Screen
	State      game.GameState
	Rates      engine.ResourceRates
	LastTick   time.Time
	CmdInput   textinput.Model
	ShowCmdRef bool
	Width      int
	Height     int
	StatusMsg  string
	tooSmall   bool
}

// New returns a freshly initialized Model ready for the boot screen.
func New() Model {
	return Model{
		Screen:   BootScreen,
		State:    game.NewGameState(),
		CmdInput: textinput.New(),
	}
}

// Init starts the boot timer. The game tick starts only after the boot screen
// completes to avoid accumulating ticks before the player sees the game.
func (m Model) Init() tea.Cmd {
	return bootTimerCmd()
}
