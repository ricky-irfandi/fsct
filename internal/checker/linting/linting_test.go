package linting

import (
	"testing"

	"github.com/ricky-irfandi/fsct/internal/checker"
)

func TestAnalysisOptionsCheck_ID(t *testing.T) {
	c := &AnalysisOptionsCheck{}
	if c.ID() != "LINT-001" {
		t.Errorf("expected LINT-001, got %s", c.ID())
	}
}

func TestAnalysisOptionsCheck_Run(t *testing.T) {
	c := &AnalysisOptionsCheck{}

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

func TestLinterRulesCheck_ID(t *testing.T) {
	c := &LinterRulesCheck{}
	if c.ID() != "LINT-002" {
		t.Errorf("expected LINT-002, got %s", c.ID())
	}
}

func TestLinterRulesCheck_Run(t *testing.T) {
	c := &LinterRulesCheck{}

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

func TestStrongModeCheck_ID(t *testing.T) {
	c := &StrongModeCheck{}
	if c.ID() != "LINT-003" {
		t.Errorf("expected LINT-003, got %s", c.ID())
	}
}

func TestStrongModeCheck_Run(t *testing.T) {
	c := &StrongModeCheck{}

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

func TestFileNamingCheck_ID(t *testing.T) {
	c := &FileNamingCheck{}
	if c.ID() != "LINT-004" {
		t.Errorf("expected LINT-004, got %s", c.ID())
	}
}

func TestFileNamingCheck_Run(t *testing.T) {
	c := &FileNamingCheck{}

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

func TestStyleGuideCheck_ID(t *testing.T) {
	c := &StyleGuideCheck{}
	if c.ID() != "LINT-005" {
		t.Errorf("expected LINT-005, got %s", c.ID())
	}
}

func TestStyleGuideCheck_Run(t *testing.T) {
	c := &StyleGuideCheck{}

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

func TestPublicAPIDocCheck_ID(t *testing.T) {
	c := &PublicAPIDocCheck{}
	if c.ID() != "LINT-006" {
		t.Errorf("expected LINT-006, got %s", c.ID())
	}
}

func TestPublicAPIDocCheck_Run(t *testing.T) {
	c := &PublicAPIDocCheck{}

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

func TestIgnoreCommentsCheck_ID(t *testing.T) {
	c := &IgnoreCommentsCheck{}
	if c.ID() != "LINT-007" {
		t.Errorf("expected LINT-007, got %s", c.ID())
	}
}

func TestIgnoreCommentsCheck_Run(t *testing.T) {
	c := &IgnoreCommentsCheck{}

	t.Run("no dart files", func(t *testing.T) {
		project := &checker.Project{}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
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
