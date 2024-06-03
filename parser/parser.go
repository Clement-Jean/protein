package parser

import (
	"fmt"

	"github.com/Clement-Jean/protein/lexer"
)

type Parser struct {
	toks    *lexer.TokenizedBuffer
	tree    ParseTree
	stack   []stateStackEntry
	errs    []error
	currTok int
}

func New(toks *lexer.TokenizedBuffer) *Parser {
	return &Parser{
		toks:  toks,
		tree:  make([]Node, 0, len(toks.TokenInfos)),
		stack: make([]stateStackEntry, 0, 31),
	}
}

func (p *Parser) pushState(st state) {
	p.stack = append(p.stack, stateStackEntry{
		st:           st,
		tokIdx:       p.currTok,
		subtreeStart: int32(len(p.tree) - 1),
	})
}

func (p *Parser) pushStateWithIdx(st state, idx int) {
	p.stack = append(p.stack, stateStackEntry{
		st:           st,
		tokIdx:       idx,
		subtreeStart: int32(len(p.tree) - 1),
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
	if p.currTok+1 >= len(p.toks.TokenInfos) {
		return lexer.TokenKindEOF
	}
	p.currTok++
	return p.toks.TokenInfos[p.currTok].Kind
}

func (p *Parser) curr() lexer.TokenKind {
	if p.currTok >= len(p.toks.TokenInfos) {
		return lexer.TokenKindEOF
	}
	return p.toks.TokenInfos[p.currTok].Kind
}

func (p *Parser) addLeafNode(hasError bool) {
	p.tree = append(p.tree, Node{
		TokIdx:      p.currTok,
		SubtreeSize: 1,
		HasError:    hasError,
	})
}

func (p *Parser) addNode(tokIdx int, state stateStackEntry) {
	p.tree = append(p.tree, Node{
		TokIdx:      tokIdx,
		SubtreeSize: int32(len(p.tree)) - state.subtreeStart + 1,
		HasError:    state.hasError,
	})
}

func (p *Parser) error(err error) {
	p.errs = append(p.errs, err)
}

func (p *Parser) expectedCurr(kind ...lexer.TokenKind) {
	p.error(fmt.Errorf("expected %v, got %s", kind, p.curr()))
}

func (p *Parser) parseTopLevel() {
	if p.curr() == lexer.TokenKindEOF {
		p.popState()
		p.addLeafNode(false)
		return
	}

	curr := p.curr()
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
		p.parseImport()
	case lexer.TokenKindPackage:
		p.parsePackage()
	case lexer.TokenKindOption:
		p.parseOption()
	}
}

func (p *Parser) Parse() (ParseTree, []error) {
	p.pushState(stateTopLevel)
	p.addLeafNode(false)
	p.next()

	for len(p.stack) != 0 && p.currTok < len(p.toks.TokenInfos) {
		switch p.topState().st {
		case stateTopLevel:
			p.parseTopLevel()
		case stateSyntaxAssign:
			p.parseSyntaxAssign()
		case stateSyntaxFinish:
			p.parseSyntaxFinish()
		case stateEditionAssign:
			p.parseEditionAssign()
		case stateEditionFinish:
			p.parseEditionFinish()
		case stateImportValue:
			p.parseImportValue()
		case stateImportFinish:
			p.parseImportFinish()
		case statePackageFinish:
			p.parsePackageFinish()
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

		case stateFullIdentifierRoot:
			p.parseFullIdentifierRoot()
		case stateFullIdentifierRest:
			p.parseFullIdentifierRest()
		}
	}
	return p.tree, p.errs
}
