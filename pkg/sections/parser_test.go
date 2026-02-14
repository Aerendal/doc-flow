package sections

import "testing"

func TestParseThreeLevels(t *testing.T) {
	content := `# Title

## A
### A1
#### A1a
### A2

## B
### B1
`
	tree := Parse(content)

	if len(tree.Root.Children) != 2 {
		t.Fatalf("root children = %d, want 2", len(tree.Root.Children))
	}

	a := tree.Root.Children[0]
	if a.Text != "A" || a.Level != 2 {
		t.Fatalf("first child = (%s,%d), want (A,2)", a.Text, a.Level)
	}

	if len(a.Children) != 2 {
		t.Fatalf("A children = %d, want 2", len(a.Children))
	}

	a1 := a.Children[0]
	if a1.Text != "A1" || a1.Level != 3 {
		t.Fatalf("A1 node mismatch: (%s,%d)", a1.Text, a1.Level)
	}
	if len(a1.Children) != 1 || a1.Children[0].Text != "A1a" || a1.Children[0].Level != 4 {
		t.Fatalf("A1a child missing or invalid")
	}

	a2 := a.Children[1]
	if a2.Text != "A2" || a2.Level != 3 {
		t.Fatalf("A2 node mismatch: (%s,%d)", a2.Text, a2.Level)
	}

	b := tree.Root.Children[1]
	if b.Text != "B" || b.Level != 2 {
		t.Fatalf("B node mismatch: (%s,%d)", b.Text, b.Level)
	}
	if len(b.Children) != 1 || b.Children[0].Text != "B1" || b.Children[0].Level != 3 {
		t.Fatalf("B1 child missing or invalid")
	}
}

func TestParseSkipsMissingLevels(t *testing.T) {
	content := `# Title
### Deep
`
	tree := Parse(content)

	if len(tree.Root.Children) != 1 {
		t.Fatalf("root children = %d, want 1", len(tree.Root.Children))
	}
	deep := tree.Root.Children[0]
	if deep.Level != 3 {
		t.Fatalf("deep level = %d, want 3", deep.Level)
	}
	if deep.Text != "Deep" {
		t.Fatalf("deep text = %s, want Deep", deep.Text)
	}
	if deep.Children != nil && len(deep.Children) != 0 {
		t.Fatalf("deep should have no children, got %d", len(deep.Children))
	}
}
