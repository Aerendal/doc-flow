package sections

import "testing"

func TestComputeMetrics(t *testing.T) {
	content := `# T
## A
## B
`
	m := ComputeMetrics(content)
	if m.Total != 2 {
		t.Fatalf("total=%d, want 2", m.Total)
	}
	if m.Completeness != 1.0 {
		t.Fatalf("completeness=%f, want 1.0", m.Completeness)
	}
}

func TestComputeMetricsEmpty(t *testing.T) {
	content := `# T
## 
`
	m := ComputeMetrics(content)
	if m.Empty != 1 {
		t.Fatalf("empty=%d, want 1", m.Empty)
	}
	if m.Completeness >= 1.0 {
		t.Fatalf("completeness should be <1.0")
	}
}
