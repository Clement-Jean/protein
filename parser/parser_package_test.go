package parser

import (
	"fmt"
	"testing"

	"github.com/Clement-Jean/protein/lexer"
)

func runPackageCheck(t *testing.T, expected string, tokens []lexer.Token) string {
	d, err := runCheck(t, tokens)

	if len(err) != 0 {
		return err
	}

	if pkg := d.GetPackage(); pkg != expected {
		t.Fatalf("package wrong. expected='%s', got='%s'", expected, pkg)
	}

	return ""
}

func TestParsePackageIdentifier(t *testing.T) {
	// package google;
	err := runPackageCheck(t, "google", []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "package", Position: lexer.Position{}},
		{Type: lexer.TokenIdentifier, Literal: "google", Position: lexer.Position{}},
		{Type: lexer.TokenSemicolon, Literal: ";", Position: lexer.Position{}},
	})

	if len(err) != 0 {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestParsePackageFullIdentifier(t *testing.T) {
	// package google.protobuf;
	err := runPackageCheck(t, "google.protobuf", []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "package", Position: lexer.Position{}},
		{Type: lexer.TokenIdentifier, Literal: "google", Position: lexer.Position{}},
		{Type: lexer.TokenDot, Literal: ".", Position: lexer.Position{}},
		{Type: lexer.TokenIdentifier, Literal: "protobuf", Position: lexer.Position{}},
		{Type: lexer.TokenSemicolon, Literal: ";", Position: lexer.Position{}},
	})

	if len(err) != 0 {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestParseUnterminedPackage(t *testing.T) {
	// package google.;
	err := runPackageCheck(t, "", []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "package", Position: lexer.Position{}},
		{Type: lexer.TokenIdentifier, Literal: "google", Position: lexer.Position{}},
		{Type: lexer.TokenDot, Literal: ".", Position: lexer.Position{}},
		{Type: lexer.TokenSemicolon, Literal: ";", Position: lexer.Position{}},
	})

	expectedErr := fmt.Sprintf(errorUnexpected, "Identifier", ";")
	if err != expectedErr {
		t.Fatalf("error wrong. expected='%s', got='%s'", expectedErr, err)
	}
}

func TestParseExpectedIdentifier(t *testing.T) {
	// package 'google';
	err := runPackageCheck(t, "", []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "package", Position: lexer.Position{}},
		{Type: lexer.TokenStr, Literal: "'google'", Position: lexer.Position{}},
		{Type: lexer.TokenSemicolon, Literal: ";", Position: lexer.Position{}},
	})

	expectedErr := fmt.Sprintf(errorUnexpected, "Identifier", "String")
	if err != expectedErr {
		t.Fatalf("error wrong. expected='%s', got='%s'", expectedErr, err)
	}
}

func TestParseExpectedSemicolon(t *testing.T) {
	// package google
	err := runPackageCheck(t, "", []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "package", Position: lexer.Position{}},
		{Type: lexer.TokenIdentifier, Literal: "google", Position: lexer.Position{}},
	})

	expectedErr := fmt.Sprintf(errorUnexpected, ";", "EOF")
	if err != expectedErr {
		t.Fatalf("error wrong. expected='%s', got='%s'", expectedErr, err)
	}
}
