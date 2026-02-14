package engine

import (
	"sort"

	"docflow/pkg/compliance"
)

func BuildComplianceReport(sum *compliance.Summary) ComplianceReport {
	if sum == nil {
		return ComplianceReport{
			ReportMeta: ReportMeta{SchemaVersion: "1.0"},
			Docs:       []ComplianceDoc{},
		}
	}

	report := ComplianceReport{
		ReportMeta: ReportMeta{
			SchemaVersion: sum.SchemaVersion,
		},
		IdentityVersion: sum.IdentityVersion,
		ReportViewMeta: ReportViewMeta{
			FailOn: sum.FailOn,
			Show:   sum.Show,
		},
		RulesPath:          sum.RulesPath,
		RulesChecksum:      sum.RulesChecksum,
		Documents:          sum.Documents,
		Passed:             sum.Passed,
		Failed:             sum.Failed,
		PassRate:           sum.PassRate,
		ViolationsCount:    cloneStringIntMap(sum.ViolationsCount),
		NewFailed:          sum.NewFailed,
		ExistingFailed:     sum.ExistingFailed,
		NewViolationsCount: cloneStringIntMap(sum.NewViolations),
		DuplicateDocIDs:    cloneDuplicateDocIDs(sum.DuplicateDocIDs),
		Docs:               make([]ComplianceDoc, 0, len(sum.Docs)),
	}
	if sum.Baseline != nil {
		report.Baseline = &BaselineMeta{
			Loaded:          sum.Baseline.Loaded,
			Path:            sum.Baseline.Path,
			SchemaVersion:   sum.Baseline.SchemaVersion,
			IdentityVersion: sum.Baseline.IdentityVersion,
		}
		report.ReportViewMeta.Against = sum.Baseline.Path
	}

	for _, doc := range sum.Docs {
		v := append([]string(nil), doc.Violations...)
		sort.Strings(v)
		report.Docs = append(report.Docs, ComplianceDoc{
			DocID:      doc.DocID,
			Path:       doc.Path,
			Status:     doc.Status,
			DocType:    doc.DocType,
			Violations: v,
			Details:    cloneStringStringMap(doc.Details),
		})
	}
	sort.Slice(report.Docs, func(i, j int) bool {
		return report.Docs[i].Path < report.Docs[j].Path
	})

	return report
}

func ApplyComplianceReport(sum *compliance.Summary, report ComplianceReport) {
	if sum == nil {
		return
	}

	sum.SchemaVersion = report.SchemaVersion
	sum.IdentityVersion = report.IdentityVersion
	sum.FailOn = report.FailOn
	sum.Show = report.Show
	sum.RulesPath = report.RulesPath
	sum.RulesChecksum = report.RulesChecksum
	sum.Documents = report.Documents
	sum.Passed = report.Passed
	sum.Failed = report.Failed
	sum.PassRate = report.PassRate
	sum.ViolationsCount = cloneStringIntMap(report.ViolationsCount)
	sum.NewFailed = report.NewFailed
	sum.ExistingFailed = report.ExistingFailed
	sum.NewViolations = cloneStringIntMap(report.NewViolationsCount)
	sum.DuplicateDocIDs = cloneDuplicateDocIDs(report.DuplicateDocIDs)

	if report.Baseline != nil {
		sum.Baseline = &compliance.BaselineMeta{
			Path:            report.Baseline.Path,
			SchemaVersion:   report.Baseline.SchemaVersion,
			IdentityVersion: report.Baseline.IdentityVersion,
			Loaded:          report.Baseline.Loaded,
		}
	} else {
		sum.Baseline = nil
	}

	docs := make([]compliance.DocResult, 0, len(report.Docs))
	for _, doc := range report.Docs {
		docs = append(docs, compliance.DocResult{
			DocID:      doc.DocID,
			Path:       doc.Path,
			Status:     doc.Status,
			DocType:    doc.DocType,
			Violations: append([]string(nil), doc.Violations...),
			Details:    cloneStringStringMap(doc.Details),
		})
	}
	sum.Docs = docs
}

func cloneStringIntMap(in map[string]int) map[string]int {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]int, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func cloneStringStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func cloneDuplicateDocIDs(in map[string][]string) map[string][]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string][]string, len(in))
	for k, v := range in {
		cp := append([]string(nil), v...)
		sort.Strings(cp)
		out[k] = cp
	}
	return out
}
