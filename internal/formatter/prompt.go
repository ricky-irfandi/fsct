package formatter

import (
	"github.com/ricky-irfandi/fsct/internal/prompt"
	"github.com/ricky-irfandi/fsct/internal/report"
)

// PromptFormatter generates AI prompt output
type PromptFormatter struct {
	templateType prompt.PromptTemplateType
	withHeader   bool
}

// NewPromptFormatter creates a new prompt formatter
func NewPromptFormatter() *PromptFormatter {
	return &PromptFormatter{
		templateType: prompt.TemplateTypeComprehensive,
		withHeader:   true,
	}
}

// SetTemplateType sets the prompt template type
func (f *PromptFormatter) SetTemplateType(t prompt.PromptTemplateType) {
	f.templateType = t
}

// SetWithHeader enables/disables the header section
func (f *PromptFormatter) SetWithHeader(withHeader bool) {
	f.withHeader = withHeader
}

// Format generates the prompt output
func (f *PromptFormatter) Format(results []report.Finding, summary report.Summary) ([]byte, error) {
	// This is a simplified version - in reality, we'd need the project
	// For now, we'll create a basic version that works with findings

	// Create a basic prompt data from findings
	data := prompt.NewPromptData()
	data.TotalChecks = summary.Passed + summary.High + summary.Warning + summary.Info
	data.PassedChecks = summary.Passed

	// Convert findings to summary format
	for _, f := range results {
		category := extractCategoryFromID(f.ID)
		data.AddFinding(prompt.FindingSummary{
			ID:         f.ID,
			Severity:   f.Severity,
			Category:   category,
			Title:      f.Title,
			Message:    f.Message,
			File:       f.File,
			Suggestion: f.Suggestion,
		})
	}

	// Set default app name
	data.AppName = "Flutter App"
	data.Version = "1.0.0"
	data.Description = "Flutter mobile application"

	// Build the prompt
	builder := prompt.NewPromptBuilder()
	builder.SetTemplateType(f.templateType)

	var output string
	var err error

	if f.withHeader {
		output, err = builder.GenerateWithHeader(data)
	} else {
		output, err = builder.Generate(data)
	}

	if err != nil {
		return nil, err
	}

	return []byte(output), nil
}

// GetExtension returns the file extension
func (f *PromptFormatter) GetExtension() string {
	return "md"
}

// Helper to extract category from finding ID
func extractCategoryFromID(id string) string {
	if len(id) < 3 {
		return "Other"
	}

	prefix := id[:3]
	switch prefix {
	case "AND":
		return "Android"
	case "IOS":
		return "iOS"
	case "FLT":
		return "Flutter"
	case "SEC":
		return "Security"
	case "POL":
		return "Policy"
	case "COD":
		return "Code Quality"
	case "TST":
		return "Testing"
	case "LIN":
		return "Linting"
	case "DOC":
		return "Documentation"
	case "PER":
		return "Performance"
	case "REV":
		return "Reviewer"
	case "AI-":
		return "AI Analysis"
	default:
		return "Other"
	}
}
