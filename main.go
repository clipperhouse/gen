package main

import (
	_ "gen/typewriters/container"
	_ "gen/typewriters/genwriter"
	"github.com/clipperhouse/typewriter"
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
