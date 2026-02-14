package validator

import (
	"fmt"
	"sort"
	"strings"

	execache "docflow/internal/exec/cache"
	"docflow/pkg/index"
	"docflow/pkg/sections"
)

// ValidateFacts validates precomputed document facts.
// Semantics are equivalent to ValidateDocs; differences should be treated as regressions.
func ValidateFacts(facts []execache.DocumentFacts, idx *index.DocumentIndex, opts *Options) (*Report, error) {
	if opts == nil {
		opts = &Options{PromoteContextFor: map[string]bool{}, SectionAliases: map[string][]string{}, RequiredDeps: map[string][]string{}}
	}
	if opts.SectionAliases == nil {
		opts.SectionAliases = map[string][]string{}
	}

	report := &Report{}
	var docs []docInfo

	if idx != nil {
		for _, rec := range idx.All() {
			if len(rec.Cycle) > 0 {
				report.add(Issue{
					Level:   LevelError,
					Type:    IssueCycleDetected,
					File:    rec.Path,
					DocID:   rec.DocID,
					Message: fmt.Sprintf("wykryto cykl: %s", strings.Join(rec.Cycle, " -> ")),
				})
			}
		}
	}

	for _, fact := range facts {
		report.Files++
		rel := fact.Path

		if fact.ParseErrorCode == "READ_ERROR" {
			report.add(Issue{
				Level:   LevelError,
				Type:    IssueMissingField,
				File:    rel,
				Message: fact.ParseErrorMessage,
			})
			continue
		}
		if fact.ParseErrorCode == "INVALID_YAML" {
			report.add(Issue{
				Level:   LevelError,
				Type:    IssueInvalidYAML,
				File:    rel,
				Line:    fact.ParseErrorLine,
				Message: fact.ParseErrorMessage,
			})
			continue
		}
		if !fact.HasFrontmatter {
			report.add(Issue{
				Level:   LevelError,
				Type:    IssueMissingFrontmatter,
				File:    rel,
				Message: "brak frontmatter YAML",
			})
			continue
		}

		raw := fact.RawFields
		if raw == nil {
			raw = map[string]bool{}
		}

		docID := strings.TrimSpace(fact.DocID)
		if docID == "" {
			report.add(Issue{Level: LevelError, Type: IssueMissingField, File: rel, Message: "brak doc_id"})
			docID = index.NormalizeDocID(rel)
		}

		if strings.TrimSpace(fact.Title) == "" || !raw["title"] {
			report.add(Issue{Level: LevelError, Type: IssueMissingField, File: rel, DocID: docID, Message: "brak title"})
		}
		docType := strings.TrimSpace(fact.DocType)
		if docType == "" || !raw["doc_type"] {
			report.add(Issue{Level: LevelError, Type: IssueMissingField, File: rel, DocID: docID, Message: "brak doc_type"})
		}
		if strings.TrimSpace(fact.Status) == "" || !raw["status"] {
			report.add(Issue{Level: LevelError, Type: IssueMissingField, File: rel, DocID: docID, Message: "brak status"})
		}
		if strings.TrimSpace(fact.Version) == "" || !raw["version"] {
			report.add(Issue{Level: LevelError, Type: IssueMissingField, File: rel, DocID: docID, Message: "brak version"})
		}

		if fact.DocID != "" && !docIDPattern.MatchString(docID) {
			report.add(Issue{Level: LevelError, Type: IssueInvalidDocID, File: rel, DocID: docID, Message: "doc_id nie spełnia konwencji snake_case"})
		}

		expectedID := index.NormalizeDocID(rel)
		if fact.DocID != "" && fact.DocID != expectedID {
			report.add(Issue{
				Level:   LevelWarn,
				Type:    IssueDocIDFilename,
				File:    rel,
				DocID:   fact.DocID,
				Message: fmt.Sprintf("doc_id=%s nie zgadza się z nazwą pliku (%s)", fact.DocID, expectedID),
			})
		}

		if fact.DocID != "" {
			docs = append(docs, docInfo{DocID: fact.DocID, File: rel})
			report.Documents++
		}

		statusLower := strings.ToLower(fact.Status)
		relaxed := statusLower == "archived" || statusLower == "deprecated" || statusLower == "generated"
		emptyStrict, contextStrict := statusProfile(statusLower)
		if !opts.StatusAware {
			emptyStrict = opts.PublishStrict && statusLower == "published"
			contextStrict = opts.StrictMode || opts.PromoteContextFor[docType]
		}

		if len(fact.ContextSources) == 0 {
			level := LevelWarn
			if contextStrict {
				level = LevelError
			}
			if relaxed {
				level = LevelWarn
			}
			report.add(Issue{Level: level, Type: IssueMissingContext, File: rel, DocID: docID, Message: "brak context_sources"})
		} else if bad, ok := containsTraversal(fact.ContextSources); ok {
			report.add(Issue{
				Level:   LevelError,
				Type:    IssueInvalidDocID,
				File:    rel,
				DocID:   docID,
				Message: fmt.Sprintf("context_sources zawiera traversal: %s", bad),
			})
		}

		if reqDeps, ok := opts.RequiredDeps[docType]; ok && len(reqDeps) > 0 {
			missing := missingDeps(reqDeps, fact.DependsOn)
			for _, md := range missing {
				report.add(Issue{
					Level:   levelFor(opts.StrictMode),
					Type:    IssueMissingExpectedDep,
					File:    rel,
					DocID:   docID,
					Message: fmt.Sprintf("brak wymaganej zależności dla %s: %s", docType, md),
				})
			}
		}
		if bad, ok := containsTraversal(fact.DependsOn); ok {
			report.add(Issue{
				Level:   LevelError,
				Type:    IssueInvalidDocID,
				File:    rel,
				DocID:   docID,
				Message: fmt.Sprintf("depends_on zawiera traversal: %s", bad),
			})
		}

		if emptyStrict && statusLower == "published" {
			metrics := sections.ComputeMetrics(fact.Content)
			if metrics.Empty > 0 {
				report.add(Issue{
					Level:   LevelError,
					Type:    IssueMissingField,
					File:    rel,
					DocID:   docID,
					Message: fmt.Sprintf("published: %d pustych sekcji (completeness %.2f)", metrics.Empty, metrics.Completeness),
				})
			}
		} else if emptyStrict && relaxed {
			metrics := sections.ComputeMetrics(fact.Content)
			if metrics.Empty > 0 {
				report.add(Issue{
					Level:   LevelWarn,
					Type:    IssueMissingField,
					File:    rel,
					DocID:   docID,
					Message: fmt.Sprintf("%s: %d pustych sekcji (completeness %.2f)", statusLower, metrics.Empty, metrics.Completeness),
				})
			}
		}

		if len(opts.SectionAliases) > 0 {
			hits := sections.FindAliasHits(fact.Content, opts.SectionAliases)
			for _, h := range hits {
				report.add(Issue{
					Level:   LevelWarn,
					Type:    IssueLegacySectionName,
					File:    rel,
					DocID:   docID,
					Message: fmt.Sprintf("legacy section name '%s' → '%s'", h.From, h.To),
				})
			}
		}

		if opts.GovernanceRules != nil {
			checkGovernance(report, opts.GovernanceRules, rel, docID, statusLower, docType, fact.Content)
		}
	}

	seen := make(map[string][]string)
	for _, d := range docs {
		seen[d.DocID] = append(seen[d.DocID], d.File)
	}
	for docID, filesWithID := range seen {
		if len(filesWithID) <= 1 {
			continue
		}
		sort.Strings(filesWithID)
		for _, f := range filesWithID {
			report.add(Issue{
				Level:   LevelError,
				Type:    IssueDuplicateDocID,
				File:    f,
				DocID:   docID,
				Message: fmt.Sprintf("duplikat doc_id=%s (%d wystąpień)", docID, len(filesWithID)),
			})
		}
	}

	return report, nil
}
