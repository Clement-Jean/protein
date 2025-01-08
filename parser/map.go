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

func (p *Parser) parseMessageMapKeyValue() {
	state := p.popState()
	curr := p.curr()
	hasError := curr != lexer.TokenKindLeftAngle
	p.addLeafNode(hasError)

	if hasError {
		p.expectedCurr(lexer.TokenKindLeftAngle)
		p.skipPastLikelyEnd(p.currTok)
		return
	}
	curr = p.next()

	hasError = !slices.Contains(mapKeyTypes, curr)
	p.addLeafNode(hasError)

	if !hasError {
		curr = p.next()
	} else {
		p.expectedCurr(mapKeyTypes...)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightAngle)
		curr = p.curr()
	}

	hasError = curr != lexer.TokenKindComma
	commaIdx := p.currTok

	if hasError {
		p.addLeafNode(hasError)
		p.expectedCurr(lexer.TokenKindComma)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightAngle)
		curr = p.curr()
		state.subtreeStart += 3
		goto end_generic
	}
	curr = p.next()

	hasError = !curr.IsIdentifier()
	p.addLeafNode(hasError)

	if !hasError {
		curr = p.next()
	} else {
		p.expectedCurr(lexer.TokenKindIdentifier)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightAngle)
		curr = p.curr()
	}

	state.subtreeStart += 3
	p.addNode(commaIdx, state)

end_generic:
	hasError = curr != lexer.TokenKindRightAngle

	if !hasError {
		state.subtreeStart--
		p.addNode(p.currTok, state)
		p.next()
	} else {
		p.addLeafNode(hasError)
		p.expectedCurr(lexer.TokenKindRightAngle)
		p.skipPastLikelyEnd(p.currTok)
	}
}
