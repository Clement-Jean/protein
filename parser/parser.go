package parser

import "github.com/Clement-Jean/protein/lexer"

type Parser struct {
	toks    *lexer.TokenizedBuffer
	tree    ParseTree
	stack   []stateStackEntry
	errs    []error
	currTok int
}

func New(toks *lexer.TokenizedBuffer) *Parser {
	return &Parser{
		toks:  toks,
		tree:  make([]Node, 0, len(toks.TokenInfos)),
		stack: make([]stateStackEntry, 0, 31),
	}
}

func (p *Parser) pushState(st state) {
	p.stack = append(p.stack, stateStackEntry{
		st:           st,
		tokIdx:       p.currTok,
		subtreeStart: int32(len(p.tree) - 1),
	})
}

func (p *Parser) popState() stateStackEntry {
	state := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return state
}

func (p *Parser) topState() stateStackEntry {
	return p.stack[len(p.stack)-1]
}

func (p *Parser) parseTopLevel() {
	p.popState()
}

func (p *Parser) Parse() (ParseTree, []error) {
	p.pushState(stateTopLevel)

	for len(p.stack) != 0 && p.currTok < len(p.toks.TokenInfos) {
		switch p.topState().st {
		case stateTopLevel:
			p.parseTopLevel()
		}
	}
	return p.tree, p.errs
}