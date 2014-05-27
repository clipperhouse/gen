package genwriter

import (
	"fmt"
	"github.com/clipperhouse/inflect"
	"github.com/clipperhouse/typewriter"
	"io"
)

func init() {
	err := typewriter.Register(NewGenWriter())
	if err != nil {
		panic(err)
	}
}

type GenWriter struct {
	// Type is not comparable, use .String() as keys instead
	models    map[string]model
	validated map[string]bool
}

func NewGenWriter() GenWriter {
	return GenWriter{
		models:    make(map[string]model),
		validated: make(map[string]bool),
	}
}

// A convenience struct for passing data to templates.
type model struct {
	typewriter.Type
	methods     []string
	projections []Projection
}

func (m model) Plural() (result string) {
	result = inflect.Pluralize(m.Name)
	if result == m.Name {
		result += "s"
	}
	return
}

// genwriter prepares models for later use in the .Validate() method. It must be called prior.
func (g GenWriter) ensureValidation(t typewriter.Type) error {
	if !g.validated[t.String()] {
		return fmt.Errorf("Type '%s' has not been previously validated. TypeWriter.Validate() must be called on all types before using them in subsequent methods.", t.String())
	}

	return nil
}

func (g GenWriter) Name() string {
	return "genwriter"
}

func (g GenWriter) Validate(t typewriter.Type) (bool, error) {
	standardMethods, projectionMethods, err := evaluateTags(t.Tags)
	if err != nil {
		return false, err
	}

	// filter methods applicable to type
	for _, s := range standardMethods {
		tmpl, ok := standardTemplates[s]

		if !ok {
			err = fmt.Errorf("unknown method %v", s)
			return false, err
		}

		if !tmpl.ApplicableTo(t) {
			standardMethods = remove(standardMethods, s)
		}
	}

	projectionTag, found, err := t.Tags.ByName("projections")

	if err != nil {
		return false, err
	}

	m := model{
		Type:    t,
		methods: standardMethods,
	}

	if found {
		for _, s := range projectionTag.Items {
			projectionType, err := t.Package.Eval(s)

			if err != nil {
				return false, fmt.Errorf("unable to identify type %s, projected on %s (%s)", s, t.Name, err)
			}

			for _, pm := range projectionMethods {
				tmpl, ok := projectionTemplates[pm]

				if !ok {
					return false, fmt.Errorf("unknown projection method %v", pm)
				}

				if tmpl.ApplicableTo(projectionType) {
					m.projections = append(m.projections, Projection{
						Method: pm,
						Type:   s,
						Parent: &m,
					})
				}
			}
		}
	}

	g.validated[t.String()] = true

	// only add to models if we are going to do something with it
	if len(m.methods) > 0 || len(m.projections) > 0 {
		g.models[t.String()] = m
	}

	return true, nil
}

func (g GenWriter) WriteHeader(w io.Writer, t typewriter.Type) {
	err := g.ensureValidation(t)

	if err != nil {
		panic(err)
	}

	m, exists := g.models[t.String()]

	if !exists {
		return
	}

	if includeSortSupport(m.methods) {
		s := `// Sort implementation is a modification of http://golang.org/pkg/sort/#Sort
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found at http://golang.org/LICENSE.
`
		w.Write([]byte(s))
	}

	return
}

func (g GenWriter) Imports(t typewriter.Type) (result []string) {
	err := g.ensureValidation(t)

	if err != nil {
		panic(err)
	}

	m, exists := g.models[t.String()]

	if !exists {
		return
	}

	imports := make(map[string]bool)

	methodRequiresErrors := map[string]bool{
		"First":   true,
		"Single":  true,
		"Max":     true,
		"Min":     true,
		"MaxBy":   true,
		"MinBy":   true,
		"Average": true,
	}

	methodRequiresSort := map[string]bool{
		"Sort": true,
	}

	for _, s := range m.methods {
		if methodRequiresErrors[s] {
			imports["errors"] = true
		}
		if methodRequiresSort[s] {
			imports["sort"] = true
		}
	}

	for _, p := range m.projections {
		if methodRequiresErrors[p.Method] {
			imports["errors"] = true
		}
		if methodRequiresSort[p.Method] {
			imports["sort"] = true
		}
	}

	for s := range imports {
		result = append(result, s)
	}

	return
}

func (g GenWriter) Write(w io.Writer, t typewriter.Type) {
	err := g.ensureValidation(t)

	if err != nil {
		panic(err)
	}

	m, exists := g.models[t.String()]

	if !exists {
		return
	}

	tmpl, _ := standardTemplates.Get("plural")
	if err := tmpl.Execute(w, m); err != nil {
		panic(err)
	}

	for _, s := range m.methods {
		tmpl, _ := standardTemplates.Get(s) // already validated above
		err := tmpl.Execute(w, m)
		if err != nil {
			panic(err)
		}
	}

	for _, p := range m.projections {
		tmpl, _ := projectionTemplates.Get(p.Method) // already validated above
		err := tmpl.Execute(w, p)
		if err != nil {
			panic(err)
		}
	}

	if includeSortInterface(m.methods) {
		tmpl, _ := standardTemplates.Get("sortInterface") // already validated above
		err := tmpl.Execute(w, m)
		if err != nil {
			panic(err)
		}
	}

	if includeSortSupport(m.methods) {
		tmpl, _ := standardTemplates.Get("sortSupport") // already validated above
		err := tmpl.Execute(w, m)
		if err != nil {
			panic(err)
		}
	}
}
