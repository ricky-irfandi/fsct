package policy

import (
	"regexp"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

var (
	privacyPattern   = regexp.MustCompile(`(?i)(privacy|policy)[^\s]*\.?(com|org|net|io)`)
	tosPattern       = regexp.MustCompile(`(?i)(terms|service|conditions)[^\s]*\.?(com|org|net|io)`)
	deletionPatterns = []*regexp.Regexp{regexp.MustCompile(`(?i)delete.*data`), regexp.MustCompile(`(?i)data.*deletion`), regexp.MustCompile(`(?i)remove.*account`), regexp.MustCompile(`(?i)gdpr`), regexp.MustCompile(`(?i)ccpa`)}
	logoutPatterns   = []*regexp.Regexp{regexp.MustCompile(`(?i)signOut`), regexp.MustCompile(`(?i)logout`), regexp.MustCompile(`(?i)sign_?out`), regexp.MustCompile(`(?i)log_?out`)}
	recoveryPatterns = []*regexp.Regexp{regexp.MustCompile(`(?i)resetPassword`), regexp.MustCompile(`(?i)forgotPassword`), regexp.MustCompile(`(?i)recoverAccount`), regexp.MustCompile(`(?i)passwordReset`)}
)

type PrivacyPolicyCheck struct{}

func (c *PrivacyPolicyCheck) ID() string {
	return "POL-001"
}

func (c *PrivacyPolicyCheck) Name() string {
	return "Privacy Policy URL"
}

func (c *PrivacyPolicyCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	foundPrivacyURL := false

	for _, file := range project.DartFiles {
		if privacyPattern.MatchString(file) {
			foundPrivacyURL = true
			break
		}
	}

	if !foundPrivacyURL {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Privacy policy URL not found",
			"Info.plist or source files",
			"Add privacy policy URL to Info.plist (NSPrivacyPolicyURL) or in-app settings",
			report.SeverityHigh,
			0,
		))
	}

	return findings
}

type TermsOfServiceCheck struct{}

func (c *TermsOfServiceCheck) ID() string {
	return "POL-002"
}

func (c *TermsOfServiceCheck) Name() string {
	return "Terms of Service URL"
}

func (c *TermsOfServiceCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	foundTOSURL := false

	for _, file := range project.DartFiles {
		if tosPattern.MatchString(file) {
			foundTOSURL = true
			break
		}
	}

	if !foundTOSURL {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Terms of Service URL not found",
			"source files",
			"Add Terms of Service URL for compliance",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type DataDeletionCheck struct{}

func (c *DataDeletionCheck) ID() string {
	return "POL-003"
}

func (c *DataDeletionCheck) Name() string {
	return "Data Deletion Contact"
}

func (c *DataDeletionCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	hasDeletionContact := false

	for _, file := range project.DartFiles {
		for _, re := range deletionPatterns {
			if re.MatchString(file) {
				hasDeletionContact = true
				break
			}
		}
	}

	if !hasDeletionContact {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Data deletion contact URL not found",
			"source files",
			"Add data deletion contact URL (required for App Store compliance under GDPR/CCPA)",
			report.SeverityHigh,
			0,
		))
	}

	return findings
}

type LogoutCheck struct{}

func (c *LogoutCheck) ID() string {
	return "POL-004"
}

func (c *LogoutCheck) Name() string {
	return "Logout Functionality"
}

func (c *LogoutCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	hasLogout := false

	for _, file := range project.DartFiles {
		for _, re := range logoutPatterns {
			if re.MatchString(file) {
				hasLogout = true
				break
			}
		}
	}

	if !hasLogout {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Logout functionality not detected",
			"source files",
			"Implement logout functionality for user account management",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type AccountRecoveryCheck struct{}

func (c *AccountRecoveryCheck) ID() string {
	return "POL-005"
}

func (c *AccountRecoveryCheck) Name() string {
	return "Account Recovery Options"
}

func (c *AccountRecoveryCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	hasRecovery := false

	for _, file := range project.DartFiles {
		for _, re := range recoveryPatterns {
			if re.MatchString(file) {
				hasRecovery = true
				break
			}
		}
	}

	if !hasRecovery {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Account recovery options not detected",
			"source files",
			"Implement password reset or account recovery functionality",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}
