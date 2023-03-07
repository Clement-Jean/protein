package parser

import (
	"fmt"
	"testing"

	"github.com/Clement-Jean/protein/lexer"
)

func runImportCheck(t *testing.T, expected string, tokens []lexer.Token) string {
	d, err := runCheck(t, tokens)

	if len(err) != 0 {
		return err
	}

	if len(d.Dependency) == 0 || expected != d.Dependency[0] && (d.Package == nil || len(expected) == 0 || d.Dependency[0] != expected) {
		if len(d.Dependency) == 0 {
			t.Fatalf("import wrong. expected='%s', got=nil", expected)
		} else if len(expected) == 0 {
			t.Fatalf("import wrong. expected=nil, got='%s'", d.Dependency[0])
		}

		t.Fatalf("import wrong. expected='%s', got='%s'", expected, d.Dependency[0])
	}

	return ""
}

func TestParseImport(t *testing.T) {
	// import "google/protobuf/empty.proto";
	expected := "google/protobuf/empty.proto"
	err := runImportCheck(t, expected, []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "import", Position: lexer.Position{}},
		{Type: lexer.TokenStr, Literal: "\"google/protobuf/empty.proto\"", Position: lexer.Position{}},
		{Type: lexer.TokenSemicolon, Literal: ";", Position: lexer.Position{}},
	})

	if len(err) != 0 {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestParseImportExpectedStr(t *testing.T) {
	// import google;
	err := runImportCheck(t, "", []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "import", Position: lexer.Position{}},
		{Type: lexer.TokenIdentifier, Literal: "google", Position: lexer.Position{}},
		{Type: lexer.TokenSemicolon, Literal: ";", Position: lexer.Position{}},
	})

	expectedErr := fmt.Sprintf(errorUnexpectedPeek, "String", "Identifier")
	if err != expectedErr {
		t.Fatalf("error wrong. expected='%s', got='%s'", expectedErr, err)
	}
}

func TestParseImportExpectedSemicolon(t *testing.T) {
	// import google;
	err := runImportCheck(t, "", []lexer.Token{
		{Type: lexer.TokenIdentifier, Literal: "import", Position: lexer.Position{}},
		{Type: lexer.TokenStr, Literal: "'google'", Position: lexer.Position{}},
	})

	expectedErr := fmt.Sprintf(errorUnexpectedPeek, ";", "EOF")
	if err != expectedErr {
		t.Fatalf("error wrong. expected='%s', got='%s'", expectedErr, err)
	}
}
