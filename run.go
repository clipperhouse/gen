package main

import (
	"fmt"
	"os"

	"github.com/clipperhouse/gen/typewriter"
)

func run() error {
	imports := []string{
		`"fmt"`,
		`"os"`,
		`"github.com/clipperhouse/gen/typewriter"`,
	}

	return execute(runStandard, imports, runBody)
}

func runStandard() (err error) {
	app, err := typewriter.NewApp("+gen")

	if err != nil {
		return err
	}

	if len(app.Types) == 0 {
		return fmt.Errorf("No types marked with +gen were found. See http://clipperhouse.github.io/gen to get started, or type %s help.", os.Args[0])
	}

	if len(app.TypeWriters) == 0 {
		return fmt.Errorf("No typewriters were imported. See http://clipperhouse.github.io/gen to get started, or type %s help.", os.Args[0])
	}

	err = app.WriteAll()

	if err != nil {
		return err
	}

	return nil
}

const runBody string = `
func main() {
	if err := gen(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func gen() error {
	app, err := typewriter.NewApp("+gen")

	if err != nil {
		return err
	}

	if len(app.Types) == 0 {
		return fmt.Errorf("No types marked with +gen were found. See http://clipperhouse.github.io/gen to get started, or type %s help.", os.Args[0])
	}

	if len(app.TypeWriters) == 0 {
		return fmt.Errorf("No typewriters were imported. See http://clipperhouse.github.io/gen to get started, or type %s help.", os.Args[0])
	}

	if err := app.WriteAll(); err != nil {
		return err
	}

	return nil
}
`
