package main

import "github.com/clipperhouse/gen/typewriter"

func run() error {
	imports := []string{
		`"log"`,
		`"github.com/clipperhouse/gen/typewriter"`,
	}

	return execute(runStandard, imports, runBody)
}

func runStandard() error {
	app, err := typewriter.NewApp("+gen")

	if err != nil {
		return err
	}

	if err := app.WriteAll(); err != nil {
		return err
	}

	return nil
}

const runBody string = `
func main() {
	app, err := typewriter.NewApp("+gen")

	if err != nil {
		log.Fatal(err)
	}

	if err := app.WriteAll(); err != nil {
		log.Fatal(err)
	}
}
`
