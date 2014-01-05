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

type typeChecker struct {
	p        *types.Package
	typeDocs map[string]string // docs keyed by type name
}

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

func (t *typeChecker) eval(s string) (typ types.Type, err error) {
	if t.p == nil {
		err = errors.New(fmt.Sprintf("unable to evaluate type %s", s))
		return
	}

	typ, _, err = types.Eval(s, t.p, t.p.Scope())
	return typ, err
}

// Returns one type checker per package found in current directory
func getTypeCheckers() (result map[string]*typeChecker) {
	fset := token.NewFileSet()
	dir, err := parser.ParseDir(fset, "./", nil, parser.ParseComments)
	if err != nil {
		errs = append(errs, err)
	}

	result = make(map[string]*typeChecker)

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

		result[k] = &typeChecker{p, typeDocs}
	}

	return
}
