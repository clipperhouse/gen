// Run this before testing: go setup.go

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/clipperhouse/gen/typewriter"
	_ "github.com/clipperhouse/gen/typewriters/slice" // make sure typewriters folder is at top of GOPATH/src
)

func main() {
	// don't let bad test or gen files get us stuck
	filter := func(f os.FileInfo) bool {
		return !strings.HasSuffix(f.Name(), "_test.go") && !strings.HasSuffix(f.Name(), "_slice.go")
	}

	a, err := typewriter.NewAppFiltered("+test", filter)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := a.WriteAll(); err != nil {
		fmt.Println(err)
		return
	}
}
