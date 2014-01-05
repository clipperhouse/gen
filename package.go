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

const (
	tagPattern        = `([\p{L}\p{N},]+)`
	getTagPattern     = `gen:"` + tagPattern + `"`
	projectTagPattern = `project:"` + tagPattern + `"`
)

type Package struct {
	p                *types.Package    // only intended for internal use with Eval() below
	TypeNamesAndDocs map[string]string // docs keyed by type name
}

func (p *Package) GetType(s string) (result *Type, err error) {
	ts := typeString(s)

	doc, found := p.TypeNamesAndDocs[ts.Name()]

	if !found {
		err = errors.New(fmt.Sprintf("%s is not a known type in the current directory", ts))
	}

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

	result = &Type{ts.Pointer(), ts.Package(), ts.Name(), subsettedMethods, projectedTypes}
	return
}

func (p *Package) Eval(s string) (typ types.Type, err error) {
	if p.p == nil {
		err = errors.New(fmt.Sprintf("unable to evaluate type %s", s))
		return
	}

	ts := typeString(s)
	typ, _, err = types.Eval(ts.LocalName(), p.p, p.p.Scope())
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
