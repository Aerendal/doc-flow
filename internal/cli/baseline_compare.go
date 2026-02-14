package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"docflow/pkg/compliance"
)

type baselineMetaJSON struct {
	Path            string `json:"path"`
	SchemaVersion   string `json:"schema_version,omitempty"`
	IdentityVersion string `json:"identity_version,omitempty"`
	Loaded          bool   `json:"loaded"`
}

func applyValidateBaseline(report validateReportJSON, againstPath, show string) (validateReportJSON, error) {
	if report.IdentityVersion == "" {
		report.IdentityVersion = currentIdentityVersion()
	}
	path := filepath.Clean(againstPath)
	data, err := os.ReadFile(path)
	if err != nil {
		return report, err
	}

	var base validateReportJSON
	if err := json.Unmarshal(data, &base); err != nil {
		return report, err
	}
	baseIdentityVersion, err := normalizeIdentityVersion(base.IdentityVersion)
	if err != nil {
		return report, err
	}

	baseSet := map[string]struct{}{}
	for _, issue := range base.Issues {
		baseIssue := issue
		baseIssue.Details = inferValidateIssueDetails(baseIssue)
		baseSet[validateIssueIdentityForVersion(baseIssue, baseIdentityVersion)] = struct{}{}
	}

	newIssues := make([]validateIssueJSON, 0, len(report.Issues))
	existingIssues := make([]validateIssueJSON, 0, len(report.Issues))
	for _, issue := range report.Issues {
		currentIssue := issue
		currentIssue.Details = inferValidateIssueDetails(currentIssue)
		if _, ok := baseSet[validateIssueIdentityForVersion(currentIssue, baseIdentityVersion)]; ok {
			existingIssues = append(existingIssues, issue)
			continue
		}
		newIssues = append(newIssues, issue)
	}

	report.Baseline = &baselineMetaJSON{
		Path:            path,
		SchemaVersion:   base.SchemaVersion,
		IdentityVersion: baseIdentityVersion,
		Loaded:          true,
	}
	report.NewErrorCount = countValidateIssueLevel(newIssues, "error")
	report.NewWarnCount = countValidateIssueLevel(newIssues, "warn")
	report.ExistingErrorCount = countValidateIssueLevel(existingIssues, "error")
	report.ExistingWarnCount = countValidateIssueLevel(existingIssues, "warn")

	switch show {
	case "new":
		report.Issues = newIssues
	case "existing":
		report.Issues = existingIssues
	case "all":
		// keep all issues
	}

	return report, nil
}

func countValidateIssueLevel(issues []validateIssueJSON, level string) int {
	count := 0
	for _, issue := range issues {
		if issue.Level == level {
			count++
		}
	}
	return count
}

func applyComplianceBaseline(sum *compliance.Summary, againstPath, show string) error {
	if sum.IdentityVersion == "" {
		sum.IdentityVersion = currentIdentityVersion()
	}
	path := filepath.Clean(againstPath)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var base compliance.Summary
	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}
	baseIdentityVersion, err := normalizeIdentityVersion(base.IdentityVersion)
	if err != nil {
		return err
	}

	baseSet := map[string]struct{}{}
	for _, doc := range base.Docs {
		for _, violation := range doc.Violations {
			baseSet[complianceViolationIdentity(doc.Path, violation)] = struct{}{}
		}
	}

	sum.NewFailed = 0
	sum.ExistingFailed = 0
	sum.NewViolations = nil

	filtered := make([]compliance.DocResult, 0, len(sum.Docs))
	newViolations := map[string]int{}

	for _, doc := range sum.Docs {
		newDocViolations := make([]string, 0, len(doc.Violations))
		existingDocViolations := make([]string, 0, len(doc.Violations))
		for _, violation := range doc.Violations {
			if _, ok := baseSet[complianceViolationIdentity(doc.Path, violation)]; ok {
				existingDocViolations = append(existingDocViolations, violation)
				continue
			}
			newDocViolations = append(newDocViolations, violation)
		}

		sort.Strings(newDocViolations)
		sort.Strings(existingDocViolations)

		if len(newDocViolations) > 0 {
			sum.NewFailed++
			for _, violation := range newDocViolations {
				newViolations[violation]++
			}
		} else if len(existingDocViolations) > 0 {
			sum.ExistingFailed++
		}

		switch show {
		case "new":
			if len(newDocViolations) > 0 {
				cp := doc
				cp.Violations = newDocViolations
				filtered = append(filtered, cp)
			}
		case "existing":
			if len(existingDocViolations) > 0 {
				cp := doc
				cp.Violations = existingDocViolations
				filtered = append(filtered, cp)
			}
		default:
			filtered = append(filtered, doc)
		}
	}

	sum.Baseline = &compliance.BaselineMeta{
		Path:            path,
		SchemaVersion:   base.SchemaVersion,
		IdentityVersion: baseIdentityVersion,
		Loaded:          true,
	}
	sum.IdentityVersion = currentIdentityVersion()
	if len(newViolations) > 0 {
		sum.NewViolations = newViolations
	}
	sum.Docs = filtered
	return nil
}

func complianceViolationIdentity(path, violation string) string {
	return violation + "|" + path
}
