package main

import (
	"github.com/clipperhouse/typewriter"
	_ "typewriters/container"
	_ "typewriters/genwriter"
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
