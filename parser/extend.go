package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseExtend() {
	p.pushState(stateExtendFinish)
	p.pushState(stateExtendBlock)
	p.pushTypedState(NodeKindExtendDecl, stateFullIdentifierRoot)

	if p.curr() == lexer.TokenKindDot {
		p.addLeafNode(false)
		p.next()
	}
}

func (p *Parser) parseExtendBlock() {
	p.popState()

	hasError := p.curr() != lexer.TokenKindLeftBrace
	p.addLeafNode(hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindLeftBrace)
		p.skipPastLikelyEnd(p.currTok)
	}

	p.pushState(stateExtendField)
}

func (p *Parser) parseExtendField() {
	switch curr := p.curr(); curr {
	case lexer.TokenKindSemicolon, lexer.TokenKindComment:
		p.next()
	case lexer.TokenKindEOF, lexer.TokenKindRightBrace:
		p.popState()
	case lexer.TokenKindMap:
		p.pushState(stateMessageFieldFinish)
		p.pushTypedState(NodeKindExtendMapDecl, stateMessageFieldAssign)
		p.parseMessageMap()
		p.addLeafNode(false)
		p.next()
	default:
		// TODO same as message parsing field
		//      put that into a function?
		var (
			dotIdx      uint32
			modifierIdx uint32
		)
		hasDot := false
		hasModifier := false

		if curr == lexer.TokenKindOptional || curr == lexer.TokenKindRepeated || curr == lexer.TokenKindRequired {
			hasModifier = true
			modifierIdx = p.currTok
			curr = p.next()
		}

		if curr == lexer.TokenKindDot {
			hasDot = true
			dotIdx = p.currTok
			curr = p.next()
		}

		if curr.IsIdentifier() {
			p.pushState(stateMessageFieldFinish)
			p.pushState(stateMessageFieldAssign)
			p.pushTypedState(NodeKindExtendFieldDecl, stateFullIdentifierRoot)

			if hasModifier {
				p.addNode(modifierIdx, stateStackEntry{
					tokIdx:       modifierIdx,
					subtreeStart: uint32(len(p.tree)),
				})
			}

			if hasDot {
				p.addNode(dotIdx, stateStackEntry{
					tokIdx:       dotIdx,
					subtreeStart: uint32(len(p.tree)),
				})
			}
			break
		} else if hasModifier {
			// we try to create a coherent parse tree
			// even though we know there is an error

			// add all the tokens between modifierIdx
			// and currTok
			for i := modifierIdx; i <= p.currTok; i++ {
				p.addNode(i, stateStackEntry{
					tokIdx:       i,
					subtreeStart: uint32(len(p.tree)),
				})
			}
			nbElements := p.currTok - modifierIdx + 1
			p.expectedCurr(messageScopeExpected...)
			p.skipPastLikelyEnd(p.currTok)

			// after skip, we can now add the token
			// we skipped to
			p.addNode(p.currTok, stateStackEntry{
				tokIdx:       p.currTok,
				subtreeStart: uint32(len(p.tree)) - nbElements,
				hasError:     true,
			})
			break
		}
		p.expectedCurr(messageScopeExpected...)
		p.skipPastLikelyEnd(p.currTok)
	}
}

func (p *Parser) parseExtendFinish() {
	state := p.popState()
	tokIdx := p.currTok

	state.hasError = p.curr() != lexer.TokenKindRightBrace

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindRightBrace)
		tokIdx = p.skipPastLikelyEnd(tokIdx)
	}

	p.addTypedNode(tokIdx, NodeKindExtendClose, state)
}
