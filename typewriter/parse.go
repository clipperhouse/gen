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

		specs := getTaggedComments(a, directive)

		for s, c := range specs {
			pointer, tags, err := parseComment(c, directive)

			if err != nil {
				if serr, ok := err.(*SyntaxError); ok {
					// error should have Pos relative to the whole AST
					serr.Pos += c.Slash
					// Go-style syntax error with filename, line number, column
					serr.msg = fset.Position(serr.Pos).String() + ": " + serr.msg
				}
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

// getTaggedComments walks the AST and returns types which have directive comment
// returns a map of TypeSpec to directive
func getTaggedComments(pkg *ast.Package, directive string) map[*ast.TypeSpec]*ast.Comment {
	specs := make(map[*ast.TypeSpec]*ast.Comment)

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

			if found, d := findDirective(t.Doc, directive); found {
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
func findDirective(doc *ast.CommentGroup, directive string) (bool, *ast.Comment) {
	if doc == nil {
		return false, nil
	}

	// check lines of doc for directive
	for _, c := range doc.List {
		l := c.Text
		// does the line start with the directive?
		t := strings.TrimLeft(l, "/ ")
		if !strings.HasPrefix(t, directive) {
			continue
		}

		// remove the directive from the line
		t = strings.TrimPrefix(t, directive)

		// must be eof or followed by a space
		if len(t) > 0 && t[0] != ' ' {
			continue
		}

		return true, c
	}

	return false, nil
}

// identifies gen-marked types and parses tags
func parseComment(comment *ast.Comment, directive string) (Pointer, Tags, error) {
	var pointer Pointer
	var tags Tags

	l := lex(comment.Text)

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

		if item.typ == itemCommentPrefix {
			// don't care, move on
			continue
		}

		if item.typ == itemDirective {
			// is it the directive we care about?
			if item.val != directive {
				return false, nil, nil
			}
			continue
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
		if item.typ != itemTag {
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
			case itemMinus:
				if len(t.Items) > 0 {
					err := &SyntaxError{
						msg: fmt.Sprintf("negation must precede tag values"),
						Pos: item.pos,
					}
					return false, nil, err
				}
				t.Negated = true
			case itemTagValue:
				t.Items = append(t.Items, item.val)
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
	Pos token.Pos
}

func (e *SyntaxError) Error() string {
	return e.msg
}
