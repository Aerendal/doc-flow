package graph

import (
	"encoding/json"
	"reflect"
	"testing"

	execache "docflow/internal/exec/cache"
)

func TestBuildGraphDeterministicAcrossFactsOrder(t *testing.T) {
	facts := []execache.DocumentFacts{
		{
			Path:           "docs/a.md",
			DocID:          "A",
			TemplateSource: "T",
			DependsOn:      []string{"B", "MISSING"},
			Content:        "[ref](b.md)",
		},
		{
			Path:  "docs/b.md",
			DocID: "B",
		},
		{
			Path:    "templates/t.md",
			DocID:   "T",
			DocType: "template",
		},
	}

	g1, issues1 := BuildGraph(facts, BuildGraphOptions{
		IncludeLinks:     true,
		IncludeTemplates: true,
		Deterministic:    true,
		Portable:         true,
	})

	reversed := append([]execache.DocumentFacts(nil), facts...)
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}

	g2, issues2 := BuildGraph(reversed, BuildGraphOptions{
		IncludeLinks:     true,
		IncludeTemplates: true,
		Deterministic:    true,
		Portable:         true,
	})

	enc1, err := json.Marshal(g1)
	if err != nil {
		t.Fatalf("marshal graph 1: %v", err)
	}
	enc2, err := json.Marshal(g2)
	if err != nil {
		t.Fatalf("marshal graph 2: %v", err)
	}
	if string(enc1) != string(enc2) {
		t.Fatalf("graph output is not deterministic")
	}

	iss1, err := json.Marshal(issues1)
	if err != nil {
		t.Fatalf("marshal issues 1: %v", err)
	}
	iss2, err := json.Marshal(issues2)
	if err != nil {
		t.Fatalf("marshal issues 2: %v", err)
	}
	if string(iss1) != string(iss2) {
		t.Fatalf("graph issues are not deterministic")
	}
}

func TestBuildGraphCyclesCanonical(t *testing.T) {
	facts := []execache.DocumentFacts{
		{Path: "a.md", DocID: "A", DependsOn: []string{"B"}},
		{Path: "b.md", DocID: "B", DependsOn: []string{"C"}},
		{Path: "c.md", DocID: "C", DependsOn: []string{"A"}},
		{Path: "d.md", DocID: "D", DependsOn: []string{"D"}},
	}
	g, _ := BuildGraph(facts, BuildGraphOptions{Deterministic: true, Portable: true})

	expected := [][]string{
		{"doc:a.md", "doc:b.md", "doc:c.md", "doc:a.md"},
		{"doc:d.md", "doc:d.md"},
	}
	if !reflect.DeepEqual(g.Cycles, expected) {
		t.Fatalf("unexpected cycles: got=%v want=%v", g.Cycles, expected)
	}
}

func TestImpactReachableDocPaths(t *testing.T) {
	facts := []execache.DocumentFacts{
		{Path: "templates/t.md", DocID: "T", DocType: "template"},
		{Path: "docs/a.md", DocID: "A", TemplateSource: "T"},
		{Path: "docs/b.md", DocID: "B", TemplateSource: "T"},
		{Path: "docs/c.md", DocID: "C", DependsOn: []string{"A"}},
	}
	g, _ := BuildGraph(facts, BuildGraphOptions{
		IncludeTemplates: true,
		Deterministic:    true,
		Portable:         true,
	})

	start := NodeIDsForDocID(g, "T")
	docs := ReachableDocPaths(g, start, ImpactOptions{
		MaxDepth:     1,
		IncludeStart: false,
		EdgeKinds: map[EdgeKind]bool{
			EdgeUsesTemplate: true,
		},
	})
	wantDirect := []string{"docs/a.md", "docs/b.md"}
	if !reflect.DeepEqual(docs, wantDirect) {
		t.Fatalf("unexpected direct impact: got=%v want=%v", docs, wantDirect)
	}

	startA := NodeIDsForDocID(g, "A")
	reverseDeps := ReachableDocPaths(g, startA, ImpactOptions{
		MaxDepth:       -1,
		IncludeReverse: true,
		IncludeStart:   false,
		EdgeKinds: map[EdgeKind]bool{
			EdgeDependsOn: true,
		},
	})
	wantReverse := []string{"docs/c.md"}
	if !reflect.DeepEqual(reverseDeps, wantReverse) {
		t.Fatalf("unexpected reverse impact: got=%v want=%v", reverseDeps, wantReverse)
	}
}
