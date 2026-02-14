package graph

import "sort"

func NodeIDsForDocID(g Graph, docID string) []string {
	if docID == "" {
		return nil
	}
	out := make([]string, 0)
	for _, n := range g.Nodes {
		if n.DocID == docID && (n.Kind == NodeKindDoc || n.Kind == NodeKindTemplate) {
			out = append(out, n.NodeID)
		}
	}
	sort.Strings(out)
	return uniqueStrings(out)
}

func ReachableNodeIDs(g Graph, startNodeIDs []string, opts ImpactOptions) []string {
	nodeSet := make(map[string]bool, len(g.Nodes))
	for _, n := range g.Nodes {
		nodeSet[n.NodeID] = true
	}

	start := make([]string, 0, len(startNodeIDs))
	startSeen := map[string]bool{}
	for _, id := range startNodeIDs {
		if id == "" || !nodeSet[id] || startSeen[id] {
			continue
		}
		startSeen[id] = true
		start = append(start, id)
	}
	sort.Strings(start)
	if len(start) == 0 {
		return nil
	}

	adj := buildAdjacency(g, opts)
	dist := map[string]int{}
	queue := make([]string, 0, len(start))
	for _, id := range start {
		dist[id] = 0
		queue = append(queue, id)
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		currDepth := dist[curr]
		if opts.MaxDepth >= 0 && currDepth >= opts.MaxDepth {
			continue
		}

		neighbors := adj[curr]
		for _, next := range neighbors {
			if _, ok := dist[next]; ok {
				continue
			}
			dist[next] = currDepth + 1
			queue = append(queue, next)
		}
	}

	out := make([]string, 0, len(dist))
	for id := range dist {
		if !opts.IncludeStart && startSeen[id] {
			continue
		}
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

func ReachableDocPaths(g Graph, startNodeIDs []string, opts ImpactOptions) []string {
	nodeByID := map[string]Node{}
	for _, n := range g.Nodes {
		nodeByID[n.NodeID] = n
	}

	nodeIDs := ReachableNodeIDs(g, startNodeIDs, opts)
	paths := make([]string, 0, len(nodeIDs))
	for _, id := range nodeIDs {
		n, ok := nodeByID[id]
		if !ok || n.Kind != NodeKindDoc || n.Path == "" {
			continue
		}
		paths = append(paths, n.Path)
	}
	sort.Strings(paths)
	return uniqueStrings(paths)
}

func buildAdjacency(g Graph, opts ImpactOptions) map[string][]string {
	edgeKinds := opts.EdgeKinds
	if len(edgeKinds) == 0 {
		edgeKinds = map[EdgeKind]bool{
			EdgeDependsOn:    true,
			EdgeUsesTemplate: true,
		}
	}

	adj := map[string][]string{}
	for _, e := range g.Edges {
		if !edgeKinds[e.Kind] {
			continue
		}
		adj[e.From] = append(adj[e.From], e.To)
		if opts.IncludeReverse {
			adj[e.To] = append(adj[e.To], e.From)
		}
	}
	for from := range adj {
		sort.Strings(adj[from])
		adj[from] = uniqueStrings(adj[from])
	}
	return adj
}

func uniqueStrings(values []string) []string {
	if len(values) < 2 {
		return append([]string(nil), values...)
	}
	out := make([]string, 0, len(values))
	prev := ""
	for i, v := range values {
		if i == 0 || v != prev {
			out = append(out, v)
			prev = v
		}
	}
	return out
}
