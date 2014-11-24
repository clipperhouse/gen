package main

import (
	"fmt"

	"github.com/clipperhouse/typewriter"
)

func list(c config) error {
	imports := typewriter.NewImportSpecSet(
		typewriter.ImportSpec{Path: "fmt"},
		typewriter.ImportSpec{Path: "os"},
		typewriter.ImportSpec{Path: "github.com/clipperhouse/typewriter"},
	)

	listFunc := func() error {
		app, err := typewriter.NewApp("+gen")

		if err != nil {
			return err
		}

		fmt.Fprintln(c.out, "Installed typewriters:")
		for _, tw := range app.TypeWriters {
			fmt.Fprintf(c.out, "  %s\n", tw.Name())
		}

		return nil
	}

	return execute(listFunc, c, imports, listBody)
}

const listBody string = `
func main() {
	app, err := typewriter.NewApp("+gen")

	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println("Imported typewriters:")
	for _, tw := range app.TypeWriters {
		fmt.Println("  " + tw.Name())
	}
}
`
