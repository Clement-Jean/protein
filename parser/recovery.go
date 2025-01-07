package parser

import (
	"slices"

	"github.com/Clement-Jean/protein/lexer"
)

func (p *Parser) skipTo(end ...lexer.TokenKind) {
	curr := p.curr()
	for curr != lexer.TokenKindEOF && !slices.Contains(end, curr) {
		if p.skipSubscope(curr) {
			curr = p.curr()
			continue
		}

		curr = p.next()
	}
}

func (p *Parser) skipSubscope(curr lexer.TokenKind) bool {
	if !curr.IsOpeningSymbol() {
		return false
	}
	end := curr.MatchingClosingSymbol()
	for curr != lexer.TokenKindEOF && curr != end {
		curr = p.next()
	}
	return true
}

func (p *Parser) skipPastLikelyEnd(idx uint32) uint32 {
	if p.currTok >= uint32(len(p.toks.TokenInfos)-1) {
		p.currTok++
		return p.currTok - 1
	}

	rootTok := p.toks.TokenInfos[idx]
	rootLine := p.toks.LineInfos[rootTok.LineIdx]
	rootIndent := rootLine.Start + rootTok.Column
	keepSkipping := func(idx uint32) bool {
		// while we are:
		//   - on the same line
		//   - on a line with bigger indentation
		// we can skip
		tok := p.toks.TokenInfos[idx]
		line := p.toks.LineInfos[tok.LineIdx]
		return line == rootLine || line.Start+tok.Column > rootIndent
	}

	curr := p.curr()
	for {
		if curr == lexer.TokenKindRightBrace || curr == lexer.TokenKindRightAngle {
			return p.currTok
		}

		if curr == lexer.TokenKindSemicolon {
			return p.currTok
		}

		if p.skipSubscope(curr) {
			curr = p.curr()
			continue
		}

		curr = p.next()
		if p.currTok >= uint32(len(p.toks.TokenInfos)-1) || !keepSkipping(p.currTok) {
			break
		}
	}

	return p.currTok - 1
}
