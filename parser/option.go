package parser

import (
	"slices"

	"github.com/Clement-Jean/protein/lexer"
)

func (p *Parser) parseOption() {
	p.pushState(stateOptionFinish)
	p.pushState(stateOptionAssign)
	p.pushState(stateOptionName)
}

func (p *Parser) parseOptionName() {
	p.popState()

	curr := p.curr()
	hasError := curr != lexer.TokenKindIdentifier &&
		curr != lexer.TokenKindLeftParen &&
		!curr.IsIdentifier()
	p.addLeafNode(hasError)

	switch curr {
	case lexer.TokenKindIdentifier:
		p.next()
	case lexer.TokenKindLeftParen:
		curr = p.next()
		p.pushState(stateOptionNameParenFinish)
		p.pushState(stateFullIdentifierRoot)

		if curr == lexer.TokenKindDot {
			p.addLeafNode(false)
			p.next()
		}
	default:
		if curr.IsIdentifier() {
			p.next()
			break
		}

		p.expectedCurr(lexer.TokenKindIdentifier, lexer.TokenKindLeftParen)
	}

	if optionNameRest := p.topState(); optionNameRest.st == stateOptionNameRest {
		// we are coming back from a dot
		p.popState()
		optionNameState := p.topState()
		optionNameState.subtreeStart++
		p.addNode(optionNameRest.tokIdx, optionNameState)

		if optionNameRest.hasError {
			return
		}
	}

	if p.curr() == lexer.TokenKindDot {
		p.pushState(stateOptionNameRest)
	}
}

func (p *Parser) parseOptionNameRest() {
	p.pushState(stateOptionName)
	p.next()
}

func (p *Parser) parseOptionNameParenFinish() {
	state := p.popState()
	tokIdx := p.currTok

	state.hasError = p.curr() != lexer.TokenKindRightParen

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindRightParen)
		tokIdx = p.skipPastLikelyEnd(tokIdx)
		state.subtreeStart--
	}

	p.addNode(tokIdx, state)

	if optionNameRest := p.topState(); optionNameRest.st == stateOptionNameRest {
		// we are coming back from a dot
		p.popState()
		optionNameState := p.topState()
		optionNameState.subtreeStart++
		p.addNode(optionNameRest.tokIdx, optionNameState)

		if optionNameRest.hasError {
			return
		}
	}

	if p.curr() == lexer.TokenKindDot {
		p.pushState(stateOptionNameRest)
	}
}

func (p *Parser) parseOptionAssign() {
	p.popState()
	curr := p.curr()
	hasError := curr != lexer.TokenKindEqual

	if hasError {
		p.addLeafNode(hasError)
		p.expectedCurr(lexer.TokenKindEqual)
		p.skipPastLikelyEnd(p.currTok)
		return
	}
	p.pushState(stateEnder) // add the = as a node containing the name and value as subtree
	curr = p.next()

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
		p.skipPastLikelyEnd(p.currTok)
	}
}

func (p *Parser) parseOptionFinish() {
	state := p.popState()
	tokIdx := p.currTok

	state.hasError = p.curr() != lexer.TokenKindSemicolon

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindSemicolon)
		tokIdx = p.skipPastLikelyEnd(tokIdx)
	}

	p.addNode(tokIdx, state)
}
