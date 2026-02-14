package index

import (
	"reflect"
	"testing"
)

func TestMarkCyclesAnnotatesRecords(t *testing.T) {
	idx := New()
	a := &DocumentRecord{DocID: "a"}
	b := &DocumentRecord{DocID: "b"}
	c := &DocumentRecord{DocID: "c"}
	idx.Add(a)
	idx.Add(b)
	idx.Add(c)

	adj := map[string][]string{
		"a": {"b"},
		"b": {"c"},
		"c": {"a"},
	}

	markCycles(idx, adj)

	if len(a.Cycle) == 0 || len(b.Cycle) == 0 || len(c.Cycle) == 0 {
		t.Fatalf("expected cycle annotated on all nodes")
	}
}

func TestCanonicalizeCycleStable(t *testing.T) {
	tests := [][]string{
		{"b", "c", "a", "b"},
		{"c", "a", "b", "c"},
		{"a", "b", "c", "a"},
		{"a", "c", "b", "a"},
	}

	want := []string{"a", "b", "c", "a"}
	for _, tc := range tests {
		got := canonicalizeCycle(tc)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("canonicalizeCycle(%v) = %v, want %v", tc, got, want)
		}
	}
}

func TestMarkCyclesDeterministicAcrossRuns(t *testing.T) {
	want := []string{"a", "b", "c", "a"}

	for i := 0; i < 20; i++ {
		idx := New()
		idx.Add(&DocumentRecord{DocID: "a"})
		idx.Add(&DocumentRecord{DocID: "b"})
		idx.Add(&DocumentRecord{DocID: "c"})

		adj := map[string][]string{
			"b": {"c"},
			"c": {"a"},
			"a": {"b"},
		}
		markCycles(idx, adj)

		if got := idx.GetByID("a").Cycle; !reflect.DeepEqual(got, want) {
			t.Fatalf("run %d: cycle = %v, want %v", i, got, want)
		}
	}
}
