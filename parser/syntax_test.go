package parser

import (
	"testing"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/parser/internal"
	"github.com/Clement-Jean/protein/token"
)

func checkSyntax(t *testing.T, got ast.Syntax, expected ast.Syntax) {
	t.Helper()
	checkIDs(t, got.ID, expected.ID)
}

func TestSyntax(t *testing.T) {
	tests := []TestCase[ast.Syntax]{
		{
			name:        internal.CaseName("syntax", true),
			expectedObj: ast.Syntax{ID: 4, Value: ast.String{ID: 2}},

			content: "syntax = 'proto2';",
			indices: "a-----bcde-------fg",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindEqual,
				token.KindStr,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("syntax", false, "expect_equal"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindStr}, token.KindEqual),
			},

			content: "syntax 'proto3';",
			indices: "a-----bc-------de",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'d', 'e'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindStr,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("syntax", false, "expect_string"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindIdentifier}, token.KindStr),
			},

			content: "syntax = proto3;",
			indices: "a-----bcde-----fg",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindEqual,
				token.KindIdentifier,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("syntax", false, "expect_semicolon"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 0, Kind: token.KindEOF}, token.KindSemicolon),
			},

			content: "syntax = 'proto3'",
			indices: "a-----bcde-------f",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindEqual,
				token.KindStr,
			},
		},
	}

	runTestCases(t, tests, checkSyntax, (*impl).parseSyntax)
}
