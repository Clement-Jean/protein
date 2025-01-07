package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parsePackage() {
	p.pushState(statePackageFinish)
	p.pushState(stateFullIdentifierRoot)
}

func (p *Parser) parsePackageFinish() {
	state := p.popState()
	tokIdx := p.currTok

	state.hasError = p.curr() != lexer.TokenKindSemicolon

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindSemicolon)
		tokIdx = p.skipPastLikelyEnd(tokIdx)
	}

	p.addNode(tokIdx, state)
}
