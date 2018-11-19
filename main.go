package main

import (
	"fmt"
	"os"

	"github.com/applike/gomono/internal/dot"
	"github.com/applike/gomono/internal/search"
)

func main1() {
	fmt.Println(search.MainPackages("./..."))
}

func main() {

	mains, err := search.MainPackages("./...")
	if err != nil {
		panic(err)
	}

	printGraphs := true // TODO: move to flag
	if printGraphs {
		for _, m := range mains {
			g, nodes, err := search.Graph(m)
			if err != nil {
				panic(err)
			}
			names := make(map[int]string)
			for _, v := range nodes {
				names[int(v.ID())] = v.ImportPath
			}

			err = dot.Dot(os.Stdout, m, g, names)
			if err != nil {
				panic(err)
			}
			fmt.Println()
		}
	}
}
