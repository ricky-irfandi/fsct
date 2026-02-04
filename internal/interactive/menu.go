package interactive

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type MenuItem struct {
	ID          string
	Label       string
	Icon        string
	Description string
	Action      func() (tea.Model, tea.Cmd)
}

type MenuModel struct {
	cursor   int
	items    []MenuItem
	selected map[int]struct{}
	width    int
	height   int
	showHelp bool
}

func NewMenuModel() *MenuModel {
	return &MenuModel{
		cursor:   0,
		items:    MainMenuItems(),
		selected: make(map[int]struct{}),
		showHelp: false,
	}
}

func MainMenuItems() []MenuItem {
	return []MenuItem{
		{ID: "check", Label: "Run Compliance Check", Icon: "🚀", Description: "Analyze your Flutter project", Action: func() (tea.Model, tea.Cmd) {
			return transition(NewCheckWizard())
		}},
		{ID: "checks", Label: "View Available Checks", Icon: "📋", Description: "Browse all available compliance checks", Action: func() (tea.Model, tea.Cmd) {
			return transition(NewChecksScreen())
		}},
		{ID: "checklist", Label: "Generate Pre-Submit Checklist", Icon: "✅", Description: "Store submission readiness", Action: func() (tea.Model, tea.Cmd) {
			return transition(NewChecklistScreen())
		}},
		{ID: "hook", Label: "Install Git Pre-commit Hook", Icon: "🔧", Description: "Auto-check before commits", Action: func() (tea.Model, tea.Cmd) {
			return transition(NewHookScreen())
		}},
		{ID: "cert", Label: "Generate Compliance Certificate", Icon: "📜", Description: "Compliance documentation", Action: func() (tea.Model, tea.Cmd) {
			return transition(NewCertScreen())
		}},
		{ID: "config", Label: "Configure AI Settings", Icon: "🤖", Description: "Setup AI analysis provider", Action: func() (tea.Model, tea.Cmd) {
			return transition(NewConfigScreen())
		}},
		{ID: "exit", Label: "Exit", Icon: "👋", Description: "Quit FSCT", Action: func() (tea.Model, tea.Cmd) {
			return nil, tea.Quit
		}},
	}
}

func (m *MenuModel) Init() tea.Cmd {
	return nil
}

func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, menuKeys.Help):
			m.showHelp = !m.showHelp
		case key.Matches(msg, menuKeys.Up):
			if m.showHelp {
				break
			}
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.items) - 1
			}
		case key.Matches(msg, menuKeys.Down):
			if m.showHelp {
				break
			}
			m.cursor++
			if m.cursor >= len(m.items) {
				m.cursor = 0
			}
		case key.Matches(msg, menuKeys.Enter):
			if m.showHelp {
				m.showHelp = false
				break
			}
			if item, ok := m.getItem(); ok {
				if item.Action != nil {
					return item.Action()
				}
			}
		case key.Matches(msg, menuKeys.Quit):
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *MenuModel) View() string {
	if m.showHelp {
		return m.renderHelpOverlay()
	}

	var s string
	s += renderHeader("FSCT - Flutter Store Compliance Tool")
	s += "\n\n"

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = cursorGlyph()
			s += Styles.MenuItemSelected.Render(cursor + " " + item.Icon + " " + item.Label)
		} else {
			s += Styles.MenuItem.Render(cursor + " " + item.Icon + " " + item.Label)
		}
		s += "\n" + Styles.MenuDescription.Render("  "+item.Description) + "\n"
	}

	s += "\n"
	s += renderFooter("↑↓ Navigate • Enter Select • ? Help • q Quit")

	return padToWidth(s)
}

func (m *MenuModel) renderHelpOverlay() string {
	var s string
	s += renderHeader("Keyboard Shortcuts")
	s += "\n\n"

	shortcuts := []struct {
		key  string
		desc string
	}{
		{"↑ / k", "Move selection up"},
		{"↓ / j", "Move selection down"},
		{"Enter", "Select current item"},
		{"← / h", "Go back (in wizard)"},
		{"?", "Show/hide this help"},
		{"q", "Quit FSCT"},
	}

	for _, sc := range shortcuts {
		s += Styles.MenuItem.Render(fmt.Sprintf("  %-12s  %s", Styles.Key.Render(sc.key), sc.desc))
		s += "\n"
	}

	s += "\n"
	s += Styles.MenuDescription.Render("Press ? to close this help")
	s += "\n\n"
	s += renderFooter("Press ? to close • q Quit")

	return padToWidth(s)
}

func (m *MenuModel) getItem() (MenuItem, bool) {
	if m.cursor < 0 || m.cursor >= len(m.items) {
		return MenuItem{}, false
	}
	return m.items[m.cursor], true
}

var menuKeys = struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Quit  key.Binding
	Back  key.Binding
	Help  key.Binding
}{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("left", "h", "backspace"),
		key.WithHelp("←/h", "go back"),
	),
	Help: key.NewBinding(
		key.WithKeys("?", "f1"),
		key.WithHelp("?", "help"),
	),
}
