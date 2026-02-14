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
        "unicode"
)

type TemplateInfo struct {
        Name          string
        Path          string
        Title         string
        Status        string
        Lines         int
        Bytes         int
        Headings      []string
        HeadingDepths []int
        MaxDepth      int
        Language      string
        Category      string
}

type CategoryStats struct {
        Count    int
        Examples []string
}

type HeadingPattern struct {
        Pattern string
        Count   int
}

var polishWords = map[string]bool{
        "cel": true, "dokumentu": true, "zakres": true, "granice": true,
        "metadane": true, "wejścia": true, "wyjścia": true, "powiązania": true,
        "zależności": true, "fazy": true, "cyklu": true, "życia": true,
        "właściciel": true, "wersja": true, "kroki": true, "sekcja": true,
        "obejmuje": true, "wymaga": true, "odniesienia": true, "określ": true,
        "czy": true, "tej": true, "fazie": true, "dokument": true,
        "powstaje": true, "jest": true, "aktualizowany": true,
        "przeglądany": true, "archiwizowany": true, "uzasadnienie": true,
        "odpowiedzialnych": true, "podaj": true, "standardy": true,
        "bezpieczeństwo": true, "implementacja": true, "testowanie": true,
        "planowanie": true, "analiza": true, "wymagań": true,
        "koncepcja": true, "wizja": true, "projekt": true, "design": true,
        "dane": true, "referencyjne": true, "decyzje": true,
        "konsumpcja": true, "rezultatów": true, "zespoły": true,
        "dostawcy": true, "branżowe": true,
        "kontrole": true, "jakości": true, "walidacja": true,
        "eskalacja": true, "komunikacja": true, "monitoring": true,
}

var englishWords = map[string]bool{
        "the": true, "and": true, "for": true, "with": true, "this": true,
        "that": true, "from": true, "will": true, "are": true, "have": true,
        "should": true, "must": true, "requirements": true, "overview": true,
        "introduction": true, "scope": true, "purpose": true, "description": true,
        "implementation": true, "testing": true, "deployment": true,
        "configuration": true, "architecture": true, "design": true,
        "security": true, "performance": true, "monitoring": true,
}

func detectLanguage(content string) string {
        lower := strings.ToLower(content)
        words := strings.FieldsFunc(lower, func(r rune) bool {
                return !unicode.IsLetter(r)
        })

        plCount := 0
        enCount := 0
        for _, w := range words {
                if polishWords[w] {
                        plCount++
                }
                if englishWords[w] {
                        enCount++
                }
        }

        if plCount > enCount*2 {
                return "PL"
        } else if enCount > plCount*2 {
                return "EN"
        } else if plCount > 0 && enCount > 0 {
                return "PL+EN"
        }
        return "unknown"
}

func categorizeByName(name string) string {
        name = strings.ToLower(name)

        categories := []struct {
                keywords []string
                category string
        }{
                {[]string{"api_"}, "API"},
                {[]string{"security", "vulnerability", "penetration", "access_control", "threat"}, "Security"},
                {[]string{"test", "qa_", "quality_assurance"}, "Testing/QA"},
                {[]string{"deploy", "release", "rollback", "ci_cd", "pipeline"}, "Deployment/CI-CD"},
                {[]string{"incident", "postmortem", "root_cause", "disaster", "crisis"}, "Incident/Operations"},
                {[]string{"monitor", "alert", "observ", "logging", "metric"}, "Monitoring/Observability"},
                {[]string{"onboard", "training", "knowledge", "mentor"}, "Onboarding/Training"},
                {[]string{"architecture", "design_", "system_design", "technical_design", "adr_"}, "Architecture/Design"},
                {[]string{"requirement", "specification", "functional_", "user_story", "acceptance_criteria"}, "Requirements/Specs"},
                {[]string{"runbook", "playbook", "procedure", "sop_", "checklist"}, "Runbooks/Procedures"},
                {[]string{"review", "audit", "compliance", "governance"}, "Review/Compliance"},
                {[]string{"plan", "roadmap", "strategy", "budget"}, "Planning/Strategy"},
                {[]string{"data_", "database", "migration", "schema"}, "Data/Database"},
                {[]string{"communication", "stakeholder", "announcement", "report"}, "Communication/Reporting"},
                {[]string{"performance", "benchmark", "capacity", "load_test", "scalab"}, "Performance"},
                {[]string{"config", "environment", "infrastructure", "network"}, "Infrastructure/Config"},
                {[]string{"project_", "sprint", "retrospective", "standup", "backlog"}, "Project Management"},
                {[]string{"documentation", "template", "style_guide", "writing"}, "Documentation"},
                {[]string{"code_", "coding", "refactor", "technical_debt"}, "Code Quality"},
                {[]string{"integration", "third_party", "vendor", "external"}, "Integration"},
                {[]string{"feature_", "product_", "user_experience", "ux_", "ui_"}, "Product/Feature"},
                {[]string{"change_", "version", "changelog"}, "Change Management"},
                {[]string{"risk_", "mitigation"}, "Risk Management"},
                {[]string{"team_", "role_", "raci", "responsibility"}, "Team/Roles"},
        }

        for _, cat := range categories {
                for _, kw := range cat.keywords {
                        if strings.Contains(name, kw) {
                                return cat.category
                        }
                }
        }
        return "Other"
}

func main() {
        dir := "testdata/templates"
        if len(os.Args) > 1 {
                dir = os.Args[1]
        }

        headingRe := regexp.MustCompile(`^(#{1,6})\s+(.+)`)
        frontmatterRe := regexp.MustCompile(`^---\s*$`)

        var templates []TemplateInfo
        headingPatterns := map[string]int{}
        categoryMap := map[string]*CategoryStats{}
        langCount := map[string]int{}

        totalLines := 0
        totalBytes := 0
        maxDepthAll := 0
        headingCounts := []int{}

        files, _ := filepath.Glob(filepath.Join(dir, "*.md"))
        sort.Strings(files)

        for _, fpath := range files {
                data, err := os.ReadFile(fpath)
                if err != nil {
                        continue
                }

                name := strings.TrimSuffix(filepath.Base(fpath), ".md")
                content := string(data)
                lines := strings.Count(content, "\n") + 1

                ti := TemplateInfo{
                        Name:  name,
                        Path:  fpath,
                        Lines: lines,
                        Bytes: len(data),
                }

                scanner := bufio.NewScanner(strings.NewReader(content))
                inFrontmatter := false
                frontmatterCount := 0
                var headingList []string

                for scanner.Scan() {
                        line := scanner.Text()

                        if frontmatterRe.MatchString(line) {
                                frontmatterCount++
                                if frontmatterCount == 1 {
                                        inFrontmatter = true
                                        continue
                                } else if frontmatterCount == 2 {
                                        inFrontmatter = false
                                        continue
                                }
                        }

                        if inFrontmatter {
                                trimmed := strings.TrimSpace(line)
                                if strings.HasPrefix(trimmed, "title:") {
                                        ti.Title = strings.TrimSpace(strings.TrimPrefix(trimmed, "title:"))
                                } else if strings.HasPrefix(trimmed, "status:") {
                                        ti.Status = strings.TrimSpace(strings.TrimPrefix(trimmed, "status:"))
                                }
                                continue
                        }

                        if m := headingRe.FindStringSubmatch(line); m != nil {
                                depth := len(m[1])
                                heading := strings.TrimSpace(m[2])
                                ti.Headings = append(ti.Headings, heading)
                                ti.HeadingDepths = append(ti.HeadingDepths, depth)
                                headingList = append(headingList, fmt.Sprintf("h%d", depth))
                                if depth > ti.MaxDepth {
                                        ti.MaxDepth = depth
                                }
                        }
                }

                ti.Language = detectLanguage(content)
                ti.Category = categorizeByName(name)

                if len(headingList) > 0 {
                        pattern := strings.Join(headingList, " > ")
                        headingPatterns[pattern]++
                }

                cat := ti.Category
                if _, ok := categoryMap[cat]; !ok {
                        categoryMap[cat] = &CategoryStats{}
                }
                categoryMap[cat].Count++
                if len(categoryMap[cat].Examples) < 3 {
                        categoryMap[cat].Examples = append(categoryMap[cat].Examples, name)
                }

                langCount[ti.Language]++
                totalLines += lines
                totalBytes += len(data)
                if ti.MaxDepth > maxDepthAll {
                        maxDepthAll = ti.MaxDepth
                }
                headingCounts = append(headingCounts, len(ti.Headings))

                templates = append(templates, ti)
        }

        allHeadings := map[string]int{}
        for _, t := range templates {
                for _, h := range t.Headings {
                        allHeadings[h]++
                }
        }

        fmt.Println("========================================")
        fmt.Println("ANALIZA SZABLONÓW — WYNIKI")
        fmt.Println("========================================")
        fmt.Printf("\nLiczba szablonów: %d\n", len(templates))
        fmt.Printf("Łączna liczba linii: %d\n", totalLines)
        fmt.Printf("Łączny rozmiar: %d bajtów (%.1f MB)\n", totalBytes, float64(totalBytes)/1024/1024)
        if len(templates) > 0 {
                fmt.Printf("Średnia długość: %.0f linii\n", float64(totalLines)/float64(len(templates)))
                fmt.Printf("Średni rozmiar: %.0f bajtów\n", float64(totalBytes)/float64(len(templates)))
        }
        fmt.Printf("Maksymalna głębokość sekcji: %d\n", maxDepthAll)

        avgHeadings := 0.0
        if len(headingCounts) > 0 {
                sum := 0
                for _, c := range headingCounts {
                        sum += c
                }
                avgHeadings = float64(sum) / float64(len(headingCounts))
        }
        fmt.Printf("Średnia liczba nagłówków: %.1f\n", avgHeadings)

        fmt.Println("\n--- KATEGORIE ---")
        type catEntry struct {
                Name  string
                Stats *CategoryStats
        }
        var cats []catEntry
        for name, stats := range categoryMap {
                cats = append(cats, catEntry{name, stats})
        }
        sort.Slice(cats, func(i, j int) bool { return cats[i].Stats.Count > cats[j].Stats.Count })
        for _, c := range cats {
                fmt.Printf("  %-30s %3d  [%s]\n", c.Name, c.Stats.Count, strings.Join(c.Stats.Examples, ", "))
        }

        fmt.Println("\n--- ROZKŁAD JĘZYKÓW ---")
        for lang, count := range langCount {
                pct := float64(count) / float64(len(templates)) * 100
                fmt.Printf("  %-10s %3d (%.1f%%)\n", lang, count, pct)
        }

        fmt.Println("\n--- TOP 20 NAJCZĘSTSZYCH NAGŁÓWKÓW ---")
        type headEntry struct {
                Name  string
                Count int
        }
        var heads []headEntry
        for name, count := range allHeadings {
                heads = append(heads, headEntry{name, count})
        }
        sort.Slice(heads, func(i, j int) bool { return heads[i].Count > heads[j].Count })
        for i, h := range heads {
                if i >= 20 {
                        break
                }
                fmt.Printf("  %3d×  %s\n", h.Count, h.Name)
        }

        fmt.Println("\n--- TOP 10 WZORCÓW STRUKTURY SEKCJI ---")
        type patEntry struct {
                Pattern string
                Count   int
        }
        var pats []patEntry
        for p, c := range headingPatterns {
                pats = append(pats, patEntry{p, c})
        }
        sort.Slice(pats, func(i, j int) bool { return pats[i].Count > pats[j].Count })
        for i, p := range pats {
                if i >= 10 {
                        break
                }
                fmt.Printf("  %3d×  %s\n", p.Count, p.Pattern)
        }

        fmt.Println("\n--- FRONTMATTER STATUS ---")
        statusCount := map[string]int{}
        for _, t := range templates {
                s := t.Status
                if s == "" {
                        s = "(brak)"
                }
                statusCount[s]++
        }
        for s, c := range statusCount {
                fmt.Printf("  %-20s %3d\n", s, c)
        }

        jsonData := map[string]interface{}{
                "total_templates":    len(templates),
                "total_lines":        totalLines,
                "total_bytes":        totalBytes,
                "avg_lines":          float64(totalLines) / float64(len(templates)),
                "avg_headings":       avgHeadings,
                "max_depth":          maxDepthAll,
                "categories":         categoryMap,
                "language_dist":      langCount,
                "top_headings":       heads,
                "heading_patterns":   pats,
                "frontmatter_status": statusCount,
        }
        jsonFile, _ := os.Create("tools/scripts/analysis_results.json")
        enc := json.NewEncoder(jsonFile)
        enc.SetIndent("", "  ")
        enc.Encode(jsonData)
        jsonFile.Close()

        fmt.Println("\nJSON zapisany do tools/scripts/analysis_results.json")
}
