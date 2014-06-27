package main

import "github.com/clipperhouse/gen/typewriter"

func run() error {
	imports := []string{
		`"os"`,
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
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	if err := app.WriteAll(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
`
