// Run this before testing: go setup.go

package main

import (
	"github.com/clipperhouse/typewriter"
	_ "typewriters/genwriter" // make sure typewriters folder is at top of GOPATH/src
)

func main() {
	a, err := typewriter.NewApp("+gen")
	if err != nil {
		panic(err)
	}
	a.WriteAll()
}
