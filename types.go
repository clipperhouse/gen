package main

import (
	_ "code.google.com/p/go.tools/go/gcimporter"
	"code.google.com/p/go.tools/go/types"
	"errors"
	"fmt"
	"github.com/clipperhouse/gen/inflect"
	"regexp"
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

const ptr = "*"

func (t typeString) Pointer() string {
	if strings.HasPrefix(string(t), ptr) {
		return ptr
	}
	return ""
}

func (t typeString) Package() string {
	parts := strings.Split(string(t), ".")
	if len(parts) > 1 {
		s := parts[0]
		if strings.HasPrefix(s, ptr) {
			return s[1:]
		}
		return s
	}
	return ""
}

// name of the type only, no pointer or package
func (t typeString) Name() string {
	s := string(t)
	parts := strings.Split(s, ".")
	if len(parts) > 1 {
		return parts[1]
	}
	if strings.HasPrefix(s, ptr) {
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

func (t *typeChecker) eval(s string) (typ types.Type, err error) {
	if t.p == nil {
		err = errors.New(fmt.Sprintf("unable to evaluate type %s", s))
		return
	}

	typ, _, err = types.Eval(s, t.p, t.p.Scope())
	return typ, err
}

const (
	tagPattern        = `([\p{L}\p{N},]+)`
	getTagPattern     = `gen:"` + tagPattern + `"`
	projectTagPattern = `project:"` + tagPattern + `"`
)

func (t *typeChecker) getTypeSpec(s string) *typeSpec {
	ts := typeString(s)

	doc := t.typeDocs[ts.Name()]

	var subsettedMethods []string
	genTag := regexp.MustCompile(getTagPattern)
	genMatch := genTag.FindStringSubmatch(doc)
	if genMatch != nil && len(genMatch) > 1 {
		subsettedMethods = strings.Split(genMatch[1], ",")
	}

	var projectedTypes []string
	projectTag := regexp.MustCompile(projectTagPattern)
	projectMatch := projectTag.FindStringSubmatch(doc)
	if projectMatch != nil && len(projectMatch) > 1 {
		projectedTypes = strings.Split(projectMatch[1], ",")
	}

	result := &typeSpec{ts.Pointer(), ts.Package(), ts.Name(), subsettedMethods, projectedTypes}

	return result
}

type typeSpec struct {
	Pointer          string
	Package          string
	Name             string
	SubsettedMethods []string
	ProjectedTypes   []string
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
		Command:   fmt.Sprintf("%s %s%s.%s", "gen", ptr, pkg, name),
		FileName:  strings.ToLower(name) + "_gen.go",
	}
}

func (g genSpec) Type() string {
	return g.Pointer + g.Package + "." + g.Singular
}

func (g genSpec) String() string {
	return g.Plural
}

type projection struct {
	Method string
	Type   string
	Parent *genSpec
}

func (p projection) String() string {
	return p.Method + strings.Title(strings.Replace(p.Type, "[]", "ArrayOf", -1)) // super hacky, TODO: really handle arbitrarty types
}
