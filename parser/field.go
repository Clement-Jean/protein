package parser

import (
	"slices"

	"github.com/Clement-Jean/protein/lexer"
)

func (p *Parser) parseMessageFieldAssign() {
	state := p.popState()

	// unfortunately we have dynamic size introducers
	// so we need to calculate their len to later on
	// link against them
	introducerLen := uint32(len(p.tree)) - state.subtreeStart

	curr := p.curr()
	hasError := !curr.IsIdentifier()

	if hasError {
		p.popState()
		p.expectedCurr(lexer.TokenKindIdentifier)
		tokIdx := p.skipPastLikelyEnd(p.currTok)
		p.addNode(tokIdx, stateStackEntry{
			tokIdx:       tokIdx,
			subtreeStart: state.subtreeStart + 1,
			hasError:     true,
		})
		return
	}
	p.addLeafNode(hasError)
	curr = p.next()

	hasError = curr != lexer.TokenKindEqual
	equalTok := p.currTok

	if hasError {
		p.addLeafNode(true)
		p.popState()
		p.expectedCurr(lexer.TokenKindEqual)
		tokIdx := p.skipPastLikelyEnd(p.currTok)
		p.addNode(tokIdx, stateStackEntry{
			tokIdx:       tokIdx,
			subtreeStart: state.subtreeStart + 1,
			hasError:     false,
		})
		return
	}
	curr = p.next()

	hasError = curr != lexer.TokenKindInt
	p.addLeafNode(hasError)

	if !hasError {
		curr = p.next()
	} else {
		p.expectedCurr(lexer.TokenKindInt)
		p.skipPastLikelyEnd(p.currTok)
	}

	state.subtreeStart += introducerLen
	p.addNode(equalTok, state)

	if curr == lexer.TokenKindLeftSquare {
		p.addLeafNode(false)
		p.next()
		p.pushState(stateMessageFieldOption)
	}
}

func (p *Parser) parseMessageFieldOption() {
	curr := p.curr()
	switch curr {
	case lexer.TokenKindRightSquare:
		p.pushState(stateMessageFieldOptionFinish)
		return
	case lexer.TokenKindComma:
		p.pushState(stateEnder)
		curr = p.next()
	case lexer.TokenKindIdentifier:
	default:
		p.popState()
		p.expectedCurr(lexer.TokenKindRightSquare, lexer.TokenKindComma)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightSquare)
		return
	}

	hasError := curr != lexer.TokenKindIdentifier
	p.addLeafNode(hasError)

	if hasError {
		p.expectedCurr(lexer.TokenKindIdentifier)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightSquare)
		return
	}
	curr = p.next()

	hasError = curr != lexer.TokenKindEqual

	if !hasError {
		p.pushState(stateEnder)
		curr = p.next()
	} else {
		p.tree = append(p.tree, Node{
			TokIdx:      p.currTok,
			HasError:    hasError,
			SubtreeSize: 2,
		})
		p.expectedCurr(lexer.TokenKindEqual)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightSquare)
		return
	}

	hasError = !slices.Contains(constantTypes, curr)

	if !hasError {
		p.addLeafNode(false)
		p.next()
	} else {
		if curr == lexer.TokenKindLeftBrace || curr == lexer.TokenKindLeftAngle {
			p.addLeafNode(false)
			p.parseTextMessage()
			p.next()
			return
		}

		p.addLeafNode(true)
		p.expectedCurr(constantTypes...)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightSquare)
	}
}

func (p *Parser) parseMessageFieldOptionFinish() {
	state := p.popState()
	currTok := p.currTok

	state.hasError = p.curr() != lexer.TokenKindRightSquare

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindRightSquare)
		currTok = p.skipPastLikelyEnd(p.currTok)
	}

	top := p.popState() // stop field option loop
	p.addNode(currTok, top)
}

func (p *Parser) parseMessageFieldFinish() {
	state := p.popState()
	tokIdx := p.currTok

	state.hasError = p.curr() != lexer.TokenKindSemicolon

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindSemicolon)
		tokIdx = p.skipPastLikelyEnd(p.currTok)
	}

	state.subtreeStart++
	p.addNode(tokIdx, state)
}
