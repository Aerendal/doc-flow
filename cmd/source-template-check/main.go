package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Policy struct {
	RequiredSections []string `yaml:"required_sections"`
}

type CheckResult struct {
	File      string   `json:"file"`
	Status    string   `json:"status"`
	Missing   []string `json:"missing"`
	Headings  []string `json:"headings"`
	Forbidden []string `json:"forbidden"`
}

func loadAliases(docflowPath string) (map[string][]string, error) {
	b, err := os.ReadFile(docflowPath)
	if err != nil {
		return nil, err
	}
	var cfg struct {
		SectionAliases map[string][]string `yaml:"section_aliases"`
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return cfg.SectionAliases, nil
}

func normalize(name string, canon map[string]string) string {
	key := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(name, "#")))
	if v, ok := canon[key]; ok {
		return v
	}
	return strings.TrimSpace(strings.TrimPrefix(name, "#"))
}

func headings(md string) []string {
	lines := strings.Split(md, "\n")
	var out []string
	for _, l := range lines {
		if strings.HasPrefix(l, "#") {
			title := strings.TrimSpace(strings.TrimLeft(l, "#"))
			if title != "" {
				out = append(out, title)
			}
		}
	}
	return out
}

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "usage: %s <policy_yaml> <docflow_yaml> <root_dir> [--suggest]\n", os.Args[0])
		os.Exit(1)
	}

	suggest := false
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		if args[i] == "--suggest" {
			suggest = true
			args = append(args[:i], args[i+1:]...)
			i--
		}
	}
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <policy_yaml> <docflow_yaml> <root_dir> [--suggest]\n", os.Args[0])
		os.Exit(1)
	}

	policyPath := args[0]
	docflowPath := args[1]
	root := args[2]

	var policy Policy
	if b, err := os.ReadFile(policyPath); err == nil {
		if err := yaml.Unmarshal(b, &policy); err != nil {
			fmt.Fprintf(os.Stderr, "policy parse error: %v\n", err)
			os.Exit(2)
		}
	} else {
		fmt.Fprintf(os.Stderr, "missing policy: %v\n", err)
		os.Exit(2)
	}

	aliases, err := loadAliases(docflowPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "docflow parse error: %v\n", err)
		os.Exit(2)
	}

	canon := make(map[string]string)
	for _, r := range policy.RequiredSections {
		canon[strings.ToLower(r)] = r
		for _, a := range aliases[r] {
			canon[strings.ToLower(a)] = r
		}
	}

	var results []CheckResult
	err = filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		b, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		txt := string(b)
		if strings.HasPrefix(txt, "---") {
			parts := strings.SplitN(txt, "---", 3)
			if len(parts) >= 3 {
				fm := parts[1]
				for _, line := range strings.Split(fm, "\n") {
					if strings.HasPrefix(strings.TrimSpace(line), "doc_type:") {
						dt := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
						if dt == "tasklist" {
							return nil
						}
					}
				}
			}
		}

		hs := headings(txt)
		norm := map[string]bool{}
		for _, h := range hs {
			norm[normalize(h, canon)] = true
		}
		var missing []string
		for _, req := range policy.RequiredSections {
			if !norm[req] {
				missing = append(missing, req)
			}
		}
		status := "PASS"
		if len(missing) > 0 {
			status = "FAIL"
		}
		results = append(results, CheckResult{File: path, Status: status, Missing: missing, Headings: hs})
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "walk error: %v\n", err)
		os.Exit(2)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].File < results[j].File
	})

	fail := 0
	for _, r := range results {
		if r.Status != "PASS" {
			fail++
		}
	}

	if suggest {
		for _, r := range results {
			if len(r.Missing) == 0 {
				continue
			}
			fmt.Println("##", r.File)
			fmt.Println("```markdown")
			limit := len(r.Missing)
			if limit > 10 {
				limit = 10
			}
			for i, m := range r.Missing[:limit] {
				fmt.Printf("%d. %s\n- TODO: uzupełnij sekcję %s (krótko: cel, decyzje, status)\n\n", i+1, m, m)
			}
			if len(r.Missing) > limit {
				fmt.Printf("+%d więcej sekcji...\n", len(r.Missing)-limit)
			}
			fmt.Println("```")
		}
	} else {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(results)
	}

	if fail > 0 {
		os.Exit(2)
	}
}
