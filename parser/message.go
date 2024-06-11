package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseMessage() {
	p.pushState(stateMessageFinish)
	p.pushState(stateMessageBlock)
	p.pushState(stateMessageName)
}

func (p *Parser) parseMessageName() {
	p.popState()

	curr := p.curr()
	hasError := curr != lexer.TokenKindIdentifier && !curr.IsIdentifier()
	p.addLeafNode(hasError)

	if !hasError {
		curr = p.next()
	} else {
		p.expectedCurr(lexer.TokenKindEqual)
		p.skipPastLikelyEnd(p.currTok)
	}
}

func (p *Parser) parseMessageBlock() {
	p.popState()

	curr := p.curr()
	hasError := curr != lexer.TokenKindLeftBrace
	p.addLeafNode(hasError)

	if !hasError {
		curr = p.next()
	} else {
		p.expectedCurr(lexer.TokenKindLeftBrace)
		p.skipPastLikelyEnd(p.currTok)
	}

	switch curr {
	case lexer.TokenKindOption:
		p.addLeafNode(false)
		p.next()
		p.parseOption()
	case lexer.TokenKindRightBrace:
	default:
		panic("not implemented")
	}
}

func (p *Parser) parseMessageFinish() {
	state := p.popState()
	tokIdx := p.currTok

	state.hasError = p.curr() != lexer.TokenKindRightBrace

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindRightBrace)
		tokIdx = p.skipPastLikelyEnd(tokIdx)
	}

	p.addNode(tokIdx, state)
}
