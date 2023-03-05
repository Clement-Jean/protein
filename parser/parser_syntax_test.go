package parser

import (
	"testing"

	"github.com/Clement-Jean/protein/lexer"
)

func runSyntaxCheck(t *testing.T, expected *string, tokens []lexer.Token) {
	l := &FakeLexer{tokens: tokens}
	p := New(l)
	d := p.Parse()

	if expected != d.Syntax && (d.Syntax == nil || expected == nil || *d.Syntax != *expected) {
		if d.Syntax == nil {
			t.Fatalf("syntax wrong. expected='%s', got=nil", *expected)
		} else if expected == nil {
			t.Fatalf("syntax wrong. expected=nil, got='%s'", *d.Syntax)
		}

		t.Fatalf("syntax wrong. expected='%s', got='%s'", *expected, *d.Syntax)
	}
}

func TestParseSyntaxProto3(t *testing.T) {
	// syntax = "proto3";
	expected := "proto3"
	runSyntaxCheck(t, &expected, []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "syntax", Position: lexer.Position{}},
		{Type: lexer.TokenEqual, Literal: "=", Position: lexer.Position{}},
		{Type: lexer.TokenStr, Literal: "\"proto3\"", Position: lexer.Position{}},
		{Type: lexer.TokenSemicolon, Literal: ";", Position: lexer.Position{}},
	})
}

func TestParseSyntaxProto2(t *testing.T) {
	// syntax = 'proto2';
	expected := "proto2"
	runSyntaxCheck(t, &expected, []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "syntax", Position: lexer.Position{}},
		{Type: lexer.TokenEqual, Literal: "=", Position: lexer.Position{}},
		{Type: lexer.TokenStr, Literal: "'proto2'", Position: lexer.Position{}},
		{Type: lexer.TokenSemicolon, Literal: ";", Position: lexer.Position{}},
	})
}

func TestParseSyntaxExpectedEqual(t *testing.T) {
	// syntax 'proto2';
	runSyntaxCheck(t, nil, []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "syntax", Position: lexer.Position{}},
		{Type: lexer.TokenStr, Literal: "'proto2'", Position: lexer.Position{}},
		{Type: lexer.TokenSemicolon, Literal: ";", Position: lexer.Position{}},
	})
}

func TestParseSyntaxExpectedStr(t *testing.T) {
	// syntax = proto2;
	runSyntaxCheck(t, nil, []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "syntax", Position: lexer.Position{}},
		{Type: lexer.TokenEqual, Literal: "=", Position: lexer.Position{}},
		{Type: lexer.TokenIdentifier, Literal: "proto2", Position: lexer.Position{}},
		{Type: lexer.TokenSemicolon, Literal: ";", Position: lexer.Position{}},
	})
}

func TestParseSyntaxExpectedSemicolon(t *testing.T) {
	// syntax = 'proto2'
	runSyntaxCheck(t, nil, []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "syntax", Position: lexer.Position{}},
		{Type: lexer.TokenEqual, Literal: "=", Position: lexer.Position{}},
		{Type: lexer.TokenStr, Literal: "'proto2'", Position: lexer.Position{}},
	})
}
