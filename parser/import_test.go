package parser_test

import "testing"

func TestParseImport(t *testing.T) {
	tests := parseTestContent(t, "import.txt")
	runParseTestCase(t, tests)
}
