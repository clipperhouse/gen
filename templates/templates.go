// Package templates includes the types and functionality for including templates in the gen package
// The pattern is similar to that used in the image and database/sql packages, where the main templates package
// is imported, and specific template packages are included as image formats or sql drivers would be.

package templates

import (
	"errors"
	"fmt"
	"go/ast"
	"sort"
	"text/template"
)

var templateSets = make(map[string]TemplateSet)

// Register allows template packages to make themselves known to a 'parent' package, usually in the init() func.
// Comparable to the approach taken by builtin image pacakge for registration of image types (eg image/png)
// Your program will do something like:
//	import (
//		"github.com/clipperhouse/gen/templates"
//		_ "github.com/clipperhouse/gen/templates/projection"
//	)
func Register(name string, ts TemplateSet) {
	templateSets[name] = ts
}

// GetTemplateSet attempts to a template set from the registered templates sets, by name
// Returns error if not found
func GetTemplateSet(name string) (TemplateSet, error) {
	var err error
	ts, ok := templateSets[name]
	if !ok {
		err = errors.New(fmt.Sprintf("%s is not a known TemplateSet", name))
	}
	return ts, err
}

// Template includes the text of a tempalte as well as requirements for the types to which it can be applied
type Template struct {
	Text            string
	RequiresNumeric bool
	// A comparable type is one that supports the == operator. Map keys must be comparable, for example.
	RequiresComparable bool
	// An ordered type is one where greater-than and less-than are supported
	RequiresOrdered bool
}

// TemplateSet is a map of string names to Template
type TemplateSet map[string]*Template

// Contains returns true if the TemplateSet includes a template of a given name
func (ts TemplateSet) Contains(name string) bool {
	_, ok := ts[name]
	return ok
}

// Get attempts to 1) locate a tempalte of that name and 2) parse the template
// Returns an error if the template is not found, and panics if the template can not be parsed (per text/template.Must)
func (ts TemplateSet) Get(name string) (t *template.Template, err error) {
	if !ts.Contains(name) {
		err = errors.New(fmt.Sprintf("%s is not a known template", name))
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
