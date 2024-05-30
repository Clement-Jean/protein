package parser_test

import "testing"

func TestParseSyntax(t *testing.T) {
	tests := parseTestContent(t, "syntax.txt")
	runParseTestCase(t, tests)
}
