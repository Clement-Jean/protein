package parser_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/parser"

	"github.com/google/go-cmp/cmp"
)

var headersRegexp = *regexp.MustCompile(
	`(?m)` +
		`^(?:=+){3,}([^=\r\n][^\r\n]*)?` +
		`\r?\n` +
		`(?P<testName>(?:[^=][^\r\n]*))` +
		`\r?\n+` +
		`^(?:=+){3,}([^=\r\n][^\r\n]*)?` +
		`\r?\n$`,
)

var hyphensRegexp = *regexp.MustCompile(
	`(?m)^(?:-+){3,}([^-\r\n][^\r\n]*)?\r?\n$`,
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

type ParseTestCase struct {
	name         string
	input        string
	expectedTree string
	expectedErrs []error
}

func parseTestContent(t *testing.T, filename string) (tests []ParseTestCase) {
	testPrefix := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	filePath := filepath.Join(basepath, "testdata", filename)
	b, err := os.ReadFile(filePath)

	if err != nil {
		t.Fatal(err)
	}

	if bytes.IndexByte(b, '\t') != -1 {
		t.Fatalf("found tabs in %s. use spaces instead", filename)
	}

	content := string(b)
	headers := headersRegexp.FindAllStringSubmatchIndex(content, -1)
	hyphens := hyphensRegexp.FindAllStringSubmatchIndex(content, -1)

	if len(hyphens) != len(headers) {
		t.Fatalf("expected %d hyphens, got %d", len(headers), len(hyphens))
	}

	for i, header := range headers {
		testNameIdx := headersRegexp.SubexpIndex("testName") * 2
		startTestNameIdx := header[testNameIdx]
		endTestNameIdx := header[testNameIdx+1]
		endHeaderIdx := header[1]
		startHyphenIdx := hyphens[i][0]
		endHyphenIdx := hyphens[i][1]
		testName := content[startTestNameIdx:endTestNameIdx]
		input := strings.TrimSpace(content[endHeaderIdx:startHyphenIdx])

		var output string
		if i+1 < len(headers) {
			startNextHeaderIdx := headers[i+1][0]
			output = strings.TrimSpace(content[endHyphenIdx:startNextHeaderIdx])
		} else {
			output = strings.TrimSpace(content[endHyphenIdx:])
		}

		name := strings.ToLower(fmt.Sprintf("%s_%s", testPrefix, testName))
		tests = append(tests, ParseTestCase{
			name:         name,
			input:        input,
			expectedTree: output,
		})
	}
	return tests
}

func runParseTestCase(t *testing.T, tests []ParseTestCase) {
	t.Helper()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l, err := lexer.NewFromReader(strings.NewReader(test.input))
			if err != nil {
				t.Fatal(err)
			}

			tb, errs := l.Lex()
			if len(errs) != 0 {
				t.Fatal(errs)
			}

			p := parser.New(tb)
			pt, errs := p.Parse()

			buf := new(bytes.Buffer)
			pt.Print(buf, tb)
			if len(errs) != 0 {
				fmt.Fprintf(buf, "errs = %v", errs)
			}

			trimmed := strings.TrimSpace(buf.String())
			if diff := cmp.Diff(test.expectedTree, trimmed); diff != "" {
				t.Errorf("%s mismatch (-want +got):\n%s", t.Name(), diff)
			}
		})
	}
}

var testFiles = []string{
	"syntax.txt",
	"edition.txt",
	"import.txt",
	"package.txt",
	"option.txt",
	"text_message.txt",
	"a_bit_of_everything.txt",
}

func TestParser(t *testing.T) {
	for _, file := range testFiles {
		subTests := parseTestContent(t, file)
		runParseTestCase(t, subTests)
	}
}
