package parser

import (
	internal_bytes "github.com/Clement-Jean/protein/internal/bytes"
	"github.com/Clement-Jean/protein/token"
)

var protoTopLevelStart = map[token.Kind]bool{
	token.KindSyntax:  true,
	token.KindEdition: true,
	token.KindPackage: true,
	token.KindImport:  true,
	token.KindOption:  true,
	token.KindEnum:    true,
	token.KindMessage: true,
	token.KindService: true,
	token.KindExtend:  true,
}

var exprEnd = map[token.Kind]bool{
	token.KindComma:       true,
	token.KindColon:       true,
	token.KindSemicolon:   true,
	token.KindRightSquare: true,
	token.KindRightAngle:  true,
	token.KindRightBrace:  true,
}

func (p *impl) advanceTo(to map[token.Kind]bool) {
	curr := p.nextToken()

	for ; curr != nil && curr.Kind != token.KindEOF; curr = p.nextToken() {
		kind := curr.Kind

		if !to[kind] {
			literal := p.fm.Lookup(curr.ID)
			kind = literalToKind[internal_bytes.ToString(literal)]
			continue
		}

		if p.idx == p.syncPos && p.syncCnt < 10 { // prevent infinite loop
			p.syncCnt++
			return
		}
		if p.idx > p.syncPos {
			p.syncPos = p.idx
			p.syncCnt = 0
			return
		}
	}
}
