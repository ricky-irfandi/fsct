package interactive

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = size.Width
		m.height = size.Height
		setScreenSize(size.Width)
	}
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
	if m.width <= 0 {
		return m.screen.View()
	}
	return lipgloss.NewStyle().Width(m.width).Render(m.screen.View())
}

func NewMainMenuScreen() Screen {
	return NewMenuModel()
}

func transition(next Screen) (tea.Model, tea.Cmd) {
	if next == nil {
		return nil, nil
	}
	return next, next.Init()
}
