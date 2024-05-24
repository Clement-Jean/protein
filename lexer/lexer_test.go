package lexer_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Clement-Jean/protein/lexer"
)

type TestCase struct {
	name       string
	input      string
	tokenInfos []lexer.TokenInfo
	lineInfos  []lexer.LineInfo
	errs       []error
}

func runTestCases(t *testing.T, tests []TestCase) {
	t.Helper()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l, err := lexer.NewFromReader(strings.NewReader(test.input))
			if err != nil {
				t.Fatal(err)
			}

			tb := l.Lex()

			if !reflect.DeepEqual(test.tokenInfos, tb.TokenInfos) {
				t.Fatalf(`
expected token infos: %+v
                 got: %+v`, test.tokenInfos, tb.TokenInfos)
			}

			if !reflect.DeepEqual(test.lineInfos, tb.LineInfos) {
				t.Fatalf(`
expected line infos: %+v
                got: %+v`, test.lineInfos, tb.LineInfos)
			}
		})
	}
}

func symbolsTestCase() TestCase {
	const start = uint8(lexer.TokenKindUnderscore)
	const end = uint8(lexer.TokenKindSlash) + 1
	const symbols = "_=,:;.{}[]()<>/"

	expectedTokenInfos := make([]lexer.TokenInfo, len(symbols)+2)
	expectedTokenInfos[0].Kind = lexer.TokenKindBOF
	for i := uint8(0); i < uint8(len(symbols)); i++ {
		info := &expectedTokenInfos[i+1]
		info.Kind = lexer.TokenKind(start + i)
		info.Column = uint32(i)
	}
	info := &expectedTokenInfos[len(symbols)+1]
	info.Kind = lexer.TokenKindEOF
	info.Column = uint32(len(symbols))

	return TestCase{
		name:       "symbols",
		input:      symbols,
		tokenInfos: expectedTokenInfos,
		lineInfos:  []lexer.LineInfo{{Len: uint32(len(symbols))}},
	}
}

func TestLexer(t *testing.T) {
	tests := []TestCase{
		symbolsTestCase(),
	}

	runTestCases(t, tests)
}
