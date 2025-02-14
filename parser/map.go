package parser

import (
	"slices"

	"github.com/Clement-Jean/protein/lexer"
)

var mapKeyTypes = []lexer.TokenKind{
	lexer.TokenKindTypeInt32,
	lexer.TokenKindTypeInt64,
	lexer.TokenKindTypeUint32,
	lexer.TokenKindTypeUint64,
	lexer.TokenKindTypeSint32,
	lexer.TokenKindTypeSint64,
	lexer.TokenKindTypeFixed32,
	lexer.TokenKindTypeFixed64,
	lexer.TokenKindTypeSfixed32,
	lexer.TokenKindTypeSfixed64,
	lexer.TokenKindTypeBool,
	lexer.TokenKindTypeString,
}

func (p *Parser) parseMessageMap() {
	p.pushState(stateMessageMapFinish)
	p.pushState(stateMessageMapStart)
}

func (p *Parser) parseMessageMapStart() {
	p.popState()
	hasError := p.curr() != lexer.TokenKindLeftAngle
	p.addLeafNode(hasError)

	if hasError {
		p.expectedCurr(lexer.TokenKindLeftAngle)
		p.skipPastLikelyEnd(p.currTok)
		return
	}
	p.next()
	p.pushState(stateMessageMapKeyValue)
}

func (p *Parser) parseMessageMapKeyValue() {
	p.popState()
	curr := p.curr()

	hasError := !slices.Contains(mapKeyTypes, curr)
	p.addLeafNode(hasError)

	if !hasError {
		curr = p.next()
	} else {
		p.expectedCurr(mapKeyTypes...)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightAngle)
		curr = p.curr()
	}

	hasError = curr != lexer.TokenKindComma

	if hasError {
		p.addLeafNode(hasError)
		p.expectedCurr(lexer.TokenKindComma)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightAngle)
		return
	}
	p.pushState(stateMessageMapComma)
	curr = p.next()

	hasError = !curr.IsIdentifier()

	if hasError {
		p.addLeafNode(hasError)
		p.expectedCurr(lexer.TokenKindIdentifier)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightAngle)
		return
	}
	p.pushState(stateFullIdentifierRoot)
}

func (p *Parser) parseMessageMapComma() {
	state := p.popState()
	p.addNode(state.tokIdx, state)
}

func (p *Parser) parseMessageMapFinish() {
	state := p.popState()
	hasError := p.curr() != lexer.TokenKindRightAngle

	if hasError {
		p.addLeafNode(true)
		p.expectedCurr(lexer.TokenKindRightAngle)
		p.skipPastLikelyEnd(p.currTok)
		return
	}

	state.subtreeStart += 2 // map <
	p.addNode(p.currTok, state)
	p.next()
}
