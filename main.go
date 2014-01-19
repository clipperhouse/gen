package main

import (
	"errors"
	"fmt"
	"os"
)

var errs = make([]error, 0)

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

	writeFile(packages, opts)
}

func addError(text string) {
	errs = append(errs, errors.New(text))
}

func writeFile(packages []*Package, opts options) {
	for _, p := range packages {
		for _, t := range p.Types {
			file, err := os.Create(t.FileName())
			if err != nil {
				panic(err)
			}
			defer file.Close()

			h := getHeaderTemplate()
			h.Execute(file, t)

			for _, m := range t.StandardMethods {
				tmpl, err := getStandardTemplate(m)
				if err == nil {
					tmpl.Execute(file, t)
				} else if opts.Force {
					fmt.Printf("  skipping %v method\n", m)
				} else {
					panic(err) // shouldn't get here, should have been caught in getSubsettedMethods
				}
			}

			for _, f := range t.Projections {
				tmpl, err := getProjectionTemplate(f.Method)
				if err == nil {
					tmpl.Execute(file, f)
				} else if opts.Force {
					fmt.Printf("  skipping %v projection method\n", f.Method)
				} else {
					panic(err) // shouldn't get here, should have been caught in getProjectionSpecs
				}
			}

			if t.requiresSortSupport() {
				s := getSortSupportTemplate()
				s.Execute(file, t)
			}

			fmt.Printf("  generated %s, yay!\n", t.Plural())
		}

	}
}

const usage = `Documentation is available at http://clipperhouse.github.io/gen

Usage: gen [-force]
  -force    # forces generation to continue despite errors; voids warranty; shortcut -f
`
