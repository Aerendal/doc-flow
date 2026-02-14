package version

import (
	"fmt"
	"strings"
)

// BumpMinor increments minor component of semver-like MAJOR.MINOR.PATCH.
// Preserves presence of leading "v" if provided.
func BumpMinor(v string) string {
	hasV := strings.HasPrefix(v, "v")
	maj, min, patch := parse(v)
	min++
	if hasV {
		return fmt.Sprintf("v%d.%d.%d", maj, min, patch)
	}
	return fmt.Sprintf("%d.%d.%d", maj, min, patch)
}

// BumpPatch increments patch component.
func BumpPatch(v string) string {
	hasV := strings.HasPrefix(v, "v")
	maj, min, patch := parse(v)
	patch++
	if hasV {
		return fmt.Sprintf("v%d.%d.%d", maj, min, patch)
	}
	return fmt.Sprintf("%d.%d.%d", maj, min, patch)
}

func parse(v string) (int, int, int) {
	v = strings.TrimPrefix(v, "v")
	var maj, min, patch int
	fmt.Sscanf(v, "%d.%d.%d", &maj, &min, &patch)
	if maj == 0 && min == 0 && patch == 0 {
		maj = 1
	}
	return maj, min, patch
}
