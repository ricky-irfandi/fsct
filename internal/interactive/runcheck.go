package interactive

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func defaultOutputName(format string) string {
	switch format {
	case "prompt":
		return "fsct-prompt"
	case "html":
		return "fsct-report"
	case "json", "yaml":
		return "fsct-report"
	default:
		return "fsct-report"
	}
}

func runCheckCmd(path, platform, format, severity, aiMode string) tea.Cmd {
	return func() tea.Msg {
		exe, err := os.Executable()
		if err != nil {
			return runCheckErrorMsg{err: fmt.Errorf("resolve executable: %w", err)}
		}

		args := []string{
			"check",
			path,
			"--platform", platform,
			"--format", format,
			"--severity", severity,
		}

		switch aiMode {
		case "skip":
			args = append(args, "--offline")
		}
		if format == "prompt" {
			args = append(args, "--offline")
		}
		if format != "console" {
			args = append(args, "--output", defaultOutputName(format))
		}

		cmd := exec.Command(exe, args...)
		var buf bytes.Buffer
		cmd.Stdout = &buf
		cmd.Stderr = &buf
		if err := cmd.Run(); err != nil {
			output := strings.TrimSpace(buf.String())
			if output == "" {
				return runCheckErrorMsg{err: err}
			}
			return outputMsg{output: fmt.Sprintf("%s\n\n(exit: %v)", output, err)}
		}

		output := strings.TrimSpace(buf.String())
		if output == "" {
			output = "✓ Compliance check completed with no output."
		}
		return outputMsg{output: output}
	}
}
