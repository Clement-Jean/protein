package parser

import "github.com/Clement-Jean/protein/lexer"

var constantTypes = []lexer.TokenKind{
	lexer.TokenKindTrue, lexer.TokenKindFalse,
	lexer.TokenKindInt, lexer.TokenKindFloat,
	lexer.TokenKindStr,
}
