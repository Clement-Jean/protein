package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseService() {
	p.pushState(stateServiceFinish)
	p.pushState(stateServiceBlock)
	p.pushState(stateIdentifier)
}

func (p *Parser) parseServiceBlock() {
	p.popState()

	hasError := p.curr() != lexer.TokenKindLeftBrace
	p.addLeafNode(hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindLeftBrace)
		p.skipPastLikelyEnd(p.currTok)
	}

	p.pushState(stateServiceValue)
}

var serviceScopeExpected = []lexer.TokenKind{
	lexer.TokenKindOption,
	lexer.TokenKindRPC,
}

func (p *Parser) parseServiceValue() {
	switch curr := p.curr(); curr {
	case lexer.TokenKindSemicolon, lexer.TokenKindComment:
		p.next()
	case lexer.TokenKindEOF, lexer.TokenKindRightBrace:
		p.popState()
	case lexer.TokenKindOption:
		p.addLeafNode(false)
		p.next()
		p.parseOption()
	case lexer.TokenKindRPC:
		p.addLeafNode(false)
		p.pushState(stateRPCFinish) // optional semicolon
		p.pushState(stateRPCFinish)
		p.next()
		p.pushState(stateRPCDefinition)
	default:
		p.expectedCurr(serviceScopeExpected...)
		p.skipPastLikelyEnd(p.currTok)
	}
}

func (p *Parser) parseServiceFinish() {
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
