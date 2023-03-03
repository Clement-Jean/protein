package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Impl is the implementation for the Lexer interface.
type Impl struct {
	src             string // the input text
	start           int    // the start of a token
	startLine       int    // the line at which a token start
	startLineOffset int    // the offset of the starting line relative to beginning of file
	line            int    // the current file line being process
	pos             int    // the reading position in the file
	atEOF           bool   // tells wether the Lexer is finished
	token           Token  // the token to return
}

// New creates a new instance of the Lexer
func New(input string) Lexer {
	return &Impl{
		src:       input,
		line:      1, // lines are 1-indexed
		startLine: 1,
	}
}

type stateFn func(*Impl) stateFn

func (l *Impl) next() rune {
	if int(l.pos) >= len(l.src) {
		l.atEOF = true
		return rune(EOF)
	}

	r, w := utf8.DecodeRuneInString(l.src[l.pos:])
	l.pos += w

	if r == '\n' {
		l.line++
	}
	return r
}

func (l *Impl) emit(tt TokenType) stateFn {
	t := Token{tt, l.src[l.start:l.pos], Position{
		Offset: l.start,
		Line:   l.startLine,
		Column: l.start - l.startLineOffset,
	}}
	l.start = l.pos
	l.startLine = l.line
	if tt == TokenSpace {
		if lineStart := strings.LastIndex(t.Literal, "\n"); lineStart != -1 {
			l.startLineOffset = l.start - (len(t.Literal) - 1 - lineStart)
		}
	}
	l.token = t
	return nil
}

func (l *Impl) peek() rune {
	if int(l.pos) >= len(l.src) {
		return rune(EOF)
	}

	r, _ := utf8.DecodeRuneInString(l.src[l.pos:])
	return r
}

func (l *Impl) backup() {
	if !l.atEOF && l.pos > 0 {
		r, w := utf8.DecodeLastRuneInString(l.src[:l.pos])
		l.pos -= w

		if r == '\n' {
			l.line--
		}
	}
}

func (l *Impl) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *Impl) acceptWhile(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *Impl) errorf(format string, args ...any) stateFn {
	l.token = Token{TokenError, fmt.Sprintf(format, args...), Position{
		Offset: l.start,
		Line:   l.startLine,
		Column: l.start - l.startLineOffset,
	}}
	l.start = 0
	l.pos = 0
	l.src = l.src[:0]
	return nil
}

func lexSpaces(l *Impl) stateFn {
	var r rune

	for {
		r = l.peek()
		if !unicode.IsSpace(r) {
			break
		}

		l.next()
	}
	return l.emit(TokenSpace)
}

func lexLineComment(l *Impl) stateFn {
	var r rune

	for {
		r = l.peek()
		if r == '\n' || r == rune(EOF) {
			break
		}

		l.next()
	}
	return l.emit(TokenComment)
}

func lexMultilineComment(l *Impl) stateFn {
	var p rune
	var r rune

	for {
		p = r
		if r == rune(EOF) {
			return l.errorf(errorUnterminatedMultilineComment)
		}

		r = l.peek()
		if p == '*' && r == '/' {
			l.next()
			break
		}

		l.next()
	}
	return l.emit(TokenComment)
}

func lexIdentifier(l *Impl) stateFn {
	l.acceptWhile("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")
	return l.emit(TokenIdentifier)
}

func lexString(l *Impl) stateFn {
	open := l.src[l.pos]
	l.next()
Loop:
	for {
		switch l.next() {
		case rune(EOF):
			return l.errorf(errorUnterminatedQuotedString)
		case rune(open):
			break Loop
		}
	}
	return l.emit(TokenStr)
}

func lexNumber(l *Impl) stateFn {
	var t TokenType = TokenInt

	l.accept("+-")

	digits := "0123456789" // decimal

	if l.accept("0") { // starts with 0
		if l.accept("xX") {
			digits = "0123456789abcdefABCDEF" // hexadecimal
		} else {
			digits = "01234567" // octal
		}
	}

	l.acceptWhile(digits)

	if l.accept(".") {
		t = TokenFloat
		l.acceptWhile("0123456789")
	}

	if l.accept("eE") { // exponent
		t = TokenFloat
		l.accept("+-")
		l.acceptWhile("0123456789")
	}

	return l.emit(t)
}

func lexProto(l *Impl) stateFn {
	switch r := l.next(); {
	case l.atEOF:
		return l.emit(EOF)
	case r == '_':
		return l.emit(TokenUnderscore)
	case r == '=':
		return l.emit(TokenEqual)
	case r == ',':
		return l.emit(TokenColon)
	case r == ';':
		return l.emit(TokenSemicolon)
	case r == '.' && !unicode.IsNumber(l.peek()):
		return l.emit(TokenDot)
	case r == '{':
		return l.emit(TokenLeftBrace)
	case r == '}':
		return l.emit(TokenRightBrace)
	case r == '[':
		return l.emit(TokenLeftSquare)
	case r == ']':
		return l.emit(TokenRightSquare)
	case r == '(':
		return l.emit(TokenLeftParen)
	case r == ')':
		return l.emit(TokenRightParen)
	case r == '<':
		return l.emit(TokenLeftAngle)
	case r == '>':
		return l.emit(TokenRightAngle)
	case unicode.IsSpace(r):
		l.backup()
		return lexSpaces
	case r == '/' && l.peek() == '/':
		l.backup()
		return lexLineComment
	case r == '/' && l.peek() == '*':
		l.backup()
		return lexMultilineComment
	case unicode.IsLetter(r):
		l.backup()
		return lexIdentifier
	case r == '"' || r == '\'':
		l.backup()
		return lexString
	case r == '+' || r == '-' || r == '.' || ('0' <= r && r <= '9'):
		l.backup()
		return lexNumber
	}

	return l.emit(TokenIllegal)
}

// NextToken provides the following token in the input
func (l *Impl) NextToken() Token {
	state := lexProto
	for {
		state = state(l)
		if state == nil {
			return l.token
		}
	}
}
