package templates

import (
	"errors"
	"fmt"
	"go/ast"
	"sort"
	"text/template"
)

var templateSets = make(map[string]TemplateSet)

func Register(name string, ts TemplateSet) {
	templateSets[name] = ts
}

func GetTemplateSet(name string) (TemplateSet, error) {
	var err error
	ts, ok := templateSets[name]
	if !ok {
		err = errors.New(fmt.Sprintf("%s is not a known TemplateSet", name))
	}
	return ts, err
}

type Template struct {
	Text               string
	RequiresNumeric    bool
	RequiresComparable bool
	RequiresOrdered    bool
}

type TemplateSet map[string]*Template

func (ts TemplateSet) Contains(name string) bool {
	_, ok := ts[name]
	return ok
}

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
