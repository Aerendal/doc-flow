package deps

import (
	"fmt"
	"sort"
)

type Graph struct {
	Nodes map[string]struct{}
	Edges map[string][]string // from -> to (depends_on)
	known map[string]struct{}
}

type Problem struct {
	Type   string // "missing" | "cycle"
	Detail string
}

type Result struct {
	Order    []string
	Problems []Problem
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]struct{}),
		Edges: make(map[string][]string),
		known: make(map[string]struct{}),
	}
}

func (g *Graph) AddNode(id string) {
	g.Nodes[id] = struct{}{}
}

func (g *Graph) AddEdge(from, to string) {
	g.Edges[from] = append(g.Edges[from], to)
}

// Build constructs graph where edges go from dependency -> document.
// Only doc_ids obecne w depsMap są traktowane jako istniejące; inne są raportowane jako missing.
func Build(depsMap map[string][]string) *Graph {
	g := NewGraph()
	for id := range depsMap {
		g.known[id] = struct{}{}
	}
	for id, deps := range depsMap {
		g.AddNode(id)
		for _, d := range deps {
			if _, ok := g.known[d]; ok {
				g.AddNode(d)
			}
			g.AddEdge(d, id)
		}
	}
	return g
}

// Topo returns stable topological order and problems (missing nodes, cycles).
func (g *Graph) Topo() Result {
	order := []string{}
	problems := []Problem{}

	inDeg := make(map[string]int)
	for n := range g.Nodes {
		inDeg[n] = 0
	}
	missingReported := make(map[string]bool)
	for from, tos := range g.Edges {
		if _, ok := inDeg[from]; !ok {
			if !missingReported[from] {
				problems = append(problems, Problem{Type: "missing", Detail: fmt.Sprintf("brak doc_id: %s (wspomniany jako zależność)", from)})
				missingReported[from] = true
			}
			continue
		}
		for _, to := range tos {
			if _, ok := inDeg[to]; !ok {
				if !missingReported[to] {
					problems = append(problems, Problem{Type: "missing", Detail: fmt.Sprintf("brak doc_id: %s (wymagany przez %s)", to, from)})
					missingReported[to] = true
				}
				continue
			}
			inDeg[to]++
		}
	}

	// stable queue sorted by doc_id
	var queue []string
	for n, deg := range inDeg {
		if deg == 0 {
			queue = append(queue, n)
		}
	}
	sort.Strings(queue)

	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		order = append(order, n)

		for _, to := range g.Edges[n] {
			if _, ok := inDeg[to]; !ok {
				continue
			}
			inDeg[to]--
			if inDeg[to] == 0 {
				queue = append(queue, to)
			}
		}
		sort.Strings(queue)
	}

	if len(order) != len(inDeg) {
		problems = append(problems, Problem{Type: "cycle", Detail: "wykryto cykl w zależnościach"})
	}

	return Result{Order: order, Problems: problems}
}
