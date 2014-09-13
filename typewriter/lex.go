// Derived from http://golang.org/pkg/text/template/parse/

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typewriter

import (
	"fmt"
	"go/token"
	"unicode"
	"unicode/utf8"
)

// item represents a token or text string returned from the scanner.
type item struct {
	typ itemType  // The type of this item.
	pos token.Pos // The starting position, in bytes, of this item in the input string.
	val string    // The value of this item.
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota // error occurred; value is text of error
	itemCommentPrefix
	itemDirective
	itemPointer
	itemTag
	itemColonQuote
	itemMinus
	itemTagValue
	itemTypeParameter
	itemCloseQuote
	itemEOF
)

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	input        string    // the string being scanned
	state        stateFn   // the next lexing function to enter
	pos          token.Pos // current position in the input
	start        token.Pos // start position of this item
	width        int       // width of last rune read from input
	lastPos      token.Pos // position of most recent item returned by nextItem
	items        chan item // channel of scanned items
	bracketDepth int
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += token.Pos(l.width)
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= token.Pos(l.width)
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.pos, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
func (l *lexer) nextItem() item {
	item := <-l.items
	l.lastPos = item.pos
	return item
}

// lex creates a new scanner for the input string.
func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexComment; l.state != nil; {
		l.state = l.state(l)
	}
}

// state functions

func lexComment(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case r == eof:
			break Loop
		case r == '/':
			return lexCommentPrefix
		case r == '+':
			return lexDirective
		case isSpace(r):
			l.ignore()
		case r == '*':
			l.emit(itemPointer)
			p := l.peek()
			if !isSpace(p) && p != eof {
				return l.errorf("pointer must be followed by a space or EOL")
			}
		case isIdentifierPrefix(r):
			return lexTag
		default:
			return l.errorf("illegal leading character '%s' in tag name", string(r))
		}
	}
	l.emit(itemEOF)
	return nil
}

// lexTag scans the elements inside quotes
func lexTag(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isIdentifierPrefix(r):
			return lexIdentifier(l, itemTag)
		case r == ':':
			if l.next() != '"' {
				return l.errorf(`expected :" following tag name`)
			}
			l.emit(itemColonQuote)
			return lexTagValues
		case isSpace(r):
			l.ignore()
		case r == '"':
			l.emit(itemCloseQuote)
			return lexComment
		default:
			return l.errorf("illegal character '%s' in tag name", string(r))
		}
	}
}

func lexTagValues(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == '-':
			l.emit(itemMinus)
		case isIdentifierPrefix(r):
			return lexIdentifier(l, itemTagValue)
		case r == '[':
			l.bracketDepth++
			// parser has no use for bracket, only important as delimiter here
			l.ignore()
			return lexTypeParameters
		case isSpace(r) || r == ',':
			// parser has no use for comma, only important as delimiter here
			l.ignore()
		case r == '"':
			// let lexTag handle it
			l.backup()
			return lexTag
		case r == eof:
			// we fell off the end without a close quote
			return lexComment
		default:
			return l.errorf("illegal character '%s' in tag value", string(r))
		}
	}
}

func lexTypeParameters(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == ']':
			if l.bracketDepth == 0 {
				// closing bracket of type parameter
				return lexTagValues
			}
			return l.errorf("additional close bracket")
		case isTypeDecl(r):
			return lexTypeDeclaration
		case isSpace(r) || r == ',':
			l.ignore()
		case r == '"':
			// premature end
			return l.errorf("expected close bracket")
		default:
			return l.errorf("illegal character '%s' in type parameter", string(r))
		}
	}
}

func lexTypeDeclaration(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case r == ']':
			l.bracketDepth--
			if l.bracketDepth == 0 {
				// closing bracket of type parameter
				break Loop
			}
			// if bracket depth remains > 0, it's part of the type declaration eg []string
			// absorb
		case isTypeDecl(r):
			// absorb
			if r == '[' {
				l.bracketDepth++
			}
		case isSpace(r) || r == ',':
			// legal delimiters for multiple type parameters
			break Loop
		default:
			// anything else is illegal
			return l.errorf("illegal character '%c' in type declaration", r)
		}
	}

	// once we get here, we've absorbed the delimiter; backup as not to emit it
	l.backup()
	l.emit(itemTypeParameter)
	return lexTypeParameters
}

func lexCommentPrefix(l *lexer) stateFn {
	for l.peek() == '/' {
		l.next()
	}
	l.emit(itemCommentPrefix)
	return lexComment
}

func lexDirective(l *lexer) stateFn {
	for isAlphaNumeric(l.peek()) {
		l.next()
	}
	l.emit(itemDirective)
	return lexComment
}

// lexIdentifier scans an alphanumeric.
func lexIdentifier(l *lexer, typ itemType) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			if !isTerminator(r) {
				return l.errorf("illegal character '%c' in identifier", r)
			}
			l.backup()
			l.emit(typ)
			break Loop
		}
	}
	switch typ {
	case itemTag:
		return lexTag
	case itemTagValue:
		return lexTagValues
	default:
		return l.errorf("unknown itemType %v", typ)
	}
}

func isTerminator(r rune) bool {
	if isSpace(r) || isEndOfLine(r) {
		return true
	}
	switch r {
	case eof, ':', ',', '"', '[', ']':
		return true
	}
	return false
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isIdentifierPrefix reports whether r is an alphabetic or underscore, per http://golang.org/ref/spec#Identifiers
func isIdentifierPrefix(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

// isTypeDecl reports whether r a character legal in a type declaration, eg map[*Thing]interface{}
// brackets are a special case, handled in lexTypeParameter
func isTypeDecl(r rune) bool {
	return r == '*' || r == '{' || r == '}' || isAlphaNumeric(r)
}
