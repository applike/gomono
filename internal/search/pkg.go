package search

import (
	"go/build"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type Package struct {
	*build.Package
	graph.Node
}

// Graph returns the dependency graph of all packages named by pattern
func Graph(pattern []string) (graph.Graph, map[string]*Package, error) {
	var (
		g       *simple.DirectedGraph
		pkgs    []*build.Package
		p       *build.Package
		nodes   = make(map[string]*Package)
		visited = make(map[string]struct{})
		err     error
	)

	g = simple.NewDirectedGraph()
	pkgs, err = Packages(pattern)
	if err != nil {
		return nil, nil, err
	}

	for len(pkgs) > 0 {
		p = pkgs[0]
		pkgs = pkgs[1:]
		if _, ok := visited[p.ImportPath]; ok {
			continue // already been here
		}

		if _, ok := nodes[p.ImportPath]; !ok {
			addNode(g, p, nodes)
		}

		for _, ip := range p.Imports {
			importedPackage, err := Import(ip)
			if importedPackage.Goroot {
				continue // Don't care about stdlib
			}
			if err != nil {
				return nil, nil, err
			}

			if _, ok := nodes[importedPackage.ImportPath]; !ok {
				addNode(g, importedPackage, nodes)
			}

			from := nodes[p.ImportPath].Node
			to := nodes[importedPackage.ImportPath].Node

			e := g.NewEdge(from, to)
			g.SetEdge(e)

			pkgs = append(pkgs, importedPackage)
		}
		visited[p.ImportPath] = struct{}{}
	}
	return g, nodes, nil
}

func addNode(g graph.NodeAdder, p *build.Package, nodes map[string]*Package) {
	n := g.NewNode()
	g.AddNode(n)
	nodes[p.ImportPath] = &Package{
		Package: p,
		Node:    n,
	}
}
