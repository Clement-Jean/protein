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

	if imp := d.GetDependency(); len(imp) == 0 || imp[0] != expected {
		if len(imp) == 0 {
			t.Fatalf("import wrong. expected='%s', got=nil", expected)
		}

		t.Fatalf("import wrong. expected='%s', got='%s'", expected, imp[0])
	}

	return ""
}

func TestParseImport(t *testing.T) {
	// import "google/protobuf/empty.proto";
	err := runImportCheck(t, "google/protobuf/empty.proto", []lexer.Token{
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

	expectedErr := fmt.Sprintf(errorUnexpected, "String", "Identifier")
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

	expectedErr := fmt.Sprintf(errorUnexpected, ";", "EOF")
	if err != expectedErr {
		t.Fatalf("error wrong. expected='%s', got='%s'", expectedErr, err)
	}
}
