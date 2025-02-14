package parser

import (
	"github.com/Clement-Jean/protein/lexer"
)

type Parser struct {
	toks    *lexer.TokenizedBuffer
	tree    ParseTree
	stack   []stateStackEntry
	errs    []error
	currTok uint32
}

func New(toks *lexer.TokenizedBuffer) *Parser {
	if len(toks.TokenInfos) <= 0 {
		panic("we need at least one token to be able to create a parser")
	}

	return &Parser{
		toks: toks,
		tree: make([]Node, 0, len(toks.TokenInfos)),
	}
}

func (p *Parser) pushTypedState(kind NodeKind, st state) {
	p.stack = append(p.stack, stateStackEntry{
		st:           st,
		kind:         kind,
		tokIdx:       p.currTok,
		subtreeStart: uint32(len(p.tree)) - 1,
	})
}

func (p *Parser) pushState(st state) {
	p.stack = append(p.stack, stateStackEntry{
		st:           st,
		tokIdx:       p.currTok,
		subtreeStart: uint32(len(p.tree)) - 1,
	})
}

func (p *Parser) popState() stateStackEntry {
	state := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return state
}

func (p *Parser) topState() stateStackEntry {
	return p.stack[len(p.stack)-1]
}

func (p *Parser) next() lexer.TokenKind {
	if p.currTok+1 >= uint32(len(p.toks.TokenInfos)) {
		return lexer.TokenKindEOF
	}
	p.currTok++
	return p.toks.TokenInfos[p.currTok].Kind
}

func (p *Parser) curr() lexer.TokenKind {
	if p.currTok >= uint32(len(p.toks.TokenInfos)) {
		return lexer.TokenKindEOF
	}
	return p.toks.TokenInfos[p.currTok].Kind
}

func (p *Parser) peek() lexer.TokenKind {
	if p.currTok+1 >= uint32(len(p.toks.TokenInfos)) {
		return lexer.TokenKindEOF
	}
	return p.toks.TokenInfos[p.currTok+1].Kind
}

func (p *Parser) addTypedLeafNode(kind NodeKind, hasError bool) {
	p.tree = append(p.tree, Node{
		TokIdx:      p.currTok,
		SubtreeSize: 1,
		Kind:        kind,
		HasError:    hasError,
	})
}

func (p *Parser) addLeafNode(hasError bool) {
	p.tree = append(p.tree, Node{
		TokIdx:      p.currTok,
		SubtreeSize: 1,
		HasError:    hasError,
	})
}

func (p *Parser) addTypedNode(tokIdx uint32, kind NodeKind, state stateStackEntry) {
	p.tree = append(p.tree, Node{
		TokIdx:      tokIdx,
		SubtreeSize: uint32(len(p.tree)) - state.subtreeStart + 1,
		Kind:        kind,
		HasError:    state.hasError,
	})
}

func (p *Parser) addNode(tokIdx uint32, state stateStackEntry) {
	p.tree = append(p.tree, Node{
		TokIdx:      tokIdx,
		SubtreeSize: uint32(len(p.tree)) - state.subtreeStart + 1,
		HasError:    state.hasError,
	})
}

func (p *Parser) stackDump() {
	println("currTok:", p.currTok)
	println("stack dump:")
	for i := len(p.stack) - 1; i >= 0; i-- {
		println("  ", i, p.stack[i].st.String())
	}
	println()
}

func (p *Parser) error(err error) {
	//	p.stackDump()
	p.errs = append(p.errs, err)
}

func (p *Parser) expectedCurr(kind ...lexer.TokenKind) {
	p.error(&ExpectedError{
		Expected: kind,
		Got:      p.curr(),
		TokIdx:   p.currTok,
	})
}

func (p *Parser) parseEnderState() {
	state := p.popState()
	top := p.topState()
	top.subtreeStart++
	p.addNode(state.tokIdx, top)
}

func (p *Parser) parseTopLevel() {
	curr := p.curr()

	if curr == lexer.TokenKindEOF {
		p.popState()
		p.addLeafNode(false)
		return
	}

	p.addLeafNode(false)
	p.next()

	switch curr {
	case lexer.TokenKindComment:
		return
	case lexer.TokenKindSyntax:
		p.parseSyntax()
	case lexer.TokenKindEdition:
		p.parseEdition()
	case lexer.TokenKindImport:
		p.tree[len(p.tree)-1].Kind = NodeKindImportStmt
		p.parseImport()
	case lexer.TokenKindPackage:
		p.tree[len(p.tree)-1].Kind = NodeKindPackageStmt
		p.parsePackage()
	case lexer.TokenKindOption:
		p.parseOption()
	case lexer.TokenKindMessage:
		p.tree[len(p.tree)-1].Kind = NodeKindMessageDecl
		p.parseMessage()
	case lexer.TokenKindEnum:
		p.tree[len(p.tree)-1].Kind = NodeKindEnumDecl
		p.parseEnum()
	case lexer.TokenKindService:
		p.tree[len(p.tree)-1].Kind = NodeKindServiceDecl
		p.parseService()
	}
}

func (p *Parser) Parse() (ParseTree, []error) {
	p.pushState(stateTopLevel)
	p.addLeafNode(false)
	p.next()

	for len(p.stack) != 0 && p.currTok < uint32(len(p.toks.TokenInfos)) {
		switch p.topState().st {
		case stateTopLevel:
			p.parseTopLevel()

		// SYNTAXES
		case stateSyntaxAssign:
			p.parseSyntaxAssign()
		case stateSyntaxFinish:
			p.parseSyntaxFinish()

		// EDITIONS
		case stateEditionAssign:
			p.parseEditionAssign()
		case stateEditionFinish:
			p.parseEditionFinish()

		// IMPORTS
		case stateImportValue:
			p.parseImportValue()
		case stateImportFinish:
			p.parseImportFinish()

		// PACKAGES
		case statePackageFinish:
			p.parsePackageFinish()

		// OPTIONS
		case stateOptionName:
			p.parseOptionName()
		case stateOptionNameRest:
			p.parseOptionNameRest()
		case stateOptionNameParenFinish:
			p.parseOptionNameParenFinish()
		case stateOptionAssign:
			p.parseOptionAssign()
		case stateOptionFinish:
			p.parseOptionFinish()
		case stateTextFieldValue:
			p.parseTextFieldValue()
		case stateTextFieldAssign:
			p.parseTextFieldAssign()
		case stateTextFieldName:
			p.parseTextFieldName()
		case stateTextFieldExtensionName:
			p.parseTextFieldExtensionName()
		case stateTextFieldExtensionNameFinish:
			p.parseTextFieldExtensionNameFinish()
		case stateTextMessageValue:
			p.parseTextMessageValue()
		case stateTextMessageInsertSemicolon:
			p.parseTextMessageInsertSemicolon()
		case stateTextMessageFinishRightBrace, stateTextMessageFinishRightAngle:
			p.parseTextMessageFinish()
		case stateTextListValue:
			p.parseTextListValue()
		case stateTextListFinish:
			p.parseTextListFinish()

		// MESSAGES
		case stateMessageBlock:
			p.parseMessageBlock()
		case stateMessageFieldAssign:
			p.parseMessageFieldAssign()
		case stateMessageFieldOption:
			p.parseMessageFieldOption()
		case stateMessageFieldOptionAssign:
			p.parseMessageFieldOptionAssign()
		case stateMessageFieldOptionFinish:
			p.parseMessageFieldOptionFinish()
		case stateMessageFieldFinish:
			p.parseMessageFieldFinish()
		case stateMessageMapStart:
			p.parseMessageMapStart()
		case stateMessageMapKeyValue:
			p.parseMessageMapKeyValue()
		case stateMessageMapComma:
			p.parseMessageMapComma()
		case stateMessageMapFinish:
			p.parseMessageMapFinish()
		case stateMessageValue:
			p.parseMessageValue()
		case stateMessageFinish:
			p.parseMessageFinish()

		// RESERVEDS
		case stateReservedRange:
			p.parseReservedRange()
		case stateReservedName:
			p.parseReservedName()
		case stateReservedFinish:
			p.parseReservedFinish()

		// ONEOFS
		case stateOneofBlock:
			p.parseOneofBlock()
		case stateOneofValue:
			p.parseOneofValue()
		case stateOneofFinish:
			p.parseOneofFinish()

		// ENUMS
		case stateEnumBlock:
			p.parseEnumBlock()
		case stateEnumValue:
			p.parseEnumValue()
		case stateEnumFinish:
			p.parseEnumFinish()

		// SERVICES
		case stateServiceBlock:
			p.parseServiceBlock()
		case stateServiceValue:
			p.parseServiceValue()
		case stateServiceFinish:
			p.parseServiceFinish()

		// RPCS
		case stateRPCDefinition:
			p.parseRPCDefinition()
		case stateRPCReqRes:
			p.parseRPCReqRes()
		case stateRPCReqResFinish:
			p.parseRPCReqResFinish()
		case stateRPCValue:
			p.parseRPCValue()
		case stateRPCFinish:
			p.parseRPCFinish()

		// IDENTIFIERS
		case stateIdentifier:
			p.parseIdentifier()
		case stateFullIdentifierRoot:
			p.parseFullIdentifierRoot()
		case stateFullIdentifierRest:
			p.parseFullIdentifierRest()

		// MISC
		case stateEnder:
			p.parseEnderState()
		}
	}
	return p.tree, p.errs
}
