package fix

import (
	"sort"

	"docflow/internal/engine"
)

type Scope string

const (
	ScopeFrontmatter Scope = "frontmatter"
	ScopeSections    Scope = "sections"
)

type FileInput struct {
	Path    string
	Content string
}

type Options struct {
	Scopes     map[Scope]bool
	OnlyCodes  map[string]bool
	MaxFiles   int
	MaxChanges int
	Aliases    map[string][]string
}

type PlannedFile struct {
	Path             string   `json:"path"`
	Codes            []string `json:"codes"`
	Summaries        []string `json:"summaries"`
	Changes          int      `json:"changes"`
	Diff             string   `json:"diff,omitempty"`
	OriginalChecksum string   `json:"-"`
	UpdatedContent   string   `json:"-"`
}

type PlanStats struct {
	FilesTouched int `json:"files_touched"`
	TotalEdits   int `json:"total_edits"`
}

type Plan struct {
	SchemaVersion   string        `json:"schema_version"`
	IdentityVersion string        `json:"identity_version,omitempty"`
	Root            string        `json:"root,omitempty"`
	Files           []PlannedFile `json:"files"`
	Stats           PlanStats     `json:"stats"`
}

type ApplyOptions struct {
	Root      string
	BackupDir string
}

func defaultOptions(opts Options) Options {
	if len(opts.Scopes) == 0 {
		opts.Scopes = map[Scope]bool{
			ScopeFrontmatter: true,
			ScopeSections:    true,
		}
	}
	if opts.MaxFiles <= 0 {
		opts.MaxFiles = 200
	}
	if opts.MaxChanges <= 0 {
		opts.MaxChanges = 1000
	}
	if opts.Aliases == nil {
		opts.Aliases = map[string][]string{}
	}
	return opts
}

func allowedByCode(code string, only map[string]bool) bool {
	if len(only) == 0 {
		return true
	}
	return only[code]
}

func sortIssuesCanonical(issues []engine.Issue) {
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Path != issues[j].Path {
			return issues[i].Path < issues[j].Path
		}
		if issues[i].Code != issues[j].Code {
			return issues[i].Code < issues[j].Code
		}
		if issues[i].DocID != issues[j].DocID {
			return issues[i].DocID < issues[j].DocID
		}
		li := 0
		if issues[i].Location != nil {
			li = issues[i].Location.Line
		}
		lj := 0
		if issues[j].Location != nil {
			lj = issues[j].Location.Line
		}
		if li != lj {
			return li < lj
		}
		return issues[i].Message < issues[j].Message
	})
}
