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

const usage = `Documentation is available at http://clipperhouse.github.io/gen

Usage: gen [-force]
  -force    # forces generation to continue despite errors; voids warranty; shortcut -f
`

const deprecationUrl = `https://github.com/clipperhouse/gen/blob/projection/CHANGELOG.md`
