package parser

import (
	"testing"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/parser/internal"
	"github.com/Clement-Jean/protein/token"
)

func checkOptions(t *testing.T, got, expected []ast.Option) {
	if len(expected) != len(got) {
		t.Fatalf("expected %d options, got %d", len(expected), len(got))
	}

	for i := range expected {
		checkOption(t, got[i], expected[i])
	}
}

func checkOption(t *testing.T, got, expected ast.Option) {
	if got.Value == nil {
		t.Fatalf("expected Option.Value not nil")
	}

	checkIDs(t, got.ID, expected.ID)
	checkIDs(t, got.Name.ID, expected.Name.ID)
	checkIDs(t, got.Value.GetID(), expected.Value.GetID())
	checkIdentifier(t, got.Name, expected.Name)
	checkIdentifierParts(t, got.Name.Parts, expected.Name.Parts)
	checkConstantValue(t, got.Value, expected.Value)
}

func TestOption(t *testing.T) {
	tests := []TestCase[ast.Option]{
		{
			name:        internal.CaseName("option", true),
			expectedObj: &ast.Option{ID: 6, Name: ast.Identifier{ID: 1}, Value: ast.Boolean{ID: 3}},

			content: "option deprecated = true;",
			indices: "a-----bc---------defg---hi",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'}, {'h', 'i'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier,
				token.KindEqual,
				token.KindIdentifier,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("option", true, "extend"),
			expectedObj: &ast.Option{
				ID:    11,
				Name:  ast.Identifier{ID: 10, Parts: []token.UniqueID{1, 2, 3, 5}},
				Value: ast.Boolean{ID: 7},
			},

			content: "option (custom).deprecated = true;",
			indices: "a-----bcd-----efg---------hijk---lm",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'},
				{'f', 'g'}, {'g', 'h'}, {'i', 'j'}, {'k', 'l'},
				{'l', 'm'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindLeftParen,
				token.KindIdentifier, // custom
				token.KindRightParen,
				token.KindDot,
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("option", true, "full_extend"),
			expectedObj: &ast.Option{
				ID:    14,
				Name:  ast.Identifier{ID: 13, Parts: []token.UniqueID{1, 5, 7, 12}},
				Value: ast.Boolean{ID: 9},
			},

			content: "option (protein.custom).deprecated = true;",
			indices: "a-----bcd------ef-----ghi---------jklm---no",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'},
				{'f', 'g'}, {'g', 'h'}, {'h', 'i'}, {'i', 'j'},
				{'k', 'l'}, {'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindLeftParen,
				token.KindIdentifier, // protein
				token.KindDot,
				token.KindIdentifier, // custom
				token.KindRightParen,
				token.KindDot,
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("option", true, "dot_full_extend"),
			expectedObj: &ast.Option{
				ID:    15,
				Name:  ast.Identifier{ID: 14, Parts: []token.UniqueID{1, 6, 8, 13}},
				Value: ast.Boolean{ID: 10},
			},

			content: "option (.protein.custom).deprecated = true;",
			indices: "a-----bcde------fg-----hij---------klmn---op",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'},
				{'f', 'g'}, {'g', 'h'}, {'h', 'i'}, {'i', 'j'},
				{'j', 'k'}, {'l', 'm'}, {'n', 'o'}, {'o', 'p'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindLeftParen,
				token.KindDot,
				token.KindIdentifier, // protein
				token.KindDot,
				token.KindIdentifier, // custom
				token.KindRightParen,
				token.KindDot,
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("option", false, "expected_identifier"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindStr}, token.KindIdentifier),
			},

			content: "option 'test' = true;",
			indices: "a-----bc-----defg---hi",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'h', 'i'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindStr,
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("option", false, "expected_identifier"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindStr}, token.KindIdentifier),
			},

			content: "option ('test') = true;",
			indices: "a-----bcd-----efghi---jk",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'},
				{'g', 'h'}, {'i', 'j'}, {'j', 'k'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindLeftParen,
				token.KindStr,
				token.KindRightParen,
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("option", false, "expected_left_paren"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindEqual}, token.KindLeftParen),
			},

			content: "option (test = true;",
			indices: "a-----bcd---efgh---ij",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'f', 'g'},
				{'h', 'i'}, {'i', 'j'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindLeftParen,
				token.KindIdentifier,
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("option", false, "expected_equal"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindIdentifier}, token.KindEqual),
			},

			content: "option test true;",
			indices: "a-----bc---de---fg",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier,
				token.KindIdentifier, // true
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("option", false, "expected_constant"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindSemicolon}, token.KindInt, token.KindFloat, token.KindIdentifier, token.KindStr, token.KindLeftBrace, token.KindLeftAngle),
			},

			content: "option test = ;",
			indices: "a-----bc---defgh",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier,
				token.KindEqual,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("option", false, "expected_semicolon"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 4, Kind: token.KindEOF}, token.KindSemicolon),
			},

			content: "option test = true",
			indices: "a-----bc---defg---h",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier,
				token.KindEqual,
				token.KindIdentifier,
			},
		},
	}

	runTestCases(t, tests, checkOption, (*impl).parseOption)
}
