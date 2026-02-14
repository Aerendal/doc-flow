package recommend

import "testing"

func TestVersionValue(t *testing.T) {
	if versionValue("v1.2.3") <= versionValue("v1.0.0") {
		t.Fatalf("expected 1.2.3 > 1.0.0")
	}
	if versionValue("") != 0 {
		t.Fatalf("empty version should be 0")
	}
}
