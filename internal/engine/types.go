package engine

type Severity string

const (
	SeverityError Severity = "error"
	SeverityWarn  Severity = "warn"
	SeverityInfo  Severity = "info"
)

type IssueKind string

const (
	IssueKindValidate   IssueKind = "validate"
	IssueKindCompliance IssueKind = "compliance"
	IssueKindAnalysis   IssueKind = "analysis"
)

type Location struct {
	Line   int `json:"line,omitempty"`
	Column int `json:"column,omitempty"`
}

type Issue struct {
	Code     string         `json:"code"`
	Severity Severity       `json:"level"`
	Type     string         `json:"type"`
	Kind     IssueKind      `json:"kind"`
	Path     string         `json:"path"`
	DocID    string         `json:"doc_id,omitempty"`
	Message  string         `json:"message"`
	Details  map[string]any `json:"details,omitempty"`
	Location *Location      `json:"location,omitempty"`
}

type ReportMeta struct {
	SchemaVersion string `json:"schema_version"`
	GeneratedAt   string `json:"generated_at,omitempty"`
	Root          string `json:"root,omitempty"`
	ToolVersion   string `json:"tool_version,omitempty"`
}

type BaselineMeta struct {
	Loaded          bool   `json:"loaded"`
	Path            string `json:"path,omitempty"`
	SchemaVersion   string `json:"schema_version,omitempty"`
	IdentityVersion string `json:"identity_version,omitempty"`
}

type ReportViewMeta struct {
	Against string `json:"against,omitempty"`
	FailOn  string `json:"fail_on,omitempty"`
	Show    string `json:"show,omitempty"`
}

type ValidateReport struct {
	ReportMeta
	ReportViewMeta

	IdentityVersion string `json:"identity_version,omitempty"`

	Baseline *BaselineMeta `json:"baseline,omitempty"`

	Files     int `json:"files"`
	Documents int `json:"documents"`

	ErrorCount int `json:"error_count"`
	WarnCount  int `json:"warn_count"`

	NewErrorCount      int `json:"new_error_count,omitempty"`
	NewWarnCount       int `json:"new_warn_count,omitempty"`
	ExistingErrorCount int `json:"existing_error_count,omitempty"`
	ExistingWarnCount  int `json:"existing_warn_count,omitempty"`

	Issues []Issue `json:"issues"`
}

type ComplianceDoc struct {
	DocID      string            `json:"doc_id"`
	Path       string            `json:"path"`
	Status     string            `json:"status"`
	DocType    string            `json:"doc_type"`
	Violations []string          `json:"violations"`
	Details    map[string]string `json:"details,omitempty"`
}

type ComplianceReport struct {
	ReportMeta
	ReportViewMeta

	IdentityVersion string `json:"identity_version,omitempty"`

	Baseline      *BaselineMeta `json:"baseline,omitempty"`
	RulesPath     string        `json:"rules_path,omitempty"`
	RulesChecksum string        `json:"rules_checksum,omitempty"`

	Documents int     `json:"documents"`
	Passed    int     `json:"passed"`
	Failed    int     `json:"failed"`
	PassRate  float64 `json:"pass_rate"`

	ViolationsCount map[string]int `json:"violations_count"`

	NewFailed          int            `json:"new_failed,omitempty"`
	ExistingFailed     int            `json:"existing_failed,omitempty"`
	NewViolationsCount map[string]int `json:"new_violations_count,omitempty"`

	DuplicateDocIDs map[string][]string `json:"duplicate_doc_ids,omitempty"`
	Docs            []ComplianceDoc     `json:"docs"`
}

type RunMode string

const (
	RunModeValidate   RunMode = "validate"
	RunModeCompliance RunMode = "compliance"
	RunModeScan       RunMode = "scan"
	RunModeHealth     RunMode = "health"
)

type OutputFormat string

const (
	OutputFormatText  OutputFormat = "text"
	OutputFormatJSON  OutputFormat = "json"
	OutputFormatSARIF OutputFormat = "sarif"
	OutputFormatHTML  OutputFormat = "html"
)

type RunOptions struct {
	ConfigPath string
	Root       string

	Mode RunMode

	Deterministic bool
	Portable      bool

	AgainstPath string
	FailOn      string
	Show        string

	RulesPath string

	Format     OutputFormat
	OutputPath string
}
