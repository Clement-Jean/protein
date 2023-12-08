package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/codemap"
	"github.com/Clement-Jean/protein/internal/bytes"
	"github.com/Clement-Jean/protein/internal/span"
	"github.com/Clement-Jean/protein/parser/internal"
	"github.com/Clement-Jean/protein/token"
)

type UnderTest interface {
	ast.Syntax
}

type TestCase[T UnderTest] struct {
	name string

	expectedObj  T
	expectedErrs []error

	content string
	indices string

	locs  [][2]rune
	kinds []token.Kind

	keepFirstToken bool
}

func checkErrs(t *testing.T, errs []error, expectedErrs []error) {
	t.Helper()
	if len(errs) != len(expectedErrs) {
		t.Fatalf("expected %d errors, got %d", len(expectedErrs), len(errs))
	}

	for i, err := range errs {
		var got *Error
		var expected *Error

		if errors.As(err, &got) && errors.As(expectedErrs[i], &expected) {
			checkIDs(t, got.ID, expected.ID)
		}
		if strings.Compare(got.Error(), expected.Error()) != 0 {
			t.Fatalf("expected error '%s', got '%s'", expected.Error(), got.Error())
		}
	}
}

func checkIDs(t *testing.T, got token.UniqueID, expected token.UniqueID) {
	t.Helper()
	if got != expected {
		t.Fatalf("expected id %d, got %d", expected, got)
	}
}

func runTestCases[T UnderTest](
	t *testing.T,
	tests []TestCase[T],
	onObj func(t *testing.T, got T, expected T),
	parseFn func(p *impl) (T, error),
) {
	t.Helper()
	cm := codemap.New()
	dummy := "test.proto"
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fm := cm.Insert(dummy, bytes.FromString(test.content))
			ref := internal.ReferenceString(test.indices, '-')
			spans := internal.MakeSpansFromIndices(ref, test.locs)

			if len(spans) != len(test.kinds) {
				t.Fatalf("have %d kinds and %d locs", len(test.kinds), len(spans))
			}

			tokens := fm.RegisterTokens(test.kinds, spans)
			p := New(tokens, fm)
			i := p.(*impl)
			if !test.keepFirstToken {
				i.nextToken()
			}
			obj, errs := parseFn(i)

			fm.PrintItems()

			if test.expectedErrs != nil {
				checkErrs(t, []error{errs}, test.expectedErrs)
			} else {
				onObj(t, obj, test.expectedObj)
			}

			cm.Remove(dummy)
		})
	}
}

func TestPeekEOF(t *testing.T) {
	fm := &codemap.FileMap{}
	p := New(nil, fm).(*impl)
	tok := p.peek()

	if tok != nil {
		t.Fatalf("expected nil, got %v", tok)
	}
}

func TestPeekSkipSpacesAndComments(t *testing.T) {
	fm := &codemap.FileMap{}
	tokens := []token.Token{
		token.Token{
			ID:   1,
			Kind: token.KindComment,
		},
		token.Token{
			ID:   2,
			Kind: token.KindSpace,
		},
		token.Token{
			ID:   3,
			Kind: token.KindIdentifier,
		},
	}
	p := New(tokens, fm).(*impl)
	tok := p.peek()

	if tok.ID != 3 {
		t.Fatalf("expected ID 3, got %d", tok.ID)
	}

	if tok.Kind != token.KindIdentifier {
		t.Fatalf("expected Identifier, got %s", tok.Kind.String())
	}
}

func TestNextEOF(t *testing.T) {
	fm := &codemap.FileMap{}
	p := New(nil, fm).(*impl)
	tok := p.nextToken()

	if tok != nil {
		t.Fatalf("expected nil, got %v", tok)
	}
}

func TestNextSkipSpacesAndComments(t *testing.T) {
	fm := &codemap.FileMap{}
	tokens := []token.Token{
		token.Token{
			ID:   1,
			Kind: token.KindComment,
		},
		token.Token{
			ID:   2,
			Kind: token.KindSpace,
		},
		token.Token{
			ID:   3,
			Kind: token.KindIdentifier,
		},
	}
	p := New(tokens, fm).(*impl)
	tok := p.nextToken()

	if tok.ID != 3 {
		t.Fatalf("expected ID 3, got %d", tok.ID)
	}

	if tok.Kind != token.KindIdentifier {
		t.Fatalf("expected Identifier, got %s", tok.Kind.String())
	}
}

func TestParseHandleUnknownIdentifier(t *testing.T) {
	content := []byte("unknown")
	cm := codemap.New()
	fm := cm.Insert("test.proto", content)
	kinds := []token.Kind{token.KindIdentifier}
	spans := []span.Span{span.Span{Start: 0, End: 7}}
	tokens := fm.RegisterTokens(kinds, spans)
	p := New(tokens, fm)
	_, errs := p.Parse()

	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}

	expectErr := gotUnexpected(&tokens[0], token.KindSyntax)
	if strings.Compare(errs[0].Error(), expectErr.Error()) != 0 {
		t.Fatalf("expected error '%s', got '%s'", expectErr.Error(), errs[0].Error())
	}
}
