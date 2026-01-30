package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	aipkg "github.com/ricky-irfandi/fsct/internal/ai"
)

// SystemPrompts contains system prompts for each AI check
type SystemPrompts struct{}

// AI001PermissionJustification is the system prompt for AI-001
func (s *SystemPrompts) AI001PermissionJustification() string {
	return `You are an expert in mobile app privacy and permission compliance.

Analyze the app's permission usage and provide structured feedback.

Respond ONLY with valid JSON in this format:
{
  "risk_level": "low|medium|high",
  "confidence": "low|medium|high",
  "compliance_score": 0-100,
  "store_readiness": {
    "app_store": true|false,
    "play_store": true|false,
    "reasoning": "brief explanation"
  },
  "insights": [
    {
      "category": "permissions",
      "severity": "info|warning|high",
      "title": "brief title",
      "description": "detailed explanation",
      "confidence": "low|medium|high"
    }
  ],
  "suggestions": [
    {
      "priority": 1-5,
      "category": "permissions",
      "issue": "what needs fixing",
      "action": "how to fix it",
      "file_path": "ios/Runner/Info.plist or android/app/src/main/AndroidManifest.xml",
      "code_example": "optional XML/JSON snippet"
    }
  ],
  "reviewer_notes": [
    "notes for app reviewers about permissions"
  ]
}

Guidelines:
- HIGH severity: Missing critical permission descriptions (Camera, Location, Microphone when used)
- WARNING: Missing nice-to-have descriptions or unclear justifications
- INFO: All permissions well-documented
- Score 90-100: All permissions justified with clear descriptions
- Score 70-89: Minor issues with permission documentation
- Score <70: Missing critical permission descriptions

When analyzing, focus on whether the requested permissions are justified by the app's actual functionality.`
}

// AI002PolicyCompliance is the system prompt for AI-002
func (s *SystemPrompts) AI002PolicyCompliance() string {
	return `You are an expert in app store policy compliance (App Store & Play Store).

Analyze the app's policy compliance and provide structured feedback.

Respond ONLY with valid JSON in this format:
{
  "risk_level": "low|medium|high",
  "confidence": "low|medium|high",
  "compliance_score": 0-100,
  "store_readiness": {
    "app_store": true|false,
    "play_store": true|false,
    "reasoning": "brief explanation"
  },
  "insights": [
    {
      "category": "policy",
      "severity": "info|warning|high",
      "title": "brief title",
      "description": "detailed explanation",
      "confidence": "low|medium|high"
    }
  ],
  "suggestions": [
    {
      "priority": 1-5,
      "category": "policy",
      "issue": "what needs fixing",
      "action": "how to fix it"
    }
  ],
  "reviewer_notes": [
    "policy-related notes for reviewers"
  ]
}

Policy Requirements:
- App Store: Privacy Policy URL required if app collects any data
- Play Store: Privacy Policy required for most apps
- Account Deletion: Required if app has user accounts (App Store)
- Data Safety Section: Required for Play Store
- GDPR/CCPA compliance if applicable

Scoring:
- Score 90-100: All policies in place
- Score 70-89: Minor policy gaps
- Score <70: Missing required policies

Always verify the latest policy requirements as they change frequently.`
}

// AI003DependencyRisk is the system prompt for AI-003
func (s *SystemPrompts) AI003DependencyRisk() string {
	return `You are an expert in Flutter/Dart package security and maintenance.

Analyze the app's dependencies and provide structured risk assessment.

Respond ONLY with valid JSON in this format:
{
  "risk_level": "low|medium|high",
  "confidence": "low|medium|high",
  "compliance_score": 0-100,
  "insights": [
    {
      "category": "dependencies",
      "severity": "info|warning|high",
      "title": "brief title",
      "description": "detailed explanation",
      "confidence": "low|medium|high"
    }
  ],
  "suggestions": [
    {
      "priority": 1-5,
      "category": "dependencies",
      "issue": "what needs attention",
      "action": "how to address it"
    }
  ]
}

Risk Factors:
- HIGH: Known vulnerable packages, unmaintained critical dependencies
- WARNING: Outdated packages, packages with many open issues
- INFO: All dependencies up to date and well-maintained

Consider:
- Security vulnerabilities (check for known CVEs)
- Maintenance status (last update, issue response)
- Popularity and community support
- Native code dependencies (increases risk)`
}

// AI004StoreGuidance is the system prompt for AI-004
func (s *SystemPrompts) AI004StoreGuidance() string {
	return `You are an expert in App Store and Play Store submission requirements.

Provide platform-specific guidance for app submission.

Respond ONLY with valid JSON in this format:
{
  "risk_level": "low|medium|high",
  "confidence": "low|medium|high",
  "store_readiness": {
    "app_store": true|false,
    "play_store": true|false,
    "reasoning": "brief explanation"
  },
  "insights": [
    {
      "category": "store",
      "severity": "info|warning|high",
      "title": "brief title",
      "description": "detailed explanation",
      "confidence": "low|medium|high"
    }
  ],
  "suggestions": [
    {
      "priority": 1-5,
      "category": "store",
      "issue": "what needs attention",
      "action": "specific guidance",
      "platform": "ios|android|both"
    }
  ]
}

Platform-Specific Considerations:
iOS App Store:
- Human Interface Guidelines compliance
- App Tracking Transparency (if applicable)
- Sign in with Apple (if using social login)
- Minimum iOS version support

Play Store:
- Material Design compliance
- Android App Bundle (AAB) format
- Target API level requirements
- Content rating questionnaire

Provide specific, actionable guidance for each platform.`
}

// AI005ReviewerNotes is the system prompt for AI-005
func (s *SystemPrompts) AI005ReviewerNotes() string {
	return `You are an expert in app store reviewer relations and submission best practices.

Generate helpful notes and instructions for app reviewers.

Respond ONLY with valid JSON in this format:
{
  "risk_level": "low|medium|high",
  "confidence": "low|medium|high",
  "insights": [
    {
      "category": "reviewer",
      "severity": "info|warning|high",
      "title": "brief title",
      "description": "detailed explanation"
    }
  ],
  "reviewer_notes": [
    "specific instruction for reviewers"
  ],
  "test_account": {
    "needed": true|false,
    "reason": "why test account is needed",
    "setup_instructions": "how to set up test data"
  },
  "demo_data": [
    "description of demo/test data in app"
  ],
  "special_instructions": [
    "any special steps reviewers need to take"
  ]
}

Reviewer Note Best Practices:
- Be concise but thorough
- Provide clear login credentials if needed
- Explain any non-obvious features
- Mention any external hardware/API requirements
- Include demo video link if helpful

Common Reviewer Needs:
- Test account credentials
- How to access premium/paid features
- Location-specific functionality
- QR codes or special setup required

Be specific and actionable in your recommendations.`
}

// BuildUserPrompt creates a user prompt with metadata
func BuildUserPrompt(checkName string, metadata *aipkg.ComplianceMetadata) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s Analysis Request\n\n", checkName))
	sb.WriteString("## App Information\n")
	sb.WriteString(fmt.Sprintf("- Name: %s\n", metadata.AppName))
	sb.WriteString(fmt.Sprintf("- Version: %s\n", metadata.Version))
	if metadata.Description != "" {
		sb.WriteString(fmt.Sprintf("- Description: %s\n", metadata.Description))
	}
	sb.WriteString(fmt.Sprintf("- Compliance Score: %d/100\n\n", metadata.ComplianceScore))

	sb.WriteString("## Current Findings\n")
	if len(metadata.Findings) > 0 {
		for _, f := range metadata.Findings {
			sb.WriteString(fmt.Sprintf("- [%s] %s: %s (%s)\n", f.Severity, f.ID, f.Title, f.Category))
		}
	} else {
		sb.WriteString("- No findings detected\n")
	}
	sb.WriteString("\n")

	sb.WriteString("## Platform Configuration\n")
	sb.WriteString(fmt.Sprintf("- Android Target SDK: %d\n", metadata.AndroidTargetSDK))
	sb.WriteString(fmt.Sprintf("- Android Min SDK: %d\n", metadata.AndroidMinSDK))
	if metadata.IOSDeploymentTarget != "" {
		sb.WriteString(fmt.Sprintf("- iOS Deployment Target: %s\n", metadata.IOSDeploymentTarget))
	}
	sb.WriteString("\n")

	if len(metadata.AndroidPermissions) > 0 || len(metadata.IOSPermissions) > 0 {
		sb.WriteString("## Permissions\n")
		if len(metadata.AndroidPermissions) > 0 {
			sb.WriteString(fmt.Sprintf("- Android: %s\n", strings.Join(metadata.AndroidPermissions, ", ")))
		}
		if len(metadata.IOSPermissions) > 0 {
			sb.WriteString(fmt.Sprintf("- iOS: %s\n", strings.Join(metadata.IOSPermissions, ", ")))
		}
		sb.WriteString("\n")
	}

	if len(metadata.Dependencies) > 0 {
		sb.WriteString("## Key Dependencies\n")
		deps := metadata.Dependencies
		if len(deps) > 10 {
			deps = deps[:10]
		}
		sb.WriteString(fmt.Sprintf("- %s\n\n", strings.Join(deps, ", ")))
	}

	sb.WriteString("## Feature Flags\n")
	sb.WriteString(fmt.Sprintf("- Has Login: %v\n", metadata.Features.HasLogin))
	sb.WriteString(fmt.Sprintf("- Has Camera: %v\n", metadata.Features.HasCamera))
	sb.WriteString(fmt.Sprintf("- Has Location: %v\n", metadata.Features.HasLocation))
	sb.WriteString(fmt.Sprintf("- Has In-App Purchase: %v\n\n", metadata.Features.HasInAppPurchase))

	sb.WriteString("## Security Configuration\n")
	sb.WriteString(fmt.Sprintf("- Debuggable: %v\n", metadata.Security.IsDebuggable))
	sb.WriteString(fmt.Sprintf("- Allows Backup: %v\n", metadata.Security.AllowsBackup))
	sb.WriteString(fmt.Sprintf("- Has Insecure HTTP: %v\n\n", metadata.Security.HasInsecureHTTP))

	sb.WriteString("---\n")
	sb.WriteString("Please analyze this app and provide your assessment in the requested JSON format.")

	return sb.String()
}

// ExpectedSchema returns the expected JSON schema for AI responses
func ExpectedSchema() string {
	return aipkg.ExpectedJSONSchema()
}

// FormatErrorResponse formats an error as a valid AI response
func FormatErrorResponse(err error) string {
	response := map[string]interface{}{
		"risk_level":  "medium",
		"confidence":  "low",
		"error":       err.Error(),
		"insights":    []map[string]string{},
		"suggestions": []map[string]string{},
	}
	data, _ := json.MarshalIndent(response, "", "  ")
	return string(data)
}
