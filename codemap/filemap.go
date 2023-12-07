package codemap

import (
	"log"

	"github.com/Clement-Jean/protein/internal/bytes"
	"github.com/Clement-Jean/protein/internal/span"
	"github.com/Clement-Jean/protein/token"
)

type FileMap struct {
	kinds   []token.Kind
	spans   []span.Span
	content []byte
}

func (fm *FileMap) Lookup(id token.UniqueID) []byte {
	s, _ := fm.SpanOf(id)
	if s.End > len(fm.content) {
		log.Panicf("FileMap thinks it contains %d, but the range (%s) doesn't point to anything valid!", id, s)
	}

	return fm.content[s.Start:s.End]
}

func (fm *FileMap) SpanOf(id token.UniqueID) (span.Span, bool) {
	if id >= len(fm.kinds) {
		return span.Span{}, false
	}
	s := fm.spans[id]
	return s, true
}

func (fm *FileMap) Merge(kind token.Kind, ids ...token.UniqueID) token.UniqueID {
	if len(ids) == 0 {
		log.Panic("FileMap cannot merge 0 id, it expects 1+ ids")
	}

	def, _ := fm.SpanOf(ids[0])
	start := def.Start
	end := def.End

	for _, id := range ids {
		r, ok := fm.SpanOf(id)

		if !ok {
			continue
		}

		start = min(start, r.Start)
		end = max(end, r.End)
	}

	if end < start {
		log.Panicf("FileMap seem to have wrong data (end < start): %v", ids)
	}

	id := len(fm.kinds)
	fm.kinds = append(fm.kinds, kind)
	fm.spans = append(fm.spans, span.Span{Start: start, End: end})
	return id
}

func (fm *FileMap) RegisterTokens(kinds []token.Kind, spans []span.Span) []token.Token {
	if len(kinds) != len(spans) {
		return nil
	}

	fm.kinds = kinds
	fm.spans = spans
	r := make([]token.Token, len(kinds))

	for i := 0; i < len(kinds); i++ {
		r[i] = token.Token{ID: i, Kind: kinds[i]}
	}

	return r
}

func (fm *FileMap) PrintItems() {
	for i := 0; i < len(fm.kinds); i++ {
		if literal := fm.Lookup(i); literal != nil {
			log.Println(i, bytes.ToString(literal))
		}
	}
}
