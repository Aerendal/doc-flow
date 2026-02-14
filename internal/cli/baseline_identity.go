package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	identityVersionV1 = "1"
	identityVersionV2 = "2"
)

var (
	missingExpectedDepRe     = regexp.MustCompile(`^brak wymaganej zależności dla ([^:]+):\s*(.+)$`)
	docIDFilenameMismatchRe  = regexp.MustCompile(`^doc_id=([^\s]+)\s+nie zgadza się z nazwą pliku\s+\(([^)]+)\)$`)
	legacySectionAliasRe     = regexp.MustCompile(`^legacy section name '(.+)' → '(.+)'$`)
	governanceFieldMissingRe = regexp.MustCompile(`^governance: missing required field '([^']+)' for status (.+)$`)
	governanceSectionMissRe  = regexp.MustCompile(`^governance: missing required section '([^']+)' for family (.+)$`)
	duplicateDocIDCountRe    = regexp.MustCompile(`\((\d+)\s+wystąpień\)$`)
	emptySectionsMetricRe    = regexp.MustCompile(`^([a-z_]+):\s+(\d+)\s+pustych sekcji`)
)

func normalizeIdentityVersion(version string) (string, error) {
	switch strings.TrimSpace(version) {
	case "", identityVersionV1:
		return identityVersionV1, nil
	case identityVersionV2:
		return identityVersionV2, nil
	default:
		return "", fmt.Errorf("unsupported identity_version: %q", version)
	}
}

func currentIdentityVersion() string {
	return identityVersionV2
}

func validateIssueIdentityForVersion(issue validateIssueJSON, identityVersion string) string {
	switch identityVersion {
	case identityVersionV2:
		return validateIssueIdentityV2(issue)
	default:
		return validateIssueIdentityV1(issue)
	}
}

func validateIssueIdentityV1(issue validateIssueJSON) string {
	docID := issue.DocID
	if docID == "" {
		docID = "-"
	}
	location := "-"
	if issue.Line > 0 {
		location = fmt.Sprintf("L%d", issue.Line)
	}
	detail := normalizeIdentityDetail(issue.Message)
	return strings.Join([]string{issue.Code, issue.Path, docID, location, detail}, "|")
}

func validateIssueIdentityV2(issue validateIssueJSON) string {
	docID := issue.DocID
	if docID == "" {
		docID = "-"
	}
	location := "-"
	if issue.Line > 0 || issue.Column > 0 {
		location = fmt.Sprintf("L%d:C%d", issue.Line, issue.Column)
	}
	detailHash := detailsFingerprint(issue.Details)
	return strings.Join([]string{issue.Code, issue.Path, docID, location, detailHash}, "|")
}

func normalizeIdentityDetail(msg string) string {
	msg = strings.ToLower(strings.TrimSpace(msg))
	if msg == "" {
		return "-"
	}
	return strings.Join(strings.Fields(msg), " ")
}

func detailsFingerprint(details map[string]any) string {
	if len(details) == 0 {
		return "-"
	}
	data, err := json.Marshal(details)
	if err != nil {
		return "-"
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:8])
}

func inferValidateIssueDetails(issue validateIssueJSON) map[string]any {
	if len(issue.Details) > 0 {
		return issue.Details
	}
	msg := strings.TrimSpace(issue.Message)
	lower := strings.ToLower(msg)
	out := map[string]any{}

	switch issue.Type {
	case "missing_field", "empty_field":
		if strings.HasPrefix(lower, "brak ") {
			field := strings.TrimSpace(strings.TrimPrefix(lower, "brak "))
			if field != "" {
				out["field"] = field
			}
		}
		if matches := emptySectionsMetricRe.FindStringSubmatch(lower); len(matches) == 3 {
			out["status"] = matches[1]
			if n, err := strconv.Atoi(matches[2]); err == nil {
				out["empty_sections"] = n
			}
			out["reason"] = "empty_sections"
		}
	case "invalid_yaml":
		out["stage"] = "frontmatter"
	case "doc_id_filename_mismatch":
		if matches := docIDFilenameMismatchRe.FindStringSubmatch(msg); len(matches) == 3 {
			out["actual_doc_id"] = matches[1]
			out["expected_doc_id"] = matches[2]
		}
	case "missing_expected_dependency":
		if matches := missingExpectedDepRe.FindStringSubmatch(msg); len(matches) == 3 {
			out["doc_type"] = strings.TrimSpace(matches[1])
			out["dependency"] = strings.TrimSpace(matches[2])
		}
	case "invalid_doc_id":
		if strings.HasPrefix(lower, "context_sources zawiera traversal:") {
			out["field"] = "context_sources"
			out["reason"] = "path_traversal"
			out["value"] = strings.TrimSpace(strings.TrimPrefix(msg, "context_sources zawiera traversal:"))
		} else if strings.HasPrefix(lower, "depends_on zawiera traversal:") {
			out["field"] = "depends_on"
			out["reason"] = "path_traversal"
			out["value"] = strings.TrimSpace(strings.TrimPrefix(msg, "depends_on zawiera traversal:"))
		}
	case "legacy_section_name":
		if matches := legacySectionAliasRe.FindStringSubmatch(msg); len(matches) == 3 {
			out["from"] = matches[1]
			out["to"] = matches[2]
		}
	case "governance_violation":
		if matches := governanceFieldMissingRe.FindStringSubmatch(msg); len(matches) == 3 {
			out["rule"] = "missing_required_field"
			out["field"] = strings.TrimSpace(matches[1])
			out["status"] = strings.TrimSpace(matches[2])
		} else if matches := governanceSectionMissRe.FindStringSubmatch(msg); len(matches) == 3 {
			out["rule"] = "missing_required_section"
			out["section"] = strings.TrimSpace(matches[1])
			out["doc_type"] = strings.TrimSpace(matches[2])
		}
	case "duplicate_doc_id":
		if matches := duplicateDocIDCountRe.FindStringSubmatch(msg); len(matches) == 2 {
			if n, err := strconv.Atoi(matches[1]); err == nil {
				out["duplicate_count"] = n
			}
		}
	case "cycle_detected":
		if strings.HasPrefix(lower, "wykryto cykl: ") {
			out["cycle"] = strings.TrimSpace(strings.TrimPrefix(msg, "wykryto cykl: "))
		}
	}

	if len(out) == 0 {
		return nil
	}
	return out
}
