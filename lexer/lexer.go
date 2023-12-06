package lexer

import (
	"github.com/Clement-Jean/protein/internal/span"
	"github.com/Clement-Jean/protein/token"
)

type Lexer interface {
	Tokenize() ([]token.Kind, []span.Span)
}
