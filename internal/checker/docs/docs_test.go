package docs

import (
	"testing"

	"github.com/ricky-irfandi/fsct/internal/checker"
)

func TestReadmePresenceCheck_ID(t *testing.T) {
	c := &ReadmePresenceCheck{}
	if c.ID() != "DOC-001" {
		t.Errorf("expected DOC-001, got %s", c.ID())
	}
}

func TestReadmePresenceCheck_Run(t *testing.T) {
	c := &ReadmePresenceCheck{}

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

func TestReadmeContentCheck_ID(t *testing.T) {
	c := &ReadmeContentCheck{}
	if c.ID() != "DOC-002" {
		t.Errorf("expected DOC-002, got %s", c.ID())
	}
}

func TestReadmeContentCheck_Run(t *testing.T) {
	c := &ReadmeContentCheck{}

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

func TestChangelogPresenceCheck_ID(t *testing.T) {
	c := &ChangelogPresenceCheck{}
	if c.ID() != "DOC-003" {
		t.Errorf("expected DOC-003, got %s", c.ID())
	}
}

func TestChangelogPresenceCheck_Run(t *testing.T) {
	c := &ChangelogPresenceCheck{}

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

func TestLicensePresenceCheck_ID(t *testing.T) {
	c := &LicensePresenceCheck{}
	if c.ID() != "DOC-004" {
		t.Errorf("expected DOC-004, got %s", c.ID())
	}
}

func TestLicensePresenceCheck_Run(t *testing.T) {
	c := &LicensePresenceCheck{}

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

func TestApiDocumentationCheck_ID(t *testing.T) {
	c := &ApiDocumentationCheck{}
	if c.ID() != "DOC-005" {
		t.Errorf("expected DOC-005, got %s", c.ID())
	}
}

func TestApiDocumentationCheck_Run(t *testing.T) {
	c := &ApiDocumentationCheck{}

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

func TestCodeCommentsCheck_ID(t *testing.T) {
	c := &CodeCommentsCheck{}
	if c.ID() != "DOC-006" {
		t.Errorf("expected DOC-006, got %s", c.ID())
	}
}

func TestCodeCommentsCheck_Run(t *testing.T) {
	c := &CodeCommentsCheck{}

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
