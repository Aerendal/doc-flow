package graph

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"docflow/internal/engine"
	execache "docflow/internal/exec/cache"
	"docflow/pkg/index"
)

var markdownLinkRe = regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)

func BuildGraph(facts []execache.DocumentFacts, opts BuildGraphOptions) (Graph, []engine.Issue) {
	g := Graph{
		SchemaVersion: "1.0",
		Root:          ".",
		Nodes:         []Node{},
		Edges:         []Edge{},
		Cycles:        [][]string{},
	}
	if !opts.Portable && len(facts) > 0 && facts[0].Root != "" {
		g.Root = facts[0].Root
	}

	nodeByID := map[string]Node{}
	docNodeByPath := map[string]string{}
	docNodesByDocID := map[string][]string{}
	pathSet := map[string]bool{}
	issues := []engine.Issue{}

	for _, f := range facts {
		path := index.NormalizePath(f.Path)
		if path == "" {
			continue
		}
		pathSet[path] = true
		nodeID := docNodeID(path)
		docNodeByPath[path] = nodeID
		docNodesByDocID[f.DocID] = append(docNodesByDocID[f.DocID], nodeID)
		nodeByID[nodeID] = Node{
			NodeID: nodeID,
			Kind:   NodeKindDoc,
			Path:   path,
			DocID:  f.DocID,
			Meta: NodeMeta{
				DocType:  f.DocType,
				Status:   f.Status,
				Language: f.Language,
			},
		}
	}

	edgeSet := map[string]Edge{}
	addEdge := func(e Edge) {
		key := string(e.Kind) + "|" + e.From + "|" + e.To + "|" + e.Evidence + "|" + e.Source
		if _, ok := edgeSet[key]; ok {
			return
		}
		edgeSet[key] = e
	}

	for _, f := range facts {
		srcPath := index.NormalizePath(f.Path)
		srcID := docNodeByPath[srcPath]
		if srcID == "" {
			continue
		}

		for _, depID := range f.DependsOn {
			depID = strings.TrimSpace(depID)
			if depID == "" {
				continue
			}
			targets := docNodesByDocID[depID]
			if len(targets) == 0 {
				extID := "external:doc_id:" + depID
				nodeByID[extID] = Node{
					NodeID: extID,
					Kind:   NodeKindExternal,
					DocID:  depID,
				}
				addEdge(Edge{
					From:     srcID,
					To:       extID,
					Kind:     EdgeDependsOn,
					Source:   "frontmatter",
					Evidence: depID,
				})
				issues = append(issues, engine.Issue{
					Code:     "DOCFLOW.GRAPH.MISSING_DEPENDENCY_TARGET",
					Severity: engine.SeverityWarn,
					Type:     "missing_dependency_target",
					Kind:     engine.IssueKindAnalysis,
					Path:     srcPath,
					DocID:    f.DocID,
					Message:  "missing depends_on target: " + depID,
					Details:  map[string]any{"depends_on": depID},
				})
				continue
			}
			sort.Strings(targets)
			for _, tgt := range targets {
				addEdge(Edge{
					From:     srcID,
					To:       tgt,
					Kind:     EdgeDependsOn,
					Source:   "frontmatter",
					Evidence: depID,
				})
			}
		}

		if opts.IncludeTemplates {
			templateID := strings.TrimSpace(f.TemplateSource)
			if templateID != "" {
				targets := docNodesByDocID[templateID]
				if len(targets) == 0 {
					tNodeID := "template:" + templateID
					nodeByID[tNodeID] = Node{
						NodeID: tNodeID,
						Kind:   NodeKindTemplate,
						DocID:  templateID,
					}
					addEdge(Edge{
						From:     tNodeID,
						To:       srcID,
						Kind:     EdgeUsesTemplate,
						Source:   "frontmatter",
						Evidence: templateID,
					})
				} else {
					sort.Strings(targets)
					for _, tNodeID := range targets {
						addEdge(Edge{
							From:     tNodeID,
							To:       srcID,
							Kind:     EdgeUsesTemplate,
							Source:   "frontmatter",
							Evidence: templateID,
						})
					}
				}
			}
		}

		if opts.IncludeLinks {
			for _, target := range extractLinks(f.Content) {
				if strings.HasPrefix(target, "doc_id:") {
					docID := strings.TrimPrefix(target, "doc_id:")
					targets := docNodesByDocID[docID]
					if len(targets) == 0 {
						extID := "external:doc_id:" + docID
						nodeByID[extID] = Node{
							NodeID: extID,
							Kind:   NodeKindExternal,
							DocID:  docID,
						}
						addEdge(Edge{
							From:     srcID,
							To:       extID,
							Kind:     EdgeLinksTo,
							Source:   "markdown",
							Evidence: target,
						})
						continue
					}
					sort.Strings(targets)
					for _, tgt := range targets {
						addEdge(Edge{
							From:     srcID,
							To:       tgt,
							Kind:     EdgeLinksTo,
							Source:   "markdown",
							Evidence: target,
						})
					}
					continue
				}

				normalized := normalizeLinkPath(srcPath, target)
				if normalized == "" {
					continue
				}
				if targetID, ok := docNodeByPath[normalized]; ok {
					addEdge(Edge{
						From:     srcID,
						To:       targetID,
						Kind:     EdgeLinksTo,
						Source:   "markdown",
						Evidence: target,
					})
					continue
				}
				if pathSet[normalized] {
					targetID := docNodeID(normalized)
					addEdge(Edge{
						From:     srcID,
						To:       targetID,
						Kind:     EdgeLinksTo,
						Source:   "markdown",
						Evidence: target,
					})
					continue
				}
				extID := "external:path:" + normalized
				nodeByID[extID] = Node{
					NodeID: extID,
					Kind:   NodeKindExternal,
					Path:   normalized,
				}
				addEdge(Edge{
					From:     srcID,
					To:       extID,
					Kind:     EdgeLinksTo,
					Source:   "markdown",
					Evidence: target,
				})
			}
		}
	}

	g.Nodes = make([]Node, 0, len(nodeByID))
	for _, n := range nodeByID {
		g.Nodes = append(g.Nodes, n)
	}
	sort.Slice(g.Nodes, func(i, j int) bool {
		if g.Nodes[i].Kind != g.Nodes[j].Kind {
			return g.Nodes[i].Kind < g.Nodes[j].Kind
		}
		if g.Nodes[i].Path != g.Nodes[j].Path {
			return g.Nodes[i].Path < g.Nodes[j].Path
		}
		if g.Nodes[i].DocID != g.Nodes[j].DocID {
			return g.Nodes[i].DocID < g.Nodes[j].DocID
		}
		return g.Nodes[i].NodeID < g.Nodes[j].NodeID
	})

	g.Edges = make([]Edge, 0, len(edgeSet))
	for _, e := range edgeSet {
		g.Edges = append(g.Edges, e)
	}
	sort.Slice(g.Edges, func(i, j int) bool {
		if g.Edges[i].Kind != g.Edges[j].Kind {
			return g.Edges[i].Kind < g.Edges[j].Kind
		}
		if g.Edges[i].From != g.Edges[j].From {
			return g.Edges[i].From < g.Edges[j].From
		}
		if g.Edges[i].To != g.Edges[j].To {
			return g.Edges[i].To < g.Edges[j].To
		}
		if g.Edges[i].Evidence != g.Edges[j].Evidence {
			return g.Edges[i].Evidence < g.Edges[j].Evidence
		}
		return g.Edges[i].Source < g.Edges[j].Source
	})

	g.Cycles = findDocCycles(g)
	g.Stats = computeStats(g)
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
		return issues[i].Message < issues[j].Message
	})
	return g, issues
}

func docNodeID(path string) string {
	return "doc:" + path
}

func extractLinks(content string) []string {
	matches := markdownLinkRe.FindAllStringSubmatch(content, -1)
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		target := strings.TrimSpace(m[1])
		if target == "" {
			continue
		}
		target = stripLinkAnchor(target)
		lower := strings.ToLower(target)
		if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") || strings.HasPrefix(lower, "mailto:") {
			continue
		}
		if strings.HasPrefix(target, "#") {
			continue
		}
		out = append(out, target)
	}
	sort.Strings(out)
	return out
}

func stripLinkAnchor(target string) string {
	if idx := strings.Index(target, "#"); idx >= 0 {
		target = target[:idx]
	}
	if idx := strings.Index(target, "?"); idx >= 0 {
		target = target[:idx]
	}
	return strings.TrimSpace(target)
}

func normalizeLinkPath(srcPath, target string) string {
	if target == "" {
		return ""
	}
	target = strings.TrimSpace(target)
	if target == "" {
		return ""
	}
	if strings.HasPrefix(target, "/") {
		target = strings.TrimPrefix(target, "/")
	}
	srcDir := filepath.Dir(srcPath)
	joined := filepath.Clean(filepath.Join(srcDir, target))
	return index.NormalizePath(joined)
}

func computeStats(g Graph) Stats {
	outDeg := map[string]int{}
	inDeg := map[string]int{}
	maxOut := 0
	maxIn := 0
	for _, e := range g.Edges {
		outDeg[e.From]++
		inDeg[e.To]++
		if outDeg[e.From] > maxOut {
			maxOut = outDeg[e.From]
		}
		if inDeg[e.To] > maxIn {
			maxIn = inDeg[e.To]
		}
	}
	return Stats{
		NodeCount:    len(g.Nodes),
		EdgeCount:    len(g.Edges),
		CyclesCount:  len(g.Cycles),
		MaxOutDegree: maxOut,
		MaxInDegree:  maxIn,
	}
}
