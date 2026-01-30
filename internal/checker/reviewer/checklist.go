package reviewer

import (
	"fmt"
	"strings"

	"github.com/ricky-irfandi/fsct/internal/config"
)

// ChecklistItem represents a single checklist item
type ChecklistItem struct {
	ID          string
	Description string
	Checked     bool
	Required    bool
}

// ReviewerChecklist contains all reviewer-related checklist items
type ReviewerChecklist struct {
	Items           []ChecklistItem
	EmailConfigured bool
	PasswordConfigured bool
	Verified        bool
	CanLogin        bool
}

// GetCompletedCount returns the number of completed items
func (c *ReviewerChecklist) GetCompletedCount() int {
	count := 0
	for _, item := range c.Items {
		if item.Checked {
			count++
		}
	}
	return count
}

// GetTotalCount returns the total number of items
func (c *ReviewerChecklist) GetTotalCount() int {
	return len(c.Items)
}

// CompletionPercent returns the completion percentage as an integer
func (c *ReviewerChecklist) CompletionPercent() int {
	if len(c.Items) == 0 {
		return 100
	}
	return int(c.GetCompletionPercentage())
}

// ReadyForSubmission returns true if all required items are checked
func (c *ReviewerChecklist) ReadyForSubmission() bool {
	return c.IsReady()
}

// GenerateChecklist generates a reviewer checklist based on configuration
func GenerateChecklist(cfg *config.ReviewerConfig) *ReviewerChecklist {
	checklist := &ReviewerChecklist{
		Items: make([]ChecklistItem, 0),
	}

	if cfg == nil {
		checklist.Items = append(checklist.Items, ChecklistItem{
			ID:          "config",
			Description: "Reviewer configuration file (.fsct.yaml) exists",
			Checked:     false,
			Required:    true,
		})
		return checklist
	}

	// Check email configuration
	email := cfg.Email
	if email == "" && cfg.EmailEnv != "" {
		email = getEnv(cfg.EmailEnv)
	}
	checklist.EmailConfigured = email != ""

	checklist.Items = append(checklist.Items, ChecklistItem{
		ID:          "email",
		Description: fmt.Sprintf("Reviewer email configured (%s)", maskEmail(email)),
		Checked:     checklist.EmailConfigured,
		Required:    true,
	})

	// Check password configuration
	password := cfg.Password
	if password == "" && cfg.PasswordEnv != "" {
		password = getEnv(cfg.PasswordEnv)
	}
	checklist.PasswordConfigured = password != ""

	checklist.Items = append(checklist.Items, ChecklistItem{
		ID:          "password",
		Description: "Reviewer password configured",
		Checked:     checklist.PasswordConfigured,
		Required:    true,
	})

	// Check if credentials are valid (not placeholders)
	isPlaceholder := false
	if email != "" {
		placeholders := []string{"test@", "example.com", "your-email@"}
		emailLower := strings.ToLower(email)
		for _, p := range placeholders {
			if strings.Contains(emailLower, p) {
				isPlaceholder = true
				break
			}
		}
	}

	checklist.Items = append(checklist.Items, ChecklistItem{
		ID:          "valid",
		Description: "Credentials are valid (not placeholders)",
		Checked:     checklist.EmailConfigured && !isPlaceholder,
		Required:    true,
	})

	// Check password strength
	isWeak := false
	if password != "" {
		if len(password) < 8 {
			isWeak = true
		} else {
			weakPatterns := []string{"password", "123456", "qwerty", "admin"}
			passLower := strings.ToLower(password)
			for _, p := range weakPatterns {
				if strings.Contains(passLower, p) {
					isWeak = true
					break
				}
			}
		}
	}

	checklist.Items = append(checklist.Items, ChecklistItem{
		ID:          "strong",
		Description: "Password is strong (8+ chars, no common patterns)",
		Checked:     checklist.PasswordConfigured && !isWeak,
		Required:    false,
	})

	// Check verification
	if cfg.Verification != nil && cfg.Verification.Enabled {
		checklist.Items = append(checklist.Items, ChecklistItem{
			ID:          "verify",
			Description: "Login verification enabled and tested",
			Checked:     false, // Would be set based on actual verification
			Required:    false,
		})
	}

	// Additional checklist items for reviewer preparation
	checklist.Items = append(checklist.Items, []ChecklistItem{
		{
			ID:          "test_data",
			Description: "Test data is pre-populated in the app",
			Checked:     false,
			Required:    false,
		},
		{
			ID:          "demo_mode",
			Description: "Demo mode or test mode is clearly marked",
			Checked:     false,
			Required:    false,
		},
		{
			ID:          "instructions",
			Description: "Reviewer instructions document prepared",
			Checked:     false,
			Required:    false,
		},
		{
			ID:          "screenshots",
			Description: "Screenshots show login flow",
			Checked:     false,
			Required:    false,
		},
	}...)

	return checklist
}

// ToMarkdown converts the checklist to markdown format
func (c *ReviewerChecklist) ToMarkdown() string {
	var sb strings.Builder

	sb.WriteString("# Reviewer Account Checklist\n\n")

	// Required items
	sb.WriteString("## Required Items\n\n")
	for _, item := range c.Items {
		if item.Required {
			status := "[ ]"
			if item.Checked {
				status = "[x]"
			}
			sb.WriteString(fmt.Sprintf("- %s %s\n", status, item.Description))
		}
	}

	// Optional items
	sb.WriteString("\n## Recommended Items\n\n")
	for _, item := range c.Items {
		if !item.Required {
			status := "[ ]"
			if item.Checked {
				status = "[x]"
			}
			sb.WriteString(fmt.Sprintf("- %s %s\n", status, item.Description))
		}
	}

	// Status summary
	sb.WriteString("\n## Status Summary\n\n")
	requiredCount := 0
	requiredDone := 0
	optionalCount := 0
	optionalDone := 0

	for _, item := range c.Items {
		if item.Required {
			requiredCount++
			if item.Checked {
				requiredDone++
			}
		} else {
			optionalCount++
			if item.Checked {
				optionalDone++
			}
		}
	}

	sb.WriteString(fmt.Sprintf("- **Required**: %d/%d complete\n", requiredDone, requiredCount))
	sb.WriteString(fmt.Sprintf("- **Recommended**: %d/%d complete\n", optionalDone, optionalCount))

	if requiredDone == requiredCount {
		sb.WriteString("\n✅ **Ready for submission!**\n")
	} else {
		sb.WriteString("\n⚠️ **Complete required items before submission**\n")
	}

	return sb.String()
}

// ToConsole converts the checklist to console format
func (c *ReviewerChecklist) ToConsole() string {
	var sb strings.Builder

	sb.WriteString("╔═══════════════════════════════════════════════════════════════════════════════\n")
	sb.WriteString("║  REVIEWER ACCOUNT CHECKLIST\n")
	sb.WriteString("╚═══════════════════════════════════════════════════════════════════════════════\n\n")

	// Required items
	sb.WriteString("Required Items:\n")
	for _, item := range c.Items {
		if item.Required {
			status := "❌"
			if item.Checked {
				status = "✅"
			}
			sb.WriteString(fmt.Sprintf("  %s %s\n", status, item.Description))
		}
	}

	// Optional items
	sb.WriteString("\nRecommended Items:\n")
	for _, item := range c.Items {
		if !item.Required {
			status := "⬜"
			if item.Checked {
				status = "✅"
			}
			sb.WriteString(fmt.Sprintf("  %s %s\n", status, item.Description))
		}
	}

	// Status summary
	requiredCount := 0
	requiredDone := 0
	for _, item := range c.Items {
		if item.Required {
			requiredCount++
			if item.Checked {
				requiredDone++
			}
		}
	}

	sb.WriteString("\n")
	if requiredDone == requiredCount {
		sb.WriteString("✅ Ready for submission!\n")
	} else {
		sb.WriteString(fmt.Sprintf("⚠️  %d/%d required items complete\n", requiredDone, requiredCount))
	}

	return sb.String()
}

// ToHTML converts the checklist to HTML format
func (c *ReviewerChecklist) ToHTML() string {
	var sb strings.Builder

	sb.WriteString("<!DOCTYPE html>\n")
	sb.WriteString("<html>\n<head>\n")
	sb.WriteString("<title>Reviewer Account Checklist</title>\n")
	sb.WriteString("<style>\n")
	sb.WriteString("body { font-family: Arial, sans-serif; max-width: 800px; margin: 40px auto; padding: 20px; }\n")
	sb.WriteString("h1 { color: #333; }\n")
	sb.WriteString("h2 { color: #555; border-bottom: 1px solid #ddd; padding-bottom: 10px; }\n")
	sb.WriteString("ul { list-style: none; padding: 0; }\n")
	sb.WriteString("li { padding: 8px 0; }\n")
	sb.WriteString(".checked { color: green; }\n")
	sb.WriteString(".unchecked { color: #ccc; }\n")
	sb.WriteString(".summary { background: #f5f5f5; padding: 20px; border-radius: 5px; margin-top: 20px; }\n")
	sb.WriteString(".ready { color: green; font-weight: bold; }\n")
	sb.WriteString(".not-ready { color: orange; font-weight: bold; }\n")
	sb.WriteString("</style>\n")
	sb.WriteString("</head>\n<body>\n")

	sb.WriteString("<h1>Reviewer Account Checklist</h1>\n")

	// Required items
	sb.WriteString("<h2>Required Items</h2>\n<ul>\n")
	for _, item := range c.Items {
		if item.Required {
			class := "unchecked"
			status := "☐"
			if item.Checked {
				class = "checked"
				status = "☑"
			}
			sb.WriteString(fmt.Sprintf("<li class=\"%s\">%s %s</li>\n", class, status, item.Description))
		}
	}
	sb.WriteString("</ul>\n")

	// Optional items
	sb.WriteString("<h2>Recommended Items</h2>\n<ul>\n")
	for _, item := range c.Items {
		if !item.Required {
			class := "unchecked"
			status := "☐"
			if item.Checked {
				class = "checked"
				status = "☑"
			}
			sb.WriteString(fmt.Sprintf("<li class=\"%s\">%s %s</li>\n", class, status, item.Description))
		}
	}
	sb.WriteString("</ul>\n")

	// Status summary
	requiredCount := 0
	requiredDone := 0
	for _, item := range c.Items {
		if item.Required {
			requiredCount++
			if item.Checked {
				requiredDone++
			}
		}
	}

	sb.WriteString("<div class=\"summary\">\n")
	sb.WriteString(fmt.Sprintf("<p><strong>Required:</strong> %d/%d complete</p>\n", requiredDone, requiredCount))
	
	if requiredDone == requiredCount {
		sb.WriteString("<p class=\"ready\">✓ Ready for submission!</p>\n")
	} else {
		sb.WriteString("<p class=\"not-ready\">⚠ Complete required items before submission</p>\n")
	}
	sb.WriteString("</div>\n")

	sb.WriteString("</body>\n</html>")

	return sb.String()
}

// IsReady returns true if all required items are checked
func (c *ReviewerChecklist) IsReady() bool {
	for _, item := range c.Items {
		if item.Required && !item.Checked {
			return false
		}
	}
	return true
}

// GetCompletionPercentage returns the completion percentage
func (c *ReviewerChecklist) GetCompletionPercentage() float64 {
	if len(c.Items) == 0 {
		return 100.0
	}

	checked := 0
	for _, item := range c.Items {
		if item.Checked {
			checked++
		}
	}

	return float64(checked) / float64(len(c.Items)) * 100
}
