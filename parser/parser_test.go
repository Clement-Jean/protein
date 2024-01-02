package parser

import (
	"fmt"
	"log"
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
	ast.Identifier | ast.Syntax | ast.Edition | ast.Package | ast.Import | ast.Option | []ast.Option |
		ast.TextField | ast.TextMessage | ast.TextScalarList | ast.TextMessageList | ast.Enum | ast.EnumValue |
		ast.ReservedTags | ast.ReservedNames | ast.Message | ast.Field | ast.Oneof | ast.ExtensionRange |
		ast.Service | ast.Rpc | ast.Extend
}

type TestCase[T UnderTest] struct {
	name string

	expectedObj  *T
	expectedErrs []error

	content string
	indices string

	locs  [][2]rune
	kinds []token.Kind

	keepFirstToken bool
}

func checkErrs(t *testing.T, errs, expectedErrs []error) {
	t.Helper()
	if len(errs) != len(expectedErrs) {
		fmt.Printf("%v\n", errs)
		t.Fatalf("expected %d errors, got %d", len(expectedErrs), len(errs))
	}

	for i := range errs {
		got := errs[i].(*Error)
		expected := expectedErrs[i].(*Error)

		checkIDs(t, got.ID, expected.ID)
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
	checkObj func(t *testing.T, got T, expected T),
	parseFn func(p *impl) (T, []error),
) {
	t.Helper()
	cm := codemap.New()
	dummy := "test.proto"
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log.Println("----", test.name, "----")
			fm := cm.Insert(dummy, bytes.FromString(test.content))
			ref := internal.ReferenceString(test.indices, '-')
			spans := internal.MakeSpansFromIndices(ref, test.locs)

			if len(test.kinds) != 0 && test.kinds[len(test.kinds)-1] != token.KindEOF {
				// add EOF if not in kinds
				test.kinds = append(test.kinds, token.KindEOF)
			}

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

			checkErrs(t, errs, test.expectedErrs)
			if test.expectedObj != nil {
				checkObj(t, obj, *test.expectedObj)
			}
			cm.Remove(dummy)
		})
	}
}

func TestParseHandleUnknownIdentifier(t *testing.T) {
	content := []byte("unknown")
	cm := codemap.New()
	fm := cm.Insert("test.proto", content)
	kinds := []token.Kind{token.KindIdentifier}
	spans := []span.Span{{Start: 0, End: 7}}
	tokens := fm.RegisterTokens(kinds, spans)
	p := New(tokens, fm)
	_, errs := p.Parse()

	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}

	expectedKinds := []token.Kind{
		token.KindSyntax, token.KindEdition,
		token.KindPackage, token.KindImport, token.KindOption,
		token.KindMessage, token.KindEnum, token.KindService, token.KindExtend,
	}
	expectErr := gotUnexpected(&tokens[0], expectedKinds...)
	if strings.Compare(errs[0].Error(), expectErr.Error()) != 0 {
		t.Fatalf("expected error '%s', got '%s'", expectErr.Error(), errs[0].Error())
	}
}
