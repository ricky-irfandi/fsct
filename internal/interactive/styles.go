package interactive

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	ColorBackground     = "#1a1a1a"
	ColorSurface        = "#24283b"
	ColorSurfaceLight   = "#2a2f45"
	ColorPrimary        = "#d4a574"
	ColorPrimaryLight   = "#e6b88a"
	ColorPrimaryDark    = "#b8935f"
	ColorSecondary      = "#7aa2f7"
	ColorSecondaryLight = "#8ab4f7"
	ColorSuccess        = "#9ece6a"
	ColorWarning        = "#e0af68"
	ColorError          = "#f7768e"
	ColorInfo           = "#7aa2f7"
	ColorTextPrimary    = "#c0caf5"
	ColorTextSecondary  = "#565f89"
	ColorTextMuted      = "#414868"
)

var Styles = struct {
	App              lipgloss.Style
	Header           lipgloss.Style
	Footer           lipgloss.Style
	Menu             lipgloss.Style
	MenuItem         lipgloss.Style
	MenuItemSelected lipgloss.Style
	MenuItemActive   lipgloss.Style
	MenuDescription  lipgloss.Style
	Wizard           lipgloss.Style
	WizardStep       lipgloss.Style
	WizardProgress   lipgloss.Style
	WizardTitle      lipgloss.Style
	Input            lipgloss.Style
	InputFocused     lipgloss.Style
	Label            lipgloss.Style
	Selector         lipgloss.Style
	SelectorSelected lipgloss.Style
	Button           lipgloss.Style
	ButtonPrimary    lipgloss.Style
	ButtonSecondary  lipgloss.Style
	Success          lipgloss.Style
	Warning          lipgloss.Style
	Error            lipgloss.Style
	Info             lipgloss.Style
	Help             lipgloss.Style
	Key              lipgloss.Style
	Separator        lipgloss.Style
	Title            lipgloss.Style
	Subtitle         lipgloss.Style
}{
	App: lipgloss.NewStyle().
		Background(lipgloss.Color(ColorBackground)).
		Padding(1, 2),

	Header: lipgloss.NewStyle().
		Background(lipgloss.Color(ColorSurface)).
		Foreground(lipgloss.Color(ColorSecondary)).
		Padding(0, 2).
		Bold(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color(ColorTextMuted)),

	Footer: lipgloss.NewStyle().
		Background(lipgloss.Color(ColorSurface)).
		Foreground(lipgloss.Color(ColorTextSecondary)).
		Padding(0, 2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderForeground(lipgloss.Color(ColorTextMuted)),

	Menu: lipgloss.NewStyle().
		Padding(1, 0),

	MenuItem: lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2).
		PaddingTop(1).
		PaddingBottom(1).
		Foreground(lipgloss.Color(ColorTextPrimary)),

	MenuItemSelected: lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2).
		PaddingTop(1).
		PaddingBottom(1).
		Foreground(lipgloss.Color(ColorBackground)).
		Background(lipgloss.Color(ColorPrimary)).
		Bold(true),

	MenuItemActive: lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2).
		PaddingTop(1).
		PaddingBottom(1).
		Foreground(lipgloss.Color(ColorBackground)).
		Background(lipgloss.Color(ColorPrimaryDark)).
		Bold(true),

	MenuDescription: lipgloss.NewStyle().
		PaddingLeft(4).
		Foreground(lipgloss.Color(ColorTextSecondary)).
		Italic(true),

	Wizard: lipgloss.NewStyle().
		Padding(1, 2),

	WizardStep: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSecondary)).
		Bold(true),

	WizardProgress: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorPrimary)),

	WizardTitle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextPrimary)).
		Bold(true).
		MarginBottom(1),

	Input: lipgloss.NewStyle().
		Background(lipgloss.Color(ColorSurface)).
		Foreground(lipgloss.Color(ColorTextPrimary)).
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(ColorTextMuted)),

	InputFocused: lipgloss.NewStyle().
		Background(lipgloss.Color(ColorSurface)).
		Foreground(lipgloss.Color(ColorTextPrimary)).
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(ColorSecondary)),

	Label: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextPrimary)).
		Bold(true).
		MarginBottom(1),

	Selector: lipgloss.NewStyle().
		PaddingLeft(2).
		Foreground(lipgloss.Color(ColorTextSecondary)),

	SelectorSelected: lipgloss.NewStyle().
		PaddingLeft(2).
		Foreground(lipgloss.Color(ColorPrimary)).
		Bold(true),

	Button: lipgloss.NewStyle().
		Background(lipgloss.Color(ColorSurface)).
		Foreground(lipgloss.Color(ColorTextPrimary)).
		Padding(0, 2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(ColorTextMuted)),

	ButtonPrimary: lipgloss.NewStyle().
		Background(lipgloss.Color(ColorPrimary)).
		Foreground(lipgloss.Color(ColorBackground)).
		Padding(0, 2).
		Bold(true),

	ButtonSecondary: lipgloss.NewStyle().
		Background(lipgloss.Color(ColorSurface)).
		Foreground(lipgloss.Color(ColorSecondary)).
		Padding(0, 2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(ColorSecondary)),

	Success: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSuccess)).
		Bold(true),

	Warning: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorWarning)).
		Bold(true),

	Error: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorError)).
		Bold(true),

	Info: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorInfo)).
		Bold(true),

	Help: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextSecondary)),

	Key: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSecondary)).
		Bold(true),

	Separator: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextMuted)),

	Title: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextPrimary)).
		Bold(true).
		MarginBottom(1),

	Subtitle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextSecondary)).
		Italic(true),
}
