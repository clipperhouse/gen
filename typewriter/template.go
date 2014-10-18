package typewriter

import (
	"fmt"

	"text/template"
)

// Template includes the text of a template as well as requirements for the types to which it can be applied.
type Template struct {
	Text           string
	TypeConstraint Constraint
	// Indicates both the number of required type parameters, and the constraints of each (if any)
	TypeParameterConstraints []Constraint
}

func (tmpl *Template) tryTypeAndValue(t Type, v TagValue) error {
	if err := tmpl.TypeConstraint.tryType(t); err != nil {
		return fmt.Errorf("cannot implement %s: %s", v, err)
	}

	if len(tmpl.TypeParameterConstraints) != len(v.TypeParameters) {
		return fmt.Errorf("%s requires %d type parameters", v.Name, len(tmpl.TypeParameterConstraints))
	}

	for i := range v.TypeParameters {
		c := tmpl.TypeParameterConstraints[i]
		tp := v.TypeParameters[i]
		if err := c.tryType(tp); err != nil {
			return fmt.Errorf("cannot implement %s on %s: %s", v, t, err)
		}
	}

	return nil
}

// TemplateSet is a map of string names to Template.
type TemplateSet map[string]*Template

// Get attempts to 1) locate a template of that name and 2) parse the template
func (ts TemplateSet) Get(v TagValue) (*template.Template, error) {
	return ts.ByName(v.TemplateKey())
}

// Get attempts to 1) locate a template of that name and 2) parse the template
func (ts TemplateSet) ByName(name string) (*template.Template, error) {
	tmpl, found := ts[name]
	if !found {
		err := fmt.Errorf("%s is not a known template", name)
		return nil, err
	}
	return template.New(name).Parse(tmpl.Text)
}
