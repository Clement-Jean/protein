package lexer_test

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"
	"runtime"
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

func findFiles(root, ext string) (files []string) {
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			files = append(files, s)
		}
		return nil
	})
	return files
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

func rangeTestCase(name string, start, end lexer.TokenKind) TestCase {
	var b strings.Builder

	realStart := uint8(start)
	realEnd := uint8(end) + 1
	diff := realEnd - realStart

	expectedTokenInfos := make([]lexer.TokenInfo, diff+2)
	expectedTokenInfos[0].Kind = lexer.TokenKindBOF
	offset := 0
	for i := uint8(0); i < diff; i++ {
		info := &expectedTokenInfos[i+1]
		info.Kind = lexer.TokenKind(realStart + i)
		info.Offset = uint32(offset)

		b.WriteString(info.Kind.String())
		b.WriteString(" ")

		offset += len(info.Kind.String()) + 1
	}
	info := &expectedTokenInfos[diff+1]
	info.Kind = lexer.TokenKindEOF
	info.Offset = uint32(offset)

	return TestCase{
		name:       name,
		input:      b.String(),
		tokenInfos: expectedTokenInfos,
		lineInfos:  []lexer.LineInfo{{}},
	}
}

var tests = []TestCase{
	rangeTestCase("symbols", lexer.TokenKindUnderscore, lexer.TokenKindSlash),
	rangeTestCase("keywords", lexer.TokenKindTypeBool, lexer.TokenKindWeak),
	{
		name:  "invalid",
		input: "&",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindError},
			{Kind: lexer.TokenKindEOF, Offset: 1},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
		},
		errs: []error{errors.New("invalid char '&'")},
	},
	{
		name:  "invalid_utf8",
		input: "🙈",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindError},
			{Kind: lexer.TokenKindError, Offset: 1},
			{Kind: lexer.TokenKindError, Offset: 2},
			{Kind: lexer.TokenKindError, Offset: 3},
			{Kind: lexer.TokenKindEOF, Offset: 4},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
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
			{Kind: lexer.TokenKindBOF, Offset: 3},
			{Kind: lexer.TokenKindEOF, Offset: 3},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 3},
		},
	},
	{
		name:  "spaces",
		input: "\t\n\v\f\n ",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindEOF, Offset: 6},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
			{Start: 2},
			{Start: 5},
		},
	},
	{
		name:  "double_newline",
		input: "\n\n",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindEOF, Offset: 2},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
			{Start: 1},
			{Start: 2},
		},
	},
	{
		name:  "line_comment_eof",
		input: "//this is a comment",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindComment},
			{Kind: lexer.TokenKindEOF, Offset: 19},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
		},
	},
	{
		name:  "line_comment_newline",
		input: "//this is a comment\n",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindComment},
			{Kind: lexer.TokenKindEOF, Offset: 20},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
			{Start: 20},
		},
	},
	{
		name:  "multiline_comment",
		input: "/*this is a comment*/",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindComment},
			{Kind: lexer.TokenKindEOF, Offset: 21},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
		},
	},
	{
		name:  "multiline_comment_inner_asterisk",
		input: "/*this is * a comment*/",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindComment},
			{Kind: lexer.TokenKindEOF, Offset: 23},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
		},
	},
	{
		name:  "multiline_comment_eof",
		input: "/*this is a comment",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindError},
			{Kind: lexer.TokenKindEOF, Offset: 19},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
		},
		errs: []error{errors.New("unclosed multiline comment")},
	},
	{
		name:  "identifier",
		input: "hello_world2024 HelloWorld2024",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindIdentifier},
			{Kind: lexer.TokenKindIdentifier, Offset: 16},
			{Kind: lexer.TokenKindEOF, Offset: 30},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
		},
	},
	{
		name:  "string",
		input: "'test' \"test\"",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindStr},
			{Kind: lexer.TokenKindStr, Offset: 7},
			{Kind: lexer.TokenKindEOF, Offset: 13},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
		},
	},
	{
		name:  "string_utf8",
		input: "'🙈🙉🙊'",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindStr},
			{Kind: lexer.TokenKindEOF, Offset: 14},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
		},
	},
	{
		name:  "escaped_string",
		input: "'this is a \\\"123string\\\"'",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindStr},
			{Kind: lexer.TokenKindEOF, Offset: 25},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
		},
	},
	{
		name:  "unterminated_string_eof",
		input: "'test",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindError},
			{Kind: lexer.TokenKindEOF, Offset: 5},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
		},
		errs: []error{errors.New("unclosed string")},
	},
	{
		name:  "unterminated_string_newline",
		input: "'test\n'",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindError},
			{Kind: lexer.TokenKindError, Offset: 6},
			{Kind: lexer.TokenKindEOF, Offset: 7},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
			{Start: 6},
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
			{Kind: lexer.TokenKindEOF, Offset: 6},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
		},
		errs: []error{errors.New("unclosed string")},
	},
	{
		name:  "decimal",
		input: "5 0 -5 +5",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindInt},
			{Kind: lexer.TokenKindInt, Offset: 2},
			{Kind: lexer.TokenKindInt, Offset: 4},
			{Kind: lexer.TokenKindInt, Offset: 7},
			{Kind: lexer.TokenKindEOF, Offset: 9},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
		},
	},
	{
		name:  "hexadecimal",
		input: "0xff 0XFF",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindInt},
			{Kind: lexer.TokenKindInt, Offset: 5},
			{Kind: lexer.TokenKindEOF, Offset: 9},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
		},
	},
	{
		name:  "octal",
		input: "056",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindInt},
			{Kind: lexer.TokenKindEOF, Offset: 3},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
		},
	},
	{
		name:  "float",
		input: "-8.8 +0.8 -.8 +.8 .8 .8e8 .8e+8 .8e-8 8e8",
		tokenInfos: []lexer.TokenInfo{
			{Kind: lexer.TokenKindBOF},
			{Kind: lexer.TokenKindFloat},
			{Kind: lexer.TokenKindFloat, Offset: 5},
			{Kind: lexer.TokenKindFloat, Offset: 10},
			{Kind: lexer.TokenKindFloat, Offset: 14},
			{Kind: lexer.TokenKindFloat, Offset: 18},
			{Kind: lexer.TokenKindFloat, Offset: 21},
			{Kind: lexer.TokenKindFloat, Offset: 26},
			{Kind: lexer.TokenKindFloat, Offset: 32},
			{Kind: lexer.TokenKindFloat, Offset: 38},
			{Kind: lexer.TokenKindEOF, Offset: 41},
		},
		lineInfos: []lexer.LineInfo{
			{Start: 0},
		},
	},
}

func TestLexer(t *testing.T) {
	runTestCases(t, tests)
}

var toks *lexer.TokenizedBuffer

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func BenchmarkLexer(b *testing.B) {
	corpusPath := filepath.Join(basepath, "../corpus/")

	for _, s := range findFiles(corpusPath, ".proto") {
		l, err := lexer.NewFromFile(s)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(s, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				toks, _ = l.Lex()
			}
		})
	}
}

func FuzzLexer(f *testing.F) {
	for _, tc := range tests {
		f.Add(tc.input)
	}

	f.Fuzz(func(t *testing.T, s string) {
		l, err := lexer.NewFromReader(strings.NewReader(s))
		if err != nil {
			f.Fatal(err)
		}

		toks, _ = l.Lex()
	})
}
