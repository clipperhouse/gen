// Run this before testing: go setup.go

package main

import (
	"github.com/clipperhouse/typewriter"
	"os"
	"strings"
	_ "typewriters/container" // make sure typewriters folder is at top of GOPATH/src
)

func main() {
	// don't let bad test or gen'd files get us stuck
	filter := func(f os.FileInfo) bool {
		return !strings.HasSuffix(f.Name(), "_test.go") && !strings.HasSuffix(f.Name(), "_container.go")
	}

	a, err := typewriter.NewAppFiltered("+gen", filter)
	if err != nil {
		panic(err)
	}
	a.WriteAll()
}
