package parser

import (
	"testing"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/parser/internal"
	"github.com/Clement-Jean/protein/token"
)

func checkOneof(t *testing.T, got ast.Oneof, expected ast.Oneof) {
	checkIDs(t, got.ID, expected.ID)
	checkIDs(t, got.Name.ID, expected.Name.ID)

	if len(got.Options) != len(expected.Options) {
		t.Fatalf("expected %d options, got %d", len(expected.Options), len(got.Options))
	}
	for i, opt := range got.Options {
		checkOption(t, opt, expected.Options[i])
	}

	if len(got.Fields) != len(expected.Fields) {
		t.Fatalf("expected %d field, got %d", len(expected.Fields), len(got.Fields))
	}
	for i, opt := range got.Fields {
		checkField(t, opt, expected.Fields[i])
	}
}

func TestParseOneof(t *testing.T) {
	tests := []TestCase[ast.Oneof]{
		{
			keepFirstToken: true,
			name:           internal.CaseName("oneof", true, "empty_statement"),
			expectedObj: &ast.Oneof{
				ID:   6,
				Name: ast.Identifier{ID: 1},
			},

			content: "oneof Test { ; }",
			indices: "a----bc---defghij",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // oneof
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("oneof", true),
			expectedObj: &ast.Oneof{
				ID:   17,
				Name: ast.Identifier{ID: 1},
				Fields: []ast.Field{
					{ID: 15, Type: ast.FieldTypeUint64, TypeID: 3, Name: ast.Identifier{ID: 4}, Tag: ast.Integer{ID: 6}},
					{ID: 16, Type: ast.FieldTypeString, TypeID: 8, Name: ast.Identifier{ID: 9}, Tag: ast.Integer{ID: 11}},
				},
			},

			content: "oneof Test { uint64 id = 1; string uuid = 2; }",
			indices: "a----bc---defg-----hi-jklmnop-----qr---stuvwxyz",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'n', 'o'},
				{'p', 'q'}, {'r', 's'}, {'t', 'u'}, {'v', 'w'},
				{'w', 'x'}, {'y', 'z'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // oneof
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // uint64
				token.KindIdentifier, // id
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
				token.KindIdentifier, // string
				token.KindIdentifier, // uuid
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("oneof", true, "option"),
			expectedObj: &ast.Oneof{
				ID: 11, Name: ast.Identifier{ID: 1},
				Options: []ast.Option{{
					ID: 10, Name: ast.Identifier{ID: 4}, Value: ast.Boolean{ID: 6},
				}},
			},

			content: "oneof Test { option deprecated = true; }",
			indices: "a----bc---defg-----hi---------jklm---nopq",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'n', 'o'},
				{'p', 'q'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // oneof
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
			name: internal.CaseName("oneof", false, "expected_identifier"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindLeftBrace}, token.KindIdentifier),
			},

			content: "oneof {}",
			indices: "a----bcde",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // oneof
				token.KindLeftBrace,
				token.KindRightBrace,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("oneof", false, "expected_left_brace"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindLeftSquare}, token.KindLeftBrace),
			},

			content: "oneof Test [}",
			indices: "a----bc---defg",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'f', 'g'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // oneof
				token.KindIdentifier, // Test
				token.KindLeftSquare,
				token.KindRightBrace,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("oneof", false, "expected_right_brace"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindEOF}, token.KindRightBrace),
			},

			content: "oneof Test {",
			indices: "a----bc---def",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // oneof
				token.KindIdentifier, // Test
				token.KindLeftBrace,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("oneof", false, "unexpected_int"),
			expectedObj: &ast.Oneof{
				ID:   6,
				Name: ast.Identifier{ID: 1},
			},
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindInt}, token.KindOption, token.KindIdentifier),
			},

			content: "oneof Test { 2 }",
			indices: "a----bc---defghij",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // oneof
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindInt,
				token.KindRightBrace,
			},
		},
	}

	runTestCases(t, tests, checkOneof, (*impl).parseOneof)
}
