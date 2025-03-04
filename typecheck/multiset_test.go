package typecheck

import (
	"testing"
)

func TestSplitAndMerge(t *testing.T) {
	tests := []struct {
		id, scope, expected string
	}{
		{".A", "", ".A"},
		{"google.protobuf.Empty", "", ".google.protobuf.Empty"},
		{"B", ".A", ".A]B"},
		{"A.B", ".A.B", ".A.B]A.B"},
		{"google.protobuf.Empty", "com.google.A", ".com[google.A]google.protobuf.Empty"},
		{"google.protobuf.Empty", "com.google.protobuf.A", ".com[google.protobuf.A]google.protobuf.Empty"},
		{"google.protobuf.Empty", ".com.google", ".com[google]google.protobuf.Empty"},
		{"google.Empty", ".com", ".com]google.Empty"},
	}

	for _, test := range tests {
		if res := splitAndMerge(test.id, test.scope); test.expected != res {
			t.Errorf("expected %q, got %q", test.expected, res)
		}
	}
}

func FuzzAlgo(f *testing.F) {
	tests := []struct {
		id, scope string
	}{
		{".A", ""},
		{"google.protobuf.Empty", ""},
		{"B", ".A"},
		{"A.B", ".A.B"},
		{"google.protobuf.Empty", "com.google.A"},
		{"google.protobuf.Empty", "com.google.protobuf.A"},
		{"google.protobuf.Empty", ".com.google"},
		{"google.Empty", ".com"},
	}

	for _, test := range tests {
		f.Add(test.id, test.scope)
	}

	f.Fuzz(func(t *testing.T, a, b string) {
		splitAndMerge(a, b)
	})
}
