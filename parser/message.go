package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseMessage() {
	p.pushState(stateMessageFinish)
	p.pushState(stateMessageBlock)
	p.pushState(stateMessageName)
}

func (p *Parser) parseMessageName() {
	p.popState()

	curr := p.curr()
	hasError := curr != lexer.TokenKindIdentifier && !curr.IsIdentifier()
	p.addLeafNode(hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindEqual)
		p.skipPastLikelyEnd(p.currTok)
	}
}

func (p *Parser) parseMessageBlock() {
	p.popState()

	hasError := p.curr() != lexer.TokenKindLeftBrace
	p.addLeafNode(hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindLeftBrace)
		p.skipPastLikelyEnd(p.currTok)
	}

	p.pushState(stateMessageValue)
}

func (p *Parser) parseMessageField() {
	state := p.popState()
	curr := p.curr()
	hasError := curr != lexer.TokenKindIdentifier
	p.addLeafNode(hasError)

	if !hasError {
		curr = p.next()
	} else {
		p.expectedCurr(lexer.TokenKindIdentifier)
		p.skipPastLikelyEnd(p.currTok)
		return
	}

	hasError = curr != lexer.TokenKindEqual
	equalTok := p.currTok

	if !hasError {
		curr = p.next()
	} else {
		p.addLeafNode(hasError)
		p.expectedCurr(lexer.TokenKindEqual)
		p.skipPastLikelyEnd(p.currTok)
		return
	}

	hasError = curr != lexer.TokenKindInt
	p.addLeafNode(hasError)

	if !hasError {
		curr = p.next()
	} else {
		p.expectedCurr(lexer.TokenKindInt)
		p.skipPastLikelyEnd(p.currTok)
		return
	}

	state.subtreeStart++
	p.addNode(equalTok, state)

	state.subtreeStart--
	state.hasError = curr != lexer.TokenKindSemicolon
	p.addNode(p.currTok, state)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindSemicolon)
		p.skipPastLikelyEnd(p.currTok)
	}
}

func (p *Parser) parseMessageValue() {
	switch p.curr() {
	case lexer.TokenKindOption:
		p.addLeafNode(false)
		p.next()
		p.parseOption()
	case lexer.TokenKindTypeFloat,
		lexer.TokenKindTypeDouble,
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
		lexer.TokenKindTypeBytes:
		p.addLeafNode(false)
		p.next()
		p.pushState(stateMessageField)
	case lexer.TokenKindSemicolon:
		p.next()
	case lexer.TokenKindRightBrace:
		p.popState()
	default:
		panic("not implemented")
	}
}

func (p *Parser) parseMessageFinish() {
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
