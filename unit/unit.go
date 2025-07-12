package unit

import (
	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/source"
)

type Unit struct {
	File   string
	Buffer *source.Buffer
	Toks   *lexer.TokenizedBuffer
	Tree   parser.ParseTree
}
