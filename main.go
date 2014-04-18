package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/clipperhouse/gen/templates"
	_ "github.com/clipperhouse/gen/templates/container"
	_ "github.com/clipperhouse/gen/templates/projection"
	_ "github.com/clipperhouse/gen/templates/standard"
)

var errs = make([]error, 0)
var standardTemplates, projectionTemplates, containerTemplates templates.TemplateSet

func init() {
	if ts, err := templates.Get("standard"); err == nil {
		standardTemplates = ts
	} else {
		errs = append(errs, err)
	}

	if ts, err := templates.Get("projection"); err == nil {
		projectionTemplates = ts
	} else {
		errs = append(errs, err)
	}

	if ts, err := templates.Get("container"); err == nil {
		containerTemplates = ts
	} else {
		errs = append(errs, err)
	}
}

func main() {
	args := os.Args[1:]
	opts, err := parseArgs(args)

	if err != nil {
		fmt.Println(err)
		return // command-line errors are fatal, other errors can be forced
	}

	if opts.Help {
		fmt.Println(usage)
		return
	}

	packages := getPackages()

	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Printf("  error: %v\n", e)
		}
		if opts.Force {
			fmt.Println("  forced...")
		} else {
			fmt.Println("  operation canceled")
			fmt.Println("  use the -f flag if you wish to force generation (i.e., ignore errors)")
			return
		}
	}

	writeFiles(packages, opts)
}

func addError(text string) {
	errs = append(errs, errors.New(text))
}

const usage = `Documentation is available at http://clipperhouse.github.io/gen

Usage: gen [-force]
  -force    # forces generation to continue despite errors; voids warranty; shortcut -f
`

const deprecationUrl = `http://clipperhouse.github.io/gen/#Changelog`
