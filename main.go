package main

import (
	"encoding/json"
	"fmt"

	"github.com/applike/gomono/internal/search"
)

func main() {
	pkgs, err := search.Packages(".")
	if err != nil {
		panic(err)
	}
	for _, p := range pkgs {
		b, _ := json.MarshalIndent(p, "", "  ")
		fmt.Println(string(b))
	}
}
