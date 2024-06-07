package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseTextMessage() {
	if p.curr() == lexer.TokenKindLeftBrace {
		p.pushState(stateTextMessageFinishRightBrace)
	} else {
		p.pushState(stateTextMessageFinishRightAngle)
	}

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
	expected := lexer.TokenKindRightBrace
	if state.st == stateTextMessageFinishRightAngle {
		expected = lexer.TokenKindRightAngle
	}

	state.hasError = curr != expected

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(expected)
		tokIdx = p.skipPastLikelyEnd(tokIdx)
		p.next()
	}

	p.addNode(tokIdx, state)
}
