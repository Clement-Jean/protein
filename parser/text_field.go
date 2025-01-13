package parser

import (
	"slices"

	"github.com/Clement-Jean/protein/lexer"
)

func (p *Parser) parseTextField() {
	p.pushState(stateTextFieldAssign)
	p.pushState(stateTextFieldName)
}

func (p *Parser) parseTextFieldName() {
	p.popState()

	curr := p.curr()
	hasError := curr != lexer.TokenKindIdentifier &&
		curr != lexer.TokenKindLeftSquare &&
		!curr.IsIdentifier()
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
		p.popState() // skip the whole field
		p.expectedCurr(lexer.TokenKindIdentifier, lexer.TokenKindLeftSquare)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindSemicolon, lexer.TokenKindRightBrace)
	}
}

func (p *Parser) parseTextFieldExtensionName() {
	p.popState()
	p.pushState(stateFullIdentifierRoot)
}

func (p *Parser) parseTextFieldExtensionNameFinish() {
	if p.curr() == lexer.TokenKindSlash {
		p.pushState(stateEnder)
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
		p.popState() // skip the whole field
		p.expectedCurr(lexer.TokenKindRightSquare)
		tokIdx = p.currTok
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindSemicolon, lexer.TokenKindRightBrace)
	}

	p.addNode(tokIdx, state)
}

func (p *Parser) parseTextFieldAssign() {
	p.popState()

	curr := p.curr()
	if curr != lexer.TokenKindColon {
		if curr != lexer.TokenKindLeftBrace && curr != lexer.TokenKindLeftAngle {
			p.addLeafNode(true)
			p.expectedCurr(lexer.TokenKindLeftBrace, lexer.TokenKindLeftAngle)
			p.skipTo(lexer.TokenKindComma, lexer.TokenKindSemicolon, lexer.TokenKindRightBrace)
			return
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
			p.pushState(stateEnder)
			p.next() // skip :
			p.addLeafNode(false)
			p.parseTextMessage()
			p.next() // skip { or <
		}
	}
}

func (p *Parser) parseTextFieldValue() {
	state := p.popState()
	curr := p.curr()
	isConstant := slices.Contains(constantTypes, curr)

	hasError := !isConstant && !curr.IsIdentifier()

	if !hasError {
		p.addLeafNode(false)
		p.next()
	} else if curr == lexer.TokenKindLeftSquare {
		p.addLeafNode(false)
		p.pushState(stateTextListFinish)
		p.pushState(stateTextListValue)
		p.next()
		return
	} else {
		p.addLeafNode(true)
		p.expectedCurr(append(constantTypes, lexer.TokenKindLeftSquare)...)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindSemicolon, lexer.TokenKindRightBrace)
	}

	top := p.topState()
	top.subtreeStart++
	p.addNode(state.tokIdx, top)
}
