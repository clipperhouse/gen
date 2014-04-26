package main

import (
	"fmt"
	"gen/typewriter"
	_ "gen/typewriters/standard"
	"os"
)

func main() {
	pkgs, err := typewriter.GetPackages()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, p := range pkgs {
		for _, t := range p.Types {
			writeFile(t)
		}
	}
}

func writeFile(t typewriter.Type) {
	f, err := os.Create(t.LocalName() + "_gen.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	typewriter.Write(f, t)
}

// +gen
type Silly int
