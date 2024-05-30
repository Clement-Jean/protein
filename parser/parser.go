package parser

import "github.com/Clement-Jean/protein/lexer"

type Parser struct {
	toks    *lexer.TokenizedBuffer
	tree    ParseTree
	errs    []error
	currTok int
}

func New(toks *lexer.TokenizedBuffer) *Parser {
	return &Parser{
		toks: toks,
		tree: make([]Node, 0, len(toks.TokenInfos)),
	}
}

func (p *Parser) Parse() (ParseTree, []error) {
	return p.tree, p.errs
}
