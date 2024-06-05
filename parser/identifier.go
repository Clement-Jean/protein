package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseFullIdentifierRoot() {
	p.popState()

	hasError := p.curr() != lexer.TokenKindIdentifier && !p.curr().IsIdentifier()
	p.addLeafNode(hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindIdentifier)
	}

	if state := p.topState(); state.st == stateFullIdentifierRest {
		// we are coming back from a dot
		p.popState()
		top := p.topState()
		top.subtreeStart++
		p.addNode(state.tokIdx, top)

		if state.hasError {
			return
		}
	}

	if p.curr() == lexer.TokenKindDot {
		p.pushState(stateFullIdentifierRest)
	}
}

func (p *Parser) parseFullIdentifierRest() {
	p.pushState(stateFullIdentifierRoot)
	p.next()
}
