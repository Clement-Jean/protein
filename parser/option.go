package parser

import (
	"slices"

	"github.com/Clement-Jean/protein/lexer"
)

func (p *Parser) parseOption() {
	p.pushState(stateOptionFinish)
	p.pushState(stateOptionAssign)
}

func (p *Parser) parseOptionAssign() {
	p.popState()

	// TODO option names can be more complex
	hasError := p.curr() != lexer.TokenKindIdentifier
	p.addLeafNode(hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindIdentifier)
		p.skipPastLikelyEnd(p.currTok)
		return
	}

	hasError = p.curr() != lexer.TokenKindEqual
	p.addLeafNode(hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindEqual)
		p.skipPastLikelyEnd(p.currTok)
		return
	}

	accepted := []lexer.TokenKind{
		lexer.TokenKindTrue, lexer.TokenKindFalse,
		lexer.TokenKindInt, lexer.TokenKindFloat,
		lexer.TokenKindStr,
	}
	hasError = !slices.Contains(accepted, p.curr())
	p.addLeafNode(hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(accepted...)
		p.skipPastLikelyEnd(p.currTok)
		// TODO check for text message
	}
}

func (p *Parser) parseOptionFinish() {
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
