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

	hasError := p.curr() != lexer.TokenKindIdentifier &&
		p.curr() != lexer.TokenKindLeftParen &&
		!p.curr().IsIdentifier()
	p.addLeafNode(hasError)

	switch curr := p.curr(); curr {
	case lexer.TokenKindIdentifier:
		p.next()
	case lexer.TokenKindLeftParen:
		p.next()
		p.pushState(stateOptionNameParenFinish)
		p.pushState(stateFullIdentifierRoot)

		if p.curr() == lexer.TokenKindDot {
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
	state := p.popState()
	tok := p.currTok
	hasError := p.curr() != lexer.TokenKindEqual

	if !hasError {
		p.next()
	} else {
		p.addLeafNode(hasError)
		p.expectedCurr(lexer.TokenKindEqual)
		p.skipPastLikelyEnd(p.currTok)
		return
	}

	accepted := []lexer.TokenKind{
		lexer.TokenKindTrue, lexer.TokenKindFalse,
		lexer.TokenKindInt, lexer.TokenKindFloat,
		lexer.TokenKindStr,
	}
	hasError = !slices.Contains(accepted, p.curr())

	if !hasError {
		p.addLeafNode(false)
		p.next()
	} else {
		if p.curr() == lexer.TokenKindLeftBrace || p.curr() == lexer.TokenKindLeftAngle {
			p.addLeafNode(false)
			p.next()
			p.parseTextMessage()
			return
		} else {
			p.addLeafNode(true)
			p.expectedCurr(accepted...)
			p.skipPastLikelyEnd(p.currTok)
		}
	}

	state.subtreeStart++
	p.addNode(tok, state)
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
