package main

import (
	"fmt"
	"io"
	"os"
)

func writeType(w io.Writer, t *Type, opts options) {
	h := getHeaderTemplate()
	h.Execute(w, t)

	for _, m := range t.StandardMethods {
		tmpl, err := getStandardTemplate(m)
		if err == nil {
			err := tmpl.Execute(w, t)
			if err != nil {
				panic(err)
			}
		} else if opts.Force {
			fmt.Printf("  skipping %v method\n", m)
		} else {
			panic(err)
		}
	}

	for _, f := range t.Projections {
		tmpl, err := getProjectionTemplate(f.Method)
		if err == nil {
			err := tmpl.Execute(w, f)
			if err != nil {
				panic(err)
			}
		} else if opts.Force {
			fmt.Printf("  skipping %v projection method\n", f.Method)
		} else {
			panic(err)
		}
	}

	if t.requiresSortSupport() {
		s := getSortSupportTemplate()
		err := s.Execute(w, t)
		if err != nil {
			panic(err)
		}
	}
}

func writeFiles(packages []*Package, opts options) {
	for _, p := range packages {
		for _, t := range p.Types {
			file, err := os.Create(t.FileName())
			if err != nil {
				panic(err)
			}
			defer file.Close()

			writeType(file, t, opts)

			fmt.Printf("  generated %s, yay!\n", t.Plural())
		}
	}
}
