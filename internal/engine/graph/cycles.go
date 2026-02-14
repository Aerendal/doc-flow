package graph

import (
	"sort"
	"strings"
)

func findDocCycles(g Graph) [][]string {
	docKind := map[string]bool{}
	for _, n := range g.Nodes {
		if n.Kind == NodeKindDoc {
			docKind[n.NodeID] = true
		}
	}

	adj := map[string][]string{}
	for _, e := range g.Edges {
		if e.Kind != EdgeDependsOn {
			continue
		}
		if !docKind[e.From] || !docKind[e.To] {
			continue
		}
		adj[e.From] = append(adj[e.From], e.To)
	}
	for k := range adj {
		sort.Strings(adj[k])
	}

	sccs := tarjanSCC(adj)
	cycles := make([][]string, 0)
	for _, scc := range sccs {
		if len(scc) <= 1 {
			node := scc[0]
			if !hasSelfLoop(adj[node], node) {
				continue
			}
		}
		sort.Strings(scc)
		allowed := map[string]bool{}
		for _, n := range scc {
			allowed[n] = true
		}
		cycle := findCycleInSCC(scc[0], adj, allowed)
		if len(cycle) == 0 {
			cycle = append([]string{}, scc...)
			cycle = append(cycle, scc[0])
		}
		cycles = append(cycles, canonicalizeCycle(cycle))
	}

	sort.Slice(cycles, func(i, j int) bool {
		return strings.Join(cycles[i], "\x00") < strings.Join(cycles[j], "\x00")
	})
	return cycles
}

func hasSelfLoop(neighbors []string, node string) bool {
	for _, n := range neighbors {
		if n == node {
			return true
		}
	}
	return false
}

func tarjanSCC(adj map[string][]string) [][]string {
	indexVal := 0
	stack := []string{}
	onStack := map[string]bool{}
	indexMap := map[string]int{}
	low := map[string]int{}
	result := [][]string{}

	nodes := make([]string, 0, len(adj))
	for n := range adj {
		nodes = append(nodes, n)
	}
	sort.Strings(nodes)

	var strongConnect func(v string)
	strongConnect = func(v string) {
		indexMap[v] = indexVal
		low[v] = indexVal
		indexVal++
		stack = append(stack, v)
		onStack[v] = true

		for _, w := range adj[v] {
			if _, ok := indexMap[w]; !ok {
				strongConnect(w)
				if low[w] < low[v] {
					low[v] = low[w]
				}
			} else if onStack[w] {
				if indexMap[w] < low[v] {
					low[v] = indexMap[w]
				}
			}
		}

		if low[v] == indexMap[v] {
			comp := []string{}
			for {
				n := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[n] = false
				comp = append(comp, n)
				if n == v {
					break
				}
			}
			result = append(result, comp)
		}
	}

	for _, n := range nodes {
		if _, ok := indexMap[n]; !ok {
			strongConnect(n)
		}
	}
	return result
}

func findCycleInSCC(start string, adj map[string][]string, allowed map[string]bool) []string {
	path := []string{start}
	seenOnPath := map[string]bool{start: true}

	var dfs func(curr string) bool
	dfs = func(curr string) bool {
		neighbors := append([]string(nil), adj[curr]...)
		sort.Strings(neighbors)
		for _, next := range neighbors {
			if !allowed[next] {
				continue
			}
			if next == start && len(path) > 1 {
				path = append(path, start)
				return true
			}
			if seenOnPath[next] {
				continue
			}
			seenOnPath[next] = true
			path = append(path, next)
			if dfs(next) {
				return true
			}
			path = path[:len(path)-1]
			delete(seenOnPath, next)
		}
		return false
	}

	if dfs(start) {
		return path
	}
	return nil
}

func canonicalizeCycle(cycle []string) []string {
	if len(cycle) == 0 {
		return nil
	}
	base := append([]string(nil), cycle...)
	if len(base) > 1 && base[0] == base[len(base)-1] {
		base = base[:len(base)-1]
	}
	if len(base) == 0 {
		return nil
	}

	forward := minLexRotation(base)
	reversed := append([]string(nil), base...)
	for l, r := 0, len(reversed)-1; l < r; l, r = l+1, r-1 {
		reversed[l], reversed[r] = reversed[r], reversed[l]
	}
	backward := minLexRotation(reversed)

	best := forward
	if strings.Join(backward, "\x00") < strings.Join(forward, "\x00") {
		best = backward
	}

	out := append([]string(nil), best...)
	out = append(out, best[0])
	return out
}

func minLexRotation(items []string) []string {
	best := rotate(items, 0)
	for i := 1; i < len(items); i++ {
		candidate := rotate(items, i)
		if strings.Join(candidate, "\x00") < strings.Join(best, "\x00") {
			best = candidate
		}
	}
	return best
}

func rotate(items []string, start int) []string {
	out := make([]string, 0, len(items))
	out = append(out, items[start:]...)
	out = append(out, items[:start]...)
	return out
}

