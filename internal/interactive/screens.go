package interactive

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const tickInterval = 200 * time.Millisecond

type wizardStep int

const (
	stepPath wizardStep = iota
	stepPlatform
	stepFormat
	stepSeverity
	stepAI
	stepConfirm
)

type CheckWizard struct {
	step        wizardStep
	stepCursor  int
	stepCursors []int
	path        string
	platform    string
	format      string
	severity    string
	aiMode      string
	pathInput   textinput.Model
}

func NewCheckWizard() *CheckWizard {
	pathInput := textinput.New()
	pathInput.Placeholder = "/path/to/project"
	pathInput.CharLimit = 256
	pathInput.Width = 48
	pathInput.Prompt = ""
	pathInput.TextStyle = Styles.Input
	pathInput.PlaceholderStyle = Styles.MenuDescription
	pathInput.CursorStyle = Styles.Info

	return &CheckWizard{
		step:        stepPath,
		stepCursor:  0,
		stepCursors: make([]int, 6),
		path:        ".",
		platform:    "both",
		format:      "console",
		severity:    "info",
		aiMode:      "auto",
		pathInput:   pathInput,
	}
}

func (m *CheckWizard) Init() tea.Cmd {
	return nil
}

func (m *CheckWizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return transition(NewMenuModel())
		}
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return transition(NewMenuModel())
		case tea.KeyLeft, tea.KeyBackspace:
			if !m.pathInputActive() {
				m.retreatStep()
			}
		case tea.KeyEnter:
			if m.step == stepConfirm {
				if m.stepCursor == 1 {
					m.retreatStep()
					return m, nil
				}
				next := NewRunCheckScreen(m.path, m.platform, m.format, m.severity, m.aiMode)
				return transition(next)
			}
			if m.step == stepAI && m.aiMode == "config" {
				return transition(NewConfigScreen())
			}
			if m.step == stepPath && m.stepCursor == 1 {
				customPath := strings.TrimSpace(m.pathInput.Value())
				if customPath == "" {
					return m, nil
				}
				m.path = customPath
			}
			m.advanceStep()
		case tea.KeyUp:
			m.handleUp()
		case tea.KeyDown:
			m.handleDown()
		}
		if msg.String() == "h" {
			if !m.pathInputActive() {
				m.retreatStep()
			}
		}
	}
	if m.pathInputActive() {
		var cmd tea.Cmd
		m.pathInput, cmd = m.pathInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *CheckWizard) handleUp() {
	if m.stepCursor > 0 {
		m.stepCursor--
	}
	m.updateSelection()
}

func (m *CheckWizard) handleDown() {
	max := m.getMaxCursor()
	if m.stepCursor < max {
		m.stepCursor++
	}
	m.updateSelection()
}

func (m *CheckWizard) getMaxCursor() int {
	switch m.step {
	case stepPath:
		return 1
	case stepPlatform:
		return 2
	case stepFormat:
		return 4
	case stepSeverity:
		return 2
	case stepAI:
		return 2
	case stepConfirm:
		return 1
	}
	return 0
}

func (m *CheckWizard) advanceStep() {
	if m.step >= stepConfirm {
		return
	}
	m.stepCursors[m.step] = m.stepCursor
	if m.step == stepSeverity && m.format == "prompt" {
		m.step = stepConfirm
		m.stepCursor = m.stepCursors[m.step]
		m.updateSelection()
		return
	}
	m.step++
	m.stepCursor = m.stepCursors[m.step]
	m.updateSelection()
}

func (m *CheckWizard) retreatStep() {
	if m.step <= stepPath {
		return
	}
	m.stepCursors[m.step] = m.stepCursor
	m.step--
	m.stepCursor = m.stepCursors[m.step]
	m.updateSelection()
}

func (m *CheckWizard) pathInputActive() bool {
	return m.step == stepPath && m.stepCursor == 1
}

func (m *CheckWizard) updateSelection() {
	switch m.step {
	case stepPath:
		if m.stepCursor == 0 {
			m.path = "."
			m.pathInput.Blur()
			m.pathInput.TextStyle = Styles.Input
		} else {
			m.pathInput.Focus()
			m.pathInput.TextStyle = Styles.InputFocused
		}
	case stepPlatform:
		platforms := []string{"both", "android", "ios"}
		m.platform = platforms[m.stepCursor]
	case stepFormat:
		formats := []string{"console", "json", "yaml", "html", "prompt"}
		m.format = formats[m.stepCursor]
		if m.format == "prompt" {
			m.aiMode = "skip"
		}
	case stepSeverity:
		severities := []string{"info", "warning", "high"}
		m.severity = severities[m.stepCursor]
	case stepAI:
		aiModes := []string{"auto", "skip", "config"}
		m.aiMode = aiModes[m.stepCursor]
	}
}

func (m *CheckWizard) View() string {
	var s string
	s += renderHeader("FSCT Check Wizard")
	s += "\n\n"
	s += m.renderProgress()
	s += "\n\n"
	s += m.renderStep()
	s += "\n\n"
	s += renderFooter("↑↓ Choose • Enter Continue • ← Back • q Menu")
	return padToWidth(s)
}

func (m *CheckWizard) renderProgress() string {
	steps := []string{"Path", "Platform", "Format", "Severity", "AI", "Run"}
	stepLabel := steps[int(m.step)]
	return fmt.Sprintf("%s %s",
		Styles.WizardProgress.Render(fmt.Sprintf("Step %d of %d", m.step+1, len(steps))),
		Styles.WizardStep.Render("• "+stepLabel),
	)
}

func (m *CheckWizard) renderStep() string {
	switch m.step {
	case stepPath:
		return m.renderPathStep()
	case stepPlatform:
		return m.renderPlatformStep()
	case stepFormat:
		return m.renderFormatStep()
	case stepSeverity:
		return m.renderSeverityStep()
	case stepAI:
		return m.renderAIStep()
	case stepConfirm:
		return m.renderConfirmStep()
	}
	return ""
}

func (m *CheckWizard) renderPathStep() string {
	var s string
	s += Styles.WizardTitle.Render("Project Path")
	s += "\n\n"

	options := []struct {
		label string
		path  string
		desc  string
	}{
		{"Current Directory", ".", "Use current working directory"},
		{"Enter Custom Path", "custom", "Type a specific path"},
	}

	for i, opt := range options {
		cursor := " "
		if m.stepCursor == i {
			cursor = cursorGlyph()
			s += Styles.MenuItemSelected.Render(fmt.Sprintf("%s %s", cursor, opt.label))
		} else {
			s += Styles.MenuItem.Render(fmt.Sprintf("%s %s", cursor, opt.label))
		}
		s += "\n" + Styles.MenuDescription.Render(fmt.Sprintf("  %s", opt.desc)) + "\n"
	}

	if m.stepCursor == 0 {
		m.path = "."
	}

	if m.stepCursor == 1 {
		s += "\n"
		s += Styles.Label.Render("Path")
		s += "\n"
		if m.pathInputActive() {
			s += Styles.InputFocused.Render(m.pathInput.View())
		} else {
			s += Styles.Input.Render(m.pathInput.View())
		}
	}

	s += "\n\n"
	selected := m.path
	if m.stepCursor == 1 && strings.TrimSpace(m.pathInput.Value()) != "" {
		selected = m.pathInput.Value()
	}
	s += Styles.Info.Render(fmt.Sprintf("Selected: %s", selected))
	return s
}

func (m *CheckWizard) renderPlatformStep() string {
	var s string
	s += Styles.WizardTitle.Render("Platform")
	s += "\n\n"

	options := []struct {
		label string
		desc  string
	}{
		{"Both (Android & iOS)", "Check both Android and iOS platforms"},
		{"Android only", "Check only Android platform"},
		{"iOS only", "Check only iOS platform"},
	}

	for i, opt := range options {
		cursor := " "
		if m.stepCursor == i {
			cursor = cursorGlyph()
			s += Styles.MenuItemSelected.Render(fmt.Sprintf("%s %s", cursor, opt.label))
		} else {
			s += Styles.MenuItem.Render(fmt.Sprintf("%s %s", cursor, opt.label))
		}
		s += "\n" + Styles.MenuDescription.Render(fmt.Sprintf("  %s", opt.desc)) + "\n"
	}

	s += "\n"
	s += Styles.Info.Render(fmt.Sprintf("Selected: %s", strings.Title(m.platform)))
	return s
}

func (m *CheckWizard) renderFormatStep() string {
	var s string
	s += Styles.WizardTitle.Render("Output Format")
	s += "\n\n"

	options := []struct {
		label string
		desc  string
	}{
		{"Console", "Display results in terminal"},
		{"JSON", "Output as JSON format"},
		{"YAML", "Output as YAML format"},
		{"HTML", "Generate HTML report"},
		{"AI Prompt", "Generate AI analysis prompt"},
	}

	for i, opt := range options {
		cursor := " "
		if m.stepCursor == i {
			cursor = cursorGlyph()
			s += Styles.MenuItemSelected.Render(fmt.Sprintf("%s %s", cursor, opt.label))
		} else {
			s += Styles.MenuItem.Render(fmt.Sprintf("%s %s", cursor, opt.label))
		}
		s += "\n" + Styles.MenuDescription.Render(fmt.Sprintf("  %s", opt.desc)) + "\n"
	}

	s += "\n"
	s += Styles.Info.Render(fmt.Sprintf("Selected: %s", strings.Title(m.format)))
	return s
}

func (m *CheckWizard) renderSeverityStep() string {
	var s string
	s += Styles.WizardTitle.Render("Severity Filter")
	s += "\n\n"

	options := []struct {
		label string
		desc  string
	}{
		{"All severities", "Show Info, Warning, and High severity issues"},
		{"Warning & High only", "Filter out Info level issues"},
		{"High only", "Show only critical issues"},
	}

	for i, opt := range options {
		cursor := " "
		if m.stepCursor == i {
			cursor = cursorGlyph()
			s += Styles.MenuItemSelected.Render(fmt.Sprintf("%s %s", cursor, opt.label))
		} else {
			s += Styles.MenuItem.Render(fmt.Sprintf("%s %s", cursor, opt.label))
		}
		s += "\n" + Styles.MenuDescription.Render(fmt.Sprintf("  %s", opt.desc)) + "\n"
	}

	s += "\n"
	s += Styles.Info.Render(fmt.Sprintf("Selected: %s", strings.Title(m.severity)))
	return s
}

func (m *CheckWizard) renderAIStep() string {
	var s string
	s += Styles.WizardTitle.Render("AI Analysis")
	s += "\n\n"

	options := []struct {
		label string
		desc  string
	}{
		{"Use AI if configured", "Enable AI analysis when API key is set"},
		{"Skip AI analysis", "Run only static checks"},
		{"Configure AI Settings", "Set up AI provider configuration"},
	}

	for i, opt := range options {
		cursor := " "
		if m.stepCursor == i {
			cursor = cursorGlyph()
			s += Styles.MenuItemSelected.Render(fmt.Sprintf("%s %s", cursor, opt.label))
		} else {
			s += Styles.MenuItem.Render(fmt.Sprintf("%s %s", cursor, opt.label))
		}
		s += "\n" + Styles.MenuDescription.Render(fmt.Sprintf("  %s", opt.desc)) + "\n"
	}

	s += "\n"
	aiLabel := "Use if configured"
	if m.aiMode == "skip" {
		aiLabel = "Skip"
	} else if m.aiMode == "config" {
		aiLabel = "Configure"
	}
	s += Styles.Info.Render(fmt.Sprintf("Selected: %s", aiLabel))
	return s
}

func (m *CheckWizard) renderConfirmStep() string {
	var s string
	s += Styles.WizardTitle.Render("Ready to Run")
	s += "\n\n"

	s += Styles.MenuDescription.Render("Configuration Summary:")
	s += "\n\n"

	summary := fmt.Sprintf("%s • %s • %s • %s • %s",
		m.path,
		strings.Title(m.platform),
		strings.Title(m.format),
		strings.Title(m.severity),
		strings.Title(m.aiMode),
	)
	s += Styles.Info.Render("Summary: ")
	s += Styles.MenuItem.Render(summary)
	s += "\n\n"

	s += Styles.MenuItem.Render(fmt.Sprintf("  Path:     %s", m.path))
	s += "\n"
	s += Styles.MenuItem.Render(fmt.Sprintf("  Platform: %s", strings.Title(m.platform)))
	s += "\n"
	s += Styles.MenuItem.Render(fmt.Sprintf("  Format:   %s", strings.Title(m.format)))
	s += "\n"
	s += Styles.MenuItem.Render(fmt.Sprintf("  Severity: %s", strings.Title(m.severity)))
	s += "\n"
	s += Styles.MenuItem.Render(fmt.Sprintf("  AI Mode:  %s", strings.Title(m.aiMode)))
	s += "\n\n"

	options := []struct {
		label string
		desc  string
	}{
		{"Run Check", "Execute compliance check with these settings"},
		{"Back", "Return to previous steps"},
	}

	for i, opt := range options {
		cursor := " "
		if m.stepCursor == i {
			cursor = cursorGlyph()
			s += Styles.MenuItemSelected.Render(fmt.Sprintf("%s %s", cursor, opt.label))
		} else {
			s += Styles.MenuItem.Render(fmt.Sprintf("%s %s", cursor, opt.label))
		}
		s += "\n"
		_ = opt.desc
	}

	return s
}

type RunCheckScreen struct {
	path     string
	platform string
	format   string
	severity string
	aiMode   string
	done     bool
	output   string
	tick     int
}

func NewRunCheckScreen(path, platform, format, severity, aiMode string) *RunCheckScreen {
	return &RunCheckScreen{
		path:     path,
		platform: platform,
		format:   format,
		severity: severity,
		aiMode:   aiMode,
		done:     false,
		output:   "",
		tick:     0,
	}
}

func (m *RunCheckScreen) Init() tea.Cmd {
	return tea.Batch(tickCmd(), runCheckCmd(m.path, m.platform, m.format, m.severity, m.aiMode))
}

func tickCmd() tea.Cmd {
	return tea.Tick(tickInterval, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

type tickMsg struct{}

type outputMsg struct{ output string }
type runCheckErrorMsg struct{ err error }

func (m *RunCheckScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if !m.done {
			m.tick++
			return m, tickCmd()
		}
	case outputMsg:
		m.done = true
		m.output = msg.output
	case runCheckErrorMsg:
		m.done = true
		m.output = fmt.Sprintf("Error: %v", msg.err)
	case tea.KeyMsg:
		if m.done {
			return transition(NewMenuModel())
		}
	}
	return m, nil
}

func (m *RunCheckScreen) View() string {
	var s string
	s += renderHeader("Running Compliance Check")
	s += "\n\n"
	s += Styles.Subtitle.Render(fmt.Sprintf("Analyzing: %s", m.path))
	s += "\n\n"

	if !m.done {
		spinner := []string{"-", "\\", "|", "/"}
		if supportsUnicode() {
			spinner = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠸"}
		}
		frame := spinner[m.tick%len(spinner)]
		s += Styles.Info.Render(fmt.Sprintf("%s Running compliance checks...", frame))
		s += "\n\n"
		s += Styles.MenuDescription.Render("This may take a few moments...")
	} else {
		if len(m.output) > 500 {
			s += Styles.MenuItem.Render(m.output[:500])
			s += "\n"
			s += Styles.MenuDescription.Render("... (output truncated)")
		} else {
			s += Styles.MenuItem.Render(m.output)
		}
		s += "\n\n"
	}

	s += renderFooter("Press any key to return to menu")
	return padToWidth(s)
}

type MessageScreen struct {
	message   string
	subtitle  string
	done      bool
	showAbout bool
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
			return transition(NewMenuModel())
		case key.Matches(msg, menuKeys.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *MessageScreen) View() string {
	var s string
	s += renderHeader(m.message)
	s += "\n\n"
	s += Styles.Subtitle.Render(m.subtitle)
	s += "\n\n"
	s += renderFooter("Enter to continue • q Quit")
	return padToWidth(s)
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
