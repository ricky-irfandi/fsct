package interactive

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type MessageScreen struct {
	message  string
	subtitle string
	done     bool
}

func NewMessageScreen(message, subtitle string) *MessageScreen {
	return &MessageScreen{
		message:  message,
		subtitle: subtitle,
		done:     false,
	}
}

func (m *MessageScreen) Init() tea.Cmd {
	return nil
}

func (m *MessageScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, menuKeys.Enter):
			m.done = true
		case key.Matches(msg, menuKeys.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *MessageScreen) View() string {
	var s string
	s += Styles.Header.Render(m.message)
	s += "\n\n"
	s += Styles.Subtitle.Render(m.subtitle)
	s += "\n\n"
	s += Styles.Footer.Render("Enter to continue • q Quit")
	return s
}

type CheckWizard struct {
	step     int
	path     string
	platform string
	format   string
	severity string
	aiMode   string
}

func NewCheckWizard() *CheckWizard {
	return &CheckWizard{
		step:     0,
		path:     ".",
		platform: "both",
		format:   "console",
		severity: "info",
		aiMode:   "auto",
	}
}

func (m *CheckWizard) Init() tea.Cmd {
	return nil
}

func (m *CheckWizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, menuKeys.Up):
			m.step--
			if m.step < 0 {
				m.step = 0
			}
		case key.Matches(msg, menuKeys.Down):
			m.step++
			if m.step > 4 {
				m.step = 4
			}
		case key.Matches(msg, menuKeys.Enter):
			if m.step == 4 {
				return NewRunCheckScreen(m.path, m.platform, m.format, m.severity, m.aiMode), nil
			}
			m.step++
		case key.Matches(msg, menuKeys.Quit):
			return NewMenuModel(), nil
		}
	}
	return m, nil
}

func (m *CheckWizard) View() string {
	steps := []string{"Project Path", "Platform", "Output Format", "Severity", "AI Options"}
	stepNames := []string{"Select project to analyze", "Choose platforms to check", "Choose output format", "Filter by severity", "Configure AI analysis"}

	var s string
	s += Styles.Header.Render("FSCT Check Wizard")
	s += "\n\n"
	s += Styles.WizardStep.Render(fmt.Sprintf("Step %d/5: %s", m.step+1, steps[m.step]))
	s += "\n\n"
	s += Styles.Label.Render(stepNames[m.step])
	s += "\n\n"

	switch m.step {
	case 0:
		s += Styles.Input.Render("📁 " + m.path)
		s += "\n\n"
		s += Styles.Selector.Render("Current directory")
	case 1:
		platforms := []string{"Both (Android & iOS)", "Android only", "iOS only"}
		for i, p := range platforms {
			prefix := "  "
			if m.platform == strings.ToLower(strings.Split(p, " ")[0]) {
				prefix = "❯"
			}
			if m.step == 1 && i == 0 && m.platform == "both" {
				prefix = "❯"
			}
			if m.step == 1 && i == 1 && m.platform == "android" {
				prefix = "❯"
			}
			if m.step == 1 && i == 2 && m.platform == "ios" {
				prefix = "❯"
			}
			if prefix == "❯" {
				s += Styles.MenuItemSelected.Render(prefix + " " + p)
			} else {
				s += Styles.MenuItem.Render(prefix + " " + p)
			}
			s += "\n"
		}
	case 2:
		formats := []string{"Console", "JSON", "YAML", "HTML", "AI Prompt"}
		for j, f := range formats {
			prefix := "  "
			if m.format == strings.ToLower(f) || (m.format == "console" && f == "Console") {
				prefix = "❯"
			}
			if prefix == "❯" {
				s += Styles.MenuItemSelected.Render(prefix + " " + f)
			} else {
				s += Styles.MenuItem.Render(prefix + " " + f)
			}
			s += "\n"
			_ = j
		}
	case 3:
		severities := []string{"All severities", "Warning & High only", "High only"}
		for i, sev := range severities {
			prefix := "  "
			if (m.severity == "info" && i == 0) || (m.severity == "warning" && i == 1) || (m.severity == "high" && i == 2) {
				prefix = "❯"
			}
			if prefix == "❯" {
				s += Styles.MenuItemSelected.Render(prefix + " " + sev)
			} else {
				s += Styles.MenuItem.Render(prefix + " " + sev)
			}
			s += "\n"
		}
	case 4:
		aiOpts := []string{"Use AI if configured", "Skip AI analysis", "Configure AI..."}
		for i, opt := range aiOpts {
			prefix := "  "
			if (m.aiMode == "auto" && i == 0) || (m.aiMode == "skip" && i == 1) || (m.aiMode == "config" && i == 2) {
				prefix = "❯"
			}
			if prefix == "❯" {
				s += Styles.MenuItemSelected.Render(prefix + " " + opt)
			} else {
				s += Styles.MenuItem.Render(prefix + " " + opt)
			}
			s += "\n"
		}
	}

	s += "\n"
	s += Styles.Footer.Render("↑↓ Navigate • Enter Next/Select • q Back to Menu")

	return s
}

type RunCheckScreen struct {
	path     string
	platform string
	format   string
	severity string
	aiMode   string
	done     bool
}

func NewRunCheckScreen(path, platform, format, severity, aiMode string) *RunCheckScreen {
	return &RunCheckScreen{
		path:     path,
		platform: platform,
		format:   format,
		severity: severity,
		aiMode:   aiMode,
		done:     false,
	}
}

func (m *RunCheckScreen) Init() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("go", "run", "./cmd/fsct", "check", m.path,
			"--platform", m.platform,
			"--format", m.format,
			"--severity", m.severity)
		output, _ := cmd.CombinedOutput()
		return outputMsg{output: string(output)}
	}
}

type outputMsg struct{ output string }

func (m *RunCheckScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case outputMsg:
		m.done = true
	case tea.KeyMsg:
		if key.Matches(msg, menuKeys.Enter) || key.Matches(msg, menuKeys.Quit) {
			return NewMenuModel(), nil
		}
	}
	return m, nil
}

func (m *RunCheckScreen) View() string {
	var s string
	s += Styles.Header.Render("Running Compliance Check")
	s += "\n\n"
	s += Styles.Subtitle.Render(fmt.Sprintf("Analyzing: %s", m.path))
	s += "\n\n"
	s += Styles.Info.Render("Running checks... (output would appear here)")
	s += "\n\n"
	s += Styles.Footer.Render("Press any key to return to menu")
	return s
}

func NewChecksScreen() *MessageScreen {
	return NewMessageScreen("Available Checks", "Lists all 98 compliance checks")
}

func NewChecklistScreen() *MessageScreen {
	return NewMessageScreen("Pre-Submit Checklist", "Generates store submission readiness checklist")
}

func NewHookScreen() *MessageScreen {
	return NewMessageScreen("Git Pre-commit Hook", "Installs FSCT as pre-commit hook")
}

func NewCertScreen() *MessageScreen {
	return NewMessageScreen("Compliance Certificate", "Generates compliance documentation")
}

func NewConfigScreen() *MessageScreen {
	return NewMessageScreen("AI Configuration", "Configure AI analysis provider settings")
}
