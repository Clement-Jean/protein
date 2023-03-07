package parser

import (
	"fmt"
	"testing"

	"github.com/Clement-Jean/protein/lexer"
)

func runSyntaxCheck(t *testing.T, expected *string, tokens []lexer.Token) string {
	d, err := runCheck(t, tokens)

	if len(err) != 0 {
		return err
	}

	if expected != d.Syntax && (d.Syntax == nil || expected == nil || *d.Syntax != *expected) {
		if d.Syntax == nil {
			t.Fatalf("syntax wrong. expected='%s', got=nil", *expected)
		} else if expected == nil {
			t.Fatalf("syntax wrong. expected=nil, got='%s'", *d.Syntax)
		}

		t.Fatalf("syntax wrong. expected='%s', got='%s'", *expected, *d.Syntax)
	}

	return ""
}

func TestParseSyntaxProto3(t *testing.T) {
	// syntax = "proto3";
	expected := "proto3"
	err := runSyntaxCheck(t, &expected, []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "syntax", Position: lexer.Position{}},
		{Type: lexer.TokenEqual, Literal: "=", Position: lexer.Position{}},
		{Type: lexer.TokenStr, Literal: "\"proto3\"", Position: lexer.Position{}},
		{Type: lexer.TokenSemicolon, Literal: ";", Position: lexer.Position{}},
	})

	if len(err) != 0 {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestParseSyntaxProto2(t *testing.T) {
	// syntax = 'proto2';
	expected := "proto2"
	err := runSyntaxCheck(t, &expected, []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "syntax", Position: lexer.Position{}},
		{Type: lexer.TokenEqual, Literal: "=", Position: lexer.Position{}},
		{Type: lexer.TokenStr, Literal: "'proto2'", Position: lexer.Position{}},
		{Type: lexer.TokenSemicolon, Literal: ";", Position: lexer.Position{}},
	})

	if len(err) != 0 {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestParseSyntaxExpectedEqual(t *testing.T) {
	// syntax 'proto2';
	err := runSyntaxCheck(t, nil, []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "syntax", Position: lexer.Position{}},
		{Type: lexer.TokenStr, Literal: "'proto2'", Position: lexer.Position{}},
		{Type: lexer.TokenSemicolon, Literal: ";", Position: lexer.Position{}},
	})

	expectedErr := fmt.Sprintf(errorUnexpectedPeek, "=", "String")
	if err != expectedErr {
		t.Fatalf("error wrong. expected='%s', got='%s'", expectedErr, err)
	}
}

func TestParseSyntaxExpectedStr(t *testing.T) {
	// syntax = proto2;
	err := runSyntaxCheck(t, nil, []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "syntax", Position: lexer.Position{}},
		{Type: lexer.TokenEqual, Literal: "=", Position: lexer.Position{}},
		{Type: lexer.TokenIdentifier, Literal: "proto2", Position: lexer.Position{}},
		{Type: lexer.TokenSemicolon, Literal: ";", Position: lexer.Position{}},
	})

	expectedErr := fmt.Sprintf(errorUnexpectedPeek, "String", "Identifier")
	if err != expectedErr {
		t.Fatalf("error wrong. expected='%s', got='%s'", expectedErr, err)
	}
}

func TestParseSyntaxExpectedSemicolon(t *testing.T) {
	// syntax = 'proto2'
	err := runSyntaxCheck(t, nil, []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "syntax", Position: lexer.Position{}},
		{Type: lexer.TokenEqual, Literal: "=", Position: lexer.Position{}},
		{Type: lexer.TokenStr, Literal: "'proto2'", Position: lexer.Position{}},
	})

	expectedErr := fmt.Sprintf(errorUnexpectedPeek, ";", "EOF")
	if err != expectedErr {
		t.Fatalf("error wrong. expected='%s', got='%s'", expectedErr, err)
	}
}
