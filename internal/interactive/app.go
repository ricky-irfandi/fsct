package interactive

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Screen interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

type AppModel struct {
	width  int
	height int
	screen Screen
}

func Run() error {
	m := AppModel{
		screen: NewMainMenuScreen(),
	}

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if err := p.Start(); err != nil {
		return err
	}

	return nil
}

func (m AppModel) Init() tea.Cmd {
	return m.screen.Init()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.Type == tea.KeyCtrlC || key.Type == tea.KeyEsc {
			return m, tea.Quit
		}
	}

	newScreen, cmd := m.screen.Update(msg)
	m.screen = newScreen
	return m, cmd
}

func (m AppModel) View() string {
	return m.screen.View()
}

func NewMainMenuScreen() Screen {
	return NewMenuModel()
}
