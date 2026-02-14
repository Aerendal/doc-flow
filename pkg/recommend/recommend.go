package recommend

import (
	"sort"
	"strings"
)

type TemplateInfo struct {
	Path       string
	DocType    string
	Language   string
	Quality    int
	Usage      int
	Deprecated bool
	Version    string
	Status     string
	CodeBlocks int
	Tables     int
	Links      int
	Images     int
}

type Params struct {
	PreferredDocType string
	PreferredLang    string
	Weights          Weights
	TopN             int
}

type Weights struct {
	DocType float64
	Lang    float64
	Quality float64
	Usage   float64
}

type Recommendation struct {
	Template TemplateInfo
	Score    float64
	Reason   string
}

func Recommend(templates []TemplateInfo, params Params) []Recommendation {
	if params.TopN <= 0 {
		params.TopN = 5
	}
	w := params.Weights
	if w.DocType == 0 && w.Lang == 0 && w.Quality == 0 && w.Usage == 0 {
		w = Weights{DocType: 0.4, Lang: 0.2, Quality: 0.25, Usage: 0.15}
	}

	var recs []Recommendation
	for _, t := range templates {
		status := strings.ToLower(t.Status)
		if t.Deprecated || status == "deprecated" || status == "archived" {
			continue
		}
		score := 0.0
		reasonParts := []string{}

		if params.PreferredDocType != "" && t.DocType == params.PreferredDocType {
			score += w.DocType * 100
			reasonParts = append(reasonParts, "doc_type match")
		}
		if params.PreferredLang != "" && t.Language == params.PreferredLang {
			score += w.Lang * 100
			reasonParts = append(reasonParts, "lang match")
		}
		score += w.Quality * float64(t.Quality)
		reasonParts = append(reasonParts, "quality")

		score += w.Usage * float64(t.Usage)
		if t.Usage > 0 {
			reasonParts = append(reasonParts, "usage")
		}

		recs = append(recs, Recommendation{
			Template: t,
			Score:    score,
			Reason:   joinReasons(reasonParts),
		})
	}

	sort.Slice(recs, func(i, j int) bool {
		if recs[i].Score == recs[j].Score {
			// prefer higher semver version if same base name
			vi := versionValue(recs[i].Template.Version)
			vj := versionValue(recs[j].Template.Version)
			if vi == vj {
				return recs[i].Template.Path < recs[j].Template.Path
			}
			return vi > vj
		}
		return recs[i].Score > recs[j].Score
	})

	if len(recs) > params.TopN {
		recs = recs[:params.TopN]
	}
	return recs
}

func joinReasons(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, ", ")
}
