package recommend

import (
	"strconv"
	"strings"
)

// versionValue converts semver-ish string like "v1.2.3" to numeric weight.
func versionValue(v string) int {
	if v == "" {
		return 0
	}
	v = strings.TrimPrefix(v, "v")
	parts := strings.Split(v, ".")
	val := 0
	for i := 0; i < len(parts) && i < 3; i++ {
		n, _ := strconv.Atoi(parts[i])
		shift := 16 * (2 - i)
		val += n << shift
	}
	return val
}
