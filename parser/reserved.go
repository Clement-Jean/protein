package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseReserved() {
	switch p.curr() {
	case lexer.TokenKindInt:
		p.pushState(stateReservedFinish)
		p.pushState(stateReservedRange)
	case lexer.TokenKindStr:
		p.pushState(stateReservedFinish)
		p.pushState(stateReservedName)
	default:
		p.expectedCurr(lexer.TokenKindInt, lexer.TokenKindStr)
		tokIdx := p.skipPastLikelyEnd(p.currTok)
		p.addNode(tokIdx, stateStackEntry{
			tokIdx:       tokIdx,
			subtreeStart: uint32(len(p.tree)) - 1,
			hasError:     true,
		})
	}
}

func (p *Parser) parseReservedName() {
	p.popState()

	curr := p.curr()
	hasError := curr != lexer.TokenKindStr
	p.addLeafNode(hasError)

	if !hasError {
		curr = p.next()
	} else {
		p.popState()
		p.expectedCurr(lexer.TokenKindStr)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindSemicolon)
	}

	if curr == lexer.TokenKindComma {
		p.pushState(stateEnder)
		p.next()
		p.pushState(stateReservedName)
	}
}

func (p *Parser) parseReservedRange() {
	p.popState()

	curr := p.curr()
	hasError := curr != lexer.TokenKindInt && curr != lexer.TokenKindMax
	p.addLeafNode(hasError)

	if !hasError {
		curr = p.next()
	} else {
		p.popState()
		p.expectedCurr(lexer.TokenKindInt, lexer.TokenKindMax)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindSemicolon)
	}

	if curr == lexer.TokenKindTo {
		p.pushState(stateEnder)
		p.next()
		p.pushState(stateReservedRange)
	}
}

func (p *Parser) parseReservedFinish() {
	curr := p.curr()
	if curr == lexer.TokenKindComma {
		p.pushState(stateEnder)
		p.next()
		p.pushState(stateReservedRange)
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
