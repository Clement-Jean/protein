package lexer

import (
	"errors"
	"fmt"
)

type stateFn func() stateFn

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isLetter(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z'
}

func isIdentifier(b byte) bool {
	return isLetter(b) || b == '_' || isDigit(b)
}

func isQuote(b byte) bool {
	return b == '"' || b == '\''
}

func (l *Lexer) lexLineComment() (state stateFn) {
	const prefixLen = 2

	len := l.goToEndOfLineComment()
	state = l.emit(TokenKindComment, l.tokPos)
	l.tokPos += len + prefixLen
	return state
}

func (l *Lexer) goToEndOfLineComment() (len int) {
	var ch byte

	for ch = l.next(); ch != 0 && ch != '\n'; ch = l.next() {
		len++
	}

	if ch == '\n' {
		l.backup()
	}
	return len
}

func (l *Lexer) lexMultilineComment() (state stateFn) {
	const (
		prefixLen = 2
		suffixLen = 2
	)

	len, ok := l.goToEndOfMultilineComment()
	if ok {
		state = l.emit(TokenKindComment, l.tokPos)
		l.tokPos += len + prefixLen + suffixLen
		return state
	}

	state = l.error(errors.New("unclosed multiline comment"))
	l.tokPos += len + prefixLen
	return state
}

func (l *Lexer) goToEndOfMultilineComment() (len int, ok bool) {
	for ch := l.next(); ch != 0; ch = l.next() {
		if ch == '*' {
			if peek := l.next(); peek == '/' {
				return len, true
			}
			len++
		}
		len++
	}
	return len, false
}

func (l *Lexer) lexIdentifier() (state stateFn) {
	len := l.acceptWhile(isIdentifier)
	state = l.emit(TokenKindIdentifier, l.tokPos)
	l.tokPos += len
	return state
}

func (l *Lexer) lexString() (state stateFn) {
	const (
		escape    = '\\'
		prefixLen = 1
		suffixLen = 1
	)

	inEscape := false
	len := 0
	open := l.next()

	ch := l.next()
	for ; ch != 0 && ch != '\n'; ch = l.next() {
		switch {
		case inEscape:
			inEscape = false
		case ch == escape:
			inEscape = true
		case ch == open: // open and not escaped
			state = l.emit(TokenKindStr, l.tokPos)
			l.tokPos += len + prefixLen + suffixLen
			return state
		}
		len++
	}

	if ch == '\n' {
		l.backup()
	}

	state = l.error(errors.New("unclosed string"))
	l.tokPos += len + prefixLen
	return state
}

func (l *Lexer) lexProto() (state stateFn) {
	switch ch := l.next(); ch {
	case 0:
		return l.emit(TokenKindEOF, l.tokPos)
	case '\v', '\f', '\r', '\t', ' ', 0x85, 0xA0:
		break // skip
	case '\n':
		l.currLineIdx++
		if int(l.currLineIdx) >= len(l.toks.LineInfos) {
			l.currLineIdx--
			return l.emit(TokenKindEOF, l.tokPos)
		}
		l.tokPos = l.toks.LineInfos[l.currLineIdx].Start
		return nil
	case '_':
		state = l.emit(TokenKindUnderscore, l.tokPos)
	case '=':
		state = l.emit(TokenKindEqual, l.tokPos)
	case ',':
		state = l.emit(TokenKindComma, l.tokPos)
	case ':':
		state = l.emit(TokenKindColon, l.tokPos)
	case ';':
		state = l.emit(TokenKindSemicolon, l.tokPos)
	case '.':
		state = l.emit(TokenKindDot, l.tokPos)
	case '{':
		state = l.emit(TokenKindLeftBrace, l.tokPos)
	case '}':
		state = l.emit(TokenKindRightBrace, l.tokPos)
	case '[':
		state = l.emit(TokenKindLeftSquare, l.tokPos)
	case ']':
		state = l.emit(TokenKindRightSquare, l.tokPos)
	case '(':
		state = l.emit(TokenKindLeftParen, l.tokPos)
	case ')':
		state = l.emit(TokenKindRightParen, l.tokPos)
	case '<':
		state = l.emit(TokenKindLeftAngle, l.tokPos)
	case '>':
		state = l.emit(TokenKindRightAngle, l.tokPos)
	default:
		switch {
		case isLetter(ch):
			l.backup()
			return l.lexIdentifier
		case isQuote(ch):
			l.backup()
			return l.lexString
		case ch == '/':
			if l.readPos >= len(l.src) {
				state = l.emit(TokenKindSlash, l.tokPos)
				break
			}

			switch l.src[l.readPos] {
			case '/':
				l.next()
				return l.lexLineComment
			case '*':
				l.next()
				return l.lexMultilineComment
			default:
				state = l.emit(TokenKindSlash, l.tokPos)
			}
		default:
			state = l.error(fmt.Errorf("invalid char %q", ch))
		}
	}

	l.tokPos++
	return state
}
