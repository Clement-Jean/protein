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

	// the name of a field can be a keyword (e.g) message
	// so we override the kind to be an identifier in order
	// to make sure we don't treat this as something else
	// than an identifier
	p.toks.TokenInfos[p.currTok].Kind = lexer.TokenKindIdentifier

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
	p.popState()

	curr := p.curr()
	hasError := curr != lexer.TokenKindIdentifier && curr != lexer.TokenKindLeftParen

	if !hasError {
		p.pushState(stateMessageFieldOptionFinish)
		p.pushState(stateMessageFieldOptionAssign)
		p.pushState(stateOptionName)
	} else {
		p.popState()
		p.expectedCurr(lexer.TokenKindIdentifier, lexer.TokenKindLeftParen)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightSquare)
	}
}

func (p *Parser) parseMessageFieldOptionAssign() {
	p.popState()

	curr := p.curr()
	hasError := curr != lexer.TokenKindEqual

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
		return
	}
}

func (p *Parser) parseMessageFieldOptionFinish() {
	curr := p.curr()
	if curr == lexer.TokenKindComma {
		p.pushState(stateEnder)
		p.next()
		p.pushState(stateMessageFieldOptionAssign)
		p.pushState(stateOptionName)
		return
	}

	state := p.popState()
	currTok := p.currTok

	state.hasError = p.curr() != lexer.TokenKindRightSquare

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindRightSquare)
		currTok = p.skipPastLikelyEnd(p.currTok)
	}

	p.addNode(currTok, state)
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
