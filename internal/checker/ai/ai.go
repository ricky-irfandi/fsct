package ai

import (
	"context"
	"fmt"
	"time"

	aipkg "github.com/ricky-irfandi/fsct/internal/ai"
	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

// Check represents an AI-powered compliance check
type Check struct {
	id       string
	name     string
	client   *aipkg.Client
	prompts  *SystemPrompts
	category string
}

// NewCheck creates a new AI check
func NewCheck(id, name, category string, client *aipkg.Client) *Check {
	return &Check{
		id:       id,
		name:     name,
		client:   client,
		prompts:  &SystemPrompts{},
		category: category,
	}
}

// ID returns the check ID
func (c *Check) ID() string {
	return c.id
}

// Name returns the check name
func (c *Check) Name() string {
	return c.name
}

// Category returns the check category
func (c *Check) Category() string {
	return c.category
}

// Run executes the AI check
func (c *Check) Run(project *checker.Project) []report.Finding {
	// Extract metadata
	findings := []report.Finding{}
	metadata := aipkg.ExtractMetadata(project, findings)

	// Get system prompt based on check type
	systemPrompt := c.getSystemPrompt()

	// Build user prompt
	userPrompt := BuildUserPrompt(c.name, metadata)

	// Call AI with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := c.client.Complete(ctx, &aipkg.CompletionRequest{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		MaxTokens:    2000,
		Temperature:  0.3, // Lower temperature for more consistent results
	})

	if err != nil {
		// Return offline finding if AI fails
		return []report.Finding{c.createOfflineFinding(err)}
	}

	// Parse response
	analysis, err := aipkg.ParseResponse(resp.Content)
	if err != nil {
		// Try to parse as plain text if JSON parsing fails
		analysis, _ = aipkg.ParseResponse(resp.Content)
	}

	// Convert to findings
	return c.convertToFindings(analysis)
}

// getSystemPrompt returns the appropriate system prompt
func (c *Check) getSystemPrompt() string {
	switch c.id {
	case "AI-001":
		return c.prompts.AI001PermissionJustification()
	case "AI-002":
		return c.prompts.AI002PolicyCompliance()
	case "AI-003":
		return c.prompts.AI003DependencyRisk()
	case "AI-004":
		return c.prompts.AI004StoreGuidance()
	case "AI-005":
		return c.prompts.AI005ReviewerNotes()
	default:
		return c.prompts.AI001PermissionJustification()
	}
}

// convertToFindings converts AI analysis to report findings
func (c *Check) convertToFindings(analysis *aipkg.AIAnalysis) []report.Finding {
	findings := make([]report.Finding, 0)

	// Convert insights to findings
	for i, insight := range analysis.Insights {
		severity := report.SeverityInfo
		switch insight.Severity {
		case "high":
			severity = report.SeverityHigh
		case "warning":
			severity = report.SeverityWarning
		}

		// Find matching suggestion
		suggestion := ""
		for _, s := range analysis.Suggestions {
			if s.Priority == i+1 {
				suggestion = s.Action
				break
			}
		}

		finding := report.Finding{
			ID:         fmt.Sprintf("%s-%03d", c.id, i+1),
			Severity:   severity,
			Title:      fmt.Sprintf("[%s] %s", c.name, insight.Title),
			Message:    insight.Description,
			Suggestion: suggestion,
		}
		findings = append(findings, finding)
	}

	// If no insights but we have suggestions, create generic findings
	if len(findings) == 0 && len(analysis.Suggestions) > 0 {
		for i, suggestion := range analysis.Suggestions {
			finding := report.Finding{
				ID:         fmt.Sprintf("%s-%03d", c.id, i+1),
				Severity:   report.SeverityInfo,
				Title:      fmt.Sprintf("[%s] Suggestion %d", c.name, i+1),
				Message:    suggestion.Issue,
				Suggestion: suggestion.Action,
			}
			findings = append(findings, finding)
		}
	}

	// If still no findings, create a summary finding
	if len(findings) == 0 {
		findings = append(findings, report.Finding{
			ID:       fmt.Sprintf("%s-001", c.id),
			Severity: report.SeverityInfo,
			Title:    fmt.Sprintf("[%s] Analysis Complete", c.name),
			Message:  fmt.Sprintf("AI analysis completed with risk level: %s", analysis.RiskLevel),
			Suggestion: fmt.Sprintf("Store readiness - App Store: %v, Play Store: %v",
				analysis.StoreReadiness.AppStore, analysis.StoreReadiness.PlayStore),
		})
	}

	return findings
}

// createOfflineFinding creates a finding when AI is unavailable
func (c *Check) createOfflineFinding(err error) report.Finding {
	return report.Finding{
		ID:       c.id,
		Severity: report.SeverityInfo,
		Title:    fmt.Sprintf("[%s] AI Analysis Unavailable", c.name),
		Message:  fmt.Sprintf("Could not perform AI analysis: %v", err),
		Suggestion: "To enable AI analysis, set AI_API_KEY environment variable or use --ai-key flag. " +
			"Alternatively, use --format prompt to generate a prompt for manual AI analysis.",
	}
}

// AI001PermissionJustificationCheck checks permission justification
func AI001PermissionJustificationCheck(client *aipkg.Client) checker.Check {
	return NewCheck(
		"AI-001",
		"AI Permission Justification Analysis",
		"AI Analysis",
		client,
	)
}

// AI002PolicyComplianceCheck checks policy compliance
func AI002PolicyComplianceCheck(client *aipkg.Client) checker.Check {
	return NewCheck(
		"AI-002",
		"AI Policy Compliance Analysis",
		"AI Analysis",
		client,
	)
}

// AI003DependencyRiskCheck checks dependency risks
func AI003DependencyRiskCheck(client *aipkg.Client) checker.Check {
	return NewCheck(
		"AI-003",
		"AI Dependency Risk Analysis",
		"AI Analysis",
		client,
	)
}

// AI004StoreGuidanceCheck provides store-specific guidance
func AI004StoreGuidanceCheck(client *aipkg.Client) checker.Check {
	return NewCheck(
		"AI-004",
		"AI Store-Specific Guidance",
		"AI Analysis",
		client,
	)
}

// AI005ReviewerNotesCheck generates reviewer notes
func AI005ReviewerNotesCheck(client *aipkg.Client) checker.Check {
	return NewCheck(
		"AI-005",
		"AI Reviewer Notes Generation",
		"AI Analysis",
		client,
	)
}

// AllChecks returns all AI checks
func AllChecks(client *aipkg.Client) []checker.Check {
	if client == nil || !client.IsAvailable() {
		return nil
	}

	return []checker.Check{
		AI001PermissionJustificationCheck(client),
		AI002PolicyComplianceCheck(client),
		AI003DependencyRiskCheck(client),
		AI004StoreGuidanceCheck(client),
		AI005ReviewerNotesCheck(client),
	}
}
