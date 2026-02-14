package fix

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func UnifiedDiff(path, before, after string) string {
	if before == after {
		return ""
	}

	oldLines := toLinesForDiff(before)
	newLines := toLinesForDiff(after)

	oldStart := 1
	if len(oldLines) == 0 {
		oldStart = 0
	}
	newStart := 1
	if len(newLines) == 0 {
		newStart = 0
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("--- a/%s\n", path))
	b.WriteString(fmt.Sprintf("+++ b/%s\n", path))
	b.WriteString(fmt.Sprintf("@@ -%s +%s @@\n", formatRange(oldStart, len(oldLines)), formatRange(newStart, len(newLines))))
	for _, line := range oldLines {
		b.WriteString("-")
		b.WriteString(line)
		b.WriteString("\n")
	}
	for _, line := range newLines {
		b.WriteString("+")
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}

func checksum(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:8])
}

func toLinesForDiff(content string) []string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	if content == "" {
		return nil
	}
	lines := strings.Split(content, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func formatRange(start, count int) string {
	if count == 1 {
		return fmt.Sprintf("%d", start)
	}
	return fmt.Sprintf("%d,%d", start, count)
}
