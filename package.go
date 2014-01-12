package main

import (
	_ "code.google.com/p/go.tools/go/gcimporter"
	"code.google.com/p/go.tools/go/types"
	"errors"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

var (
	genTag     = regexp.MustCompile(`gen:"(.+)"`)
	projectTag = regexp.MustCompile(`project:"(.+)"`)
)

type Package struct {
	p                *types.Package    // only intended for internal use with Eval() below
	TypeNamesAndDocs map[string]string // docs keyed by type name
}

func (p *Package) GetType(t *typeArg) (result *Type, err error) {
	doc, found := p.TypeNamesAndDocs[t.Name]

	if !found {
		err = errors.New(fmt.Sprintf("%s is not a known type in the current directory", t))
		return
	}

	var subsettedMethods []string

	if matches := genTag.FindStringSubmatch(doc); matches != nil && len(matches) > 1 {
		subsettedMethods = strings.Split(matches[1], ",")
	}

	var projectedTypes []string

	if matches := projectTag.FindStringSubmatch(doc); matches != nil && len(matches) > 1 {
		projectedTypes = strings.Split(matches[1], ",")
	}

	result = &Type{t, subsettedMethods, projectedTypes}
	return
}

func (p *Package) Eval(s string) (typ types.Type, err error) {
	scope := types.Universe
	if p.p != nil {
		scope = p.p.Scope()
	}

	typ, _, err = types.Eval(s, p.p, scope) // calling with nil p.p will assume Universe scope but being defensive
	return typ, err
}

// Returns one gen Package per Go package found in current directory, keyed by name
func getPackages() (result map[string]*Package) {
	fset := token.NewFileSet()
	dir, err := parser.ParseDir(fset, "./", nil, parser.ParseComments)
	if err != nil {
		errs = append(errs, err)
	}

	result = make(map[string]*Package)

	for k, v := range dir {
		files := make([]*ast.File, 0)
		for _, f := range v.Files {
			files = append(files, f)
		}

		p, err := types.Check(k, fset, files)
		if err != nil {
			errs = append(errs, err)
		}

		d := doc.New(v, k, doc.AllDecls)
		typeDocs := make(map[string]string)
		for _, t := range d.Types {
			typeDocs[t.Name] = t.Doc
		}

		result[k] = &Package{p, typeDocs}
	}

	return
}
