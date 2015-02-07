package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/clipperhouse/typewriter"
)

func run(c config) error {
	imports := typewriter.NewImportSpecSet(
		typewriter.ImportSpec{Path: "fmt"},
		typewriter.ImportSpec{Path: "os"},
		typewriter.ImportSpec{Path: "regexp"},
		typewriter.ImportSpec{Path: "github.com/clipperhouse/typewriter"},
	)

	return execute(runStandard, c, imports, runTmpl)
}

func runStandard(c config) (err error) {
	app, err := c.Config.NewApp("+gen")

	if err != nil {
		return err
	}

	if len(app.Packages) == 0 {
		return fmt.Errorf("No packages were found. See http://clipperhouse.github.io/gen to get started, or type %s help.", os.Args[0])
	}

	found := false

	for _, p := range app.Packages {
		found = found || len(p.Types) > 0
	}

	if !found {
		return fmt.Errorf("No types marked with +gen were found. See http://clipperhouse.github.io/gen to get started, or type %s help.", os.Args[0])
	}

	if len(app.TypeWriters) == 0 {
		return fmt.Errorf("No typewriters were imported. See http://clipperhouse.github.io/gen to get started, or type %s help.", os.Args[0])
	}

	if _, err := app.WriteAll(); err != nil {
		return err
	}

	return nil
}

var runTmpl = template.Must(template.New("run").Parse(`

var exitStatusMsg = regexp.MustCompile("^exit status \\d+$")

func main() {
	var err error

	defer func() {
		if err != nil {
			if !exitStatusMsg.MatchString(err.Error()) {
				os.Stderr.WriteString(err.Error() + "\n")
			}
			os.Exit(1)
		}
	}()

	err = run()
}

func run() error {
	config := {{ printf "%#v" .Config }}
	app, err := config.NewApp("+gen")

	if err != nil {
		return err
	}

	if len(app.Packages) == 0 {
		return fmt.Errorf("No packages were found. See http://clipperhouse.github.io/gen to get started, or type %s help.", os.Args[0])
	}

	found := false

	for _, p := range app.Packages {
		found = found || len(p.Types) > 0
	}

	if !found {
		return fmt.Errorf("No types marked with +gen were found. See http://clipperhouse.github.io/gen to get started, or type %s help.", os.Args[0])
	}

	if len(app.TypeWriters) == 0 {
		return fmt.Errorf("No typewriters were imported. See http://clipperhouse.github.io/gen to get started, or type %s help.", os.Args[0])
	}

	if _, err := app.WriteAll(); err != nil {
		return err
	}

	return nil
}
`))
