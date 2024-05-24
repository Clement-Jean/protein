package lexer_test

import (
	"errors"
	"fmt"
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

			tb, errs := l.Lex()

			if !reflect.DeepEqual(test.errs, errs) {
				t.Fatalf(`
expected errors: %+v
            got: %+v`, test.errs, errs)
			}

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
		{
			name:  "invalid",
			input: "&",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindError},
				{Kind: lexer.TokenKindEOF, Column: 1},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 1},
			},
			errs: []error{errors.New("invalid char '&'")},
		},
		{
			name:  "invalid_utf8",
			input: "🙈",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindError},
				{Kind: lexer.TokenKindError, Column: 1},
				{Kind: lexer.TokenKindError, Column: 2},
				{Kind: lexer.TokenKindError, Column: 3},
				{Kind: lexer.TokenKindEOF, Column: 4},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 4},
			},
			errs: []error{
				fmt.Errorf("invalid char %q", "🙈"[0]),
				fmt.Errorf("invalid char %q", "🙈"[1]),
				fmt.Errorf("invalid char %q", "🙈"[2]),
				fmt.Errorf("invalid char %q", "🙈"[3]),
			},
		},
		{
			name:  "skip_utf8_bom",
			input: string([]byte{0xEF, 0xBB, 0xBF}),
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF, Column: 3},
				{Kind: lexer.TokenKindEOF, Column: 3},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 3, Len: 0},
			},
		},
		{
			name:  "spaces",
			input: "\t\n\v\f\n ",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindEOF, LineIdx: 2, Column: 1},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 1},
				{Start: 2, Len: 2},
				{Start: 5, Len: 1},
			},
		},
		{
			name:  "double_newline",
			input: "\n\n",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindEOF, LineIdx: 2},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 0},
				{Start: 1, Len: 0},
				{Start: 2, Len: 0},
			},
		},
		{
			name:  "line_comment_eof",
			input: "//this is a comment",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindComment},
				{Kind: lexer.TokenKindEOF, Column: 19},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 19},
			},
		},
		{
			name:  "line_comment_newline",
			input: "//this is a comment\n",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindComment},
				{Kind: lexer.TokenKindEOF, LineIdx: 1},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 19},
				{Start: 20, Len: 0},
			},
		},
		{
			name:  "multiline_comment",
			input: "/*this is a comment*/",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindComment},
				{Kind: lexer.TokenKindEOF, Column: 21},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 21},
			},
		},
		{
			name:  "multiline_comment_inner_asterisk",
			input: "/*this is * a comment*/",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindComment},
				{Kind: lexer.TokenKindEOF, Column: 23},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 23},
			},
		},
		{
			name:  "multiline_comment_eof",
			input: "/*this is a comment",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindError},
				{Kind: lexer.TokenKindEOF, Column: 19},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 19},
			},
			errs: []error{errors.New("unclosed multiline comment")},
		},
		{
			name:  "identifier",
			input: "hello_world2024 HelloWorld2024",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindIdentifier},
				{Kind: lexer.TokenKindIdentifier, Column: 16},
				{Kind: lexer.TokenKindEOF, Column: 30},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 30},
			},
		},
		{
			name:  "string",
			input: "'test' \"test\"",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindStr},
				{Kind: lexer.TokenKindStr, Column: 7},
				{Kind: lexer.TokenKindEOF, Column: 13},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 13},
			},
		},
		{
			name:  "string_utf8",
			input: "'🙈🙉🙊'",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindStr},
				{Kind: lexer.TokenKindEOF, Column: 14},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 14},
			},
		},
		{
			name:  "escaped_string",
			input: "'this is a \\\"123string\\\"'",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindStr},
				{Kind: lexer.TokenKindEOF, Column: 25},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 25},
			},
		},
		{
			name:  "unterminated_string_eof",
			input: "'test",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindError},
				{Kind: lexer.TokenKindEOF, Column: 5},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 5},
			},
			errs: []error{errors.New("unclosed string")},
		},
		{
			name:  "unterminated_string_newline",
			input: "'test\n'",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindError},
				{Kind: lexer.TokenKindError, LineIdx: 1},
				{Kind: lexer.TokenKindEOF, LineIdx: 1, Column: 1},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 5},
				{Start: 6, Len: 1},
			},
			errs: []error{
				errors.New("unclosed string"),
				errors.New("unclosed string"),
			},
		},
		{
			name:  "unterminated_string_mismatch",
			input: "\"test'",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindError},
				{Kind: lexer.TokenKindEOF, Column: 6},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 6},
			},
			errs: []error{errors.New("unclosed string")},
		},
		{
			name:  "decimal",
			input: "5 0 -5 +5",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindInt},
				{Kind: lexer.TokenKindInt, Column: 2},
				{Kind: lexer.TokenKindInt, Column: 4},
				{Kind: lexer.TokenKindInt, Column: 7},
				{Kind: lexer.TokenKindEOF, Column: 9},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 9},
			},
		},
		{
			name:  "hexadecimal",
			input: "0xff 0XFF",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindInt},
				{Kind: lexer.TokenKindInt, Column: 5},
				{Kind: lexer.TokenKindEOF, Column: 9},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 9},
			},
		},
		{
			name:  "octal",
			input: "056",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindInt},
				{Kind: lexer.TokenKindEOF, Column: 3},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 3},
			},
		},
		{
			name:  "float",
			input: "-8.8 +0.8 -.8 +.8 .8 .8e8 .8e+8 .8e-8 8e8",
			tokenInfos: []lexer.TokenInfo{
				{Kind: lexer.TokenKindBOF},
				{Kind: lexer.TokenKindFloat},
				{Kind: lexer.TokenKindFloat, Column: 5},
				{Kind: lexer.TokenKindFloat, Column: 10},
				{Kind: lexer.TokenKindFloat, Column: 14},
				{Kind: lexer.TokenKindFloat, Column: 18},
				{Kind: lexer.TokenKindFloat, Column: 21},
				{Kind: lexer.TokenKindFloat, Column: 26},
				{Kind: lexer.TokenKindFloat, Column: 32},
				{Kind: lexer.TokenKindFloat, Column: 38},
				{Kind: lexer.TokenKindEOF, Column: 41},
			},
			lineInfos: []lexer.LineInfo{
				{Start: 0, Len: 41},
			},
		},
	}

	runTestCases(t, tests)
}
