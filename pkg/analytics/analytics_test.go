package analytics

import (
	"os"
	"path/filepath"
	"testing"

	"docflow/pkg/recommend"
)

func TestComputeAndExport(t *testing.T) {
	templates := []recommend.TemplateInfo{
		{Path: "a.md", Quality: 80, Usage: 5, Status: "active", Version: "v1", DocType: "guide", Language: "pl"},
		{Path: "b.md", Quality: 50, Usage: 1, Status: "active", Version: "v1", DocType: "spec", Language: "pl"},
	}
	stats := toStats(templates)
	summary := Compute(templates)
	if summary.Total != 2 {
		t.Fatalf("total=%d, want 2", summary.Total)
	}

	dir := t.TempDir()
	csvPath := filepath.Join(dir, "out.csv")
	if err := ExportCSV(csvPath, stats); err != nil {
		t.Fatalf("ExportCSV error: %v", err)
	}
	if _, err := os.Stat(csvPath); err != nil {
		t.Fatalf("csv not written")
	}
	htmlPath := filepath.Join(dir, "out.html")
	if err := ExportHTML(htmlPath, stats, summary); err != nil {
		t.Fatalf("ExportHTML error: %v", err)
	}
}

func toStats(tpls []recommend.TemplateInfo) []TemplateStats {
	stats := make([]TemplateStats, 0, len(tpls))
	for _, t := range tpls {
		stats = append(stats, TemplateStats{
			Path:     t.Path,
			Quality:  t.Quality,
			Usage:    t.Usage,
			Status:   t.Status,
			Version:  t.Version,
			DocType:  t.DocType,
			Language: t.Language,
		})
	}
	return stats
}
