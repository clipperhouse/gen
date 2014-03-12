package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func writeType(w io.Writer, t *Type, opts options) {
	if h, err := standardTemplates.Get("header"); err != nil {
		panic(err)
	} else {
		err := h.Execute(w, t)
		if err != nil {
			panic(err)
		}
	}

	for _, m := range t.StandardMethods {
		tmpl, err := standardTemplates.Get(m)
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
		tmpl, err := projectionTemplates.Get(f.Method)
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

	if t.requiresSortInterface() {
		tmpl, err := standardTemplates.Get("sortInterface")
		if err == nil {
			err := tmpl.Execute(w, t)
			if err != nil {
				panic(err)
			}
		} else if opts.Force {
			fmt.Printf("  skipping sortInterface\n")
		} else {
			panic(err)
		}
	}

	if t.requiresSortSupport() {
		tmpl, err := standardTemplates.Get("sortSupport")
		if err == nil {
			err := tmpl.Execute(w, t)
			if err != nil {
				panic(err)
			}
		} else if opts.Force {
			fmt.Printf("  skipping sortSupport\n")
		} else {
			panic(err)
		}
	}

	for _, c := range t.Containers {
		tmpl, err := containerTemplates.Get(c)
		if err == nil {
			err := tmpl.Execute(w, t)
			if err != nil {
				panic(err)
			}
		} else if opts.Force {
			fmt.Printf("  skipping %v container\n", c)
		} else {
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

func getExistingGenFiles() (result map[string]bool) {
	result = make(map[string]bool)

	d, err := ioutil.ReadDir("./")

	if err != nil {
		panic(err)
	}

	for _, f := range d {
		if strings.HasSuffix(f.Name(), "_gen.go") {
			result[f.Name()] = true
		}
	}

	return
}

func promptDeletions(packages []*Package, existing map[string]bool, input io.Reader, output io.Writer) (deletions []string, ok bool) {
	for _, p := range packages {
		for _, t := range p.Types {
			delete(existing, t.FileName())
		}
	}

	if len(existing) > 0 {
		for f := range existing {
			deletions = append(deletions, f)
		}

		fmt.Fprintf(output, "  This will delete previously-generated files: %s\n  Continue [y/n]? ", strings.Join(deletions, ", "))

		var confirm string
		fmt.Fscan(input, &confirm)
		ok = confirm == "y"
	} else {
		ok = true
	}

	return
}

func writeFiles(packages []*Package, opts options) {
	existing := getExistingGenFiles()
	deletions, ok := promptDeletions(packages, existing, os.Stdin, os.Stdout)

	if !ok {
		fmt.Println("  operation cancelled")
		return
	}

	for _, p := range packages {
		for _, t := range p.Types {
			file, err := os.Create(t.FileName())
			if err != nil {
				panic(err)
			}
			defer file.Close()

			var b bytes.Buffer
			writeType(&b, t, opts)
			byts, err := formatToBytes(&b)

			if err == nil || opts.Force {
				file.Write(byts)
			} else {
				panic(err)
			}

			fmt.Printf("  generated %s, yay!\n", t.Plural())
		}
	}

	for _, f := range deletions {
		os.Remove(f)
		fmt.Printf("  deleted %s\n", f)
	}
}
