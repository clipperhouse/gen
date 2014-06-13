package main

import (
	"github.com/clipperhouse/gen/typewriter"
	_ "github.com/clipperhouse/gen/typewriters/container"
	_ "github.com/clipperhouse/gen/typewriters/genwriter"
)

func main() {
	app, err := typewriter.NewApp("+gen")

	if err != nil {
		panic(err)
	}

	if err := app.WriteAll(); err != nil {
		panic(err)
	}
}
