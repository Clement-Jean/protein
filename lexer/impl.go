package lexer

import (
	"bytes"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/Clement-Jean/protein/internal/span"
	"github.com/Clement-Jean/protein/token"
)

type impl struct {
	buf        []byte
	location   span.Span
	start, pos int
	kind       token.Kind
	atEOF      bool
}

var utf8Bom = []byte{0xEF, 0xBB, 0xBF}

func New(buf []byte) Lexer {
	start := 0
	pos := 0

	if len(buf) >= 3 && bytes.Equal(buf[0:3], utf8Bom) {
		start = 3
		pos = 3
	}

	return &impl{buf: buf, start: start, pos: pos}
}

func (l *impl) next() rune {
	if l.pos >= len(l.buf) {
		l.atEOF = true
		return 0
	}

	r, w := utf8.DecodeRune(l.buf[l.pos:])
	l.pos += w
	return r
}

func (l *impl) peek() rune {
	if l.pos >= len(l.buf) {
		return 0
	}

	r, _ := utf8.DecodeRune(l.buf[l.pos:])
	return r
}

func (l *impl) backup() {
	if l.atEOF || l.pos == 0 {
		return
	}

	_, w := utf8.DecodeLastRune(l.buf[:l.pos])
	l.pos -= w
}

func (l *impl) emit(kind token.Kind) stateFn {
	l.kind = kind
	l.location = span.Span{Start: l.start, End: l.pos}
	l.start = l.pos
	return nil
}

func (l *impl) errorf(kind token.Kind) stateFn {
	return l.emit(kind)
}

func (l *impl) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *impl) acceptWhile(fn func(rune) bool) {
	for fn(l.next()) {
	}
	l.backup()
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
