package parser

import (
	"slices"

	"github.com/Clement-Jean/protein/lexer"
)

func (p *Parser) parseMessage() {
	p.pushState(stateMessageFinish)
	p.pushState(stateMessageBlock)
	p.pushState(stateMessageName)
}

func (p *Parser) parseMessageName() {
	p.popState()

	curr := p.curr()
	hasError := !curr.IsIdentifier()
	p.addLeafNode(hasError)

	if !hasError {
		p.next()
	} else {
		p.expectedCurr(lexer.TokenKindIdentifier)
		p.skipPastLikelyEnd(p.currTok)
	}
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

func (p *Parser) parseMessageFieldAssign() {
	state := p.popState()

	// unfortunately we have dynamic size introducers
	// so we need to calculate their len to later on
	// link against them
	introducerLen := int32(len(p.tree)) - state.subtreeStart

	curr := p.curr()
	hasError := !curr.IsIdentifier()

	if !hasError {
		p.addLeafNode(hasError)
		curr = p.next()
	} else {
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

	hasError = curr != lexer.TokenKindEqual
	equalTok := p.currTok

	if !hasError {
		curr = p.next()
	} else {
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

	if !hasError {
		curr = p.next()
	} else {
		p.expectedCurr(lexer.TokenKindIdentifier)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightSquare)
		return
	}

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
		} else {
			p.addLeafNode(true)
			p.expectedCurr(constantTypes...)
			p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightSquare)
		}
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

func (p *Parser) parseMessageMapKeyValue() {
	state := p.popState()
	curr := p.curr()
	hasError := curr != lexer.TokenKindLeftAngle
	p.addLeafNode(hasError)

	if !hasError {
		curr = p.next()
	} else {
		p.expectedCurr(lexer.TokenKindLeftAngle)
		p.skipPastLikelyEnd(p.currTok)
		return
	}

	hasError = !slices.Contains(mapKeyTypes, curr)
	p.addLeafNode(hasError)

	if !hasError {
		curr = p.next()
	} else {
		p.expectedCurr(mapKeyTypes...)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightAngle)
		curr = p.curr()
	}

	hasError = curr != lexer.TokenKindComma
	commaIdx := p.currTok

	if !hasError {
		curr = p.next()
	} else {
		p.addLeafNode(hasError)
		p.expectedCurr(lexer.TokenKindComma)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightAngle)
		curr = p.curr()
		state.subtreeStart += 3
		goto end_generic
	}

	hasError = !curr.IsIdentifier()
	p.addLeafNode(hasError)

	if !hasError {
		curr = p.next()
	} else {
		p.expectedCurr(lexer.TokenKindIdentifier)
		p.skipTo(lexer.TokenKindComma, lexer.TokenKindRightAngle)
		curr = p.curr()
	}

	state.subtreeStart += 3
	p.addNode(commaIdx, state)

end_generic:
	hasError = curr != lexer.TokenKindRightAngle

	if !hasError {
		state.subtreeStart--
		p.addNode(p.currTok, state)
		p.next()
	} else {
		p.addLeafNode(hasError)
		p.expectedCurr(lexer.TokenKindRightAngle)
		p.skipPastLikelyEnd(p.currTok)
	}
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
	case lexer.TokenKindMap:
		p.pushState(stateMessageFieldFinish)
		p.pushState(stateMessageFieldAssign)
		p.pushState(stateMessageMapKeyValue)
		p.addLeafNode(false)
		p.next()
	default:
		hasDot := false
		hasModifier := false
		modifierIdx := -1
		if curr == lexer.TokenKindOptional || curr == lexer.TokenKindRepeated {
			hasModifier = true
			modifierIdx = p.currTok
			curr = p.next()
		}

		if curr == lexer.TokenKindDot {
			hasDot = true
			curr = p.next()
		}

		if curr.IsIdentifier() {
			p.pushState(stateMessageFieldFinish)
			p.pushState(stateMessageFieldAssign)
			p.pushState(stateFullIdentifierRoot)
			if hasModifier {
				p.addNode(modifierIdx, stateStackEntry{
					tokIdx:       modifierIdx,
					subtreeStart: int32(len(p.tree)),
				})
			}
			if hasDot {
				p.addLeafNode(false)
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
					subtreeStart: int32(len(p.tree)),
				})
			}
			nbElements := p.currTok - modifierIdx + 1
			p.expectedCurr(messageScopeExpected...)
			p.skipPastLikelyEnd(p.currTok)

			// after skip, we can now add the token
			// we skipped to
			p.addNode(p.currTok, stateStackEntry{
				tokIdx:       p.currTok,
				subtreeStart: int32(len(p.tree) - nbElements),
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

	p.addNode(tokIdx, state)
}
