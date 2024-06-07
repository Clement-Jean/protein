package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseImport() {
	p.pushState(stateImportFinish)
	p.pushState(stateImportValue)
}

func (p *Parser) parseImportValue() {
	p.popState()

	curr := p.curr()
	if curr == lexer.TokenKindPublic || curr == lexer.TokenKindWeak {
		p.addLeafNode(false)
		curr = p.next()
	}

	hasError := curr != lexer.TokenKindStr
	p.addLeafNode(hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindStr)
		p.skipPastLikelyEnd(p.currTok)
	}
}

func (p *Parser) parseImportFinish() {
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
