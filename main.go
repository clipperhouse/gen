package main

import (
	"github.com/clipperhouse/typewriter"
	_ "github.com/clipperhouse/typewriters/container"
	_ "github.com/clipperhouse/typewriters/genwriter"
)

func main() {
	app, err := typewriter.NewApp("+gen")
	if err != nil {
		panic(err)
	}

	app.WriteAll()
}

// +gen projections:"int" containers:"Set"
type Silly int
