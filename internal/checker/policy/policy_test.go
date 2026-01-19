package policy

import (
	"testing"

	"github.com/ricky-irfandi/fsct/internal/checker"
)

func TestPrivacyPolicyCheck_ID(t *testing.T) {
	c := &PrivacyPolicyCheck{}
	if c.ID() != "POL-001" {
		t.Errorf("expected POL-001, got %s", c.ID())
	}
}

func TestPrivacyPolicyCheck_Name(t *testing.T) {
	c := &PrivacyPolicyCheck{}
	if c.Name() != "Privacy Policy URL" {
		t.Errorf("unexpected name: %s", c.Name())
	}
}

func TestPrivacyPolicyCheck_Run(t *testing.T) {
	c := &PrivacyPolicyCheck{}

	t.Run("no dart files", func(t *testing.T) {
		project := &checker.Project{}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("runs without panic", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"test.dart"},
		}
		results := c.Run(project)
		if results == nil {
			t.Error("expected non-nil results")
		}
	})
}

func TestTermsOfServiceCheck_ID(t *testing.T) {
	c := &TermsOfServiceCheck{}
	if c.ID() != "POL-002" {
		t.Errorf("expected POL-002, got %s", c.ID())
	}
}

func TestTermsOfServiceCheck_Run(t *testing.T) {
	c := &TermsOfServiceCheck{}

	t.Run("no dart files", func(t *testing.T) {
		project := &checker.Project{}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("runs without panic", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"test.dart"},
		}
		results := c.Run(project)
		if results == nil {
			t.Error("expected non-nil results")
		}
	})
}

func TestDataDeletionCheck_ID(t *testing.T) {
	c := &DataDeletionCheck{}
	if c.ID() != "POL-003" {
		t.Errorf("expected POL-003, got %s", c.ID())
	}
}

func TestDataDeletionCheck_Run(t *testing.T) {
	c := &DataDeletionCheck{}

	t.Run("no dart files", func(t *testing.T) {
		project := &checker.Project{}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("runs without panic", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"test.dart"},
		}
		results := c.Run(project)
		if results == nil {
			t.Error("expected non-nil results")
		}
	})
}

func TestLogoutCheck_ID(t *testing.T) {
	c := &LogoutCheck{}
	if c.ID() != "POL-004" {
		t.Errorf("expected POL-004, got %s", c.ID())
	}
}

func TestLogoutCheck_Run(t *testing.T) {
	c := &LogoutCheck{}

	t.Run("no dart files", func(t *testing.T) {
		project := &checker.Project{}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("runs without panic", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"test.dart"},
		}
		results := c.Run(project)
		if results == nil {
			t.Error("expected non-nil results")
		}
	})
}

func TestAccountRecoveryCheck_ID(t *testing.T) {
	c := &AccountRecoveryCheck{}
	if c.ID() != "POL-005" {
		t.Errorf("expected POL-005, got %s", c.ID())
	}
}

func TestAccountRecoveryCheck_Run(t *testing.T) {
	c := &AccountRecoveryCheck{}

	t.Run("no dart files", func(t *testing.T) {
		project := &checker.Project{}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("runs without panic", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"test.dart"},
		}
		results := c.Run(project)
		if results == nil {
			t.Error("expected non-nil results")
		}
	})
}
