package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/ricky-irfandi/fsct/internal/report"
)

type Formatter interface {
	Format(results []report.Finding, summary report.Summary) ([]byte, error)
	GetExtension() string
}

type JSONFormatter struct{}

func (f *JSONFormatter) Format(results []report.Finding, summary report.Summary) ([]byte, error) {
	type OutputReport struct {
		Version   string           `json:"version"`
		Timestamp string           `json:"timestamp"`
		Summary   report.Summary   `json:"summary"`
		Findings  []report.Finding `json:"findings"`
	}

	output := OutputReport{
		Version:   "1.0.0",
		Timestamp: time.Now().Format(time.RFC3339),
		Summary:   summary,
		Findings:  results,
	}

	return json.MarshalIndent(output, "", "  ")
}

func (f *JSONFormatter) GetExtension() string {
	return "json"
}

type YAMLFormatter struct{}

func (f *YAMLFormatter) Format(results []report.Finding, summary report.Summary) ([]byte, error) {
	output := fmt.Sprintf(`version: "1.0.0"
timestamp: "%s"
summary:
  high: %d
  warning: %d
  info: %d
  passed: %d
findings:
`, time.Now().Format(time.RFC3339), summary.High, summary.Warning, summary.Info, summary.Passed)

	for _, finding := range results {
		output += fmt.Sprintf(`  - id: "%s"
    severity: "%s"
    title: "%s"
    message: "%s"
    file: "%s"
    line: %d
    suggestion: "%s"
`, finding.ID, finding.Severity, finding.Title, finding.Message, finding.File, finding.Line, finding.Suggestion)
	}

	return []byte(output), nil
}

func (f *YAMLFormatter) GetExtension() string {
	return "yaml"
}

type HTMLFormatter struct{}

func (f *HTMLFormatter) Format(results []report.Finding, summary report.Summary) ([]byte, error) {
	const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FSCT Report</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        header { background: #2196F3; color: white; padding: 20px; border-radius: 8px 8px 0 0; }
        h1 { margin: 0; font-size: 24px; }
        .summary { display: flex; gap: 20px; padding: 20px; background: #fafafa; border-bottom: 1px solid #eee; }
        .stat { text-align: center; padding: 10px 20px; background: white; border-radius: 4px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
        .stat-value { font-size: 32px; font-weight: bold; }
        .stat-label { font-size: 12px; color: #666; text-transform: uppercase; }
        .high { color: #f44336; }
        .warning { color: #ff9800; }
        .info { color: #2196F3; }
        .passed { color: #4CAF50; }
        .findings { padding: 20px; }
        .finding { padding: 15px; margin-bottom: 10px; border-radius: 4px; border-left: 4px solid #ddd; background: #fafafa; }
        .finding.high { border-left-color: #f44336; }
        .finding.warning { border-left-color: #ff9800; }
        .finding.info { border-left-color: #2196F3; }
        .finding-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 5px; }
        .finding-id { font-weight: bold; color: #666; }
        .finding-severity { padding: 2px 8px; border-radius: 12px; font-size: 12px; font-weight: bold; text-transform: uppercase; }
        .severity-high { background: #ffebee; color: #c62828; }
        .severity-warning { background: #fff3e0; color: #ef6c00; }
        .severity-info { background: #e3f2fd; color: #1565c0; }
        .finding-title { font-weight: 600; margin-bottom: 5px; }
        .finding-message { color: #666; margin-bottom: 5px; }
        .finding-file { font-size: 12px; color: #999; }
        .finding-suggestion { margin-top: 10px; padding: 10px; background: white; border-radius: 4px; font-size: 13px; }
        .finding-suggestion strong { color: #4CAF50; }
        .passed-message { text-align: center; padding: 40px; color: #4CAF50; font-size: 18px; }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>FSCT Report</h1>
            <p>Generated at {{.Timestamp}}</p>
        </header>
        <div class="summary">
            <div class="stat high">
                <div class="stat-value">{{.Summary.High}}</div>
                <div class="stat-label">High</div>
            </div>
            <div class="stat warning">
                <div class="stat-value">{{.Summary.Warning}}</div>
                <div class="stat-label">Warning</div>
            </div>
            <div class="stat info">
                <div class="stat-value">{{.Summary.Info}}</div>
                <div class="stat-label">Info</div>
            </div>
            <div class="stat passed">
                <div class="stat-value">{{.Summary.Passed}}</div>
                <div class="stat-label">Passed</div>
            </div>
        </div>
        <div class="findings">
            {{if .Findings}}
                {{range .Findings}}
                <div class="finding {{.SeverityClass}}">
                    <div class="finding-header">
                        <span class="finding-id">{{.ID}}</span>
                        <span class="finding-severity severity-{{.SeverityClass}}">{{.Severity}}</span>
                    </div>
                    <div class="finding-title">{{.Title}}</div>
                    <div class="finding-message">{{.Message}}</div>
                    <div class="finding-file">{{.File}}{{if .Line}}:{{.Line}}{{end}}</div>
                    {{if .Suggestion}}
                    <div class="finding-suggestion"><strong>Suggestion:</strong> {{.Suggestion}}</div>
                    {{end}}
                </div>
                {{end}}
            {{else}}
                <div class="passed-message">No issues found! All checks passed.</div>
            {{end}}
        </div>
    </div>
</body>
</html>
`

	type FormattedFinding struct {
		ID            string
		Severity      string
		SeverityClass string
		Title         string
		Message       string
		File          string
		Line          int
		Suggestion    string
	}

	type OutputData struct {
		Timestamp string
		Summary   report.Summary
		Findings  []FormattedFinding
	}

	findings := make([]FormattedFinding, 0, len(results))
	for _, f := range results {
		severityClass := "info"
		if f.Severity == report.SeverityHigh {
			severityClass = "high"
		} else if f.Severity == report.SeverityWarning {
			severityClass = "warning"
		}
		findings = append(findings, FormattedFinding{
			ID:            f.ID,
			Severity:      string(f.Severity),
			SeverityClass: severityClass,
			Title:         f.Title,
			Message:       f.Message,
			File:          f.File,
			Line:          f.Line,
			Suggestion:    f.Suggestion,
		})
	}

	data := OutputData{
		Timestamp: time.Now().Format("January 2, 2006 15:04"),
		Summary:   summary,
		Findings:  findings,
	}

	tmpl, err := template.New("html").Parse(htmlTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (f *HTMLFormatter) GetExtension() string {
	return "html"
}

type ConsoleFormatter struct{}

func (f *ConsoleFormatter) Format(results []report.Finding, summary report.Summary) ([]byte, error) {
	output := "FSCT Report\n"
	output += "────────────\n"
	output += fmt.Sprintf("Summary  High %d  |  Warning %d  |  Info %d  |  Passed %d\n\n",
		summary.High,
		summary.Warning,
		summary.Info,
		summary.Passed,
	)

	if len(results) == 0 {
		output += "No issues found. All checks passed.\n"
		return []byte(output), nil
	}

	output += "Findings\n"
	output += "────────\n"
	for _, finding := range results {
		icon := "•"
		if finding.Severity == report.SeverityHigh {
			icon = "×"
		} else if finding.Severity == report.SeverityWarning {
			icon = "!"
		}
		output += fmt.Sprintf("\n%s %s  (%s)\n", icon, finding.ID, strings.ToUpper(string(finding.Severity)))
		output += fmt.Sprintf("  %s\n", finding.Title)
		if finding.Message != "" {
			output += fmt.Sprintf("  %s\n", finding.Message)
		}
		if finding.File != "" {
			output += fmt.Sprintf("  %s", finding.File)
			if finding.Line > 0 {
				output += fmt.Sprintf(":%d", finding.Line)
			}
			output += "\n"
		}
		if finding.Suggestion != "" {
			output += fmt.Sprintf("  Suggestion: %s\n", finding.Suggestion)
		}
	}

	return []byte(output), nil
}

func (f *ConsoleFormatter) GetExtension() string {
	return ""
}

func NewFormatter(format string) Formatter {
	switch format {
	case "json":
		return &JSONFormatter{}
	case "yaml":
		return &YAMLFormatter{}
	case "html":
		return &HTMLFormatter{}
	case "prompt":
		return NewPromptFormatter()
	default:
		return &ConsoleFormatter{}
	}
}

func WriteToFile(data []byte, filename string, extension string) error {
	if filename == "" {
		filename = "fsct-report"
	}

	if extension != "" {
		filename = filename + "." + extension
	}

	return os.WriteFile(filename, data, 0644)
}
