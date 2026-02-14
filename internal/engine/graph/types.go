package graph

type NodeKind string

const (
	NodeKindDoc      NodeKind = "doc"
	NodeKindTemplate NodeKind = "template"
	NodeKindSection  NodeKind = "section"
	NodeKindAsset    NodeKind = "asset"
	NodeKindExternal NodeKind = "external"
)

type EdgeKind string

const (
	EdgeDependsOn       EdgeKind = "depends_on"
	EdgeUsesTemplate    EdgeKind = "uses_template"
	EdgeRequiresSection EdgeKind = "requires_section"
	EdgeLinksTo         EdgeKind = "links_to"
	EdgeDerivedFrom     EdgeKind = "derived_from"
)

type NodeMeta struct {
	DocType  string `json:"doc_type,omitempty"`
	Status   string `json:"status,omitempty"`
	Language string `json:"language,omitempty"`
}

type Node struct {
	NodeID string   `json:"node_id"`
	Kind   NodeKind `json:"kind"`
	Path   string   `json:"path,omitempty"`
	DocID  string   `json:"doc_id,omitempty"`
	Meta   NodeMeta `json:"meta,omitempty"`
}

type Edge struct {
	From     string   `json:"from"`
	To       string   `json:"to"`
	Kind     EdgeKind `json:"kind"`
	Weight   int      `json:"weight,omitempty"`
	Evidence string   `json:"evidence,omitempty"`
	Source   string   `json:"source,omitempty"`
}

type Stats struct {
	NodeCount    int `json:"node_count"`
	EdgeCount    int `json:"edge_count"`
	CyclesCount  int `json:"cycles_count"`
	MaxOutDegree int `json:"max_out_degree"`
	MaxInDegree  int `json:"max_in_degree"`
}

type Graph struct {
	SchemaVersion string     `json:"schema_version"`
	Root          string     `json:"root"`
	Nodes         []Node     `json:"nodes"`
	Edges         []Edge     `json:"edges"`
	Stats         Stats      `json:"stats"`
	Cycles        [][]string `json:"cycles,omitempty"`
}

type BuildGraphOptions struct {
	IncludeLinks     bool
	IncludeTemplates bool
	Deterministic    bool
	Portable         bool
}

type ImpactOptions struct {
	MaxDepth       int
	EdgeKinds      map[EdgeKind]bool
	IncludeReverse bool
	IncludeStart   bool
}
