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

		typesPkg, typesErr := types.Check(name, fset, astFiles)

		if typesErr != nil {
			return typs, typesErr
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

		l := lex(line)

		// top level can be pointer or tags
		for {
			item := l.nextItem()

			if item.typ == itemError {
				err = &SyntaxError{
					msg: item.val,
					Pos: item.pos,
				}
				return
			}

			// pick up pointer if it exists
			if item.typ == itemPointer {
				// have we already seen a pointer?
				if pointer {
					err = &SyntaxError{
						msg: fmt.Sprintf("second pointer declaration"),
						Pos: item.pos,
					}
					return
				}

				// have we already seen tags? pointer must be first
				if len(tags) > 0 {
					err = &SyntaxError{
						msg: fmt.Sprintf("pointer declaration must precede tags"),
						Pos: item.pos,
					}
					return
				}

				pointer = true

				// and move on
				item = l.nextItem()
			}

			if item.typ == itemError {
				err = &SyntaxError{
					msg: item.val,
					Pos: item.pos,
				}
				return
			}

			// ok to be done at this point
			if item.typ == itemEOF {
				return
			}

			// next item needs to be an identifier (start of tag)
			if item.typ != itemIdentifier {
				err = &SyntaxError{
					msg: fmt.Sprintf("tag name required, found '%s'", item.String()),
					Pos: item.pos,
				}
				return
			}

			// we have an identifier, start a tag & move on
			t := Tag{
				Name:  item.val,
				Items: []string{},
			}

			item = l.nextItem()

			// next item must be a colonquote
			if item.typ != itemColonQuote {
				err = &SyntaxError{
					msg: fmt.Sprintf(`tag name must be followed by ':"', found '%s'`, item.String()),
					Pos: item.pos,
				}
				return
			}

		TagValues:
			// now inside a tag, loop through tag values
			for {
				item = l.nextItem()

				switch item.typ {
				case itemError:
					err = &SyntaxError{
						msg: item.val,
						Pos: item.pos,
					}
					return
				case itemEOF:
					// shouldn't happen within a tag
					err = &SyntaxError{
						msg: "expected a close quote",
						Pos: item.pos,
					}
					return
				case itemIdentifier:
					t.Items = append(t.Items, item.val)
				case itemComma:
					// absorb
					// de facto, spaces or commas as separators, but prefer commas for readability
				case itemCloseQuote:
					// we're done with this tag, get out
					break TagValues
				default:
					err = &SyntaxError{
						msg: fmt.Sprintf("unknown value '%v' in tag", item.val),
						Pos: item.pos,
					}
					return
				}
			}

			tags = append(tags, t)
		}
	}

	return
}

type SyntaxError struct {
	msg string
	Pos Pos
}

func (e *SyntaxError) Error() string {
	return e.msg
}
