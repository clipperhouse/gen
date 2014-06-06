// Run this before testing: go setup.go

package main

import (
	"github.com/clipperhouse/typewriter"
	_ "github.com/clipperhouse/typewriters/genwriter" // make sure typewriters folder is at top of GOPATH/src
	"os"
	"strings"
)

func main() {
	// don't let bad test or gen files get us stuck
	filter := func(f os.FileInfo) bool {
		return !strings.HasSuffix(f.Name(), "_test.go") && !strings.HasSuffix(f.Name(), "_genwriter.go")
	}

	a, err := typewriter.NewAppFiltered("+test", filter)
	if err != nil {
		panic(err)
	}
	a.WriteAll()
}
