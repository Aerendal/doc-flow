package rules

import (
	"os"
	"testing"
)

func TestLoadFamilyRules(t *testing.T) {
	yaml := `families:
  - doc_type: specification
    required_deps: [product_vision_statement]
    optional_deps: [system_architecture]
`
	tmp := t.TempDir() + "/rules.yaml"
	if err := os.WriteFile(tmp, []byte(yaml), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	r, err := LoadFamilyRules(tmp)
	if err != nil {
		t.Fatalf("LoadFamilyRules error: %v", err)
	}
	req := r.RequiredFor("specification")
	if len(req) != 1 || req[0] != "product_vision_statement" {
		t.Fatalf("required = %v, want product_vision_statement", req)
	}
	if r.RequiredFor("guide") != nil {
		t.Fatalf("expected nil for unknown doc_type")
	}
}
