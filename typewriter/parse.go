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
			pointer, tags, err := parse(c.Text, directive, pkg)

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

			typ.Tags = tags
			typ.test = test(strings.HasSuffix(fset.Position(s.Pos()).Filename, "_test.go"))

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

			if c := findDirective(t.Doc, directive); c != nil {
				specs[t] = c
			}
		}

		// no need to keep walking, we don't care about TypeSpec's children
		return false
	})

	return specs
}

// findDirective return the first line of a doc which contains a directive
// the directive and '//' are removed
func findDirective(doc *ast.CommentGroup, directive string) *ast.Comment {
	if doc == nil {
		return nil
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

		return c
	}

	return nil
}

type parsr struct {
	lex       *lexer
	token     [2]item // two-token lookahead for parser.
	peekCount int
}

// next returns the next token.
func (p *parsr) next() item {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.token[0] = p.lex.nextItem()
	}
	return p.token[p.peekCount]
}

// backup backs the input stream up one token.
func (p *parsr) backup() {
	p.peekCount++
}

// peek returns but does not consume the next token.
func (p *parsr) peek() item {
	if p.peekCount > 0 {
		return p.token[p.peekCount-1]
	}
	p.peekCount = 1
	p.token[0] = p.lex.nextItem()
	return p.token[0]
}

func parse(input, directive string, evaluator evaluator) (Pointer, Tags, error) {
	var pointer Pointer
	var tags Tags
	p := &parsr{
		lex: lex(input),
	}

Loop:
	for {
		item := p.next()
		switch item.typ {
		case itemEOF:
			break Loop
		case itemError:
			err := &SyntaxError{
				msg: item.val,
				Pos: item.pos,
			}
			return false, nil, err
		case itemCommentPrefix:
			// don't care, move on
			continue
		case itemDirective:
			// is it the directive we care about?
			if item.val != directive {
				return false, nil, nil
			}
			continue
		case itemPointer:
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
		case itemTag:
			// we have an identifier, start a tag
			tag := Tag{
				Name: item.val,
			}

			// expect colonquote
			if p.next().typ != itemColonQuote {
				err := &SyntaxError{
					msg: fmt.Sprintf(`tag name must be followed by ':"', found '%s'`, item.val),
					Pos: item.pos,
				}
				return false, nil, err
			}

			negated, vals, err := parseTagValues(p, evaluator)

			if err != nil {
				return false, nil, err
			}

			tag.Negated = negated
			tag.Values = vals

			tags = append(tags, tag)
		default:
			return false, nil, unexpected(item)
		}
	}

	return pointer, tags, nil
}

func parseTagValues(p *parsr, evaluator evaluator) (bool, []TagValue, error) {
	var negated bool
	var vals []TagValue

	for {
		item := p.next()

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
			if len(vals) > 0 {
				err := &SyntaxError{
					msg: fmt.Sprintf("negation must precede tag values"),
					Pos: item.pos,
				}
				return false, nil, err
			}
			negated = true
		case itemTagValue:
			val := TagValue{
				Name: item.val,
			}

			if p.peek().typ == itemTypeParameter {
				typs, err := parseTypeParameters(p, evaluator)
				if err != nil {
					serr := &SyntaxError{
						msg: err.Error(),
						Pos: item.pos,
					}
					return false, nil, serr
				}
				val.TypeParameters = typs
			}

			vals = append(vals, val)
		case itemCloseQuote:
			// we're done
			return negated, vals, nil
		default:
			return false, nil, unexpected(item)
		}
	}
}

func parseTypeParameters(p *parsr, evaluator evaluator) ([]Type, error) {
	var typs []Type

	for {
		item := p.next()

		switch item.typ {
		case itemTypeParameter:
			typ, err := evaluator.Eval(item.val)
			if err != nil {
				serr := &SyntaxError{
					msg: err.Error(),
					Pos: item.pos,
				}
				return nil, serr
			}

			typs = append(typs, typ)
		default:
			p.backup()
			return typs, nil
		}
	}
}

func unexpected(item item) error {
	return &SyntaxError{
		msg: fmt.Sprintf("unexpected '%v'", item.val),
		Pos: item.pos,
	}
}

type SyntaxError struct {
	msg string
	Pos token.Pos
}

func (e *SyntaxError) Error() string {
	return e.msg
}
