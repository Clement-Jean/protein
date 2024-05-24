package lexer

import "fmt"

type stateFn func() stateFn

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
		case ch == '/':
			state = l.emit(TokenKindSlash, l.tokPos)
		default:
			state = l.error(fmt.Errorf("invalid char %q", ch))
		}
	}

	l.tokPos++
	return state
}
