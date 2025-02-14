package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseIdentifier() {
	st := p.popState()

	curr := p.curr()
	hasError := !curr.IsIdentifier()
	p.addTypedLeafNode(st.kind, hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindIdentifier)
		p.skipPastLikelyEnd(p.currTok)
		p.popState()
	}
}

func (p *Parser) parseFullIdentifierRoot() {
	st := p.popState()

	curr := p.curr()
	hasError := !curr.IsIdentifier()
	p.addTypedLeafNode(st.kind, hasError)

	if !hasError {
		curr = p.next()
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

	if curr == lexer.TokenKindDot {
		p.pushState(stateFullIdentifierRest)
	}
}

func (p *Parser) parseFullIdentifierRest() {
	p.pushState(stateFullIdentifierRoot)
	p.next()
}
