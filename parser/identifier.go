package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseFullIdentifierRoot() {
	p.popState()

	hasError := p.curr() != lexer.TokenKindIdentifier
	p.addLeafNode(hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindIdentifier)
		p.skipPastLikelyEnd(p.currTok)
	}

	if state := p.stack[len(p.stack)-1]; state.st == stateFullIdentifierRest {
		// we are coming back from a dot
		p.popState()
		p.addNode(state.tokIdx, state)

		if state.hasError {
			return
		}
	}

	if p.curr() == lexer.TokenKindDot {
		p.pushState(stateFullIdentifierRest)
	}
}

func (p *Parser) parseFullIdentifierRest() {
	tok := p.currTok
	p.next()
	p.pushStateWithIdx(stateFullIdentifierRoot, tok)
	//	} else {
	//		p.popState()
	//		p.expectedCurr(lexer.TokenKindIdentifier)
	//		p.skipPastLikelyEnd(p.currTok)
	//	}
}
