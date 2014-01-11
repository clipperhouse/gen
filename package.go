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
	projectTagPattern = `project:"(.+)"`
)

type Package struct {
	p                *types.Package    // only intended for internal use with Eval() below
	TypeNamesAndDocs map[string]string // docs keyed by type name
}

func (p *Package) GetType(t *typeArg) (result *Type, errs []error) {
	doc, found := p.TypeNamesAndDocs[t.Name]

	if !found {
		errs = append(errs, errors.New(fmt.Sprintf("%s is not a known type in the current directory", t)))
	}

	var subsettedMethods []string
	genTag := regexp.MustCompile(getTagPattern)
	genMatch := genTag.FindStringSubmatch(doc)
	if genMatch != nil && len(genMatch) > 1 {
		subsettedMethods = strings.Split(genMatch[1], ",")
		d := findDuplicates(subsettedMethods)
		if len(d) > 0 {
			errs = append(errs, errors.New(fmt.Sprintf("duplicate subsetted method(s) found on type %s: %v", t, d)))
		}
	}

	var projectedTypes []string
	projectTag := regexp.MustCompile(projectTagPattern)
	projectMatch := projectTag.FindStringSubmatch(doc)
	if projectMatch != nil && len(projectMatch) > 1 {
		projectedTypes = strings.Split(projectMatch[1], ",")
		d := findDuplicates(projectedTypes)
		if len(d) > 0 {
			errs = append(errs, errors.New(fmt.Sprintf("duplicate projected type(s) found on type %s: %v", t, d)))
		}
	}

	result = &Type{t, subsettedMethods, projectedTypes}
	return
}

func (p *Package) Eval(s string) (typ types.Type, err error) {
	if p.p == nil {
		typ, _, err = types.Eval(s, nil, types.Universe)
		if err != nil {
			err = errors.New(fmt.Sprintf("unable to evaluate type %s", s))
		}
		return
	}

	typ, _, err = types.Eval(s, p.p, p.p.Scope())
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

func findDuplicates(a []string) (result []string) {
	found := make(map[string]bool)

	for _, s := range a {
		if found[s] {
			result = append(result, s)
		}
		found[s] = true
	}
	return
}
