package search

import (
	"go/build"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

// Graph returns the dependency graph of all packages named by pattern
func Graph(pattern string) (graph.Graph, map[string]graph.Node, error) {
	var (
		g       *simple.DirectedGraph
		pkgs    []*build.Package
		p       *build.Package
		nodes   = make(map[string]graph.Node)
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
			addNode(g, p.ImportPath, nodes)
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
				addNode(g, importedPackage.ImportPath, nodes)
			}

			from := nodes[p.ImportPath]
			to := nodes[importedPackage.ImportPath]

			e := g.NewEdge(from, to)
			g.SetEdge(e)

			pkgs = append(pkgs, importedPackage)
		}
		visited[p.ImportPath] = struct{}{}
	}
	return g, nodes, nil
}

func addNode(g graph.NodeAdder, name string, nodes map[string]graph.Node) {
	n := g.NewNode()
	g.AddNode(n)
	nodes[name] = n
}
