package index

import (
	"fmt"
	"os"
	"path/filepath"

	"docflow/pkg/parser"
	"docflow/pkg/version"

	"gopkg.in/yaml.v3"
)

// VersionBump describes a single automatic version update.
type VersionBump struct {
	DocID  string
	Path   string
	From   string
	To     string
	Reason DocChangeType
}

// AutoBumpVersions updates frontmatter version based on detected changes.
// Content change → minor bump, metadata change → patch bump.
func AutoBumpVersions(changes []DocChange, docsRoot string) ([]VersionBump, error) {
	var bumps []VersionBump

	for _, c := range changes {
		if c.Type != ChangeContent && c.Type != ChangeMetadata {
			continue
		}
		abs := filepath.Join(docsRoot, c.Path)
		data, err := os.ReadFile(abs)
		if err != nil {
			return bumps, fmt.Errorf("read %s: %w", abs, err)
		}
		res, err := parser.ParseFrontmatterString(string(data), c.Path)
		if err != nil || res == nil || !res.HasFrontmatter || res.Frontmatter == nil {
			continue
		}

		fm := res.Frontmatter
		raw := map[string]interface{}{}
		if res.RawYAML != "" {
			_ = yaml.Unmarshal([]byte(res.RawYAML), &raw)
		}

		current := fm.Version
		if current == "" {
			current = "0.1.0"
		}

		var next string
		if c.Type == ChangeContent {
			next = version.BumpMinor(current)
		} else {
			next = version.BumpPatch(current)
		}

		raw["version"] = next
		newYaml, err := yaml.Marshal(raw)
		if err != nil {
			return bumps, fmt.Errorf("marshal frontmatter %s: %w", c.Path, err)
		}

		out := fmt.Sprintf("---\n%s---\n%s", string(newYaml), res.Body)
		if err := os.WriteFile(abs, []byte(out), 0o644); err != nil {
			return bumps, fmt.Errorf("write %s: %w", abs, err)
		}

		bumps = append(bumps, VersionBump{
			DocID:  c.DocID,
			Path:   c.Path,
			From:   current,
			To:     next,
			Reason: c.Type,
		})
	}

	return bumps, nil
}
