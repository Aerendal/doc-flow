package templates

import (
	"testing"
	"time"

	"docflow/pkg/recommend"
)

func TestRecommendTransitionsDeprecate(t *testing.T) {
	now := time.Now()
	tl := TemplateLifecycle{
		Info: recommend.TemplateInfo{
			Path:    "tmpl.md",
			Status:  "active",
			Quality: 40,
			Usage:   0,
		},
		LastUsed: now.AddDate(0, 0, -100),
		Now:      now,
	}
	rule := LifecycleRule{QualityThreshold: 50, UsageDays: 90}
	tr := RecommendTransitions(tl, rule)
	if len(tr) != 1 || tr[0].To != StateDeprecated {
		t.Fatalf("expected deprecate transition, got %+v", tr)
	}
}

func TestNoDeprecateWhenUsed(t *testing.T) {
	now := time.Now()
	tl := TemplateLifecycle{
		Info: recommend.TemplateInfo{
			Path:    "tmpl.md",
			Status:  "active",
			Quality: 40,
			Usage:   1,
		},
		LastUsed: now.AddDate(0, 0, -10),
		Now:      now,
	}
	rule := LifecycleRule{QualityThreshold: 50, UsageDays: 90}
	tr := RecommendTransitions(tl, rule)
	if len(tr) != 0 {
		t.Fatalf("expected no transitions when usage>0")
	}
}

func TestDeprecateWhenStaleFlag(t *testing.T) {
	now := time.Now()
	tl := TemplateLifecycle{
		Info: recommend.TemplateInfo{
			Path:    "tmpl.md",
			Status:  "active",
			Quality: 40,
			Usage:   0,
		},
		LastUsed:           now, // recently used but flagged stale
		Now:                now,
		UsageZeroFor90Days: true,
	}
	rule := LifecycleRule{QualityThreshold: 50, UsageDays: 90}
	tr := RecommendTransitions(tl, rule)
	if len(tr) != 1 || tr[0].To != StateDeprecated {
		t.Fatalf("expected deprecate due to UsageZeroFor90Days, got %+v", tr)
	}
}
