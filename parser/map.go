package parser

import "github.com/Clement-Jean/protein/lexer"

var mapKeyTypes = []lexer.TokenKind{
	lexer.TokenKindTypeInt32,
	lexer.TokenKindTypeInt64,
	lexer.TokenKindTypeUint32,
	lexer.TokenKindTypeUint64,
	lexer.TokenKindTypeSint32,
	lexer.TokenKindTypeSint64,
	lexer.TokenKindTypeFixed32,
	lexer.TokenKindTypeFixed64,
	lexer.TokenKindTypeSfixed32,
	lexer.TokenKindTypeSfixed64,
	lexer.TokenKindTypeBool,
	lexer.TokenKindTypeString,
}
