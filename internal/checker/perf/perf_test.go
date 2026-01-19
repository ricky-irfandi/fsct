package perf

import (
	"testing"

	"github.com/ricky-irfandi/fsct/internal/checker"
)

func TestConstConstructorCheck_ID(t *testing.T) {
	c := &ConstConstructorCheck{}
	if c.ID() != "PERF-001" {
		t.Errorf("expected PERF-001, got %s", c.ID())
	}
}

func TestConstConstructorCheck_Run(t *testing.T) {
	c := &ConstConstructorCheck{}

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

func TestBuildOptimizationCheck_ID(t *testing.T) {
	c := &BuildOptimizationCheck{}
	if c.ID() != "PERF-002" {
		t.Errorf("expected PERF-002, got %s", c.ID())
	}
}

func TestBuildOptimizationCheck_Run(t *testing.T) {
	c := &BuildOptimizationCheck{}

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

func TestListBuilderCheck_ID(t *testing.T) {
	c := &ListBuilderCheck{}
	if c.ID() != "PERF-003" {
		t.Errorf("expected PERF-003, got %s", c.ID())
	}
}

func TestListBuilderCheck_Run(t *testing.T) {
	c := &ListBuilderCheck{}

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

func TestImageOptimizationCheck_ID(t *testing.T) {
	c := &ImageOptimizationCheck{}
	if c.ID() != "PERF-004" {
		t.Errorf("expected PERF-004, got %s", c.ID())
	}
}

func TestImageOptimizationCheck_Run(t *testing.T) {
	c := &ImageOptimizationCheck{}

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

func TestStateManagementCheck_ID(t *testing.T) {
	c := &StateManagementCheck{}
	if c.ID() != "PERF-005" {
		t.Errorf("expected PERF-005, got %s", c.ID())
	}
}

func TestStateManagementCheck_Run(t *testing.T) {
	c := &StateManagementCheck{}

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

func TestDependencyOptimizationCheck_ID(t *testing.T) {
	c := &DependencyOptimizationCheck{}
	if c.ID() != "PERF-006" {
		t.Errorf("expected PERF-006, got %s", c.ID())
	}
}

func TestDependencyOptimizationCheck_Run(t *testing.T) {
	c := &DependencyOptimizationCheck{}

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
