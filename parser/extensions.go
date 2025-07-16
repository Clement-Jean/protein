package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseExtensions() {
	p.pushState(stateExtensionsFinish)
	p.pushState(stateReservedRange)
}

func (p *Parser) parseExtensionsFinish() {
	curr := p.curr()
	switch curr {
	case lexer.TokenKindComma:
		p.pushState(stateEnder)
		p.next()
		p.pushState(stateReservedRange)
		return
	case lexer.TokenKindLeftSquare:
		p.addLeafNode(false)
		p.next()
		p.pushState(stateMessageFieldOption)
		return
	}

	state := p.popState()
	tokIdx := p.currTok

	state.hasError = curr != lexer.TokenKindSemicolon

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindSemicolon)
		tokIdx = p.skipPastLikelyEnd(tokIdx)
	}

	p.addNode(tokIdx, state)
}
