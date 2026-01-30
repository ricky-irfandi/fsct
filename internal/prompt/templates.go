package prompt

// MainPromptTemplate is the comprehensive AI compliance prompt template
const MainPromptTemplate = `You are an expert mobile app store compliance reviewer with 10+ years of experience helping apps pass App Store and Play Store review on the first submission.

Analyze the following Flutter app data and provide detailed, actionable compliance recommendations.

═══════════════════════════════════════════════════════════════════════════════
APP METADATA
═══════════════════════════════════════════════════════════════════════════════

• App Name: {{.AppName}}
• Version: {{.Version}}
• Description: {{.Description}}
• Flutter SDK: {{.FlutterVersion}}
• Repository: {{.Repository}}
• Homepage: {{.Homepage}}

═══════════════════════════════════════════════════════════════════════════════
COMPLIANCE SCORE: {{.ComplianceScore}}/100
═══════════════════════════════════════════════════════════════════════════════

Status: {{.Status}}
Ready for App Store: {{if .IsReadyForAppStore}}✅ YES{{else}}❌ NO{{end}}
Ready for Play Store: {{if .IsReadyForPlayStore}}✅ YES{{else}}❌ NO{{end}}

═══════════════════════════════════════════════════════════════════════════════
STATIC ANALYSIS RESULTS
═══════════════════════════════════════════════════════════════════════════════

Total Checks: {{.TotalChecks}}
Passed: {{.PassedChecks}}
High Severity: {{.HighCount}}
Warnings: {{.WarningCount}}
Info: {{.InfoCount}}

{{if .Blockers}}
───────────────────────────────────────────────────────────────────────────────
❌ CRITICAL BLOCKERS (Must Fix Before Submit)
───────────────────────────────────────────────────────────────────────────────
{{range .Blockers}}
[{{.ID}}] {{.Title}}
  Platform: {{.Platform}}
  Issue: {{.Description}}
  Fix File: {{.FixFile}}
{{end}}
{{end}}

{{if .FindingsBySeverity.HIGH}}
───────────────────────────────────────────────────────────────────────────────
HIGH SEVERITY ISSUES
───────────────────────────────────────────────────────────────────────────────
{{range .FindingsBySeverity.HIGH}}
• [{{.ID}}] {{.Title}}
  {{.Message}}
  File: {{.File}}
  Suggestion: {{.Suggestion}}
{{end}}
{{end}}

{{if .FindingsBySeverity.WARNING}}
───────────────────────────────────────────────────────────────────────────────
WARNING SEVERITY ISSUES
───────────────────────────────────────────────────────────────────────────────
{{range .FindingsBySeverity.WARNING}}
• [{{.ID}}] {{.Title}}
  {{.Message}}
  Suggestion: {{.Suggestion}}
{{end}}
{{end}}

═══════════════════════════════════════════════════════════════════════════════
PLATFORM CONFIGURATION
═══════════════════════════════════════════════════════════════════════════════

ANDROID CONFIGURATION:
• Package Name: {{.AndroidConfig.PackageName}}
• Target SDK: {{.AndroidConfig.TargetSDKVersion}} (Required: 35+)
• Min SDK: {{.AndroidConfig.MinSDKVersion}} (Recommended: 21+)
• Version Code: {{.AndroidConfig.VersionCode}}
• Version Name: {{.AndroidConfig.VersionName}}
• Debuggable: {{if .AndroidConfig.IsDebuggable}}⚠️ YES (Remove for release){{else}}✅ No{{end}}
• Allow Backup: {{if .AndroidConfig.AllowBackup}}⚠️ Enabled{{else}}✅ Disabled{{end}}

iOS CONFIGURATION:
• Bundle ID: {{.IOSConfig.BundleIdentifier}}
• Bundle Version: {{.IOSConfig.BundleVersion}}
• Short Version: {{.IOSConfig.BundleShortVersion}}
• Deployment Target: {{.IOSConfig.DeploymentTarget}}
• Requires Full Screen: {{.IOSConfig.RequiresFullScreen}}
• Encryption Declared: {{if .IOSConfig.EncryptionDeclared}}✅ Yes{{else}}⚠️ No{{end}}

FLUTTER CONFIGURATION:
• Uses Material 3: {{if .FlutterConfig.UsesMaterial3}}✅ Yes{{else}}⚠️ No{{end}}
• Has Linter: {{if .FlutterConfig.HasLinter}}✅ Yes{{else}}⚠️ No{{end}}
• Has Icon Config: {{if .FlutterConfig.HasIconConfig}}✅ Yes{{else}}⚠️ No{{end}}
• Has Splash Config: {{if .FlutterConfig.HasSplashConfig}}✅ Yes{{else}}⚠️ No{{end}}
• Has Deprecated Packages: {{if .FlutterConfig.HasDeprecatedPkg}}❌ Yes{{else}}✅ No{{end}}

═══════════════════════════════════════════════════════════════════════════════
PERMISSIONS ANALYSIS
═══════════════════════════════════════════════════════════════════════════════

{{if .AndroidPermissions}}
Android Permissions:
{{range .AndroidPermissions}}
• {{.Name}}: {{if .Present}}✅ Declared{{else}}❌ Missing{{end}}{{if .HasDescription}} (with description){{else}}{{if .Present}} ⚠️ (no description){{end}}{{end}}
{{end}}
{{end}}

{{if .IOSPermissions}}
iOS Permissions (Info.plist):
{{range .IOSPermissions}}
• {{.Name}}: {{if .Present}}✅ Present{{else}}❌ Missing{{end}}{{if .Description}} - "{{.Description}}"{{end}}
{{end}}
{{end}}

═══════════════════════════════════════════════════════════════════════════════
DEPENDENCIES AUDIT
═══════════════════════════════════════════════════════════════════════════════

{{if .HighRiskDeps}}
HIGH RISK DEPENDENCIES:
{{range .HighRiskDeps}}
• {{.Name}}: {{.Version}} - {{.RiskLevel}}
{{end}}
{{else}}
✅ No high-risk dependencies detected
{{end}}

{{if .OutdatedDeps}}
OUTDATED DEPENDENCIES:
{{range .OutdatedDeps}}
• {{.Name}}: {{.Version}} → {{.LatestVer}}
{{end}}
{{else}}
✅ All dependencies up to date
{{end}}

Main Dependencies:
{{range .Dependencies}}
• {{.Name}}: {{.Version}}
{{end}}

Dev Dependencies:
{{range .DevDependencies}}
• {{.Name}}: {{.Version}}
{{end}}

═══════════════════════════════════════════════════════════════════════════════
SECURITY & POLICY FLAGS
═══════════════════════════════════════════════════════════════════════════════

{{if .SecurityFlags}}
Security Flags:
{{range .SecurityFlags}}
• {{.Type}}: {{if .Present}}⚠️ {{.Description}}{{else}}✅ OK{{end}}
{{end}}
{{end}}

{{if .PolicyFlags}}
Policy Flags:
{{range .PolicyFlags}}
• {{.Type}}: {{if .Present}}✅ Found{{else}}❌ Missing{{end}} - {{.Description}}
{{end}}
{{end}}

{{if .MissingPolicies}}
Missing Policy Elements:
{{range .MissingPolicies}}
• {{.}}
{{end}}
{{end}}

{{if .ReviewerAccount}}
═══════════════════════════════════════════════════════════════════════════════
REVIEWER TEST ACCOUNT
═══════════════════════════════════════════════════════════════════════════════

Configured: {{if .ReviewerAccount.Configured}}✅ Yes{{else}}⚠️ No{{end}}
{{if .ReviewerAccount.Configured}}
Email: {{.ReviewerAccount.Email}}
Password: {{.ReviewerAccount.PasswordEnv}} (environment variable)
Weak Password: {{if .ReviewerAccount.HasWeakPassword}}⚠️ Yes{{else}}✅ No{{end}}
Can Verify: {{if .ReviewerAccount.CanVerify}}✅ Yes{{else}}⚠️ No{{end}}
{{if .ReviewerAccount.LoginEndpoint}}Login Endpoint: {{.ReviewerAccount.LoginEndpoint}}{{end}}
{{end}}
{{end}}

═══════════════════════════════════════════════════════════════════════════════
COMPLIANCE QUESTIONS
═══════════════════════════════════════════════════════════════════════════════

Please analyze and provide detailed answers:

1. STORE SUBMISSION READINESS
   - Is this app ready for App Store submission? Why or why not?
   - Is this app ready for Play Store submission? Why or why not?
   - What's the estimated approval probability for each store?
   - How long will fixes likely take?

2. CRITICAL FIXES REQUIRED (Top 5)
   For each blocker, provide:
   - Issue description
   - Exact file to edit
   - Code fix or configuration change
   - Time estimate

3. PERMISSIONS REVIEW
   - Are all requested permissions justified?
   - Are the permission descriptions adequate for reviewers?
   - Any permissions that might trigger extra review scrutiny?

4. APP STORE CONNECT / PLAY CONSOLE SETUP
   - What app category fits best?
   - What content rating is appropriate?
   - What compliance documents are needed?
   - Any special declarations required?

5. PRIVACY & POLICY REQUIREMENTS
   - What privacy policy clauses are needed?
   - Is a terms of service required?
   - Data deletion capability requirements?
   - GDPR/CCPA compliance considerations?

6. REVIEWER CREDENTIALS & TESTING
   - Should we create test accounts for reviewers?
   - What login credentials should we provide?
   - What demo data should be pre-populated?
   - Any special instructions for reviewers?

7. COMMON REJECTION PREVENTION
   - What are the most likely rejection reasons?
   - How can we prevent them?
   - Any "red flags" in the current configuration?

═══════════════════════════════════════════════════════════════════════════════
RESPONSE FORMAT
═══════════════════════════════════════════════════════════════════════════════

Please format your response as:

## Executive Summary
(Brief 2-3 sentence overview)

## Store Readiness Assessment
(App Store: ✅/❌ with reasoning)
(Play Store: ✅/❌ with reasoning)

## Critical Fixes (In Priority Order)
1. [Priority: HIGH/MEDIUM/LOW]
   - Issue: ...
   - Fix: ...
   - File: ...
   - Code: ...

## Permissions Review
(Analysis and recommendations)

## Setup Recommendations
(App Store Connect & Play Console guidance)

## Privacy & Legal Requirements
(What documents/policies are needed)

## Reviewer Instructions
(What to tell reviewers, test account setup)

## Final Checklist
- [ ] Fix 1
- [ ] Fix 2
...`

// SummaryPromptTemplate is a shorter prompt for quick analysis
const SummaryPromptTemplate = `Analyze this Flutter app's compliance data for App Store and Play Store submission:

APP: {{.AppName}} v{{.Version}}
SCORE: {{.ComplianceScore}}/100
STATUS: {{.Status}}

FINDINGS:
- High: {{.HighCount}}
- Warnings: {{.WarningCount}}
- Info: {{.InfoCount}}
- Passed: {{.PassedChecks}}/{{.TotalChecks}}

{{if .Blockers}}BLOCKERS:{{range .Blockers}}
- {{.ID}}: {{.Title}} ({{.Platform}}){{end}}{{end}}

Provide:
1. Is it ready for submission? (App Store + Play Store)
2. Top 3 fixes needed
3. Estimated time to fix
4. Approval probability`

// PromptTemplateType represents different prompt template types
type PromptTemplateType string

const (
	TemplateTypeComprehensive PromptTemplateType = "comprehensive"
	TemplateTypeSummary       PromptTemplateType = "summary"
)

// GetTemplate returns the template string for the given type
func GetTemplate(t PromptTemplateType) string {
	switch t {
	case TemplateTypeSummary:
		return SummaryPromptTemplate
	case TemplateTypeComprehensive:
	default:
		return MainPromptTemplate
	}
	return MainPromptTemplate
}
