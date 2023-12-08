package parser

import (
	"testing"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/parser/internal"
	"github.com/Clement-Jean/protein/token"
)

func TestParseFullyQualifiedIdentifier(t *testing.T) {
	tests := []TestCase[ast.Identifier]{
		{
			keepFirstToken: true,
			name:           internal.CaseName("fully_qualified_identifier", true),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindStr}, token.KindIdentifier),
			},

			content: "test.'test'.test",
			indices: "a---bc-----de---f",
			locs:    [][2]rune{{'a', 'b'}, {'b', 'c'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'}},
			kinds: []token.Kind{
				token.KindIdentifier, // test
				token.KindDot,
				token.KindStr,
				token.KindDot,
				token.KindIdentifier, // test
			},
		},
	}

	runTestCases(t, tests, checkIdentifier, (*impl).parseFullyQualifiedIdentifier)
}
