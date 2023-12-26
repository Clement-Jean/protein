package parser

import (
	"testing"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/parser/internal"
	"github.com/Clement-Jean/protein/token"
)

func checkPackage(t *testing.T, got ast.Package, expected ast.Package) {
	t.Helper()
	checkIDs(t, got.ID, expected.ID)
	checkIdentifier(t, got.Value, expected.Value)
}

func TestPackage(t *testing.T) {
	tests := []TestCase[ast.Package]{
		{
			name:        internal.CaseName("package", true, "identifier"),
			expectedObj: &ast.Package{ID: 4, Value: ast.Identifier{ID: 1}},

			content: "package google;",
			indices: "a------bc-----de",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'d', 'e'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier,
				token.KindSemicolon,
			},
		},
		{
			name:        internal.CaseName("package", true, "full_identifier"),
			expectedObj: &ast.Package{ID: 7, Value: ast.Identifier{ID: 6, Parts: []token.UniqueID{1, 3}}},

			content: "package google.protobuf;",
			indices: "a------bc-----de-------fg",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'}, {'f', 'g'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier,
				token.KindDot,
				token.KindIdentifier,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("package", false, "expect_identifier"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindStr}, token.KindIdentifier),
			},

			content: "package 'google';",
			indices: "a------bc-------de",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'d', 'e'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindStr,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("package", false, "expect_semicolon"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindEOF}, token.KindSemicolon),
			},

			content: "package 'google'",
			indices: "a------bc-------d",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier,
			},
		},
	}

	wrap := func(p *impl) (ast.Package, []error) {
		pkg, err := p.parsePackage()
		return pkg, internal.EmptyErrorSliceIfNil(err)
	}
	runTestCases(t, tests, checkPackage, wrap)
}
