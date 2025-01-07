package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseSyntax() {
	p.pushState(stateSyntaxFinish)
	p.pushState(stateSyntaxAssign)
}

func (p *Parser) parseSyntaxAssign() {
	p.popState()

	curr := p.curr()
	hasError := curr != lexer.TokenKindEqual
	p.addLeafNode(hasError)

	if hasError {
		p.expectedCurr(lexer.TokenKindEqual)
		p.skipPastLikelyEnd(p.currTok)
		return
	}
	curr = p.next()

	hasError = curr != lexer.TokenKindStr
	p.addLeafNode(hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindStr)
		p.skipPastLikelyEnd(p.currTok)
	}
}

func (p *Parser) parseSyntaxFinish() {
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
