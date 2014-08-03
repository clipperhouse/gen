// Derived from http://golang.org/pkg/text/template/parse/

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typewriter

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// item represents a token or text string returned from the scanner.
type item struct {
	typ itemType // The type of this item.
	pos int      // The starting position, in bytes, of this item in the input string.
	val string   // The value of this item.
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota // error occurred; value is text of error
	itemPointer
	itemIdentifier
	itemColonQuote
	itemCloseQuote
	itemEOF
	itemSpace
	itemComma
	itemMinus
	itemLeftBracket
	itemRightBracket
)

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	input   string    // the string being scanned
	state   stateFn   // the next lexing function to enter
	pos     int       // current position in the input
	start   int       // start position of this item
	width   int       // width of last rune read from input
	lastPos int       // position of most recent item returned by nextItem
	items   chan item // channel of scanned items
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
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
	l.pos -= l.width
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
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
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
	for l.state = lexDirective; l.state != nil; {
		l.state = l.state(l)
	}
}

// state functions

func lexDirective(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case r == eof:
			break Loop
		case isSpace(r):
			l.ignore()
			continue
		case r == '*':
			l.emit(itemPointer)
			p := l.peek()
			if !isSpace(p) && p != eof {
				return l.errorf("pointer must be followed by a space or EOL")
			}
		case isIdentifierPrefix(r):
			l.backup()
			return lexInsideTag
		default:
			return l.errorf("illegal leading character '%s' in identifier", string(r))
		}
	}
	l.emit(itemEOF)
	return nil
}

// lexInsideTag scans the elements inside quotes
func lexInsideTag(l *lexer) stateFn {
	switch r := l.next(); {
	case isIdentifierPrefix(r):
		l.backup()
		return lexIdentifier(l, lexInsideTag)
	case r == ':':
		if l.next() != '"' {
			return l.errorf(`expected :" following tag name`)
		}
		l.emit(itemColonQuote)
		return lexInsideTagValue
	case isSpace(r):
		l.ignore()
	case r == '"':
		l.emit(itemCloseQuote)
		return lexDirective
	default:
		return l.errorf("illegal character '%s' in tag name", string(r))
	}
	return lexDirective
}

func lexInsideTagValue(l *lexer) stateFn {
	switch r := l.next(); {
	case isIdentifierPrefix(r):
		l.backup()
		return lexIdentifier(l, lexInsideTagValue)
	case r == '-':
		l.emit(itemMinus)
	case r == ',':
		l.emit(itemComma)
	case r == '"':
		// defer back up
		l.backup()
		return lexInsideTag
	case isSpace(r):
		l.ignore()
	case r == eof:
		// we fell off the end without a close quote
		return lexDirective
	default:
		return l.errorf("illegal character '%s' in tag value", string(r))
	}
	return lexInsideTagValue
}

// lexSpace scans a run of space characters.
// One space has already been seen.
func lexSpace(l *lexer) stateFn {
	for isSpace(l.peek()) {
		l.next()
	}
	l.emit(itemSpace)
	return lexInsideTag
}

// lexIdentifier scans an alphanumeric.
func lexIdentifier(l *lexer, fn stateFn) stateFn {
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
			l.emit(itemIdentifier)
			break Loop
		}
	}
	return fn
}

func isTerminator(r rune) bool {
	if isSpace(r) || isEndOfLine(r) {
		return true
	}
	switch r {
	case eof, ':', ',', '"':
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
