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

func TestParseExtensions(t *testing.T) {
	tests := []TestCase[ast.ExtensionRange]{
		{
			name: internal.CaseName("extension_range", true, "start"),
			expectedObj: ast.ExtensionRange{
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
			expectedObj: ast.ExtensionRange{
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
			expectedObj: ast.ExtensionRange{
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
			expectedObj: ast.ExtensionRange{
				ID: 14,
				Ranges: []ast.Range{
					{ID: 11, Start: ast.Integer{ID: 1}, End: ast.Integer{ID: 3}},
				},
				OptionsID: 13,
				Options: []ast.Option{
					{ID: 12, Name: ast.Identifier{ID: 5}, Value: ast.Identifier{ID: 7}},
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
