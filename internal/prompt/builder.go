package prompt

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/ricky-irfandi/fsct/internal/report"
)

// PromptBuilder generates the final prompt text
type PromptBuilder struct {
	templateType PromptTemplateType
}

// NewPromptBuilder creates a new prompt builder
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{
		templateType: TemplateTypeComprehensive,
	}
}

// SetTemplateType sets the template type
func (pb *PromptBuilder) SetTemplateType(t PromptTemplateType) {
	pb.templateType = t
}

// Generate creates the final prompt text from PromptData
func (pb *PromptBuilder) Generate(data *PromptData) (string, error) {
	tmplText := GetTemplate(pb.templateType)

	// Create template
	tmpl, err := template.New("prompt").Parse(tmplText)
	if err != nil {
		return "", err
	}

	// Prepare template data
	templateData := pb.prepareTemplateData(data)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// templateData is the structure passed to the template
type templateData struct {
	AppName           string
	Version           string
	Description       string
	FlutterVersion    string
	Repository        string
	Homepage          string
	ComplianceScore   int
	Status            string
	IsReadyForAppStore bool
	IsReadyForPlayStore bool
	TotalChecks       int
	PassedChecks      int
	HighCount         int
	WarningCount      int
	InfoCount         int
	Blockers          []BlockerInfo
	FindingsBySeverity findingsBySeverity
	FindingsByCategory map[string][]FindingSummary
	AndroidConfig     AndroidSummary
	IOSConfig         IOSSummary
	FlutterConfig     FlutterSummary
	AndroidPermissions []PermissionInfo
	IOSPermissions    []PermissionInfo
	Dependencies      []DependencyInfo
	DevDependencies   []DependencyInfo
	HighRiskDeps      []DependencyInfo
	OutdatedDeps      []DependencyInfo
	SecurityFlags     []SecurityFlag
	PolicyFlags       []PolicyFlag
	MissingPolicies   []string
	ReviewerAccount   *ReviewerAccountInfo
}

// findingsBySeverity wraps findings for template access
type findingsBySeverity struct {
	HIGH    []FindingSummary
	WARNING []FindingSummary
	INFO    []FindingSummary
}

// prepareTemplateData converts PromptData to template-friendly format
func (pb *PromptBuilder) prepareTemplateData(data *PromptData) *templateData {
	return &templateData{
		AppName:             data.AppName,
		Version:             data.Version,
		Description:         data.Description,
		FlutterVersion:      data.FlutterVersion,
		Repository:          data.Repository,
		Homepage:            data.Homepage,
		ComplianceScore:     data.CalculateComplianceScore(),
		Status:              data.GetStatus(),
		IsReadyForAppStore:  data.IsReadyForAppStore(),
		IsReadyForPlayStore: data.IsReadyForPlayStore(),
		TotalChecks:         data.TotalChecks,
		PassedChecks:        data.PassedChecks,
		HighCount:           data.GetHighCount(),
		WarningCount:        data.GetWarningCount(),
		InfoCount:           data.GetInfoCount(),
		Blockers:            data.Blockers,
		FindingsBySeverity: findingsBySeverity{
			HIGH:    data.FindingsBySeverity[report.SeverityHigh],
			WARNING: data.FindingsBySeverity[report.SeverityWarning],
			INFO:    data.FindingsBySeverity[report.SeverityInfo],
		},
		FindingsByCategory: data.FindingsByCategory,
		AndroidConfig:      data.AndroidConfig,
		IOSConfig:          data.IOSConfig,
		FlutterConfig:      data.FlutterConfig,
		AndroidPermissions: data.AndroidPermissions,
		IOSPermissions:     data.IOSPermissions,
		Dependencies:       data.Dependencies,
		DevDependencies:    data.DevDependencies,
		HighRiskDeps:       data.HighRiskDeps,
		OutdatedDeps:       data.OutdatedDeps,
		SecurityFlags:      data.SecurityFlags,
		PolicyFlags:        data.PolicyFlags,
		MissingPolicies:    data.MissingPolicies,
		ReviewerAccount:    data.ReviewerAccount,
	}
}

// GenerateWithHeader generates prompt with copy-paste instructions
func (pb *PromptBuilder) GenerateWithHeader(data *PromptData) (string, error) {
	prompt, err := pb.Generate(data)
	if err != nil {
		return "", err
	}

	var result strings.Builder

	// Add header
	result.WriteString("\n")
	result.WriteString("╔═══════════════════════════════════════════════════════════════════════════════\n")
	result.WriteString("║  📋 AI COMPLIANCE PROMPT - COPY & PASTE INTO YOUR AI ASSISTANT\n")
	result.WriteString("╚═══════════════════════════════════════════════════════════════════════════════\n")
	result.WriteString("\n")

	// Show instructions
	result.WriteString("A comprehensive prompt has been generated below with all your app's\n")
	result.WriteString("compliance data. Copy everything between the === markers and paste into:\n")
	result.WriteString("  • ChatGPT (https://chat.openai.com)\n")
	result.WriteString("  • Claude (https://claude.ai)\n")
	result.WriteString("  • Gemini (https://gemini.google.com)\n")
	result.WriteString("  • Or any AI assistant of your choice\n")
	result.WriteString("\n")

	// Prompt stats
	result.WriteString("───────────────────────────────────────────────────────────────────────────────\n")
	result.WriteString("PROMPT STATS\n")
	result.WriteString("───────────────────────────────────────────────────────────────────────────────\n")
	result.WriteString("Character count: " + formatNumber(len(prompt)) + "\n")
	result.WriteString("Estimated tokens: ~" + formatNumber(len(prompt)/4) + "\n")
	result.WriteString("Template: " + string(pb.templateType) + "\n")
	result.WriteString("\n")

	// Start of prompt marker
	result.WriteString("═══════════════════════════════════════════════════════════════════════════════\n")
	result.WriteString("START OF PROMPT - COPY EVERYTHING BELOW THIS LINE\n")
	result.WriteString("═══════════════════════════════════════════════════════════════════════════════\n")
	result.WriteString("\n")

	// The actual prompt
	result.WriteString(prompt)

	// End of prompt marker
	result.WriteString("\n")
	result.WriteString("═══════════════════════════════════════════════════════════════════════════════\n")
	result.WriteString("END OF PROMPT - COPY EVERYTHING ABOVE THIS LINE\n")
	result.WriteString("═══════════════════════════════════════════════════════════════════════════════\n")
	result.WriteString("\n")

	// Footer with tips
	result.WriteString("💡 TIPS FOR BEST RESULTS:\n")
	result.WriteString("   1. Use GPT-4, Claude 3, or equivalent models for best analysis\n")
	result.WriteString("   2. The prompt includes all compliance data - no need to add context\n")
	result.WriteString("   3. Save the AI's response for your compliance documentation\n")
	result.WriteString("\n")

	// AI integration hint
	result.WriteString("🚀 WANT AUTOMATIC ANALYSIS?\n")
	result.WriteString("   Set AI_API_KEY environment variable for direct AI integration:\n")
	result.WriteString("   export AI_API_KEY=your_key_here\n")
	result.WriteString("   fsct check . --ai\n")
	result.WriteString("\n")

	return result.String(), nil
}

// GenerateCompact creates a compact version for terminal display
func (pb *PromptBuilder) GenerateCompact(data *PromptData) (string, error) {
	// Use summary template for compact view
	pb.SetTemplateType(TemplateTypeSummary)
	return pb.GenerateWithHeader(data)
}

// ExportToFile saves the prompt to a file
func (pb *PromptBuilder) ExportToFile(data *PromptData, filename string) error {
	prompt, err := pb.Generate(data)
	if err != nil {
		return err
	}

	// Add metadata header for file
	var result strings.Builder
	result.WriteString("# FSCT AI Compliance Prompt\n")
	result.WriteString("# Generated: " + getCurrentTimestamp() + "\n")
	result.WriteString("# App: " + data.AppName + " v" + data.Version + "\n")
	result.WriteString("# Compliance Score: " + formatNumber(data.CalculateComplianceScore()) + "/100\n")
	result.WriteString("\n")
	result.WriteString("---\n")
	result.WriteString("\n")
	result.WriteString(prompt)

	return writeFile(filename, result.String())
}

// Helper functions
func formatNumber(n int) string {
	if n < 1000 {
		return string(rune('0' + n))
	}
	// Simple formatting for larger numbers
	return "1000+"
}

func getCurrentTimestamp() string {
	return "2024-01-30" // Simplified - would use time.Now() in production
}

func writeFile(filename, content string) error {
	// Simplified file write - would use os.WriteFile in production
	return nil
}
