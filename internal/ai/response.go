package ai

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/ricky-irfandi/fsct/internal/report"
)

// AIAnalysis represents the structured analysis from AI
type AIAnalysis struct {
	// Overall assessment
	RiskLevel       string           `json:"risk_level"`       // low, medium, high
	Confidence      string           `json:"confidence"`       // low, medium, high
	ComplianceScore int              `json:"compliance_score"` // 0-100

	// Store readiness
	StoreReadiness StoreReadiness `json:"store_readiness"`

	// Detailed findings from AI
	Insights []AIInsight `json:"insights"`

	// Actionable suggestions
	Suggestions []AISuggestion `json:"suggestions"`

	// Reviewer-specific information
	ReviewerNotes []string `json:"reviewer_notes"`

	// Raw response for debugging
	RawContent string `json:"-"`
}

// StoreReadiness indicates readiness for each store
type StoreReadiness struct {
	AppStore   bool   `json:"app_store"`
	PlayStore  bool   `json:"play_store"`
	Reasoning  string `json:"reasoning,omitempty"`
}

// AIInsight represents a single insight from AI analysis
type AIInsight struct {
	Category    string `json:"category"`              // permissions, policy, security, etc.
	Severity    string `json:"severity"`              // info, warning, high
	Title       string `json:"title"`
	Description string `json:"description"`
	Confidence  string `json:"confidence,omitempty"`  // low, medium, high
}

// AISuggestion represents an actionable suggestion
type AISuggestion struct {
	Priority    int    `json:"priority"`    // 1-5 (1 = highest)
	Category    string `json:"category"`
	Issue       string `json:"issue"`
	Action      string `json:"action"`
	CodeExample string `json:"code_example,omitempty"`
	FilePath    string `json:"file_path,omitempty"` // Generic path only, no absolute paths
}

// ParseResponse parses an AI response into structured analysis
func ParseResponse(content string) (*AIAnalysis, error) {
	// Try to extract JSON from the response
	jsonContent := extractJSON(content)
	if jsonContent == "" {
		// Try to parse as plain text if no JSON found
		return parsePlainTextResponse(content), nil
	}

	var analysis AIAnalysis
	if err := json.Unmarshal([]byte(jsonContent), &analysis); err != nil {
		// If JSON parsing fails, try to extract structured info from text
		return parsePlainTextResponse(content), nil
	}

	analysis.RawContent = content
	return &analysis, nil
}

// ParseResponseStrict parses response and returns error if JSON is invalid
func ParseResponseStrict(content string) (*AIAnalysis, error) {
	jsonContent := extractJSON(content)
	if jsonContent == "" {
		return nil, fmt.Errorf("%w: no JSON found in response", ErrInvalidResponse)
	}

	var analysis AIAnalysis
	if err := json.Unmarshal([]byte(jsonContent), &analysis); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidResponse, err)
	}

	analysis.RawContent = content
	return &analysis, nil
}

// extractJSON extracts JSON content from a string
func extractJSON(content string) string {
	// Look for JSON between triple backticks
	jsonBlockRegex := regexp.MustCompile("```(?:json)?\\s*([\\s\\S]*?)```")
	matches := jsonBlockRegex.FindStringSubmatch(content)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// Look for JSON between curly braces
	startIdx := strings.Index(content, "{")
	endIdx := strings.LastIndex(content, "}")
	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		return content[startIdx : endIdx+1]
	}

	return ""
}

// parsePlainTextResponse parses a plain text response into structured analysis
func parsePlainTextResponse(content string) *AIAnalysis {
	analysis := &AIAnalysis{
		RawContent:    content,
		Insights:      make([]AIInsight, 0),
		Suggestions:   make([]AISuggestion, 0),
		ReviewerNotes: make([]string, 0),
		StoreReadiness: StoreReadiness{
			AppStore:  false,
			PlayStore: false,
		},
	}

	// Extract risk level
	contentLower := strings.ToLower(content)
	if strings.Contains(contentLower, "high risk") || strings.Contains(contentLower, "critical") {
		analysis.RiskLevel = "high"
	} else if strings.Contains(contentLower, "medium risk") || strings.Contains(contentLower, "moderate") {
		analysis.RiskLevel = "medium"
	} else {
		analysis.RiskLevel = "low"
	}

	// Extract store readiness
	if strings.Contains(contentLower, "ready for app store") || strings.Contains(contentLower, "app store: yes") {
		analysis.StoreReadiness.AppStore = true
	}
	if strings.Contains(contentLower, "ready for play store") || strings.Contains(contentLower, "play store: yes") {
		analysis.StoreReadiness.PlayStore = true
	}

	// Extract insights from sections
	analysis.Insights = extractInsightsFromText(content)
	analysis.Suggestions = extractSuggestionsFromText(content)
	analysis.ReviewerNotes = extractReviewerNotes(content)

	return analysis
}

// extractInsightsFromText extracts insights from text content
func extractInsightsFromText(content string) []AIInsight {
	insights := make([]AIInsight, 0)

	// Look for common insight patterns
	lines := strings.Split(content, "\n")
	currentCategory := "general"

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Detect category from headers
		if strings.HasPrefix(line, "##") || strings.HasPrefix(line, "**") {
			lower := strings.ToLower(line)
			switch {
			case strings.Contains(lower, "permission"):
				currentCategory = "permissions"
			case strings.Contains(lower, "privacy") || strings.Contains(lower, "policy"):
				currentCategory = "policy"
			case strings.Contains(lower, "security"):
				currentCategory = "security"
			case strings.Contains(lower, "store"):
				currentCategory = "store"
			case strings.Contains(lower, "reviewer"):
				currentCategory = "reviewer"
			}
			continue
		}

		// Look for list items with issues
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "•") {
			line = strings.TrimLeft(line, "-*• ")

			// Determine severity
			severity := "info"
			lowerLine := strings.ToLower(line)
			if strings.Contains(lowerLine, "critical") || strings.Contains(lowerLine, "blocker") ||
				strings.Contains(lowerLine, "❌") || strings.Contains(lowerLine, "must fix") {
				severity = "high"
			} else if strings.Contains(lowerLine, "warning") || strings.Contains(lowerLine, "should") ||
				strings.Contains(lowerLine, "⚠️") {
				severity = "warning"
			}

			if len(line) > 10 {
				insights = append(insights, AIInsight{
					Category:    currentCategory,
					Severity:    severity,
					Title:       truncateString(line, 100),
					Description: line,
				})
			}
		}
	}

	return insights
}

// extractSuggestionsFromText extracts suggestions from text content
func extractSuggestionsFromText(content string) []AISuggestion {
	suggestions := make([]AISuggestion, 0)
	priority := 1

	// Look for numbered suggestions - simple approach without lookahead
	numberedRegex := regexp.MustCompile(`(?i)^\s*(\d+)[:.\)]\s*(.+)$`)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		matches := numberedRegex.FindStringSubmatch(line)
		if len(matches) > 2 {
			suggestion := AISuggestion{
				Priority: priority,
				Issue:    strings.TrimSpace(matches[2]),
				Action:   "Review and address this issue",
			}
			suggestions = append(suggestions, suggestion)
			priority++
		}
	}

	// If no numbered suggestions, look for bullet points with actions
	if len(suggestions) == 0 {
		for _, line := range lines {
			line = strings.TrimSpace(line)
			lower := strings.ToLower(line)

			// Look for action items
			if (strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*")) &&
				(strings.Contains(lower, "fix") || strings.Contains(lower, "add") ||
					strings.Contains(lower, "update") || strings.Contains(lower, "remove")) {
				suggestions = append(suggestions, AISuggestion{
					Priority: priority,
					Issue:    strings.TrimLeft(line, "-* "),
					Action:   "Address this item",
				})
				priority++
			}
		}
	}

	return suggestions
}

// extractReviewerNotes extracts reviewer-specific notes
func extractReviewerNotes(content string) []string {
	notes := make([]string, 0)

	// Look for reviewer section
	reviewerSection := false
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		lower := strings.ToLower(line)

		// Detect reviewer section
		if strings.Contains(lower, "reviewer") && strings.Contains(lower, "instruction") {
			reviewerSection = true
			continue
		}

		// End of section
		if reviewerSection && (strings.HasPrefix(line, "##") || strings.HasPrefix(line, "---")) {
			break
		}

		// Collect notes
		if reviewerSection && len(line) > 10 {
			notes = append(notes, line)
		}
	}

	return notes
}

// ConvertToFindings converts AI analysis to report findings
func (a *AIAnalysis) ConvertToFindings(checkPrefix string) []report.Finding {
	findings := make([]report.Finding, 0)

	// Convert insights to findings
	for i, insight := range a.Insights {
		severity := report.SeverityInfo
		switch insight.Severity {
		case "high":
			severity = report.SeverityHigh
		case "warning":
			severity = report.SeverityWarning
		}

		finding := report.Finding{
			ID:         fmt.Sprintf("%s-%03d", checkPrefix, i+1),
			Severity:   severity,
			Title:      insight.Title,
			Message:    insight.Description,
			Suggestion: a.findSuggestionForInsight(insight),
		}
		findings = append(findings, finding)
	}

	return findings
}

// findSuggestionForInsight finds a matching suggestion for an insight
func (a *AIAnalysis) findSuggestionForInsight(insight AIInsight) string {
	for _, suggestion := range a.Suggestions {
		if strings.Contains(strings.ToLower(suggestion.Issue), strings.ToLower(insight.Title)) ||
			strings.Contains(strings.ToLower(insight.Description), strings.ToLower(suggestion.Issue)) {
			return suggestion.Action
		}
	}
	return ""
}

// IsReady returns true if the app is ready for submission
func (a *AIAnalysis) IsReady() bool {
	return a.StoreReadiness.AppStore && a.StoreReadiness.PlayStore && a.RiskLevel != "high"
}

// GetPriorityFindings returns findings sorted by priority
func (a *AIAnalysis) GetPriorityFindings(max int) []AIInsight {
	// Filter high severity first
	highPriority := make([]AIInsight, 0)
	mediumPriority := make([]AIInsight, 0)
	lowPriority := make([]AIInsight, 0)

	for _, insight := range a.Insights {
		switch insight.Severity {
		case "high":
			highPriority = append(highPriority, insight)
		case "warning":
			mediumPriority = append(mediumPriority, insight)
		default:
			lowPriority = append(lowPriority, insight)
		}
	}

	// Combine and limit
	result := make([]AIInsight, 0)
	result = append(result, highPriority...)
	result = append(result, mediumPriority...)
	result = append(result, lowPriority...)

	if max > 0 && len(result) > max {
		return result[:max]
	}
	return result
}

// Helper functions

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ExpectedJSONSchema returns the expected JSON schema for AI responses
func ExpectedJSONSchema() string {
	return `{
  "risk_level": "low|medium|high",
  "confidence": "low|medium|high",
  "compliance_score": 0-100,
  "store_readiness": {
    "app_store": true|false,
    "play_store": true|false,
    "reasoning": "explanation"
  },
  "insights": [
    {
      "category": "permissions|policy|security|store|reviewer|general",
      "severity": "info|warning|high",
      "title": "brief title",
      "description": "detailed description",
      "confidence": "low|medium|high"
    }
  ],
  "suggestions": [
    {
      "priority": 1-5,
      "category": "category",
      "issue": "what needs fixing",
      "action": "how to fix it",
      "code_example": "optional code snippet",
      "file_path": "optional generic path"
    }
  ],
  "reviewer_notes": [
    "note for app reviewers"
  ]
}`
}
