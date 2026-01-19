package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/ricky-irfandi/fsct/internal/report"
)

type DiffResult struct {
	Added   []report.Finding `json:"added"`
	Removed []report.Finding `json:"removed"`
	Changed []ChangedFinding `json:"changed"`
	Summary DiffSummary      `json:"summary"`
}

type ChangedFinding struct {
	Before report.Finding `json:"before"`
	After  report.Finding `json:"after"`
}

type DiffSummary struct {
	Added   int `json:"added"`
	Removed int `json:"removed"`
	Changed int `json:"changed"`
	Same    int `json:"same"`
}

func Diff(oldReport, newReport string) (*DiffResult, error) {
	oldFindings, err := loadFindings(oldReport)
	if err != nil {
		return nil, fmt.Errorf("failed to load old report: %w", err)
	}

	newFindings, err := loadFindings(newReport)
	if err != nil {
		return nil, fmt.Errorf("failed to load new report: %w", err)
	}

	return Compare(oldFindings, newFindings), nil
}

func loadFindings(path string) ([]report.Finding, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var report struct {
		Findings []report.Finding `json:"findings"`
	}

	if err := json.Unmarshal(data, &report); err != nil {
		return nil, err
	}

	return report.Findings, nil
}

func Compare(oldFindings, newFindings []report.Finding) *DiffResult {
	result := &DiffResult{
		Added:   make([]report.Finding, 0),
		Removed: make([]report.Finding, 0),
		Changed: make([]ChangedFinding, 0),
		Summary: DiffSummary{},
	}

	oldMap := make(map[string]report.Finding)
	newMap := make(map[string]report.Finding)

	for _, f := range oldFindings {
		oldMap[generateKey(f)] = f
	}

	for _, f := range newFindings {
		newMap[generateKey(f)] = f
	}

	for key, newF := range newMap {
		if oldF, exists := oldMap[key]; !exists {
			result.Added = append(result.Added, newF)
		} else if oldF.Severity != newF.Severity {
			result.Changed = append(result.Changed, ChangedFinding{
				Before: oldF,
				After:  newF,
			})
		}
	}

	for key, oldF := range oldMap {
		if _, exists := newMap[key]; !exists {
			result.Removed = append(result.Removed, oldF)
		}
	}

	sort.Slice(result.Added, func(i, j int) bool {
		return result.Added[i].ID < result.Added[j].ID
	})

	sort.Slice(result.Removed, func(i, j int) bool {
		return result.Removed[i].ID < result.Removed[j].ID
	})

	sort.Slice(result.Changed, func(i, j int) bool {
		return result.Changed[i].Before.ID < result.Changed[j].Before.ID
	})

	result.Summary.Added = len(result.Added)
	result.Summary.Removed = len(result.Removed)
	result.Summary.Changed = len(result.Changed)

	return result
}

func generateKey(f report.Finding) string {
	return fmt.Sprintf("%s|%s|%s", f.ID, f.File, f.Message)
}

func (d *DiffResult) Format(style string) string {
	switch style {
	case "json":
		return d.formatJSON()
	case "console":
		return d.formatConsole()
	default:
		return d.formatConsole()
	}
}

func (d *DiffResult) formatJSON() string {
	data, _ := json.MarshalIndent(d, "", "  ")
	return string(data)
}

func (d *DiffResult) formatConsole() string {
	output := fmt.Sprintf("\nFSCT Diff Report\n")
	output += fmt.Sprintf("================\n")
	output += fmt.Sprintf("Generated: %s\n\n", time.Now().Format(time.RFC3339))

	output += fmt.Sprintf("Summary:\n")
	output += fmt.Sprintf("  Added:    +%d\n", d.Summary.Added)
	output += fmt.Sprintf("  Removed:  -%d\n", d.Summary.Removed)
	output += fmt.Sprintf("  Changed:  ~%d\n", d.Summary.Changed)

	if len(d.Added) > 0 {
		output += fmt.Sprintf("\nAdded Findings (+%d):\n", len(d.Added))
		for _, f := range d.Added {
			output += fmt.Sprintf("  [+] %s (%s): %s\n", f.ID, f.Severity, f.Title)
		}
	}

	if len(d.Removed) > 0 {
		output += fmt.Sprintf("\nRemoved Findings (-%d):\n", len(d.Removed))
		for _, f := range d.Removed {
			output += fmt.Sprintf("  [-] %s (%s): %s\n", f.ID, f.Severity, f.Title)
		}
	}

	if len(d.Changed) > 0 {
		output += fmt.Sprintf("\nChanged Findings (~%d):\n", len(d.Changed))
		for _, c := range d.Changed {
			output += fmt.Sprintf("  [~] %s: %s -> %s\n", c.Before.ID, c.Before.Severity, c.After.Severity)
		}
	}

	return output
}

func (d *DiffResult) Save(path string) error {
	data, _ := json.MarshalIndent(d, "", "  ")
	return os.WriteFile(path, data, 0644)
}
