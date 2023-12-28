package parser

import (
	"testing"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/parser/internal"
	"github.com/Clement-Jean/protein/token"
)

func checkExtensionRange(t *testing.T, got ast.ExtensionRange, expected ast.ExtensionRange) {
	checkIDs(t, got.ID, expected.ID)
	checkIDs(t, got.OptionsID, expected.OptionsID)

	if len(got.Options) != len(expected.Options) {
		t.Fatalf("expected %d options, got %d", len(expected.Options), len(got.Options))
	}
	for i, opt := range got.Options {
		checkOption(t, opt, expected.Options[i])
	}

	if len(got.Ranges) != len(expected.Ranges) {
		t.Fatalf("expected %d range, got %d", len(expected.Ranges), len(got.Ranges))
	}
	for j, item := range got.Ranges {
		checkIDs(t, item.ID, expected.Ranges[j].ID)
		checkIDs(t, item.Start.ID, expected.Ranges[j].Start.ID)
		checkIDs(t, item.End.ID, expected.Ranges[j].End.ID)
	}
}

func checkExtend(t *testing.T, got ast.Extend, expected ast.Extend) {
	checkIDs(t, got.ID, expected.ID)
	checkIDs(t, got.Name.ID, expected.Name.ID)

	if len(got.Fields) != len(expected.Fields) {
		t.Fatalf("expected %d options, got %d", len(expected.Fields), len(got.Fields))
	}
	for i, field := range got.Fields {
		checkField(t, field, expected.Fields[i])
	}
}

func TestParseExtensionRanges(t *testing.T) {
	tests := []TestCase[ast.ExtensionRange]{
		{
			name: internal.CaseName("extension_range", true, "start"),
			expectedObj: &ast.ExtensionRange{
				ID: 4,
				Ranges: []ast.Range{
					{ID: 1, Start: ast.Integer{ID: 1}, End: ast.Integer{ID: 1}},
				},
			},

			content: "extensions 1000;",
			indices: "a---------bc---de",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'d', 'e'}},
			kinds: []token.Kind{
				token.KindIdentifier, // extensions
				token.KindInt,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("extension_range", true),
			expectedObj: &ast.ExtensionRange{
				ID: 7,
				Ranges: []ast.Range{
					{ID: 6, Start: ast.Integer{ID: 1}, End: ast.Integer{ID: 3}},
				},
			},

			content: "extensions 1000 to max;",
			indices: "a---------bc---de-fg--hi",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'h', 'i'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // extensions
				token.KindInt,
				token.KindIdentifier, // to
				token.KindIdentifier, // max
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("extension_range", true, "multiple"),
			expectedObj: &ast.ExtensionRange{
				ID: 12,
				Ranges: []ast.Range{
					{ID: 10, Start: ast.Integer{ID: 1}, End: ast.Integer{ID: 3}},
					{ID: 11, Start: ast.Integer{ID: 5}, End: ast.Integer{ID: 7}},
				},
			},

			content: "extensions 1 to 10, 10 to max;",
			indices: "a---------bcde-fg-hij-kl-mn--op",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'h', 'i'}, {'j', 'k'}, {'l', 'm'}, {'n', 'o'},
				{'o', 'p'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // extensions
				token.KindInt,
				token.KindIdentifier, // to
				token.KindInt,
				token.KindComma,
				token.KindInt,
				token.KindIdentifier, // to
				token.KindIdentifier, // max
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("extension_range", true, "option"),
			expectedObj: &ast.ExtensionRange{
				ID: 14,
				Ranges: []ast.Range{
					{ID: 11, Start: ast.Integer{ID: 1}, End: ast.Integer{ID: 3}},
				},
				OptionsID: 13,
				Options: []ast.Option{
					{ID: 12, Name: ast.Identifier{ID: 5}, Value: ast.Boolean{ID: 7}},
				},
			},

			content: "extensions 1 to 10 [deprecated = true];",
			indices: "a---------bcde-fg-hij---------klmn---opq",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'j', 'k'}, {'l', 'm'}, {'n', 'o'},
				{'o', 'p'}, {'p', 'q'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // extensions
				token.KindInt,
				token.KindIdentifier, // to
				token.KindInt,
				token.KindLeftSquare,
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindRightSquare,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("extension_range", false, "expected_semicolon"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2}, token.KindSemicolon),
			},

			content: "extensions 1000",
			indices: "a---------bc---d",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}},
			kinds: []token.Kind{
				token.KindIdentifier, // extensions
				token.KindInt,
			},
		},
		{
			name: internal.CaseName("extension_range", false, "unexpected_string"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindStr}, token.KindMax, token.KindInt),
			},

			content: "extensions 1000 to 'max';",
			indices: "a---------bc---de-fg----hi",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'h', 'i'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // extensions
				token.KindInt,
				token.KindIdentifier, // to
				token.KindStr,
				token.KindSemicolon,
			},
		},
	}

	runTestCases(t, tests, checkExtensionRange, (*impl).parseExtensionRange)
}

func TestParseExtends(t *testing.T) {
	tests := []TestCase[ast.Extend]{
		{
			name: internal.CaseName("extend", true),
			expectedObj: &ast.Extend{
				ID: 10, Name: ast.Identifier{ID: 9},
			},

			content: "extend google.protobuf.Empty {}",
			indices: "a-----bc-----de-------fg----hijk",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'},
				{'f', 'g'}, {'g', 'h'}, {'i', 'j'}, {'j', 'k'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // extend
				token.KindIdentifier, // google
				token.KindDot,
				token.KindIdentifier, // protobuf
				token.KindDot,
				token.KindIdentifier, // Empty
				token.KindLeftBrace,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("extend", true, "empty"),
			expectedObj: &ast.Extend{
				ID: 11, Name: ast.Identifier{ID: 10},
			},

			content: "extend google.protobuf.Empty {;}",
			indices: "a-----bc-----de-------fg----hijkl",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'},
				{'f', 'g'}, {'g', 'h'}, {'i', 'j'}, {'j', 'k'},
				{'k', 'l'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // extend
				token.KindIdentifier, // google
				token.KindDot,
				token.KindIdentifier, // protobuf
				token.KindDot,
				token.KindIdentifier, // Empty
				token.KindLeftBrace,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("extend", true, "field"),
			expectedObj: &ast.Extend{
				ID: 16, Name: ast.Identifier{ID: 14},
				Fields: []ast.Field{
					{
						ID:     15,
						TypeID: 7, Type: ast.FieldTypeUint64,
						Name: ast.Identifier{ID: 8},
						Tag:  ast.Integer{ID: 10},
					},
				},
			},

			content: "extend google.protobuf.Empty { uint64 id = 1; }",
			indices: "a-----bc-----de-------fg----hijk-----lm-nopqrstu",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'},
				{'f', 'g'}, {'g', 'h'}, {'i', 'j'}, {'k', 'l'},
				{'m', 'n'}, {'o', 'p'}, {'q', 'r'}, {'r', 's'},
				{'t', 'u'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // extend
				token.KindIdentifier, // google
				token.KindDot,
				token.KindIdentifier, // protobuf
				token.KindDot,
				token.KindIdentifier, // Empty
				token.KindLeftBrace,
				token.KindIdentifier, // uint64
				token.KindIdentifier, // id
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("extend", true, "option"),
			expectedObj: &ast.Extend{
				ID: 16, Name: ast.Identifier{ID: 14},
				Options: []ast.Option{
					{
						ID:    15,
						Name:  ast.Identifier{ID: 13},
						Value: ast.Identifier{ID: 15},
					},
				},
			},

			content: "extend google.protobuf.Empty { option deprecated = true; }",
			indices: "a-----bc-----de-------fg----hijk-----lm---------nopq---rstu",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'},
				{'f', 'g'}, {'g', 'h'}, {'i', 'j'}, {'k', 'l'},
				{'m', 'n'}, {'o', 'p'}, {'q', 'r'}, {'r', 's'},
				{'t', 'u'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // extend
				token.KindIdentifier, // google
				token.KindDot,
				token.KindIdentifier, // protobuf
				token.KindDot,
				token.KindIdentifier, // Empty
				token.KindLeftBrace,
				token.KindIdentifier, // option
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("extend", false, "expected_identifier"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindLeftBrace}, token.KindIdentifier),
			},

			content: "extend {}",
			indices: "a-----bcde",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'d', 'e'}},
			kinds: []token.Kind{
				token.KindIdentifier, // extend
				token.KindLeftBrace,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("extend", false, "expected_left_brace"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindLeftSquare}, token.KindLeftBrace),
			},

			content: "extend Test [}",
			indices: "a-----bc---defg",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'}},
			kinds: []token.Kind{
				token.KindIdentifier, // extend
				token.KindIdentifier, // Test
				token.KindLeftSquare,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("extend", false, "expected_right_brace"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindEOF}, token.KindRightBrace),
			},

			content: "extend Test {",
			indices: "a-----bc---def",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}},
			kinds: []token.Kind{
				token.KindIdentifier, // extend
				token.KindIdentifier, // Test
				token.KindLeftBrace,
			},
		},
		{
			name: internal.CaseName("extend", false, "expected_field"),
			expectedObj: &ast.Extend{
				ID:   6,
				Name: ast.Identifier{ID: 1},
			},
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindInt}, token.KindOption, token.KindField),
			},

			content: "extend Test {2}",
			indices: "a-----bc---defgh",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'}, {'g', 'h'}},
			kinds: []token.Kind{
				token.KindIdentifier, // extend
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindInt,
				token.KindRightBrace,
			},
		},
	}

	runTestCases(t, tests, checkExtend, (*impl).parseExtend)
}
