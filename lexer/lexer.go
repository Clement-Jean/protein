package lexer

import (
	"bytes"
	"io"
	"math"
	"strings"

	"github.com/Clement-Jean/protein/source"
)

type Lexer struct {
	src         *source.Buffer
	toks        *TokenizedBuffer
	errs        []error
	currLineIdx LineIdx
	srcPos      uint32 // the idx at which the file content really starts
	tokPos      uint32 // the begining of a token
	readPos     uint32 // the idx we are reading at in src
}

func newLexer(src *source.Buffer) (*Lexer, error) {
	var srcPos uint32 = 0
	if bytes.Equal(src.Range(0, 3), []byte{0xEF, 0xBB, 0xBF}) {
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

func NewFromFile(filename string) (*Lexer, error) {
	src, err := source.NewFromFile(filename)
	if err != nil {
		return nil, err
	}
	return newLexer(src)
}

func NewFromReader(r io.Reader) (*Lexer, error) {
	src, err := source.NewFromReader(r)
	if err != nil {
		return nil, err
	}
	return newLexer(src)
}

func (l *Lexer) next() byte {
	if l.readPos >= l.src.Len() {
		return 0
	}
	ch := l.src.At(l.readPos)
	l.readPos++
	return ch
}

func (l *Lexer) backup() {
	if l.readPos <= 0 {
		return
	}
	l.readPos--
}

func (l *Lexer) accept(valid string) bool {
	ch := l.next()
	if strings.IndexByte(valid, ch) != -1 {
		return true
	}

	if ch != 0 {
		l.backup()
	}
	return false
}

func (l *Lexer) acceptWhile(fn func(byte) bool) {
	ch := l.next()
	for fn(ch) {
		ch = l.next()
	}

	if ch != 0 {
		l.backup()
	}
}

func (l *Lexer) computeColumn(position uint32) uint32 {
	if int(l.currLineIdx) >= len(l.toks.LineInfos) {
		return math.MaxUint32
	}
	return l.srcPos + position - l.toks.LineInfos[l.currLineIdx].Start
}

func (l *Lexer) emit(kind TokenKind, position uint32) stateFn {
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
	nbLines := bytes.Count(l.src.Bytes(), []byte{'\n'}) + 1
	l.toks.LineInfos = make([]LineInfo, nbLines)

	i := 0
	start := l.srcPos
	for i = range nbLines {
		rest := l.src.From(start)
		idx := uint32(bytes.IndexByte(rest, '\n'))
		newlineIdx := start + idx
		info := &l.toks.LineInfos[i]
		info.Start = start
		info.Len = newlineIdx - start
		start = newlineIdx + 1
	}

	info := &l.toks.LineInfos[i]
	info.Start = start
	info.Len = l.src.Len() - start
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
