package version

import "testing"

func TestBumpMinorPreservesPrefix(t *testing.T) {
	if got := BumpMinor("v1.2.3"); got != "v1.3.3" {
		t.Fatalf("expected v1.3.3, got %s", got)
	}
	if got := BumpMinor("1.2.3"); got != "1.3.3" {
		t.Fatalf("expected 1.3.3, got %s", got)
	}
}

func TestBumpPatchPreservesPrefix(t *testing.T) {
	if got := BumpPatch("v0.1.0"); got != "v0.1.1" {
		t.Fatalf("expected v0.1.1, got %s", got)
	}
	if got := BumpPatch("0.1.0"); got != "0.1.1" {
		t.Fatalf("expected 0.1.1, got %s", got)
	}
}
