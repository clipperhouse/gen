package genwriter

import (
	"fmt"
	"github.com/clipperhouse/inflect"
	"github.com/clipperhouse/typewriter"
	"io"
	"regexp"
	"strings"
)

func init() {
	typewriter.Register(GenWriter{})
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
	return "gen"
}

func (s GenWriter) Validate(t typewriter.Type) (bool, error) {
	standardMethods, projectionMethods, err := determineMethods(t.Tags)
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

// This business exists because I overload the methods tag to specify both standard and projection methods.
// Kind of a mess, but for the end user, arguably simpler. And arguably not.
func determineMethods(tags typewriter.Tags) (standardMethods, projectionMethods []string, err error) {
	var nilMethods, nilProjections bool

	methods, found, methodsErr := tags.ByName("methods")

	if methodsErr != nil {
		err = methodsErr
		return
	}

	nilMethods = !found // non-existent methods tag is different than empty

	_, found, projectionsErr := tags.ByName("projections")

	if projectionsErr != nil {
		err = projectionsErr
		return
	}

	nilProjections = !found

	if nilMethods || methods.Negated {
		// default to all
		standardMethods = standardTemplates.GetAllKeys()
		if !nilProjections {
			projectionMethods = projectionTemplates.GetAllKeys()
		}
	}

	if !nilMethods {
		// categorize subsetted methods as standard or projection
		std := make([]string, 0)
		prj := make([]string, 0)

		for _, m := range methods.Items {
			isStd := standardTemplates.Contains(m)
			if isStd {
				std = append(std, m)
			}

			// only consider projection methods in presence of projected types
			isPrj := !nilProjections && projectionTemplates.Contains(m)
			if isPrj {
				prj = append(prj, m)
			}

			if !isStd && !isPrj {
				err = fmt.Errorf("method %s is unknown", m)
				return
			}
		}

		if methods.Negated {
			standardMethods = remove(standardMethods, std...)
			projectionMethods = remove(projectionMethods, prj...)
		} else {
			standardMethods = std
			projectionMethods = prj
		}
	}

	return
}

func includeSortSupport(standardMethods []string) bool {
	for _, m := range standardMethods {
		if strings.HasPrefix(m, "SortBy") {
			return true
		}
	}
	return false
}

func includeSortInterface(standardMethods []string) bool {
	reg := regexp.MustCompile(`^Sort(Desc)?$`)
	for _, m := range standardMethods {
		if reg.MatchString(m) {
			return true
		}
	}
	return false
}
