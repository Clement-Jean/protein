package parser

import (
	"math"

	"github.com/Clement-Jean/protein/lexer"
)

func (p *Parser) parseRPCDefinition() {
	p.popState()
	p.pushState(stateRPCReqRes)
	p.pushState(stateIdentifier)
}

func (p *Parser) parseRPCReqRes() {
	state := p.popState()

	curr := p.curr()
	hasError := curr != lexer.TokenKindLeftParen

	if hasError {
		p.popState()
		p.expectedCurr(lexer.TokenKindLeftParen)
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

	p.pushState(stateRPCReqResFinish)

	if curr == lexer.TokenKindDot {
		p.addLeafNode(false)
		curr = p.next()
	}

	p.pushTypedState(NodeKindRPCInputOutput, stateFullIdentifierRoot)

	if curr == lexer.TokenKindStream {
		p.addLeafNode(hasError)
		p.next()
	}
}

func (p *Parser) parseRPCReqResFinish() {
	state := p.popState()
	curr := p.curr()
	hasError := curr != lexer.TokenKindRightParen

	if hasError {
		p.popState()
		p.expectedCurr(lexer.TokenKindRightParen)
		tokIdx := p.skipPastLikelyEnd(p.currTok)
		p.addNode(tokIdx, stateStackEntry{
			tokIdx:       tokIdx,
			subtreeStart: state.subtreeStart + 1,
			hasError:     true,
		})
		return
	}
	p.addNode(p.currTok, state)
	curr = p.next()

	if p.topState().st == stateEnder { // coming back from stateReqResFinish
		return
	}

	p.pushState(stateEnder)

	hasError = curr != lexer.TokenKindReturns

	if hasError {
		p.popState()
		p.expectedCurr(lexer.TokenKindReturns)
		tokIdx := p.skipPastLikelyEnd(p.currTok)
		p.addNode(tokIdx, stateStackEntry{
			tokIdx:       tokIdx,
			subtreeStart: state.subtreeStart + 1,
			hasError:     true,
		})
		return
	}
	p.next()

	p.pushState(stateRPCReqRes)
}

func (p *Parser) parseRPCFinish() {
	state := p.popState()
	curr := p.curr()

	if p.topState().st == stateServiceValue { // coming back after rpc option
		if curr == lexer.TokenKindSemicolon {
			p.addNode(state.tokIdx, state)
			p.next()
			return
		}

		p.tree = append(p.tree, Node{
			TokIdx:      math.MaxUint32,
			SubtreeSize: uint32(len(p.tree)) - state.subtreeStart + 1,
		})
		return
	}

	hasError := curr != lexer.TokenKindSemicolon &&
		curr != lexer.TokenKindLeftBrace

	if hasError {
		p.popState()
		p.expectedCurr(lexer.TokenKindSemicolon, lexer.TokenKindLeftBrace)
		tokIdx := p.skipPastLikelyEnd(p.currTok)
		p.addNode(tokIdx, stateStackEntry{
			tokIdx:       tokIdx,
			subtreeStart: state.subtreeStart + 1,
			hasError:     true,
		})
		return
	}

	switch curr {
	case lexer.TokenKindSemicolon:
		if p.topState().st == stateRPCFinish {
			p.popState() // remove need for checking optional semicolon
		}
		p.addNode(p.currTok, state)
		p.next()
	case lexer.TokenKindLeftBrace:
		p.addLeafNode(false)
		p.next()
		p.pushState(stateMessageFinish) // HACK: for now using MessageFinish
		p.pushState(stateRPCValue)
	}
}

func (p *Parser) parseRPCValue() {
	switch p.curr() {
	case lexer.TokenKindSemicolon, lexer.TokenKindComment:
		p.next()
	case lexer.TokenKindEOF, lexer.TokenKindRightBrace:
		p.popState()
	case lexer.TokenKindOption:
		p.addLeafNode(false)
		p.next()
		p.parseOption()
	default:
		p.expectedCurr(lexer.TokenKindOption)
		p.skipPastLikelyEnd(p.currTok)
	}
}
