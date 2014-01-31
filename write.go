package main

import (
	"bytes"
	"fmt"
	"go/format"
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

func formatToBytes(b *bytes.Buffer) ([]byte, error) {
	byts := b.Bytes()
	formatted, err := format.Source(byts)
	if err != nil {
		return byts, err
	}
	return formatted, nil
}

func writeFiles(packages []*Package, opts options) {
	for _, p := range packages {
		for _, t := range p.Types {
			file, err := os.Create(t.FileName())
			if err != nil {
				panic(err)
			}
			defer file.Close()

			b := bytes.NewBufferString("")
			writeType(b, t, opts)

			byts, err := formatToBytes(b)

			if err == nil || opts.Force {
				file.Write(byts)
			} else {
				panic(err)
			}

			fmt.Printf("  generated %s, yay!\n", t.Plural())
		}
	}
}
