package parser_test

import "testing"

func TestParseOption(t *testing.T) {
	tests := parseTestContent(t, "option.txt")
	runParseTestCase(t, tests)
}
