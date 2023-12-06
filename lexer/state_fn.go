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

func (l *impl) lexSpaces() stateFn {
	for isSpace(l.buf.peek()) {
		l.buf.next()
	}
	return l.emit(token.KindSpace)
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
		return l.emit(token.KindDot)
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
	case '/':
		return l.emit(token.KindSlash)
	default:
		switch {
		case isSpace(next):
			l.buf.backup()
			return l.lexSpaces
		case l.buf.atEOF:
			return l.emit(token.KindEOF)
		}
	}

	return l.emit(token.KindIllegal)
}
