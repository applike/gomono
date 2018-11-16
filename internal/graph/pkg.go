package graph

import (
	"github.com/applike/gomono/internal/search"
	"github.com/applike/graph"
	"github.com/applike/graph/directed"
)

// Graph returns the dependency graph of all packages named by pattern
func Graph(pattern string) (graph.Graph, map[string]graph.Node, error) {
	var g = directed.New()
	var m = make(map[string]graph.Node)
	pkgs, err := search.Packages(pattern)
	if err != nil {
		return nil, nil, err
	}
	for len(pkgs) > 0 {
		p, pkgs := pkgs[0], pkgs[1:]
		_, ok := m[p.ImportPath]
		if ok {
			// We already added this node, thus
			// we also added it's dependencies
			continue
		}
		from := g.NewNode()
		m[p.ImportPath] = from
		g.AddNode(from)

		for _, imp := range p.Imports {
			to, ok := m[imp]
			if !ok {
				to := g.NewNode()
				g.AddNode(to)
			}
			g.AddEdge(from, to)
			impPkg, err := search.Import(imp)
			if err != nil {
				return nil, nil, err
			}
			pkgs = append(pkgs, impPkg)
		}
	}
	return g, m, nil
}
