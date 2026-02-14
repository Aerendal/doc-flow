package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

var defaultCanon = map[string][]string{
	"api":          {"Przegląd", "Endpoints", "Authentication", "Errors"},
	"architecture": {"Przegląd", "Decyzje architektoniczne", "Komponenty", "Ryzyka"},
	"guide":        {"Przegląd", "Kroki", "FAQ"},
	"proposal":     {"Wstęp"},
}

func loadCanon(path string) (map[string][]string, string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return defaultCanon, fmt.Sprintf("fallback:cannot_read_governance(%v)", err)
	}

	var cfg struct {
		Families map[string]struct {
			RequiredSections []string `yaml:"required_sections"`
		} `yaml:"families"`
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return defaultCanon, fmt.Sprintf("fallback:cannot_parse_governance(%v)", err)
	}

	out := make(map[string][]string)
	for k, v := range cfg.Families {
		if len(v.RequiredSections) > 0 {
			out[k] = v.RequiredSections
		}
	}
	if len(out) == 0 {
		return defaultCanon, "fallback:empty_governance_families"
	}
	return out, "governance"
}

type LintResult struct {
	File        string   `json:"file"`
	DocType     string   `json:"doc_type"`
	Canon       []string `json:"canon"`
	CanonSource string   `json:"canon_source"`
	Headings    []string `json:"headings"`
	Violations  []string `json:"violations"`
}

func main() {
	roots := os.Args[1:]
	if len(roots) == 0 {
		roots = []string{"."}
	}

	canon, canonSource := loadCanon("docs/_meta/GOVERNANCE_RULES.yaml")
	if canonSource != "governance" {
		fmt.Fprintf(os.Stderr, "section-order-lint: %s\n", canonSource)
	}
	results := []LintResult{}

	for _, root := range roots {
		_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			docType, headings := parseFile(path)
			r := LintResult{
				File:        path,
				DocType:     docType,
				CanonSource: canonSource,
				Headings:    headings,
			}
			if c, ok := canon[docType]; ok {
				r.Canon = c
				r.Violations = checkOrder(c, headings)
			} else if docType != "" {
				r.Violations = []string{"unknown_doc_type"}
			}
			results = append(results, r)
			return nil
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].File < results[j].File
	})

	exit := 0
	for _, r := range results {
		if len(r.Violations) > 0 {
			exit = 1
			fmt.Printf("%s: %v\n", r.File, r.Violations)
		}
	}

	out, _ := json.MarshalIndent(results, "", "  ")
	_ = os.MkdirAll(".docflow", 0o755)
	_ = os.WriteFile(".docflow/section_order_lint.json", out, 0o644)
	os.Exit(exit)
}

func parseFile(path string) (string, []string) {
	f, err := os.Open(path)
	if err != nil {
		return "", nil
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	docType := ""
	inFront := false
	headings := []string{}
	reH := regexp.MustCompile(`^##\s+(.*)$`)
	for sc.Scan() {
		line := sc.Text()
		if strings.TrimSpace(line) == "---" {
			inFront = !inFront
			continue
		}
		if inFront && strings.HasPrefix(line, "doc_type:") {
			docType = strings.Trim(strings.TrimSpace(strings.SplitN(line, ":", 2)[1]), "\"")
		}
		if !inFront {
			if m := reH.FindStringSubmatch(line); m != nil {
				headings = append(headings, strings.TrimSpace(m[1]))
			}
		}
	}
	return docType, headings
}

func checkOrder(canon, headings []string) []string {
	if len(canon) == 0 {
		return nil
	}

	idx := make([]int, len(canon))
	for i := range idx {
		idx[i] = -1
	}
	for i, h := range headings {
		for j, c := range canon {
			if h == c && idx[j] == -1 {
				idx[j] = i
			}
		}
	}

	viol := []string{}
	for j, pos := range idx {
		if pos == -1 {
			viol = append(viol, fmt.Sprintf("missing:%s", canon[j]))
		}
	}

	last := -1
	for j, pos := range idx {
		if pos == -1 {
			continue
		}
		if pos < last {
			viol = append(viol, fmt.Sprintf("out_of_order:%s", canon[j]))
		}
		last = pos
	}

	return viol
}
