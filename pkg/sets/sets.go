package sets

import "sort"

type TemplateUsage struct {
	SetID    string
	Template string
}

type Set struct {
	ID        string   `yaml:"id"`
	Name      string   `yaml:"name"`
	Templates []string `yaml:"templates"`
}

// CoOccurrence builds simple frequency of pairs and returns top pairs above threshold count.
func CoOccurrence(usages []TemplateUsage, minCount int) map[string][]string {
	type pair struct {
		a, b string
	}
	count := make(map[pair]int)

	bySet := make(map[string][]string)
	for _, u := range usages {
		bySet[u.SetID] = append(bySet[u.SetID], u.Template)
	}

	for _, ts := range bySet {
		// unique templates per set
		unique := map[string]bool{}
		for _, t := range ts {
			unique[t] = true
		}
		var list []string
		for t := range unique {
			list = append(list, t)
		}
		sort.Strings(list)
		for i := 0; i < len(list); i++ {
			for j := i + 1; j < len(list); j++ {
				p := pair{list[i], list[j]}
				count[p]++
			}
		}
	}

	result := make(map[string][]string)
	for p, c := range count {
		if c >= minCount {
			result[p.a] = append(result[p.a], p.b)
			result[p.b] = append(result[p.b], p.a)
		}
	}

	for k := range result {
		sort.Strings(result[k])
	}
	return result
}
