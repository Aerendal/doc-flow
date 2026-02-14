package cache

import (
	"crypto/sha256"
	"fmt"
	"os"
)

// FileHash returns short SHA256 (8 bytes hex) of the file content.
func FileHash(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum[:8]), nil
}
