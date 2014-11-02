// Run this before testing: go setup.go

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/clipperhouse/gen/typewriter"
	_ "github.com/clipperhouse/gen/typewriters/slice"
)

func main() {
	// don't let bad test or gen files get us stuck
	filter := func(f os.FileInfo) bool {
		return !strings.HasSuffix(f.Name(), "_slice.go") && !strings.HasSuffix(f.Name(), "_test.go")
	}

	a, err := typewriter.NewAppFiltered("+test", filter)
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := a.WriteAll(); err != nil {
		fmt.Println(err)
		return
	}
}
