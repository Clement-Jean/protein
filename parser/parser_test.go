package parser

import (
	"io"
	"os"
	"testing"

	"github.com/Clement-Jean/protein/lexer"
	"google.golang.org/protobuf/types/descriptorpb"
)

type FakeLexer struct {
	i      int
	tokens []lexer.Token
}

func (l *FakeLexer) NextToken() lexer.Token {
	if l.i >= len(l.tokens) {
		return lexer.Token{Type: lexer.EOF, Position: lexer.Position{}}
	}

	token := l.tokens[l.i]
	l.i++
	return token
}

func runCheck(t *testing.T, tokens []lexer.Token) (*descriptorpb.FileDescriptorProto, string) {
	r, w, err := os.Pipe()

	if err != nil {
		t.Fatalf("couldn't create pipe")
	}

	oldStderr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = oldStderr }()

	l := &FakeLexer{tokens: tokens}
	p := New(l)
	d := p.Parse()

	w.Close()
	errors, err := io.ReadAll(r)

	if err != nil {
		t.Fatalf("couldn't read stderr")
	}

	return &d, string(errors)
}
