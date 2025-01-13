package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseTextListValue() {
	curr := p.curr()

	if curr == lexer.TokenKindComma {
		p.pushState(stateEnder)
		curr = p.next()
	}

	switch curr {
	case lexer.TokenKindRightSquare:
		p.popState()
	case lexer.TokenKindLeftBrace, lexer.TokenKindLeftAngle:
		p.addLeafNode(false)
		p.parseTextMessage()
		p.next()
	case lexer.TokenKindIdentifier:
		p.parseTextField()
	case lexer.TokenKindStr, lexer.TokenKindFloat, lexer.TokenKindInt:
		p.addLeafNode(false)
		p.next()
	default:
		if curr != lexer.TokenKindComma {
			state := p.popState()
			p.expectedCurr(lexer.TokenKindRightSquare)
			p.skipPastLikelyEnd(p.currTok)
			p.addNode(p.currTok, state)
			return
		}
	}
}

func (p *Parser) parseTextListFinish() {
	state := p.popState()
	tokIdx := p.currTok

	state.hasError = p.curr() != lexer.TokenKindRightSquare

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindRightSquare)
		tokIdx = p.skipPastLikelyEnd(tokIdx)
	}

	p.addNode(tokIdx, state)
}
