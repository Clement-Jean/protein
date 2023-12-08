package parser

import (
	"testing"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/parser/internal"
	"github.com/Clement-Jean/protein/token"
)

func checkEdition(t *testing.T, got, expected ast.Edition) {
	t.Helper()
	checkIDs(t, got.ID, expected.ID)
	checkIDs(t, got.Value.ID, expected.Value.ID)
}

func TestEdition(t *testing.T) {
	tests := []TestCase[ast.Edition]{
		{
			name:        internal.CaseName("edition", true),
			expectedObj: ast.Edition{ID: 5, Value: ast.String{ID: 2}},

			content: "edition = '2023';",
			indices: "a------bcde-----fg",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'}, {'g', 'g'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindEqual,
				token.KindStr,
				token.KindSemicolon,
				token.KindEOF,
			},
		},
		{
			name: internal.CaseName("edition", false, "expect_equal"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindStr}, token.KindEqual),
			},

			content: "edition '2023';",
			indices: "a------bc-----de",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'e', 'e'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindStr,
				token.KindSemicolon,
				token.KindEOF,
			},
		},
		{
			name: internal.CaseName("edition", false, "expect_string"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindIdentifier}, token.KindStr),
			},

			content: "edition = 2023;",
			indices: "a------bcde---fg",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'}, {'g', 'g'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindEqual,
				token.KindIdentifier,
				token.KindSemicolon,
				token.KindEOF,
			},
		},
		{
			name: internal.CaseName("edition", false, "expect_semicolon"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindEOF}, token.KindSemicolon),
			},

			content: "edition = '2023'",
			indices: "a------bcde-----f",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'f'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindEqual,
				token.KindStr,
				token.KindEOF,
			},
		},
	}

	runTestCases(t, tests, checkEdition, (*impl).parseEdition)
}
