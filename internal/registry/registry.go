package registry

import (
	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/checker/android"
	"github.com/ricky-irfandi/fsct/internal/checker/code"
	"github.com/ricky-irfandi/fsct/internal/checker/docs"
	"github.com/ricky-irfandi/fsct/internal/checker/flutter"
	"github.com/ricky-irfandi/fsct/internal/checker/ios"
	"github.com/ricky-irfandi/fsct/internal/checker/linting"
	"github.com/ricky-irfandi/fsct/internal/checker/perf"
	"github.com/ricky-irfandi/fsct/internal/checker/policy"
	"github.com/ricky-irfandi/fsct/internal/checker/security"
	"github.com/ricky-irfandi/fsct/internal/checker/testing"
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
	r.registerCodeChecks()
	r.registerTestingChecks()
	r.registerLintingChecks()
	r.registerDocsChecks()
	r.registerPerfChecks()
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
	r.checks["FLT-002"] = &flutter.Material3Check{}
	r.checks["FLT-003"] = &flutter.MinSDKVersionCheck{}
	r.checks["FLT-004"] = &flutter.PackageNameCheck{}
	r.checks["FLT-005"] = &flutter.VersionCheck{}
	r.checks["FLT-006"] = &flutter.DependencyConstraintCheck{}
	r.checks["FLT-007"] = &flutter.DeprecatedPackageCheck{}
	r.checks["FLT-008"] = &flutter.ProjectStructureCheck{}
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

func (r *CheckerRegistry) registerCodeChecks() {
	r.checks["COD-001"] = &code.FileLengthCheck{}
	r.checks["COD-002"] = &code.ClassLengthCheck{}
	r.checks["COD-003"] = &code.MethodComplexityCheck{}
	r.checks["COD-004"] = &code.NamingConventionCheck{}
	r.checks["COD-005"] = &code.ImportOrganizationCheck{}
	r.checks["COD-006"] = &code.CommentQualityCheck{}
	r.checks["COD-007"] = &code.CyclomaticComplexityCheck{}
	r.checks["COD-008"] = &code.DuplicateCodeCheck{}
}

func (r *CheckerRegistry) registerTestingChecks() {
	r.checks["TST-001"] = &testing.TestDirectoryCheck{}
	r.checks["TST-002"] = &testing.TestFileNamingCheck{}
	r.checks["TST-003"] = &testing.TestCoverageCheck{}
	r.checks["TST-004"] = &testing.WidgetTestCheck{}
	r.checks["TST-005"] = &testing.MockDependenciesCheck{}
	r.checks["TST-006"] = &testing.GoldenTestCheck{}
}

func (r *CheckerRegistry) registerLintingChecks() {
	r.checks["LINT-001"] = &linting.AnalysisOptionsCheck{}
	r.checks["LINT-002"] = &linting.LinterRulesCheck{}
	r.checks["LINT-003"] = &linting.StrongModeCheck{}
	r.checks["LINT-004"] = &linting.FileNamingCheck{}
	r.checks["LINT-005"] = &linting.StyleGuideCheck{}
	r.checks["LINT-006"] = &linting.PublicAPIDocCheck{}
	r.checks["LINT-007"] = &linting.IgnoreCommentsCheck{}
}

func (r *CheckerRegistry) registerDocsChecks() {
	r.checks["DOC-001"] = &docs.ReadmePresenceCheck{}
	r.checks["DOC-002"] = &docs.ReadmeContentCheck{}
	r.checks["DOC-003"] = &docs.ChangelogPresenceCheck{}
	r.checks["DOC-004"] = &docs.LicensePresenceCheck{}
	r.checks["DOC-005"] = &docs.ApiDocumentationCheck{}
	r.checks["DOC-006"] = &docs.CodeCommentsCheck{}
}

func (r *CheckerRegistry) registerPerfChecks() {
	r.checks["PERF-001"] = &perf.ConstConstructorCheck{}
	r.checks["PERF-002"] = &perf.BuildOptimizationCheck{}
	r.checks["PERF-003"] = &perf.ListBuilderCheck{}
	r.checks["PERF-004"] = &perf.ImageOptimizationCheck{}
	r.checks["PERF-005"] = &perf.StateManagementCheck{}
	r.checks["PERF-006"] = &perf.DependencyOptimizationCheck{}
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
		"Code Quality", "Testing", "Linting", "Documentation", "Performance",
	}
}
