package filter

import (
	"strings"

	"github.com/ricky-irfandi/fsct/internal/report"
)

type Filter struct {
	IgnorePatterns []string
	MinSeverity    report.Severity
	AllowedChecks  map[string]bool
}

func NewFilter() *Filter {
	return &Filter{
		IgnorePatterns: make([]string, 0),
		MinSeverity:    report.SeverityInfo,
		AllowedChecks:  make(map[string]bool),
	}
}

func (f *Filter) ShouldIgnore(file string) bool {
	for _, pattern := range f.IgnorePatterns {
		if strings.Contains(file, pattern) {
			return true
		}
	}
	return false
}

func (f *Filter) ShouldInclude(severity report.Severity, checkID string) bool {
	if f.MinSeverity == report.SeverityHigh && severity != report.SeverityHigh {
		return false
	}

	if len(f.AllowedChecks) > 0 {
		return f.AllowedChecks[checkID]
	}

	return true
}

func (f *Filter) SetIgnorePatterns(patterns []string) {
	f.IgnorePatterns = patterns
}

func (f *Filter) SetMinSeverity(severity string) {
	switch strings.ToLower(severity) {
	case "high":
		f.MinSeverity = report.SeverityHigh
	case "warning":
		f.MinSeverity = report.SeverityWarning
	default:
		f.MinSeverity = report.SeverityInfo
	}
}

func (f *Filter) SetAllowedChecks(checkIDs []string) {
	f.AllowedChecks = make(map[string]bool)
	for _, id := range checkIDs {
		f.AllowedChecks[id] = true
	}
}

type SeverityFilter struct {
	MinSeverity report.Severity
}

func NewSeverityFilter(minSeverity report.Severity) *SeverityFilter {
	return &SeverityFilter{MinSeverity: minSeverity}
}

func (s *SeverityFilter) Filter(findings []report.Finding) []report.Finding {
	filtered := make([]report.Finding, 0)
	for _, finding := range findings {
		if s.shouldInclude(finding.Severity) {
			filtered = append(filtered, finding)
		}
	}
	return filtered
}

func (s *SeverityFilter) shouldInclude(severity report.Severity) bool {
	if s.MinSeverity == report.SeverityHigh {
		return severity == report.SeverityHigh
	}
	if s.MinSeverity == report.SeverityWarning {
		return severity == report.SeverityHigh || severity == report.SeverityWarning
	}
	return true
}
