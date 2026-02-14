package deps

import "testing"

func TestTopoStable(t *testing.T) {
	depsMap := map[string][]string{
		"c": {"a"},
		"b": {"a"},
		"a": {},
	}
	g := Build(depsMap)
	res := g.Topo()
	if len(res.Problems) != 0 {
		t.Fatalf("unexpected problems: %+v", res.Problems)
	}
	expected := []string{"a", "b", "c"}
	for i := range expected {
		if res.Order[i] != expected[i] {
			t.Fatalf("order[%d]=%s, want %s", i, res.Order[i], expected[i])
		}
	}
}

func TestTopoCycle(t *testing.T) {
	depsMap := map[string][]string{
		"a": {"b"},
		"b": {"a"},
	}
	g := Build(depsMap)
	res := g.Topo()
	if len(res.Problems) == 0 {
		t.Fatalf("expected cycle problem")
	}
	found := false
	for _, p := range res.Problems {
		if p.Type == "cycle" {
			found = true
		}
	}
	if !found {
		t.Fatalf("cycle not reported: %+v", res.Problems)
	}
}

func TestTopoMissing(t *testing.T) {
	depsMap := map[string][]string{
		"a": {"missing"},
	}
	g := Build(depsMap)
	res := g.Topo()
	if len(res.Problems) == 0 {
		t.Fatalf("expected missing problem")
	}
	if res.Problems[0].Type != "missing" {
		t.Fatalf("problem type = %s, want missing", res.Problems[0].Type)
	}
}
