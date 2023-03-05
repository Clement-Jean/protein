package parser

import (
	"github.com/Clement-Jean/protein/lexer"
)

type FakeLexer struct {
	i      int
	tokens []lexer.Token
}

func (l *FakeLexer) NextToken() lexer.Token {
	if l.i >= len(l.tokens) {
		return lexer.Token{Type: lexer.EOF, Position: lexer.Position{}}
	}

	token := l.tokens[l.i]
	l.i++
	return token
}
