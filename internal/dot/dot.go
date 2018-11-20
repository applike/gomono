package dot

import (
	"io"
	"text/template"

	"gonum.org/v1/gonum/graph"
)

type edges struct {
	Name  string
	Edges []edge
}

type edge struct {
	A, B interface{}
}

// Dot writes the graph g in dot format to w.
// names is used as a lookup-table for Nodenames by ID
func Dot(w io.Writer, name string, g graph.Graph, names map[int]string) error {
	ns := g.Nodes()

	var pairs = make([]edge, 0)

	for hasNextFrom := ns.Next(); hasNextFrom; hasNextFrom = ns.Next() {
		n := ns.Node()
		es := g.From(n.ID())
		for hasNextTo := es.Next(); hasNextTo; hasNextTo = es.Next() {
			t := es.Node()
			pairs = append(pairs, edge{A: names[int(n.ID())], B: names[int(t.ID())]})
		}
	}

	funcMap := template.FuncMap{
		"last": func(i int) bool {
			return i == len(pairs)-1
		},
	}

	tmpl, err := template.New("dotformat").Funcs(funcMap).Parse(format)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, edges{Name: name, Edges: pairs})
}

var format = `{{define "T"}}    "{{.A}}" -> "{{.B}}"
{{end}}{{define "END"}}    "{{.A}}" -> "{{.B}}"{{end}}digraph {{.Name}} {
{{range $index, $element := .Edges}}{{if last $index}}{{template "END" $element}}{{else}}{{template "T" $element}}{{end}}{{end}}
}
`
