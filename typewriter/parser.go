package typewriter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"

	_ "code.google.com/p/go.tools/go/gcimporter"
	"code.google.com/p/go.tools/go/types"
)

func getTypes(directive string, filter func(os.FileInfo) bool) ([]Type, error) {
	typs := make([]Type, 0)

	fset := token.NewFileSet()
	rootDir := "./"

	astPackages, astErr := parser.ParseDir(fset, rootDir, filter, parser.ParseComments)

	if astErr != nil {
		return typs, astErr
	}

	for name, astPackage := range astPackages {

		// collect type nodes
		var decls []*ast.GenDecl

		ast.Inspect(astPackage, func(n ast.Node) bool {
			// is it a type?
			// http://golang.org/pkg/go/ast/#GenDecl
			if d, ok := n.(*ast.GenDecl); ok && d.Tok == token.TYPE {
				decls = append(decls, d)

				// no need to keep walking, we don't care about TypeSpec's children
				return false
			}
			return true
		})

		astFiles, astErr := getAstFiles(astPackage, rootDir)

		if astErr != nil {
			return typs, astErr
		}

		conf := types.Config{}
		typesPkg, typesErr := conf.Check(name, fset, astFiles, nil)

		if typesErr != nil {
			fmt.Println("typecheck error:", typesErr, "\nattempting to continue...")
		}

		pkg := &Package{typesPkg}

		for _, decl := range decls {
			if decl.Lparen == 0 {
				// not parenthesized, copy GenDecl.Doc into TypeSpec.Doc
				decl.Specs[0].(*ast.TypeSpec).Doc = decl.Doc
			}
			for _, gspec := range decl.Specs {
				spec := gspec.(*ast.TypeSpec)

				pointer, tags, found, err := parseTags(directive, spec.Doc.Text())

				if err != nil {
					return typs, err
				}

				if !found {
					continue
				}

				typ := Type{
					Package: pkg,
					Pointer: pointer,
					Name:    spec.Name.Name,
					Tags:    tags,
				}

				t, _, err := types.Eval(typ.LocalName(), typesPkg, typesPkg.Scope())

				known := err == nil

				if !known {
					// really shouldn't happen, since the type came from the ast in the first place
					err = fmt.Errorf("failed to evaluate type %s (%s)", typ.Name, err)
					return typs, err
				}

				typ.comparable = isComparable(t)
				typ.numeric = isNumeric(t)
				typ.ordered = isOrdered(t)
				typ.test = test(strings.HasSuffix(fset.Position(spec.Pos()).Filename, "_test.go"))
				typ.Type = t

				typs = append(typs, typ)
			}
		}
	}

	return typs, nil
}

func getAstFiles(p *ast.Package, rootDir string) (result []*ast.File, err error) {
	// pull map of *ast.File into a slice
	for _, f := range p.Files {
		result = append(result, f)
	}
	return
}

// something resembling legal identifiers in Go: http://golang.org/ref/spec#Identifiers
// TODO: should probably allow underscore
var tagreg = regexp.MustCompile(`(\p{L}[\p{L}\p{N}]*):"([^\"]+?)"`)

// identifies gen-marked types and parses tags
func parseTags(directive string, doc string) (pointer Pointer, tags []Tag, found bool, err error) {
	lines := strings.Split(doc, "\n")
	for _, line := range lines {
		original := line

		// strategy is to remove meaningful tokens as they are found
		// kind of a hack, a real parser someday

		// does the line start with the directive?
		if line = strings.TrimLeft(line, "/ "); !strings.HasPrefix(line, directive) {
			continue
		}

		// remove the directive from the line
		line = strings.TrimPrefix(line, directive)

		// next character needs to be a space or end of string
		if !(len(line) == 0 || strings.HasPrefix(line, " ")) {
			err = fmt.Errorf(`the directive %s needs to be followed by a space or end of line, see source containing
%s
`, directive, original)
			return
		}

		// ok, we got something
		found = true

		// get rid of leading spaces
		line = strings.TrimLeft(line, " ")

		// is the next character a pointer?
		p := Pointer(true).String()
		if strings.HasPrefix(line, p) {
			pointer = true
			line = strings.TrimPrefix(line, p)

			// if found, next character needs to be a space or end of string
			if !(len(line) == 0 || strings.HasPrefix(line, " ")) {
				err = fmt.Errorf(`the pointer needs to be followed by a space or end of line, see source containing
%s
`, original)
				return
			}
		}

		// find all matches of tag pattern
		matches := tagreg.FindAllString(line, -1)

		for _, m := range matches {
			if tag, found := parseTag(m); found {
				// should always be found since the matches are selected
				// and substringed (substrung?) using the same regex

				// add to the results
				tags = append(tags, tag)

				// remove the tag from the parsed line
				line = strings.Replace(line, m, "", -1)
			}
		}

		// trim spaces
		line = strings.Trim(line, " ")

		// anything remaining is invalid

		if len(line) > 0 {
			err = fmt.Errorf(`unknown syntax '%s', see
%s
`, line, original)
			return
		}
	}
	return
}

func parseTag(s string) (tag Tag, found bool) {
	var matches []string
	if matches = tagreg.FindStringSubmatch(s); matches == nil || len(matches) == 0 {
		// not a match? not an error, just not a tag
		return
	}

	found = true

	var name string
	var items []string
	var negated bool

	name = matches[1]

	splitter := regexp.MustCompile(`[, ]+`)
	if match := matches[2]; len(match) > 0 {
		index := 0
		if negated = strings.HasPrefix(match, "-"); negated {
			index = 1
		}
		items = splitter.Split(match[index:], -1)
	}

	tag = Tag{
		Name:     name,
		Items:    items,
		Negated:  negated,
		original: s,
	}

	return
}
