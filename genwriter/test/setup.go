// Run this before testing: go setup.go

package main

import (
	"github.com/clipperhouse/typewriter"
	"os"
	"strings"
	_ "typewriters/genwriter" // make sure typewriters folder is at top of GOPATH/src
)

func main() {
	// don't let bad test or gen files get us stuck
	filter := func(f os.FileInfo) bool {
		return !strings.HasSuffix(f.Name(), "_test.go") && !strings.HasSuffix(f.Name(), "_gen.go")
	}

	a, err := typewriter.NewAppFiltered("+gen", filter)
	if err != nil {
		panic(err)
	}
	a.WriteAll()
}
