package code

import (
	"testing"

	"github.com/ricky-irfandi/fsct/internal/checker"
)

func TestFileLengthCheck_ID(t *testing.T) {
	c := &FileLengthCheck{}
	if c.ID() != "COD-001" {
		t.Errorf("expected COD-001, got %s", c.ID())
	}
}

func TestFileLengthCheck_Run(t *testing.T) {
	c := &FileLengthCheck{}

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

func TestClassLengthCheck_ID(t *testing.T) {
	c := &ClassLengthCheck{}
	if c.ID() != "COD-002" {
		t.Errorf("expected COD-002, got %s", c.ID())
	}
}

func TestClassLengthCheck_Run(t *testing.T) {
	c := &ClassLengthCheck{}

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

func TestMethodComplexityCheck_ID(t *testing.T) {
	c := &MethodComplexityCheck{}
	if c.ID() != "COD-003" {
		t.Errorf("expected COD-003, got %s", c.ID())
	}
}

func TestMethodComplexityCheck_Run(t *testing.T) {
	c := &MethodComplexityCheck{}

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

func TestNamingConventionCheck_ID(t *testing.T) {
	c := &NamingConventionCheck{}
	if c.ID() != "COD-004" {
		t.Errorf("expected COD-004, got %s", c.ID())
	}
}

func TestNamingConventionCheck_Run(t *testing.T) {
	c := &NamingConventionCheck{}

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

func TestImportOrganizationCheck_ID(t *testing.T) {
	c := &ImportOrganizationCheck{}
	if c.ID() != "COD-005" {
		t.Errorf("expected COD-005, got %s", c.ID())
	}
}

func TestImportOrganizationCheck_Run(t *testing.T) {
	c := &ImportOrganizationCheck{}

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

func TestCommentQualityCheck_ID(t *testing.T) {
	c := &CommentQualityCheck{}
	if c.ID() != "COD-006" {
		t.Errorf("expected COD-006, got %s", c.ID())
	}
}

func TestCommentQualityCheck_Run(t *testing.T) {
	c := &CommentQualityCheck{}

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

func TestCyclomaticComplexityCheck_ID(t *testing.T) {
	c := &CyclomaticComplexityCheck{}
	if c.ID() != "COD-007" {
		t.Errorf("expected COD-007, got %s", c.ID())
	}
}

func TestCyclomaticComplexityCheck_Run(t *testing.T) {
	c := &CyclomaticComplexityCheck{}

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

func TestDuplicateCodeCheck_ID(t *testing.T) {
	c := &DuplicateCodeCheck{}
	if c.ID() != "COD-008" {
		t.Errorf("expected COD-008, got %s", c.ID())
	}
}

func TestDuplicateCodeCheck_Run(t *testing.T) {
	c := &DuplicateCodeCheck{}

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
