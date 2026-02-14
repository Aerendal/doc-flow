package index

import (
	"sort"
	"strings"
)

// markCycles detects cycles in the dependency graph and annotates records.
// It is tolerant: missing nodes are ignored for cycle purposes.
func markCycles(idx *DocumentIndex, adj map[string][]string) {
	type state uint8
	const (
		white state = iota
		gray
		black
	)

	visited := make(map[string]state)
	stack := []string{}

	var dfs func(string)
	dfs = func(node string) {
		visited[node] = gray
		stack = append(stack, node)

		neighbors := append([]string(nil), adj[node]...)
		sort.Strings(neighbors)
		for _, nbr := range neighbors {
			if visited[nbr] == white {
				dfs(nbr)
			} else if visited[nbr] == gray {
				// cycle found: stack from nbr to end + nbr
				cycle := []string{nbr}
				for i := len(stack) - 1; i >= 0; i-- {
					cycle = append(cycle, stack[i])
					if stack[i] == nbr {
						break
					}
				}
				// reverse to start from first occurrence of nbr
				for l, r := 0, len(cycle)-1; l < r; l, r = l+1, r-1 {
					cycle[l], cycle[r] = cycle[r], cycle[l]
				}
				cycle = canonicalizeCycle(cycle)
				for _, id := range cycle {
					if rec := idx.GetByID(id); rec != nil {
						rec.Cycle = cycle
					}
				}
			}
		}

		stack = stack[:len(stack)-1]
		visited[node] = black
	}

	nodes := make([]string, 0, len(adj))
	for node := range adj {
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)
	for _, node := range nodes {
		if visited[node] == white {
			dfs(node)
		}
	}
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
