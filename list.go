package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/clipperhouse/gen/typewriter"
)

func list(customFilename string) (io.Reader, error) {
	imports := []string{
		`"fmt"`,
		`"log"`,
		`"github.com/clipperhouse/gen/typewriter"`,
	}

	return execute(listStandard, customFilename, imports, listBody)
}

func listStandard() (io.Reader, error) {
	var out bytes.Buffer

	app, err := typewriter.NewApp("+gen")

	if err != nil {
		return &out, err
	}

	fmt.Fprintln(&out, "Installed typewriters:")
	for _, tw := range app.TypeWriters {
		fmt.Fprintf(&out, "  %s\n", tw.Name())
	}

	return &out, nil
}

const listBody string = `
func main() {
	app, err := typewriter.NewApp("+gen")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Installed typewriters (custom):")
	for _, tw := range app.TypeWriters {
		fmt.Println("  " + tw.Name())
	}
}
`
