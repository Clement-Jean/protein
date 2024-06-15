package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseTextMessage() {
	if p.curr() == lexer.TokenKindLeftBrace {
		p.pushState(stateTextMessageFinishRightBrace)
	} else {
		p.pushState(stateTextMessageFinishRightAngle)
	}
	p.pushState(stateTextMessageValue)
}

func (p *Parser) parseTextMessageValue() {
	p.popState()

	curr := p.curr()
	if curr == lexer.TokenKindComment {
		curr = p.next()
	}

	if curr != lexer.TokenKindRightBrace && curr != lexer.TokenKindRightAngle {
		p.parseTextField()
	} else if curr.IsIdentifier() {
		p.parseTextField()
	}
}

func (p *Parser) parseTextMessageInsert() {
	p.popState()
	top := p.topState()
	top.subtreeStart++
	p.tree = append(p.tree, Node{
		TokIdx:      -1,
		SubtreeSize: int32(len(p.tree)) - top.subtreeStart + 1,
	})
}

func (p *Parser) parseTextMessageFinish() {
	curr := p.curr()
	if curr == lexer.TokenKindComma || curr == lexer.TokenKindSemicolon {
		p.pushState(stateEnder)
		p.next()
		p.pushState(stateTextMessageValue)
		return
	} else if curr.IsIdentifier() {
		p.pushState(stateTextMessageInsert)
		p.pushState(stateTextMessageValue)
		return
	}

	state := p.popState()
	tokIdx := p.currTok
	expected := lexer.TokenKindRightBrace
	if state.st == stateTextMessageFinishRightAngle {
		expected = lexer.TokenKindRightAngle
	}

	state.hasError = p.curr() != expected

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(expected)
		tokIdx = p.skipPastLikelyEnd(tokIdx)
		p.next()
	}

	p.addNode(tokIdx, state)
}
