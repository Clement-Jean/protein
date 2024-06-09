package lexer

import (
	"strings"
	"testing"
)

var tokenTests = []struct {
	name string
	s    string
	pred func(TokenKind) bool
}{
	{
		name: "identifier",
		s:    "myname syntax",
		pred: TokenKind.IsIdentifier,
	},
	{
		name: "opening symbol",
		s:    "{[(<",
		pred: TokenKind.IsOpeningSymbol},
	{
		name: "closing symbol",
		s:    "}])>",
		pred: TokenKind.IsClosingSymbol,
	},
}

func TestTokens(t *testing.T) {
	for _, test := range tokenTests {
		t.Run(test.name, func(t *testing.T) {
			l, err := NewFromReader(strings.NewReader(test.s))

			if err != nil {
				t.Fatal(err)
			}

			tb, errs := l.Lex()
			if len(errs) != 0 {
				t.Fatal(errs)
			}

			for _, tok := range tb.TokenInfos {
				if tok.Kind == TokenKindBOF || tok.Kind == TokenKindEOF {
					continue
				}

				if !test.pred(tok.Kind) {
					t.Fatalf("expected %s to be %s", tok.Kind, test.name)
				}
			}
		})
	}
}

func TestKeywordLiteral(t *testing.T) {
	if len(literals) != len(kinds) {
		t.Fatalf("got %d literals and %d kinds", len(literals), len(kinds))
	}

	for i, literal := range literals {
		if literal != kinds[i].String() {
			t.Fatalf("expected literal %q, got %q", literal, kinds[i].String())
		}
	}
}
