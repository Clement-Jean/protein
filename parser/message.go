package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseMessage() {
	p.pushState(stateMessageFinish)
	p.pushState(stateMessageBlock)
	p.pushState(stateIdentifier)
}

func (p *Parser) parseMessageBlock() {
	p.popState()

	hasError := p.curr() != lexer.TokenKindLeftBrace
	p.addLeafNode(hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindLeftBrace)
		p.skipPastLikelyEnd(p.currTok)
	}

	p.pushState(stateMessageValue)
}

var messageScopeExpected = []lexer.TokenKind{
	lexer.TokenKindTypeFloat,
	lexer.TokenKindOption,
	lexer.TokenKindTypeDouble,
	lexer.TokenKindTypeInt32,
	lexer.TokenKindTypeInt64,
	lexer.TokenKindTypeUint32,
	lexer.TokenKindTypeUint64,
	lexer.TokenKindTypeSint32,
	lexer.TokenKindTypeSint64,
	lexer.TokenKindTypeFixed32,
	lexer.TokenKindTypeFixed64,
	lexer.TokenKindTypeSfixed32,
	lexer.TokenKindTypeSfixed64,
	lexer.TokenKindTypeBool,
	lexer.TokenKindTypeString,
	lexer.TokenKindTypeBytes,
	lexer.TokenKindMap,
	lexer.TokenKindIdentifier,
	lexer.TokenKindReserved,
	lexer.TokenKindExtensions,
	lexer.TokenKindOneOf,
	lexer.TokenKindMessage,
	lexer.TokenKindRightBrace,
}

func (p *Parser) parseMessageValue() {
	switch curr := p.curr(); curr {
	case lexer.TokenKindSemicolon, lexer.TokenKindComment:
		p.next()
	case lexer.TokenKindEOF, lexer.TokenKindRightBrace:
		p.popState()
	case lexer.TokenKindOption:
		p.addLeafNode(false)
		p.next()
		p.parseOption()
	case lexer.TokenKindReserved:
		p.addLeafNode(false)
		p.next()
		p.parseReserved()
	case lexer.TokenKindExtensions:
		p.addLeafNode(false)
		p.next()
		p.parseExtensions()
	case lexer.TokenKindOneOf:
		p.addTypedLeafNode(NodeKindMessageOneOfDecl, false)
		p.next()
		p.parseOneof()
	case lexer.TokenKindMap:
		p.pushState(stateMessageFieldFinish)
		p.pushState(stateMessageFieldAssign)
		p.parseMessageMap()
		p.addLeafNode(false)
		p.next()
	case lexer.TokenKindMessage:
		p.addTypedLeafNode(NodeKindMessageDecl, false)
		p.next()
		p.parseMessage()
	case lexer.TokenKindEnum:
		p.addTypedLeafNode(NodeKindEnumDecl, false)
		p.next()
		p.parseEnum()
	default:
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
			p.pushTypedState(NodeKindMessageFieldDecl, stateFullIdentifierRoot)

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

func (p *Parser) parseMessageFinish() {
	state := p.popState()
	tokIdx := p.currTok

	state.hasError = p.curr() != lexer.TokenKindRightBrace

	if !state.hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindRightBrace)
		tokIdx = p.skipPastLikelyEnd(tokIdx)
	}

	p.addTypedNode(tokIdx, NodeKindMessageClose, state)
}
