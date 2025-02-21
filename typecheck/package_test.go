package typecheck_test

import (
	"testing"

	"github.com/Clement-Jean/protein/typecheck"
)

type packageTestCase struct {
	name     string
	contents []testFile
	errors   []error
}

func TestPackage(t *testing.T) {
	tests := []importTestCase{
		{
			name: "trivial package",
			contents: []testFile{
				{"a.proto", "package a.b;"},
			},
		},
		{
			name: "redefined package",
			contents: []testFile{
				{"a.proto", "package a.b; package b.a;"},
			},
			errors: []error{
				&typecheck.PackageMultipleDefError{File: "a.proto"},
			},
		},
	}

	for _, test := range tests {
		units := createUnits(t, test.contents)

		t.Run(test.name, func(t *testing.T) {
			l := typecheck.New(
				units,
				typecheck.WithSourceCreator(fakeSourceCreator(test.contents, test.unknown)),
				typecheck.WithFileCheck(fakeFileCheck(test.contents, test.unknown)),
			)
			errs := l.Check()

			if len(errs) != len(test.errors) {
				t.Fatalf("expected %d errors, got %d: %v", len(test.errors), len(errs), errs)
			}

			for i, err := range errs {
				if err.Error() != test.errors[i].Error() {
					t.Fatalf("expected %q, got %q", test.errors[i].Error(), err.Error())
				}
			}
		})
	}

}
