package parser

import (
	"testing"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/parser/internal"
	"github.com/Clement-Jean/protein/token"
)

func checkReservedTags(t *testing.T, got ast.ReservedTags, expected ast.ReservedTags) {
	if got.Items == nil {
		t.Fatalf("expected Reserved.Items not nil")
	}

	checkIDs(t, got.ID, expected.ID)

	if len(got.Items) != len(expected.Items) {
		t.Fatalf("expected %d reserved items, got %d", len(expected.Items), len(got.Items))
	}
	for j, item := range got.Items {
		checkIDs(t, item.ID, expected.Items[j].ID)
		checkIDs(t, item.Start.ID, expected.Items[j].Start.ID)
		checkIDs(t, item.End.ID, expected.Items[j].End.ID)
	}
}

func checkReservedNames(t *testing.T, got ast.ReservedNames, expected ast.ReservedNames) {

	if got.Items == nil {
		t.Fatalf("expected Reserved.Items not nil")
	}

	checkIDs(t, got.ID, expected.ID)

	if len(got.Items) != len(expected.Items) {
		t.Fatalf("expected %d reserved items, got %d", len(expected.Items), len(got.Items))
	}
	for j, item := range got.Items {
		checkIDs(t, item.ID, expected.Items[j].ID)
	}
}

func TestReservedTags(t *testing.T) {
	tests := []TestCase[ast.ReservedTags]{
		{
			name: internal.CaseName("reserved_tag", true),
			expectedObj: &ast.ReservedTags{
				ID: 1, Items: []ast.Range{{ID: 1, Start: ast.Integer{ID: 1}, End: ast.Integer{ID: 1}}},
			},

			content: "reserved 1;",
			indices: "a-------bcde",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // reserved
				token.KindInt,        // 1
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("reserved_tag", true, "multiple_range"),
			expectedObj: &ast.ReservedTags{
				ID: 12, Items: []ast.Range{
					{ID: 10, Start: ast.Integer{ID: 1}, End: ast.Integer{ID: 3}},
					{ID: 11, Start: ast.Integer{ID: 5}, End: ast.Integer{ID: 7}},
				},
			},

			content: "reserved 1 to 10, 11 to max;",
			indices: "a-------bcde-fg-hij-kl-mn--op",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'h', 'i'}, {'j', 'k'}, {'l', 'm'}, {'n', 'o'},
				{'o', 'p'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // reserved
				token.KindInt,        // 1
				token.KindIdentifier, // to
				token.KindInt,        // 10
				token.KindComma,
				token.KindInt,        // 11
				token.KindIdentifier, // to
				token.KindIdentifier, // max
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("reserved_tag", false, "expected_int"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindIdentifier}, token.KindInt),
			},

			content: "reserved test;",
			indices: "a-------bc---de",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // reserved
				token.KindIdentifier,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("reserved_tag", false, "expected_int_after_comma"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindStr}, token.KindInt),
			},

			content: "reserved 1, 'test';",
			indices: "a-------bcdef-----gh",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'f', 'g'},
				{'g', 'h'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // reserved
				token.KindInt,
				token.KindComma,
				token.KindStr,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("reserved_tag", false, "expected_int_after_to"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindStr}, token.KindMax, token.KindInt),
			},

			content: "reserved 1 to 'test';",
			indices: "a-------bcde-fg-----hi",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'h', 'i'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // reserved
				token.KindInt,
				token.KindIdentifier,
				token.KindStr,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("reserved_tag", false, "expected_max_after_to"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindIdentifier}, token.KindMax, token.KindInt),
			},

			content: "reserved 1 to min;",
			indices: "a-------bcde-fg--hi",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'h', 'i'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // reserved
				token.KindInt,
				token.KindIdentifier,
				token.KindIdentifier, // min
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("reserved_tag", false, "expected_semicolon"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindEOF}, token.KindSemicolon),
			},

			content: "reserved 1",
			indices: "a-------bcd",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // reserved
				token.KindInt,
			},
		},
	}

	wrap := func(p *impl) (ast.ReservedTags, []error) {
		tags, err := p.parseReservedTags()
		return tags, internal.EmptyErrorSliceIfNil(err)
	}
	runTestCases(t, tests, checkReservedTags, wrap)
}

func TestReservedNames(t *testing.T) {
	tests := []TestCase[ast.ReservedNames]{
		{
			name: internal.CaseName("reserved_name", true),
			expectedObj: &ast.ReservedNames{
				ID: 1, Items: []ast.String{{ID: 1}},
			},

			content: "reserved 'test';",
			indices: "a-------bc-----de",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // reserved
				token.KindStr,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("reserved_name", true, "multiple_names"),
			expectedObj: &ast.ReservedNames{
				ID: 6, Items: []ast.String{{ID: 1}, {ID: 3}},
			},

			content: "reserved 'test', 'another';",
			indices: "a-------bc-----def--------gh",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'f', 'g'},
				{'g', 'h'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // reserved
				token.KindStr,
				token.KindComma,
				token.KindStr,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("reserved_name", false, "expected_string"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindIdentifier}, token.KindStr),
			},

			content: "reserved test;",
			indices: "a-------bc---de",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // reserved
				token.KindIdentifier,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("reserved_name", false, "expected_str_after_comma"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindInt}, token.KindStr),
			},

			content: "reserved 'test', 1;",
			indices: "a-------bc----defghi",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'h', 'i'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // reserved
				token.KindStr,
				token.KindComma,
				token.KindInt,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("reserved_name", false, "expected_semicolon"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindEOF}, token.KindSemicolon),
			},

			content: "reserved 'test'",
			indices: "a-------bc-----d",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // reserved
				token.KindStr,
			},
		},
	}

	wrap := func(p *impl) (ast.ReservedNames, []error) {
		names, err := p.parseReservedNames()
		return names, internal.EmptyErrorSliceIfNil(err)
	}
	runTestCases(t, tests, checkReservedNames, wrap)
}
