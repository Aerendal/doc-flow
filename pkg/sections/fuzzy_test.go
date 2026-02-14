package sections

import (
	"strings"
	"testing"
)

func TestCanonicalFor(t *testing.T) {
	aliases := map[string][]string{
		"Przegląd":   {"Overview", "Podsumowanie"},
		"Endpoints":  {"API Endpoints", "API"},
		"Deployment": {"Rollout"},
	}

	c, ok, exact := CanonicalFor("Overview", aliases)
	if !ok || exact || c != "Przegląd" {
		t.Fatalf("expected canonical Przegląd (alias), got %s ok=%v exact=%v", c, ok, exact)
	}
	c, ok, exact = CanonicalFor("Przegląd", aliases)
	if !ok || !exact || c != "Przegląd" {
		t.Fatalf("expected exact canonical Przegląd, got %s", c)
	}
	if _, ok, _ := CanonicalFor("Unknown", aliases); ok {
		t.Fatalf("unexpected match for Unknown")
	}
}

func TestRenameAliases(t *testing.T) {
	content := "# Title\n## Overview\nText\n## API Endpoints\n"
	aliases := map[string][]string{"Przegląd": {"Overview"}, "Endpoints": {"API Endpoints"}}
	newC, hits := RenameAliases(content, aliases)
	if len(hits) != 2 {
		t.Fatalf("hits=%d, want 2", len(hits))
	}
	if newC == content {
		t.Fatalf("content not changed")
	}
	if !strings.Contains(newC, "## Przegląd") || !strings.Contains(newC, "## Endpoints") {
		t.Fatalf("replacement failed: %s", newC)
	}
}
