package main

import (
	_ "code.google.com/p/go.tools/go/gcimporter"
	"code.google.com/p/go.tools/go/types"
	"fmt"
	"github.com/clipperhouse/gen/inflect"
	"strings"
	"time"
)

type options struct {
	All          bool
	AllPointer   string
	ExportedOnly bool
	Force        bool
}

// Utility for handling various string representations of types. Does no validation, ensure that you initialize/convert with something like *package.Type. Pointer and package are optional.
type typeString string

func (t typeString) Pointer() string {
	if strings.HasPrefix(string(t), "*") {
		return "*"
	}
	return ""
}

func (t typeString) Package() string {
	parts := strings.Split(string(t), ".")
	if len(parts) > 1 {
		return parts[0]
	}
	return ""
}

func (t typeString) Name() string {
	s := string(t)
	parts := strings.Split(s, ".")
	if len(parts) > 1 {
		return parts[1]
	}
	if strings.HasPrefix(s, "*") {
		return s[1:]
	}
	return s
}

// "Package-local" name, includes pointer but ignores package.
func (t typeString) LocalName() string {
	return t.Pointer() + t.Name()
}

type typeChecker struct {
	p        *types.Package
	typeDocs map[string]string // docs keyed by type name
}

func (t *typeChecker) getTypeSpec(s string) (typeSpec, error) {
	typ, _, err := types.Eval(s, t.p, t.p.Scope())

	if err != nil {
		return typeSpec{}, err
	}

	name := typeString(typ.String()).Name()
	result := typeSpec{Type: typ, Doc: t.typeDocs[name]}

	return result, nil // err is returned above
}

type typeSpec struct {
	Type types.Type
	Doc  string
}

type genSpec struct {
	Pointer     string
	Package     string
	Singular    string
	Methods     []string
	Projections []*projection
	Imports     []string
	Plural      string
	Receiver    string
	Loop        string
	Generated   string
	Command     string
	FileName    string
}

func newGenSpec(ptr, pkg, name string) *genSpec {
	plural := inflect.Pluralize(name)
	if plural == name {
		plural += "s"
	}
	return &genSpec{
		Pointer:   ptr,
		Package:   pkg,
		Singular:  name,
		Plural:    inflect.Pluralize(name),
		Receiver:  "rcv",
		Loop:      "v",
		Generated: time.Now().UTC().Format(time.RFC1123),
		Command:   fmt.Sprintf("%s %s%s", "gen", ptr, name),
		FileName:  strings.ToLower(name) + "_gen.go",
	}
}

func (g genSpec) String() string {
	return joinName("", g.Plural) // TODO: kill this
}

type projection struct {
	Method string
	Type   string
	Parent *genSpec
}

func (p projection) String() string {
	return p.Method + strings.Title(strings.Replace(p.Type, "[]", "ArrayOf", -1)) // super hacky, TODO: really handle arbitrarty types
}
