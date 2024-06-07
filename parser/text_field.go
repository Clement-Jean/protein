package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseTextField() {
	p.pushState(stateTextFieldAssign)
	p.pushState(stateTextFieldName)
}

func (p *Parser) parseTextFieldName() {
	p.popState()

	curr := p.curr()
	hasError := curr != lexer.TokenKindIdentifier && curr != lexer.TokenKindLeftSquare
	p.addLeafNode(hasError)

	switch curr {
	case lexer.TokenKindIdentifier:
		p.next()
	case lexer.TokenKindLeftSquare:
		p.next()
		p.pushState(stateTextFieldExtensionNameFinish)
		p.pushState(stateTextFieldExtensionName)
	default:
		if curr.IsIdentifier() {
			p.next()
			break
		}
		p.expectedCurr(lexer.TokenKindIdentifier, lexer.TokenKindLeftSquare)
	}
}

func (p *Parser) parseTextFieldExtensionName() {
	p.popState()
	p.pushState(stateFullIdentifierRoot)
}

func (p *Parser) parseTextFieldExtensionNameSlash() {
	state := p.popState()
	top := p.topState()
	top.subtreeStart++
	p.addNode(state.tokIdx, top)
}

func (p *Parser) parseTextFieldExtensionNameFinish() {
	if p.curr() == lexer.TokenKindSlash {
		p.pushState(stateTextFieldExtensionNameSlash)
		p.next()
		p.pushState(stateFullIdentifierRoot)
		return
	}

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

func (p *Parser) parseTextFieldAssign() {
	p.popState()

	curr := p.curr()
	if curr != lexer.TokenKindColon {
		if curr != lexer.TokenKindLeftBrace && curr != lexer.TokenKindLeftAngle {
			// TODO error
			panic("not implemented")
		}

		p.addLeafNode(false)
		p.parseTextMessage()
		p.next() // skip { or <
	} else {
		peek := p.peek()
		if peek != lexer.TokenKindLeftBrace && peek != lexer.TokenKindLeftAngle {
			p.pushState(stateTextFieldValue)
			p.next() // skip :
		} else {
			p.pushState(stateTextFieldColon)
			p.next() // skip :
			p.addLeafNode(false)
			p.parseTextMessage()
			p.next() // skip { or <
		}
	}
}

func (p *Parser) parseTextFieldColon() {
	state := p.popState()
	top := p.topState()
	top.subtreeStart++
	p.addNode(state.tokIdx, top)
}

func (p *Parser) parseTextFieldValue() {
	state := p.popState()

	switch curr := p.curr(); curr {
	case lexer.TokenKindInt:
		p.addLeafNode(false)
		p.next()
	case lexer.TokenKindFloat:
		p.addLeafNode(false)
		p.next()
	case lexer.TokenKindStr:
		p.addLeafNode(false)
		p.next()
	case lexer.TokenKindTrue, lexer.TokenKindFalse:
		p.addLeafNode(false)
		p.next()
	default:
		if curr.IsIdentifier() {
			p.addLeafNode(false)
			p.next()
			break
		}
		// TODO error
		panic("not implemented")
	}

	top := p.topState()
	top.subtreeStart++
	p.addNode(state.tokIdx, top)
}
