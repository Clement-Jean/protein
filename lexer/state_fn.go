package lexer

import "github.com/Clement-Jean/protein/token"

type stateFn func() stateFn

func isSpace(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
		return true
	}
	return false
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func isIdentifier(r rune) bool {
	return isLetter(r) || r == '_' || isDigit(r)
}

func isHexadecimalDigit(r rune) bool {
	return isDigit(r) || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

func isOctalDigit(r rune) bool {
	return r >= '0' && r <= '7'
}

func (l *impl) lexSpaces() stateFn {
	for isSpace(l.buf.peek()) {
		l.buf.next()
	}
	return l.emit(token.KindSpace)
}

func (l *impl) lexLineComment() stateFn {
	for r := l.buf.peek(); !l.buf.atEOF && r != '\n'; r = l.buf.peek() {
		l.buf.next()
	}

	return l.emit(token.KindComment)
}

func (l *impl) lexMultilineComment() stateFn {
	next := l.buf.next()

	for !l.buf.atEOF {
		if next == '*' && l.buf.peek() == '/' {
			l.buf.next()
			return l.emit(token.KindComment)
		}

		next = l.buf.next()
	}

	return l.errorf(token.KindErrorUnterminatedMultilineComment)
}

func (l *impl) lexIdentifier() stateFn {
	l.acceptWhile(isIdentifier)
	return l.emit(token.KindIdentifier)
}

func (l *impl) lexString() stateFn {
	const escape = '\\'

	open := l.buf.next()
	next := l.buf.next()
	inEscape := false

	for ; !l.buf.atEOF && next != '\n'; next = l.buf.next() {
		switch {
		case inEscape:
			inEscape = false
		case next == escape:
			inEscape = true
		case next == open: // open and not escaped
			return l.emit(token.KindStr)
		}
	}

	if next == '\n' {
		l.buf.backup()
	}
	return l.errorf(token.KindErrorUnterminatedQuotedString)
}

func (l *impl) lexNumber() stateFn {
	kind := token.KindInt

	l.accept("+-")

	digits := isDigit

	if l.accept("0") { // starts with 0
		if l.accept("xX") {
			digits = isHexadecimalDigit
		} else {
			digits = isOctalDigit
		}
	}

	l.acceptWhile(digits)

	if l.accept(".") {
		kind = token.KindFloat
		l.acceptWhile(isDigit)
	}

	if l.accept("eE") { // exponent
		kind = token.KindFloat
		l.accept("+-")
		l.acceptWhile(isDigit)
	}

	return l.emit(kind)
}

func (l *impl) lexProto() stateFn {
	if l.buf.atEOF {
		return l.emit(token.KindEOF)
	}

	switch next := l.buf.next(); next {
	case '_':
		return l.emit(token.KindUnderscore)
	case '=':
		return l.emit(token.KindEqual)
	case ',':
		return l.emit(token.KindComma)
	case ':':
		return l.emit(token.KindColon)
	case ';':
		return l.emit(token.KindSemicolon)
	case '.':
		if !isDigit(l.buf.peek()) {
			return l.emit(token.KindDot)
		}
		l.buf.backup()
		return l.lexNumber
	case '{':
		return l.emit(token.KindLeftBrace)
	case '}':
		return l.emit(token.KindRightBrace)
	case '[':
		return l.emit(token.KindLeftSquare)
	case ']':
		return l.emit(token.KindRightSquare)
	case '(':
		return l.emit(token.KindLeftParen)
	case ')':
		return l.emit(token.KindRightParen)
	case '<':
		return l.emit(token.KindLeftAngle)
	case '>':
		return l.emit(token.KindRightAngle)
	default:
		switch {
		case isLetter(next):
			l.buf.backup()
			return l.lexIdentifier
		case isSpace(next):
			l.buf.backup()
			return l.lexSpaces
		case isDigit(next) || next == '-' || next == '+' || next == '.':
			l.buf.backup()
			return l.lexNumber
		case next == '"' || next == '\'':
			l.buf.backup()
			return l.lexString
		case next == '/':
			peek := l.buf.peek()
			if peek == '/' {
				l.buf.backup()
				return l.lexLineComment
			} else if peek == '*' {
				l.buf.backup()
				return l.lexMultilineComment
			}
			return l.emit(token.KindSlash)
		case l.buf.atEOF:
			return l.emit(token.KindEOF)
		}
	}

	return l.emit(token.KindIllegal)
}
