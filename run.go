package main

import (
	"bytes"
	"io"

	"github.com/clipperhouse/gen/typewriter"
)

func run(customFilename string) (io.Reader, error) {
	imports := []string{
		`"log"`,
		`"github.com/clipperhouse/gen/typewriter"`,
	}

	return execute(runStandard, customFilename, imports, runBody)
}

func runStandard() (io.Reader, error) {
	var out bytes.Buffer

	app, err := typewriter.NewApp("+gen")

	if err != nil {
		return &out, err
	}

	if err := app.WriteAll(); err != nil {
		return &out, err
	}

	return &out, nil
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
