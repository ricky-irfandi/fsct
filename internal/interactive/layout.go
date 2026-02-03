package interactive

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var screenWidth int

func setScreenSize(width int) {
	if width > 0 {
		screenWidth = width
	}
}

func uiWidth() int {
	return screenWidth
}

func renderHeader(text string) string {
	if uiWidth() > 0 {
		return Styles.Header.Width(uiWidth()).Render(text)
	}
	return Styles.Header.Render(text)
}

func renderFooter(text string) string {
	if uiWidth() > 0 {
		return Styles.Footer.Width(uiWidth()).Render(text)
	}
	return Styles.Footer.Render(text)
}

func padToWidth(content string) string {
	if uiWidth() <= 0 {
		return content
	}
	return lipgloss.NewStyle().Width(uiWidth()).Render(content)
}

func supportsUnicode() bool {
	env := strings.ToLower(os.Getenv("LC_ALL") + os.Getenv("LC_CTYPE") + os.Getenv("LANG"))
	return strings.Contains(env, "utf-8") || strings.Contains(env, "utf8")
}

func cursorGlyph() string {
	if supportsUnicode() {
		return "❯"
	}
	return ">"
}
