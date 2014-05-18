package genwriter

import (
	"fmt"
	"github.com/clipperhouse/inflect"
	"github.com/clipperhouse/typewriter"
	"io"
)

func init() {
	err := typewriter.Register(GenWriter{})
	if err != nil {
		panic(err)
	}
}

type GenWriter struct{}

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

// Type is not comparable, use .String() as key instead
var models = make(map[string]model)
var validated = make(map[string]bool)

// genwriter prepares models for later use in the .Validate() method. It must be called prior.
func ensureValidation(t typewriter.Type) {
	if !validated[t.String()] {
		panic(fmt.Errorf("Type %s has not been previously validated. typewriter.Validate() must be called on all types before using them in subsequent methods.", t.String()))
	}
}

func (s GenWriter) Name() string {
	return "genwriter"
}

func (s GenWriter) Validate(t typewriter.Type) (bool, error) {
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

	models[t.String()] = m
	validated[t.String()] = true

	return len(m.methods) > 0 || len(m.projections) > 0, nil
}

func (s GenWriter) WriteHeader(w io.Writer, t typewriter.Type) {
	ensureValidation(t)

	//TODO: add licenses
	return
}

func (s GenWriter) Imports(t typewriter.Type) (result []string) {
	ensureValidation(t)

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

	m := models[t.String()]

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

func (s GenWriter) Write(w io.Writer, t typewriter.Type) {
	ensureValidation(t)

	m := models[t.String()]

	tmpl, _ := standardTemplates.Get("plural")
	err := tmpl.Execute(w, m)
	if err != nil {
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
