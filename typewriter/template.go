package typewriter

import (
	"fmt"
	"go/ast"
	"sort"
	"text/template"
)

// Template includes the text of a template as well as requirements for the types to which it can be applied.
type Template struct {
	Text            string
	RequiresNumeric bool
	// A comparable type is one that supports the == operator. Map keys must be comparable, for example.
	RequiresComparable bool
	// An ordered type is one where greater-than and less-than are supported
	RequiresOrdered bool
	// Indicates that this template requires an exact number of type parameters. Default is zero.
	RequiresTypeParameters int
}

func (tmpl Template) ApplicableToType(t Type) bool {
	return (!tmpl.RequiresComparable || t.Comparable()) && (!tmpl.RequiresNumeric || t.Numeric()) && (!tmpl.RequiresOrdered || t.Ordered())
}

func (tmpl Template) ApplicableToValue(v TagValue) bool {
	return tmpl.RequiresTypeParameters == len(v.TypeParameters)
}

// TemplateSet is a map of string names to Template.
type TemplateSet map[string]*Template

// Contains returns true if the TemplateSet includes a template of a given name.
func (ts TemplateSet) Contains(name string) bool {
	_, ok := ts[name]
	return ok
}

// Get attempts to 1) locate a tempalte of that name and 2) parse the template
// Returns an error if the template is not found, and panics if the template can not be parsed (per text/template.Must)
func (ts TemplateSet) Get(name string) (t *template.Template, err error) {
	if !ts.Contains(name) {
		err = fmt.Errorf("%s is not a known template", name)
		return
	}
	t = template.Must(template.New(name).Parse(ts[name].Text))
	return
}

// GetAllKeys returns a slice of all 'exported' key names of templates in the TemplateSet
func (ts TemplateSet) GetAllKeys() (result []string) {
	for k := range ts {
		if ast.IsExported(k) {
			result = append(result, k)
		}
	}
	sort.Strings(result)
	return
}
