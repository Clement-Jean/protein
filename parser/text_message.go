package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseTextMessage() {
	p.pushState(stateTextMessageFinish)
	p.pushState(stateTextMessageValue)
}

func (p *Parser) parseTextMessageValue() {
	p.popState()

	curr := p.curr()
	if curr != lexer.TokenKindRightBrace && curr != lexer.TokenKindRightAngle {
		p.parseTextField()
	}
}

func (p *Parser) parseTextMessageComma() {
	state := p.popState()
	top := p.topState()
	top.subtreeStart++
	p.addNode(state.tokIdx, top)
}

func (p *Parser) parseTextMessageFinish() {
	curr := p.curr()
	if curr == lexer.TokenKindComma {
		p.pushState(stateTextMessageComma)
		p.next()
		p.pushState(stateTextMessageValue)
		return
	}

	state := p.popState()
	tokIdx := p.currTok

	state.hasError = curr != expected

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindRightBrace, lexer.TokenKindRightAngle)
		tokIdx = p.skipPastLikelyEnd(tokIdx)
	}

	p.addNode(tokIdx, state)
}
