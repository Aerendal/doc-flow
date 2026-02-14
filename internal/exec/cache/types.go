package cache

type Fingerprint struct {
	Size         int64
	MTimeUnix    int64
	MetaChecksum string
	BodyChecksum string
}

type RunStats struct {
	ChangedFiles int
	CacheHits    int
	CacheMisses  int
}

type Cache interface {
	Get(path string, fp Fingerprint) (*DocumentFacts, bool, error)
	Put(path string, fp Fingerprint, facts DocumentFacts) error
	BeginRun(mode string, root string) (string, error)
	EndRun(runID string, stats RunStats) error
	Close() error
}

type DocumentFacts struct {
	Path      string `json:"path"`
	AbsPath   string `json:"-"`
	Root      string `json:"-"`
	Size      int64  `json:"size"`
	MTimeUnix int64  `json:"mtime_unix"`

	ChecksumMeta string `json:"checksum_meta,omitempty"`
	ChecksumBody string `json:"checksum_body,omitempty"`
	ChecksumFull string `json:"checksum_full,omitempty"`

	DocID          string   `json:"doc_id,omitempty"`
	Title          string   `json:"title,omitempty"`
	DocType        string   `json:"doc_type,omitempty"`
	Status         string   `json:"status,omitempty"`
	Version        string   `json:"version,omitempty"`
	Priority       string   `json:"priority,omitempty"`
	Owner          string   `json:"owner,omitempty"`
	Language       string   `json:"language,omitempty"`
	DependsOn      []string `json:"depends_on,omitempty"`
	ContextSources []string `json:"context_sources,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	TemplateSource string   `json:"template_source,omitempty"`

	HeadingCount int `json:"heading_count"`
	MaxDepth     int `json:"max_depth"`
	Lines        int `json:"lines"`
	CodeBlocks   int `json:"code_blocks,omitempty"`
	Tables       int `json:"tables,omitempty"`

	HasFrontmatter bool            `json:"has_frontmatter"`
	RawFields      map[string]bool `json:"raw_fields,omitempty"`
	RawYAML        string          `json:"raw_yaml,omitempty"`
	Body           string          `json:"body,omitempty"`
	Content        string          `json:"content,omitempty"`

	ParseErrorCode    string `json:"parse_error_code,omitempty"`
	ParseErrorMessage string `json:"parse_error_message,omitempty"`
	ParseErrorLine    int    `json:"parse_error_line,omitempty"`
}
