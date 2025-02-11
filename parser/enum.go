package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseEnum() {
	p.pushState(stateEnumFinish)
	p.pushState(stateEnumBlock)
	p.pushState(stateIdentifier)
}

func (p *Parser) parseEnumBlock() {
	p.popState()

	hasError := p.curr() != lexer.TokenKindLeftBrace
	p.addLeafNode(hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindLeftBrace)
		p.skipPastLikelyEnd(p.currTok)
	}

	p.pushState(stateEnumValue)
}

var enumScopeExpected = []lexer.TokenKind{
	lexer.TokenKindOption,
	lexer.TokenKindReserved,
	lexer.TokenKindIdentifier,
	lexer.TokenKindRightBrace,
}

func (p *Parser) parseEnumValue() {
	switch curr := p.curr(); curr {
	case lexer.TokenKindSemicolon, lexer.TokenKindComment:
		p.next()
	case lexer.TokenKindEOF, lexer.TokenKindRightBrace:
		p.popState()
	case lexer.TokenKindOption:
		p.addLeafNode(false)
		p.next()
		p.parseOption()
	case lexer.TokenKindReserved:
		p.addLeafNode(false)
		p.next()
		p.parseReserved()
	default:
		if curr.IsIdentifier() {
			p.pushState(stateMessageFieldFinish)
			p.pushState(stateMessageFieldAssign)
			break
		}
		p.expectedCurr(enumScopeExpected...)
		p.skipPastLikelyEnd(p.currTok)
	}
}

func (p *Parser) parseEnumFinish() {
	state := p.popState()
	tokIdx := p.currTok

	state.hasError = p.curr() != lexer.TokenKindRightBrace

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindRightBrace)
		tokIdx = p.skipPastLikelyEnd(tokIdx)
	}

	p.addTypedNode(tokIdx, NodeKindMessageClose, state)
}
