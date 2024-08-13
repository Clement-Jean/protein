package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseOneof() {
	p.pushState(stateOneofFinish)
	p.pushState(stateOneofBlock)
	p.pushState(stateIdentifier)
}

func (p *Parser) parseOneofBlock() {
	p.popState()

	hasError := p.curr() != lexer.TokenKindLeftBrace
	p.addLeafNode(hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindLeftBrace)
		p.skipPastLikelyEnd(p.currTok)
	}

	p.pushState(stateOneofValue)
}

var oneofScopeExpected = []lexer.TokenKind{
	lexer.TokenKindOption,
	lexer.TokenKindIdentifier,
	lexer.TokenKindRightBrace,
}

func (p *Parser) parseOneofValue() {
	switch curr := p.curr(); curr {
	case lexer.TokenKindSemicolon, lexer.TokenKindComment:
		p.next()
	case lexer.TokenKindEOF, lexer.TokenKindRightBrace:
		p.popState()
	case lexer.TokenKindOption:
		p.addLeafNode(false)
		p.next()
		p.parseOption()
	default:
		hasDot := false
		if curr == lexer.TokenKindDot {
			hasDot = true
			curr = p.next()
		}

		if curr.IsIdentifier() {
			p.pushState(stateMessageFieldFinish)
			p.pushState(stateMessageFieldAssign)
			p.pushState(stateFullIdentifierRoot)
			if hasDot {
				p.addLeafNode(false)
			}
			break
		}

		p.expectedCurr(messageScopeExpected...)
		p.skipPastLikelyEnd(p.currTok)
	}
}

func (p *Parser) parseOneofFinish() {
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
