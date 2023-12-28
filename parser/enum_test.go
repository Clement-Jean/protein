package parser

import (
	"testing"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/parser/internal"
	"github.com/Clement-Jean/protein/token"
)

func checkEnumValue(t *testing.T, got ast.EnumValue, expected ast.EnumValue) {
	checkIDs(t, got.ID, expected.ID)
	checkIDs(t, got.Name.ID, expected.Name.ID)
	checkIDs(t, got.OptionsID, expected.OptionsID)

	if len(got.Options) != len(expected.Options) {
		t.Fatalf("expected %d options, got %d", len(expected.Options), len(got.Options))
	}
	for j, opt := range got.Options {
		checkOption(t, opt, expected.Options[j])
	}
}

func checkEnum(t *testing.T, got ast.Enum, expected ast.Enum) {
	checkIDs(t, got.ID, expected.ID)
	checkIDs(t, got.Name.ID, expected.Name.ID)

	if len(got.Options) != len(expected.Options) {
		t.Fatalf("expected %d options, got %d", len(expected.Options), len(got.Options))
	}
	for i, opt := range got.Options {
		checkOption(t, opt, expected.Options[i])
	}

	if len(got.ReservedTags) != len(expected.ReservedTags) {
		t.Fatalf("expected %d reserved tags, got %d", len(expected.ReservedTags), len(got.ReservedTags))
	}
	for i, opt := range got.ReservedTags {
		checkReservedTags(t, opt, expected.ReservedTags[i])
	}
	if len(got.ReservedNames) != len(expected.ReservedNames) {
		t.Fatalf("expected %d reserved names, got %d", len(expected.ReservedNames), len(got.ReservedNames))
	}
	for i, opt := range got.ReservedNames {
		checkReservedNames(t, opt, expected.ReservedNames[i])
	}

	if len(got.Values) != len(expected.Values) {
		t.Fatalf("expected %d values, got %d", len(expected.Values), len(got.Values))
	}
	for i, opt := range got.Values {
		checkEnumValue(t, opt, expected.Values[i])
	}
}

func TestEnum(t *testing.T) {
	tests := []TestCase[ast.Enum]{
		{
			name:        internal.CaseName("enum", true),
			expectedObj: &ast.Enum{ID: 5, Name: ast.Identifier{ID: 1}},

			content: "enum Test {}",
			indices: "a---bc---defg",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindRightBrace,
			},
		},
		{
			name:        internal.CaseName("enum", true, "empty_statement"),
			expectedObj: &ast.Enum{ID: 6, Name: ast.Identifier{ID: 1}},

			content: "enum Test {;}",
			indices: "a---bc---defgh",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'}, {'g', 'h'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("enum", true, "option"),
			expectedObj: &ast.Enum{
				ID: 11, Name: ast.Identifier{ID: 1},
				Options: []ast.Option{{
					ID: 10, Name: ast.Identifier{ID: 4}, Value: ast.Boolean{ID: 6},
				}},
			},

			content: "enum Test { option deprecated = true; }",
			indices: "a---bc---defg-----hi---------jklm---nopq",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'n', 'o'},
				{'p', 'q'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier,
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("enum", true, "value"),
			expectedObj: &ast.Enum{
				ID: 10, Name: ast.Identifier{ID: 1},
				Values: []ast.EnumValue{{
					ID: 9, Name: ast.Identifier{ID: 3}, Tag: ast.Integer{ID: 5},
				}},
			},

			content: "enum Test { TEST_UNSPECIFIED = 0; }",
			indices: "a---bc---defg---------------hijklmno",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'l', 'm'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // TEST_UNSPECIFIED
				token.KindEqual,
				token.KindInt, // 0
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("enum", true, "reserved_tag"),
			expectedObj: &ast.Enum{
				ID: 8, Name: ast.Identifier{ID: 1},
				ReservedTags: []ast.ReservedTags{
					{ID: 4, Items: []ast.Range{{ID: 4, Start: ast.Integer{ID: 4}, End: ast.Integer{ID: 4}}}},
				},
			},

			content: "enum Test { reserved 1; }",
			indices: "a---bc---defg-------hijklm",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'j', 'k'}, {'l', 'm'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // reserved
				token.KindInt,        // 1
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("enum", true, "reserved_name"),
			expectedObj: &ast.Enum{
				ID: 8, Name: ast.Identifier{ID: 1},
				ReservedNames: []ast.ReservedNames{
					{ID: 4, Items: []ast.String{{ID: 4}}},
				},
			},

			content: "enum Test { reserved '1'; }",
			indices: "a---bc---defg-------hi--jklm",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'j', 'k'}, {'l', 'm'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // reserved
				token.KindStr,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("enum", true, "value_option"),
			expectedObj: &ast.Enum{
				ID: 17, Name: ast.Identifier{ID: 1},
				Values: []ast.EnumValue{{
					ID: 16, Name: ast.Identifier{ID: 3}, Tag: ast.Integer{ID: 5},
					OptionsID: 15,
					Options:   []ast.Option{{ID: 14, Name: ast.Identifier{ID: 7}, Value: ast.Boolean{ID: 9}}},
				}},
			},

			content: "enum Test { TEST_UNSPECIFIED = 0 [deprecated = true]; }",
			indices: "a---bc---defg---------------hijklmn---------opqr---stuvw",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'n', 'o'},
				{'p', 'q'}, {'r', 's'}, {'s', 't'}, {'t', 'u'},
				{'v', 'w'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // TEST_UNSPECIFIED
				token.KindEqual,
				token.KindInt, // 0
				token.KindLeftSquare,
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindRightSquare,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("enum", true, "value_options"),
			expectedObj: &ast.Enum{
				ID: 22, Name: ast.Identifier{ID: 1},
				Values: []ast.EnumValue{{
					ID: 21, Name: ast.Identifier{ID: 3}, Tag: ast.Integer{ID: 5},
					OptionsID: 20,
					Options: []ast.Option{
						{ID: 18, Name: ast.Identifier{ID: 7}, Value: ast.Boolean{ID: 9}},
						{ID: 19, Name: ast.Identifier{ID: 11}, Value: ast.Boolean{ID: 13}},
					},
				}},
			},

			content: "enum Test { TEST = 0 [deprecated = true, a = true]; }",
			indices: "a---bc---defg---hijklmn---------opqr---stuvwxy---z1234",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'n', 'o'},
				{'p', 'q'}, {'r', 's'}, {'s', 't'}, {'u', 'v'},
				{'w', 'x'}, {'y', 'z'}, {'z', '1'}, {'1', '2'},
				{'3', '4'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // TEST
				token.KindEqual,
				token.KindInt, // 0
				token.KindLeftSquare,
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindComma,
				token.KindIdentifier, // a
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindRightSquare,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("enum", false, "option_expected_right_square"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 10, Kind: token.KindEOF}, token.KindRightSquare),
				gotUnexpected(&token.Token{ID: 10, Kind: token.KindEOF}, token.KindSemicolon),
				gotUnexpected(&token.Token{ID: 10, Kind: token.KindEOF}, token.KindRightBrace),
			},

			content: "enum Test { TEST_UNSPECIFIED = 0 [deprecated = true",
			indices: "a---bc---defg---------------hijklmn---------opqr---s",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'n', 'o'},
				{'p', 'q'}, {'r', 's'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // TEST_UNSPECIFIED
				token.KindEqual,
				token.KindInt, // 0
				token.KindLeftSquare,
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
			},
		},
		{
			name: internal.CaseName("enum", false, "expected_left_brace"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindLeftSquare}, token.KindLeftBrace),
			},

			content: "enum Test [ option deprecated = true; }",
			indices: "a---bc---defg-----hi---------jklm---nopq",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'n', 'o'},
				{'p', 'q'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftSquare,
				token.KindIdentifier,
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("enum", false, "expected_right_brace"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 8, Kind: token.KindRightSquare}, token.KindOption, token.KindReserved, token.KindIdentifier),
				gotUnexpected(&token.Token{ID: 9, Kind: token.KindEOF}, token.KindRightBrace),
			},

			content: "enum Test { option deprecated = true; ]",
			indices: "a---bc---defg-----hi---------jklm---nopq",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'n', 'o'},
				{'p', 'q'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier,
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindSemicolon,
				token.KindRightSquare,
			},
		},
		{
			name: internal.CaseName("enum", false, "expected_left_brace_eof"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 8, Kind: token.KindEOF}, token.KindRightBrace),
			},

			content: "enum Test { option deprecated = true;",
			indices: "a---bc---defg-----hi---------jklm---no",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier,
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("enum", false, "expected_identifier"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindLeftBrace}, token.KindIdentifier),
			},

			content: "enum {}",
			indices: "a---bcde",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindLeftBrace,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("enum", false, "value_expected_identifier"),
			expectedObj: &ast.Enum{
				ID:   9,
				Name: ast.Identifier{ID: 1},
			},
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindStr}, token.KindOption, token.KindReserved, token.KindIdentifier),
			},

			content: "enum Test { 'a' = 1; }",
			indices: "a---bc---defg--hijklmno",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'l', 'm'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindStr,
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("enum", false, "value_expected_int"),
			expectedObj: &ast.Enum{
				ID:   9,
				Name: ast.Identifier{ID: 1},
			},
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 5, Kind: token.KindStr}, token.KindInt),
			},

			content: "enum Test { a = '1'; }",
			indices: "a---bc---defghijk--lmno",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'l', 'm'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier,
				token.KindEqual,
				token.KindStr,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("enum", false, "value_expected_equal"),
			expectedObj: &ast.Enum{
				ID:   7,
				Name: ast.Identifier{ID: 1},
			},
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 4, Kind: token.KindInt}, token.KindEqual),
			},

			content: "enum Test { a 1 }",
			indices: "a---bc---defghijkl",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier,
				token.KindInt,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("enum", false, "value_expected_left_semicolon"),
			expectedObj: &ast.Enum{
				ID:   8,
				Name: ast.Identifier{ID: 1},
			},
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 6, Kind: token.KindRightBrace}, token.KindSemicolon),
			},

			content: "enum Test { a = 1 }",
			indices: "a---bc---defghijklmn",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier,
				token.KindEqual,
				token.KindInt,
				token.KindRightBrace,
			},
		},
	}

	runTestCases(t, tests, checkEnum, (*impl).parseEnum)
}
