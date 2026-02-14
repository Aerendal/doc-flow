package analytics

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"docflow/pkg/recommend"
)

type TemplateStats struct {
	Path     string
	Quality  int
	Usage    int
	Status   string
	Version  string
	DocType  string
	Language string
}

type Summary struct {
	Total      int
	AvgQuality float64
	TopByUsage []TemplateStats
	LowQuality []TemplateStats
}

func Compute(templates []recommend.TemplateInfo) Summary {
	stats := make([]TemplateStats, 0, len(templates))
	totalQ := 0
	for _, t := range templates {
		stats = append(stats, TemplateStats{
			Path:     t.Path,
			Quality:  t.Quality,
			Usage:    t.Usage,
			Status:   t.Status,
			Version:  t.Version,
			DocType:  t.DocType,
			Language: t.Language,
		})
		totalQ += t.Quality
	}
	s := Summary{Total: len(stats)}
	if len(stats) > 0 {
		s.AvgQuality = float64(totalQ) / float64(len(stats))
	}

	sort.Slice(stats, func(i, j int) bool { return stats[i].Usage > stats[j].Usage })
	if len(stats) > 5 {
		s.TopByUsage = stats[:5]
	} else {
		s.TopByUsage = stats
	}

	sort.Slice(stats, func(i, j int) bool { return stats[i].Quality < stats[j].Quality })
	if len(stats) > 5 {
		s.LowQuality = stats[:5]
	} else {
		s.LowQuality = stats
	}
	return s
}

func ExportCSV(path string, stats []TemplateStats) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create csv %s: %w", path, err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	header := []string{"path", "doc_type", "lang", "version", "status", "quality", "usage"}
	if err := w.Write(header); err != nil {
		return err
	}
	for _, s := range stats {
		row := []string{
			s.Path, s.DocType, s.Language, s.Version, s.Status,
			strconv.Itoa(s.Quality),
			strconv.Itoa(s.Usage),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func ExportHTML(path string, stats []TemplateStats, summary Summary) error {
	var b strings.Builder
	b.WriteString("<html><body>\n")
	fmt.Fprintf(&b, "<h1>Template Analytics</h1>\n")
	fmt.Fprintf(&b, "<p>Total: %d, AvgQuality: %.1f</p>\n", summary.Total, summary.AvgQuality)
	b.WriteString("<table border=\"1\"><tr><th>Path</th><th>DocType</th><th>Lang</th><th>Ver</th><th>Status</th><th>Quality</th><th>Usage</th></tr>\n")
	for _, s := range stats {
		fmt.Fprintf(&b, "<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%d</td><td>%d</td></tr>\n",
			s.Path, s.DocType, s.Language, s.Version, s.Status, s.Quality, s.Usage)
	}
	b.WriteString("</table></body></html>")
	return os.WriteFile(path, []byte(b.String()), 0o644)
}
