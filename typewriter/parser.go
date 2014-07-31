package typewriter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	_ "code.google.com/p/go.tools/go/gcimporter"
	"code.google.com/p/go.tools/go/types"
)

func getTypes(directive string, filter func(os.FileInfo) bool) ([]Type, error) {
	typs := make([]Type, 0)

	fset := token.NewFileSet()
	rootDir := "./"

	astPackages, err := parser.ParseDir(fset, rootDir, filter, parser.ParseComments)

	if err != nil {
		return typs, err
	}

	for name, astPackage := range astPackages {

		// collect type specs
		var specs []*ast.TypeSpec
		ast.Inspect(astPackage, func(n ast.Node) bool {
			// is it a type?
			// http://golang.org/pkg/go/ast/#GenDecl
			d, ok := n.(*ast.GenDecl)

			if !ok || d.Tok != token.TYPE {
				// never mind, move on
				return true
			}

			if d.Lparen == 0 {
				// not parenthesized, copy GenDecl.Doc into TypeSpec.Doc
				d.Specs[0].(*ast.TypeSpec).Doc = d.Doc
			}

			for _, s := range d.Specs {
				specs = append(specs, s.(*ast.TypeSpec))
			}

			// no need to keep walking, we don't care about TypeSpec's children
			return false
		})

		// pull map into a slice
		var files []*ast.File
		for _, f := range astPackage.Files {
			files = append(files, f)
		}

		typesPkg, err := types.Check(name, fset, files)

		if err != nil {
			return typs, err
		}

		pkg := &Package{typesPkg}

		for _, spec := range specs {
			pointer, tags, found, err := parseTags(directive, spec.Doc.Text())

			if err != nil {
				return typs, err
			}

			if !found {
				continue
			}

			typ, err := pkg.Eval(pointer.String() + spec.Name.Name)

			if err != nil {
				// really shouldn't happen, since the type came from the ast in the first place
				err = fmt.Errorf("failed to evaluate type %s (%s)", typ.Name, err)
				return typs, err
			}

			typ.test = test(strings.HasSuffix(fset.Position(spec.Pos()).Filename, "_test.go"))
			typ.Tags = tags

			typs = append(typs, typ)
		}
	}

	return typs, nil
}

// identifies gen-marked types and parses tags
func parseTags(directive string, doc string) (pointer Pointer, tags []Tag, found bool, err error) {
	lines := strings.Split(doc, "\n")
	for _, line := range lines {
		original := line

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
