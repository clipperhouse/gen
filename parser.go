package typewriter

import (
	_ "code.google.com/p/go.tools/go/gcimporter"
	"code.google.com/p/go.tools/go/types"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

func getTypes(directive string) ([]Type, error) {
	typs := make([]Type, 0)

	fset := token.NewFileSet()
	rootDir := "./"

	astPackages, astErr := parser.ParseDir(fset, rootDir, nil, parser.ParseComments)

	if astErr != nil {
		return typs, astErr
	}

	for name, astPackage := range astPackages {
		astFiles, astErr := getAstFiles(astPackage, rootDir)

		if astErr != nil {
			return typs, astErr
		}

		typesPkg, typesErr := types.Check(name, fset, astFiles)

		if typesErr != nil {
			return typs, typesErr
		}

		pkg := &Package{typesPkg}

		// doc package is handy for pulling types and their comments
		docPkg := doc.New(astPackage, name, doc.AllDecls)

		for _, docType := range docPkg.Types {

			pointer, tags, found := parseTags(directive, docType.Doc)

			if !found {
				continue
			}

			typ := Type{
				Package: pkg,
				Pointer: pointer,
				Name:    docType.Name,
				Tags:    tags,
			}

			t, _, err := types.Eval(typ.LocalName(), typesPkg, typesPkg.Scope())
			known := err == nil

			if !known {
				err = errors.New(fmt.Sprintf("failed to evaluate type %s (%s)", typ.Name, err))
				continue
			}

			typ.comparable = isComparable(t)
			typ.numeric = isNumeric(t)
			typ.ordered = isOrdered(t)

			typs = append(typs, typ)
		}
	}

	return typs, nil
}

func getAstFiles(p *ast.Package, rootDir string) (result []*ast.File, err error) {
	// pull map of *ast.File into a slice
	// and skip files who's out of compile scope (Conditional compile, for example)
	for name, f := range p.Files {
		if ok, buildErr := build.Default.MatchFile(rootDir, name); err != nil {
			err = buildErr
		} else if ok {
			result = append(result, f)
		}
	}
	return
}

const deprecationUrl = `http://clipperhouse.github.io/gen/#Changelog`

func checkDeprecatedTags(t types.Type) bool {
	// give informative errors for use of deprecated custom methods
	switch x := t.Underlying().(type) {
	case *types.Struct:
		for i := 0; i < x.NumFields(); i++ {
			_, found := parseTag(x.Tag(i))
			if found {
				return false
			}
		}
	}
	return true
}

// identifies gen-marked types and parses tags
func parseTags(directive string, doc string) (pointer Pointer, tags []Tag, found bool) {
	lines := strings.Split(doc, "\n")
	for _, line := range lines {
		if line = strings.TrimLeft(line, "/ "); !strings.HasPrefix(line, directive) {
			continue
		}

		found = true

		// parse out tags & pointer
		spaces := regexp.MustCompile(" +")
		parts := spaces.Split(line, -1)

		for _, s := range parts {
			if s == "*" {
				pointer = true
				continue
			}
			if tag, found := parseTag(s); found {
				tags = append(tags, tag)
				continue
			}
		}
	}
	return
}

func parseTag(s string) (tag Tag, found bool) {
	// same as legal identifiers in Go: http://golang.org/ref/spec#Identifiers
	r := regexp.MustCompile(`(\p{L}[\p{L}\p{N}]*):"(.*)"`)

	var matches []string

	if matches = r.FindStringSubmatch(s); matches == nil || len(matches) == 0 {
		return
	}

	found = true

	var name string
	var items []string
	var negated bool

	name = matches[1]

	if match := matches[2]; len(match) > 0 {
		index := 0
		if negated = strings.HasPrefix(match, "-"); negated {
			index = 1
		}
		items = strings.Split(match[index:], ",")
	}

	tag = Tag{
		Name:     name,
		Items:    items,
		Negated:  negated,
		original: s,
	}

	return
}
