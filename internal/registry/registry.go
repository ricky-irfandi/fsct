package registry

import (
	aipkg "github.com/ricky-irfandi/fsct/internal/ai"
	"github.com/ricky-irfandi/fsct/internal/checker"
	aichecks "github.com/ricky-irfandi/fsct/internal/checker/ai"
	"github.com/ricky-irfandi/fsct/internal/checker/android"
	"github.com/ricky-irfandi/fsct/internal/checker/flutter"
	"github.com/ricky-irfandi/fsct/internal/checker/ios"
	"github.com/ricky-irfandi/fsct/internal/checker/policy"
	"github.com/ricky-irfandi/fsct/internal/checker/reviewer"
	"github.com/ricky-irfandi/fsct/internal/checker/security"
	"github.com/ricky-irfandi/fsct/internal/config"
)

type CheckerRegistry struct {
	checks map[string]checker.Check
}

func NewRegistry() *CheckerRegistry {
	return &CheckerRegistry{
		checks: make(map[string]checker.Check),
	}
}

func (r *CheckerRegistry) RegisterAll() {
	r.registerAndroidChecks()
	r.registerIOSChecks()
	r.registerFlutterChecks()
	r.registerSecurityChecks()
	r.registerPolicyChecks()
}

func (r *CheckerRegistry) registerAndroidChecks() {
	r.checks["AND-001"] = &android.TargetSDKCheck{}
	r.checks["AND-002"] = &android.MinSDKCheck{}
	r.checks["AND-003"] = &android.InternetPermissionCheck{}
	r.checks["AND-004"] = &android.DangerousPermissionsCheck{}
	r.checks["AND-005"] = &android.DebuggableCheck{}
	r.checks["AND-006"] = &android.ExportedAttributeCheck{}
	r.checks["AND-007"] = &android.MissingAppIconCheck{}
	r.checks["AND-008"] = &android.PlaceholderIconCheck{}
	r.checks["AND-009"] = &android.ApplicationIDCheck{}
	r.checks["AND-010"] = &android.VersionCodeCheck{}
	r.checks["AND-011"] = &android.PackageVisibilityCheck{}
	r.checks["AND-012"] = &android.AllowBackupCheck{}
}

func (r *CheckerRegistry) registerIOSChecks() {
	r.checks["IOS-001"] = &ios.CameraUsageDescriptionCheck{}
	r.checks["IOS-002"] = &ios.PhotoLibraryUsageDescriptionCheck{}
	r.checks["IOS-003"] = &ios.LocationUsageDescriptionCheck{}
	r.checks["IOS-004"] = &ios.MicrophoneUsageDescriptionCheck{}
	r.checks["IOS-005"] = &ios.ContactsUsageDescriptionCheck{}
	r.checks["IOS-006"] = &ios.CalendarsUsageDescriptionCheck{}
	r.checks["IOS-007"] = &ios.EmptyUsageDescriptionCheck{}
	r.checks["IOS-008"] = &ios.MissingAppIconCheck{}
	r.checks["IOS-009"] = &ios.Missing1024IconCheck{}
	r.checks["IOS-010"] = &ios.FullScreenConflictCheck{}
	r.checks["IOS-011"] = &ios.EncryptionDeclarationCheck{}
	r.checks["IOS-012"] = &ios.DeploymentTargetCheck{}
}

func (r *CheckerRegistry) registerFlutterChecks() {
	r.checks["FLT-001"] = &flutter.FlutterSDKVersionCheck{}
	r.checks["FLT-003"] = &flutter.MinSDKVersionCheck{}
	r.checks["FLT-004"] = &flutter.PackageNameCheck{}
	r.checks["FLT-005"] = &flutter.VersionCheck{}
}

func (r *CheckerRegistry) registerSecurityChecks() {
	r.checks["SEC-001"] = &security.HardcodedCredentialsCheck{}
	r.checks["SEC-002"] = &security.DebugModeCheck{}
	r.checks["SEC-003"] = &security.InsecureHTTPCheck{}
	r.checks["SEC-004"] = &security.ExportedActivityCheck{}
	r.checks["SEC-005"] = &security.SQLInjectionCheck{}
}

func (r *CheckerRegistry) registerPolicyChecks() {
	r.checks["POL-001"] = &policy.PrivacyPolicyCheck{}
	r.checks["POL-002"] = &policy.TermsOfServiceCheck{}
	r.checks["POL-003"] = &policy.DataDeletionCheck{}
	r.checks["POL-004"] = &policy.LogoutCheck{}
	r.checks["POL-005"] = &policy.AccountRecoveryCheck{}
}

// RegisterAIChecks registers AI-powered checks if AI client is available
func (r *CheckerRegistry) RegisterAIChecks(client *aipkg.Client) {
	if client == nil || !client.IsAvailable() {
		return
	}

	r.checks["AI-001"] = aichecks.AI001PermissionJustificationCheck(client)
	r.checks["AI-002"] = aichecks.AI002PolicyComplianceCheck(client)
	r.checks["AI-003"] = aichecks.AI003DependencyRiskCheck(client)
	r.checks["AI-004"] = aichecks.AI004StoreGuidanceCheck(client)
	r.checks["AI-005"] = aichecks.AI005ReviewerNotesCheck(client)
}

func (r *CheckerRegistry) Get(id string) (checker.Check, bool) {
	check, ok := r.checks[id]
	return check, ok
}

func (r *CheckerRegistry) GetAll() []checker.Check {
	checks := make([]checker.Check, 0, len(r.checks))
	for _, check := range r.checks {
		checks = append(checks, check)
	}
	return checks
}

func (r *CheckerRegistry) GetAllByID() map[string]checker.Check {
	return r.checks
}

func (r *CheckerRegistry) Count() int {
	return len(r.checks)
}

func (r *CheckerRegistry) GetCategories() []string {
	return []string{
		"Android", "iOS", "Flutter", "Security", "Policy",
		"AI Analysis", "Reviewer",
	}
}

// HasAIChecks returns true if AI checks are registered
func (r *CheckerRegistry) HasAIChecks() bool {
	_, hasAI001 := r.checks["AI-001"]
	return hasAI001
}

// RegisterReviewerChecks registers reviewer verification checks
func (r *CheckerRegistry) RegisterReviewerChecks(cfg *config.ReviewerConfig) {
	if cfg == nil {
		cfg = reviewer.GetConfigFromEnv()
	}

	// Always register basic credential checks
	r.checks["REV-001"] = reviewer.NewNoCredentialsCheck(cfg)
	r.checks["REV-002"] = reviewer.NewPlaceholderEmailCheck(cfg)
	r.checks["REV-003"] = reviewer.NewWeakPasswordCheck(cfg)

	// Register login verification only if enabled
	if cfg != nil && cfg.Verification != nil && cfg.Verification.Enabled {
		r.checks["REV-004"] = reviewer.NewLoginVerificationCheck(cfg)
	}
}
