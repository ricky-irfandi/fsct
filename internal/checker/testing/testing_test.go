package testing

import (
	"testing"

	"github.com/ricky-irfandi/fsct/internal/checker"
)

func TestTestDirectoryCheck_ID(t *testing.T) {
	c := &TestDirectoryCheck{}
	if c.ID() != "TST-001" {
		t.Errorf("expected TST-001, got %s", c.ID())
	}
}

func TestTestDirectoryCheck_Run(t *testing.T) {
	c := &TestDirectoryCheck{}

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

func TestTestFileNamingCheck_ID(t *testing.T) {
	c := &TestFileNamingCheck{}
	if c.ID() != "TST-002" {
		t.Errorf("expected TST-002, got %s", c.ID())
	}
}

func TestTestFileNamingCheck_Run(t *testing.T) {
	c := &TestFileNamingCheck{}

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

func TestTestCoverageCheck_ID(t *testing.T) {
	c := &TestCoverageCheck{}
	if c.ID() != "TST-003" {
		t.Errorf("expected TST-003, got %s", c.ID())
	}
}

func TestTestCoverageCheck_Run(t *testing.T) {
	c := &TestCoverageCheck{}

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

func TestWidgetTestCheck_ID(t *testing.T) {
	c := &WidgetTestCheck{}
	if c.ID() != "TST-004" {
		t.Errorf("expected TST-004, got %s", c.ID())
	}
}

func TestWidgetTestCheck_Run(t *testing.T) {
	c := &WidgetTestCheck{}

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

func TestMockDependenciesCheck_ID(t *testing.T) {
	c := &MockDependenciesCheck{}
	if c.ID() != "TST-005" {
		t.Errorf("expected TST-005, got %s", c.ID())
	}
}

func TestMockDependenciesCheck_Run(t *testing.T) {
	c := &MockDependenciesCheck{}

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

func TestGoldenTestCheck_ID(t *testing.T) {
	c := &GoldenTestCheck{}
	if c.ID() != "TST-006" {
		t.Errorf("expected TST-006, got %s", c.ID())
	}
}

func TestGoldenTestCheck_Run(t *testing.T) {
	c := &GoldenTestCheck{}

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
