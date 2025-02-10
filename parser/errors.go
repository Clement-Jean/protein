package parser

import (
	"fmt"

	"github.com/Clement-Jean/protein/lexer"
)

type ExpectedError struct {
	Expected []lexer.TokenKind
	Got      lexer.TokenKind
	TokIdx   uint32
}

func (e *ExpectedError) Error() string {
	return fmt.Sprintf("expected %v, got %s", e.Expected, e.Got)
}
