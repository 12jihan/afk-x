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
