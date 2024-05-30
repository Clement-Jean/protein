package parser

import "github.com/Clement-Jean/protein/lexer"

func (p *Parser) skipPastLikelyEnd(idx int) int {
	if p.currTok >= len(p.toks.TokenInfos)-1 {
		p.currTok++
		return p.currTok - 1
	}

	rootTok := p.toks.TokenInfos[idx]
	rootLine := p.toks.LineInfos[rootTok.LineIdx]
	rootIndent := rootLine.Start + int(rootTok.Column)
	keepSkipping := func(idx int) bool {
		// while we are:
		//   - on the same line
		//   - on a line with bigger indentation
		// we can skip
		tok := p.toks.TokenInfos[idx]
		line := p.toks.LineInfos[tok.LineIdx]
		return line == rootLine || line.Start+int(tok.Column) > rootIndent
	}

	for {
		if p.curr() == lexer.TokenKindSemicolon {
			return p.currTok
		}

		p.next()
		if p.currTok >= len(p.toks.TokenInfos)-1 || !keepSkipping(p.currTok) {
			break
		}
	}

	return p.currTok - 1
}
