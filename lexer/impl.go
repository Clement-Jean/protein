package lexer

import (
	"io"
	"slices"
	"strings"

	"github.com/Clement-Jean/protein/internal/span"
	"github.com/Clement-Jean/protein/token"
)

type impl struct {
	buf      *inMemoryRuneReader
	location span.Span
	kind     token.Kind
}

func New(reader io.Reader) Lexer {
	return &impl{
		buf: newInMemoryRuneReader(reader),
	}
}

func (l *impl) emit(kind token.Kind) stateFn {
	l.kind = kind
	l.location = span.Span{Start: l.buf.start, End: l.buf.pos}
	l.buf.start = l.buf.pos
	return nil
}

func (l *impl) errorf(kind token.Kind) stateFn {
	return l.emit(kind)
}

func (l *impl) accept(valid string) bool {
	if strings.ContainsRune(valid, l.buf.next()) {
		return true
	}
	l.buf.backup()
	return false
}

func (l *impl) acceptWhile(fn func(rune) bool) {
	for fn(l.buf.next()) {
	}
	l.buf.backup()
}

func (l *impl) nextToken() (token.Kind, span.Span) {
	for state := l.lexProto(); state != nil; state = state() {
	}
	return l.kind, l.location
}

func (l *impl) Tokenize() ([]token.Kind, []span.Span) {
	var kinds []token.Kind
	var spans []span.Span

	kind, loc := l.nextToken()

	for kind != token.KindEOF {
		kinds = append(kinds, kind)
		spans = append(spans, loc)
		kind, loc = l.nextToken()
	}

	kinds = append(kinds, kind)
	spans = append(spans, loc)
	kinds = slices.Clip(kinds)
	spans = slices.Clip(spans)
	return kinds, spans
}
