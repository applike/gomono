package main

import (
	"fmt"

	"github.com/applike/gomono/internal/search"
)

func main() {
	// pkgs, err := search.Packages(".")
	// if err != nil {
	// 	panic(err)
	// }
	// for _, p := range pkgs {
	// 	b, _ := json.MarshalIndent(p, "", "  ")
	// 	fmt.Println(string(b))
	// }

	g, m, err := search.Graph(".")
	if err != nil {
		panic(err)
	}

	for k, v := range m {
		fmt.Printf("%v: %v\n", v, k)
	}

	fmt.Printf("%v map entries\n", len(m))
	fmt.Printf("%v graph entries\n", g.Nodes().Len())
}
