package main

import (
	"github.com/clipperhouse/gen/inflect"
	"regexp"
	"strings"
)

type Type struct {
	Package         *Package
	Pointer         string
	Name            string
	StandardMethods []string
	Projections     []*Projection
	Containers      []string
	Imports         []string
}

func (t *Type) LocalName() (result string) {
	return t.Pointer + t.Name
}

func (t *Type) Plural() (result string) {
	result = inflect.Pluralize(t.Name)
	if result == t.Name {
		result += "s"
	}
	return
}

func (t *Type) FileName() string {
	return strings.ToLower(t.Name) + "_gen.go"
}

func (t *Type) AddProjection(methodName, typeName string) {
	t.Projections = append(t.Projections, &Projection{methodName, typeName, t})
}

type Projection struct {
	Method string
	Type   string
	Parent *Type
}

func (p *Projection) MethodName() string {
	name := p.Type

	pointer := regexp.MustCompile(`^\**`)
	pointers := len(pointer.FindAllString(name, -1)[0])
	name = strings.Replace(name, "*", "", -1) + strings.Repeat("Pointer", pointers)

	slice := regexp.MustCompile(`(\[\])`)
	slices := len(slice.FindAllString(name, -1))
	name = strings.Replace(name, "[]", "", -1) + strings.Repeat("Slice", slices)

	illegal := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	name = illegal.ReplaceAllString(name, " ")

	name = strings.Title(name)
	name = strings.Replace(name, " ", "", -1)

	return p.Method + strings.Title(name)
}

func (p *Projection) String() string {
	return p.MethodName()
}
