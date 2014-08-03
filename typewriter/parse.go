package typewriter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func getTypes(directive string, filter func(os.FileInfo) bool) ([]Type, error) {
	// get the AST
	fset := token.NewFileSet()
	astPkgs, err := parser.ParseDir(fset, "./", filter, parser.ParseComments)

	if err != nil {
		return nil, err
	}

	var typs []Type

	for _, a := range astPkgs {
		pkg, err := getPackage(fset, a)

		if err != nil {
			return nil, err
		}

		specs := getTaggedSpecs(a, directive)

		for s, d := range specs {
			pointer, tags, err := parseTags(d)

			if err != nil {
				return nil, err
			}

			typ, err := pkg.Eval(pointer.String() + s.Name.Name)

			if err != nil {
				// really shouldn't happen, since the type came from the ast in the first place
				err = fmt.Errorf("failed to evaluate type %s (%s)", typ.Name, err)
				return nil, err
			}

			typ.test = test(strings.HasSuffix(fset.Position(s.Pos()).Filename, "_test.go"))
			typ.Tags = tags

			typs = append(typs, typ)
		}
	}

	return typs, nil
}

// getTaggedSpecs walks the AST and returns types which have directive comment
// returns a map of TypeSpec to directive
func getTaggedSpecs(pkg *ast.Package, directive string) map[*ast.TypeSpec]string {
	specs := make(map[*ast.TypeSpec]string)

	ast.Inspect(pkg, func(n ast.Node) bool {
		g, ok := n.(*ast.GenDecl)

		// is it a type?
		// http://golang.org/pkg/go/ast/#GenDecl
		if !ok || g.Tok != token.TYPE {
			// never mind, move on
			return true
		}

		if g.Lparen == 0 {
			// not parenthesized, copy GenDecl.Doc into TypeSpec.Doc
			g.Specs[0].(*ast.TypeSpec).Doc = g.Doc
		}

		for _, s := range g.Specs {
			t := s.(*ast.TypeSpec)

			if found, d := findDirective(t.Doc.Text(), directive); found {
				specs[t] = d
			}
		}

		// no need to keep walking, we don't care about TypeSpec's children
		return false
	})

	return specs
}

// findDirective return the first line of a doc which contains a directive
// the directive and '//' are removed
func findDirective(doc, directive string) (bool, string) {
	// check lines of doc for directive
	for _, l := range strings.Split(doc, "\n") {
		// does the line start with the directive?
		if l = strings.TrimLeft(l, "/ "); !strings.HasPrefix(l, directive) {
			continue
		}

		// remove the directive from the line
		l = strings.TrimPrefix(l, directive)

		// must be eof or followed by a space
		if len(l) > 0 && l[0] != ' ' {
			continue
		}

		return true, strings.TrimSpace(l)
	}

	return false, ""
}

// identifies gen-marked types and parses tags
func parseTags(d string) (Pointer, Tags, error) {
	var pointer Pointer
	var tags Tags

	l := lex(d)

	// top level can be pointer or tags
	for {
		item := l.nextItem()

		if item.typ == itemError {
			err := &SyntaxError{
				msg: item.val,
				Pos: item.pos,
			}
			return false, nil, err
		}

		// pick up pointer if it exists
		if item.typ == itemPointer {
			// have we already seen a pointer?
			if pointer {
				err := &SyntaxError{
					msg: fmt.Sprintf("second pointer declaration"),
					Pos: item.pos,
				}
				return false, nil, err
			}

			// have we already seen tags? pointer must be first
			if len(tags) > 0 {
				err := &SyntaxError{
					msg: fmt.Sprintf("pointer declaration must precede tags"),
					Pos: item.pos,
				}
				return false, nil, err
			}

			pointer = true

			// and move on
			item = l.nextItem()
		}

		if item.typ == itemError {
			err := &SyntaxError{
				msg: item.val,
				Pos: item.pos,
			}
			return false, nil, err
		}

		// ok to be done at this point
		if item.typ == itemEOF {
			return pointer, tags, nil
		}

		// next item needs to be an identifier (start of tag)
		if item.typ != itemIdentifier {
			err := &SyntaxError{
				msg: fmt.Sprintf("tag name required, found '%s'", item.val),
				Pos: item.pos,
			}
			return false, nil, err
		}

		// we have an identifier, start a tag & move on
		t := Tag{
			Name:  item.val,
			Items: []string{},
		}

		item = l.nextItem()

		// next item must be a colonquote
		if item.typ != itemColonQuote {
			err := &SyntaxError{
				msg: fmt.Sprintf(`tag name must be followed by ':"', found '%s'`, item.val),
				Pos: item.pos,
			}
			return false, nil, err
		}

	TagValues:
		// now inside a tag, loop through tag values
		for {
			item = l.nextItem()

			switch item.typ {
			case itemError:
				err := &SyntaxError{
					msg: item.val,
					Pos: item.pos,
				}
				return false, nil, err
			case itemEOF:
				// shouldn't happen within a tag
				err := &SyntaxError{
					msg: "expected a close quote",
					Pos: item.pos,
				}
				return false, nil, err
			case itemIdentifier:
				t.Items = append(t.Items, item.val)
			case itemComma:
				// absorb
				// de facto, spaces or commas as separators, but prefer commas for readability
			case itemCloseQuote:
				// we're done with this tag, get out
				break TagValues
			default:
				err := &SyntaxError{
					msg: fmt.Sprintf("unknown value '%v' in tag", item.val),
					Pos: item.pos,
				}
				return false, nil, err
			}
		}

		tags = append(tags, t)
	}
}

type SyntaxError struct {
	msg string
	Pos int
}

func (e *SyntaxError) Error() string {
	return e.msg
}
