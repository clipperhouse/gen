package main

import (
	"fmt"

	"github.com/clipperhouse/gen/typewriter"
)

func list(customFilename string) error {
	imports := []string{
		`"fmt"`,
		`"log"`,
		`"github.com/clipperhouse/gen/typewriter"`,
	}

	return execute(listStandard, customFilename, imports, listBody)
}

func listStandard() error {
	app, err := typewriter.NewApp("+gen")

	if err != nil {
		return err
	}

	fmt.Println("Installed typewriters:")
	for _, tw := range app.TypeWriters {
		fmt.Println("  " + tw.Name())
	}

	return nil
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
