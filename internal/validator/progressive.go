package validator

import "strings"

// statusProfile returns severity knobs based on document status.
func statusProfile(status string) (emptyStrict bool, missingContextStrict bool) {
	s := strings.ToLower(status)
	switch s {
	case "published":
		return true, true
	case "review":
		return true, false
	default: // draft or unknown
		return false, false
	}
}
