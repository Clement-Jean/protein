package lexer

import (
	"io"
	"unicode/utf8"
)

type inMemoryRuneReader struct {
	buf        []byte
	start, pos uint64
	atEOF      bool
}

func newInMemoryRuneReader(rd io.Reader) *inMemoryRuneReader {
	content, err := io.ReadAll(rd)

	if err != nil {
		panic(err)
	}

	return &inMemoryRuneReader{
		buf: content,
	}
}

func (l *inMemoryRuneReader) next() rune {
	if l.pos >= uint64(len(l.buf)) {
		l.atEOF = true
		return 0
	}

	r, w := utf8.DecodeRune(l.buf[l.pos:])
	l.pos += uint64(w)
	return r
}

func (l *inMemoryRuneReader) peek() rune {
	if l.pos >= uint64(len(l.buf)) {
		return 0
	}

	r, _ := utf8.DecodeRune(l.buf[l.pos:])
	return r
}

func (l *inMemoryRuneReader) backup() {
	if l.atEOF || l.pos == 0 {
		return
	}

	_, w := utf8.DecodeLastRune(l.buf[:l.pos])
	l.pos -= uint64(w)
}
