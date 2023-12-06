package lexer

import "github.com/Clement-Jean/protein/token"

type stateFn func() stateFn

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
		if l.buf.atEOF {
			return l.emit(token.KindEOF)
		}
	}

	return l.emit(token.KindIllegal)
}
