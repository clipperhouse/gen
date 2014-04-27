package main

import (
	"gen/typewriter"
	_ "gen/typewriters/standard"
)

func main() {
	app, err := typewriter.NewApp("+gen")
	if err != nil {
		panic(err)
	}

	app.WriteAll()
}

// +gen projections:"int"
type Silly int
