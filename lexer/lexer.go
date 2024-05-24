package lexer

import (
	"bytes"
	"io"
	"math"
)

type Lexer struct {
	src         []byte
	errs        []error
	toks        *TokenizedBuffer
	currLineIdx LineIdx
	srcPos      int // the idx at which the file content really starts
	tokPos      int // the begining of a token
	readPos     int // the idx we are reading at in src
}

func NewFromReader(r io.Reader) (*Lexer, error) {
	src, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	srcPos := 0
	if len(src) >= 3 && bytes.Equal(src[:3], []byte{0xEF, 0xBB, 0xBF}) {
		// skip UTF8 BOM
		srcPos = 3
	}

	return &Lexer{
		src:     src,
		toks:    &TokenizedBuffer{},
		srcPos:  srcPos,
		readPos: srcPos,
	}, nil
}

func (l *Lexer) next() byte {
	if l.readPos >= len(l.src) {
		return 0
	}
	ch := l.src[l.readPos]
	l.readPos++
	return ch
}

func (l *Lexer) backup() {
	if l.readPos <= 0 {
		return
	}
	l.readPos--
}

func (l *Lexer) acceptWhile(fn func(byte) bool) (length int) {
	var ch byte

	for ch = l.next(); fn(ch); ch = l.next() {
		length++
	}

	if ch != 0 {
		l.backup()
	}
	return length
}

func (l *Lexer) computeColumn(position int) uint32 {
	if int(l.currLineIdx) >= len(l.toks.LineInfos) {
		return math.MaxUint32
	}
	return uint32(l.srcPos + position - l.toks.LineInfos[l.currLineIdx].Start)
}

func (l *Lexer) emit(kind TokenKind, position int) stateFn {
	l.toks.TokenInfos = append(l.toks.TokenInfos, TokenInfo{
		Kind:    kind,
		LineIdx: l.currLineIdx,
		Column:  l.computeColumn(position),
	})
	return nil
}

func (l *Lexer) error(err error) stateFn {
	l.errs = append(l.errs, err)
	return l.emit(TokenKindError, l.tokPos)
}

func (l *Lexer) makeLines() {
	nbLines := bytes.Count(l.src, []byte{'\n'}) + 1
	l.toks.LineInfos = make([]LineInfo, nbLines)

	i := 0
	start := l.srcPos
	for i = range nbLines {
		rest := l.src[start:]
		idx := bytes.IndexByte(rest, '\n')
		newlineIdx := start + idx
		info := &l.toks.LineInfos[i]
		info.Start = start
		info.Len = uint32(newlineIdx - start)
		start = newlineIdx + 1
	}

	info := &l.toks.LineInfos[i]
	info.Start = start
	info.Len = uint32(len(l.src) - start)
}

func (l *Lexer) start() {
	l.tokPos = l.srcPos
	l.emit(TokenKindBOF, l.tokPos)
}

func (l *Lexer) Lex() (*TokenizedBuffer, []error) {
	l.makeLines()
	l.start()

	lastKind := TokenKindBOF
	for lastKind != TokenKindEOF {
		for state := l.lexProto(); state != nil; {
			state = state()
		}
		lastKind = l.toks.TokenInfos[len(l.toks.TokenInfos)-1].Kind
	}
	return l.toks, l.errs
}
