package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseTextMessage() {
	p.pushState(stateTextMessageFinish)
	p.pushState(stateTextMessageValue)
}

func (p *Parser) parseTextMessageValue() {
	p.popState()

	if p.curr() != lexer.TokenKindRightBrace && p.curr() != lexer.TokenKindRightAngle {
		p.parseTextField()
	}
}

func (p *Parser) parseTextMessageFinish() {
	state := p.popState()
	tokIdx := p.currTok

	state.hasError = p.curr() != lexer.TokenKindRightBrace &&
		p.curr() != lexer.TokenKindRightAngle

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindRightBrace, lexer.TokenKindRightAngle)
		tokIdx = p.skipPastLikelyEnd(tokIdx)
	}

	p.addNode(tokIdx, state)
}
