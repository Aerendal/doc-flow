package templates

import (
	"time"

	"docflow/pkg/recommend"
)

type LifecycleState string

const (
	StateDraft      LifecycleState = "draft"
	StateActive     LifecycleState = "active"
	StateDeprecated LifecycleState = "deprecated"
	StateArchived   LifecycleState = "archived"
)

type TemplateLifecycle struct {
	Info               recommend.TemplateInfo
	LastUsed           time.Time
	Now                time.Time
	UsageZeroFor90Days bool
}

type Transition struct {
	From   LifecycleState
	To     LifecycleState
	Reason string
}

type LifecycleRule struct {
	QualityThreshold int
	UsageDays        int // days with zero usage to deprecate
}

func RecommendTransitions(t TemplateLifecycle, rule LifecycleRule) []Transition {
	var tr []Transition
	// draft -> active manual (not auto)
	if t.Info.Status == string(StateActive) && shouldDeprecate(t, rule) {
		tr = append(tr, Transition{
			From:   StateActive,
			To:     StateDeprecated,
			Reason: "quality below threshold and no usage",
		})
	}
	if t.Info.Status == string(StateDeprecated) {
		tr = append(tr, Transition{
			From:   StateDeprecated,
			To:     StateArchived,
			Reason: "manual archive after migration",
		})
	}
	return tr
}

func shouldDeprecate(t TemplateLifecycle, rule LifecycleRule) bool {
	noUsage := t.Info.Usage == 0 &&
		(t.UsageZeroFor90Days || t.Now.Sub(t.LastUsed).Hours() >= float64(rule.UsageDays*24))
	lowQuality := t.Info.Quality < rule.QualityThreshold
	return lowQuality && noUsage
}
