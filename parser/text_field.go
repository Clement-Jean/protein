package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) parseTextField() {
	p.pushState(stateTextFieldAssign)
	p.pushState(stateTextFieldName)
}

func (p *Parser) parseTextFieldName() {
	p.popState()

	hasError := p.curr() != lexer.TokenKindIdentifier && p.curr() != lexer.TokenKindLeftAngle
	p.addLeafNode(hasError)

	switch p.curr() {
	case lexer.TokenKindIdentifier:
		p.next()
	case lexer.TokenKindLeftAngle:
		// TODO pushState(stateTextFieldExtensionNameFinish)
		//      pushState(stateTextFieldExtensionName)
		//      next()
	default:
		p.expectedCurr(lexer.TokenKindIdentifier, lexer.TokenKindLeftSquare)
	}
}

func (p *Parser) parseTextFieldExtensionName() {
	panic("not implemented")
}

func (p *Parser) parseTextFieldExtensionNameFinish() {
	panic("not implemented")
}

func (p *Parser) parseTextFieldAssign() {
	p.popState()

	if p.peek() != lexer.TokenKindLeftBrace && p.peek() != lexer.TokenKindLeftAngle {
		if p.curr() != lexer.TokenKindColon {
			// TODO error
			panic("not implemented")
		}
		p.pushState(stateTextFieldValue)
		p.next()
	} else {
		// here ':' is optional, however we insert it anyway as node
	}
}

func (p *Parser) parseTextFieldValue() {
	state := p.popState()

	switch curr := p.curr(); curr {
	case lexer.TokenKindInt:
		p.addLeafNode(false)
		p.next()
	case lexer.TokenKindFloat:
		p.addLeafNode(false)
		p.next()
	case lexer.TokenKindStr:
		p.addLeafNode(false)
		p.next()
	case lexer.TokenKindTrue, lexer.TokenKindFalse:
		p.addLeafNode(false)
		p.next()
	case lexer.TokenKindLeftBrace, lexer.TokenKindLeftAngle:
		p.addLeafNode(false)
		p.next()
		p.parseTextMessage()
	default:
		if curr.IsIdentifier() {
			p.addLeafNode(false)
			p.next()
			break
		}
		// TODO error
		panic("not implemented")
	}

	p.addNode(state.tokIdx, state)
}
