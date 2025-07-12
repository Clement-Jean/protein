package typecheck_test

import "github.com/Clement-Jean/protein/typecheck"

var enumTests = []typecheckTestCase{
	{
		name: "enum redefined",
		contents: []testFile{
			{"a.proto", "enum A { A_UNSPECIFIED = 0; } enum A { A_UNSPECIFIED = 0; }"},
		},
		errors: []error{
			&typecheck.TypeRedefinedError{
				Files: []string{"a.proto", "a.proto"},
				Lines: []int{1, 1},
				Cols:  []int{36, 6},
				Name:  ".A",
			},
		},
	},
	{
		name: "enum value redefined",
		contents: []testFile{
			{"a.proto", "enum A { A_UNSPECIFIED = 0; A_UNSPECIFIED = 1; }"},
		},
		errors: []error{
			&typecheck.FieldNameReusedError{
				File:       "a.proto",
				Line:       1,
				Col:        29,
				ParentName: ".A",
				Name:       "A_UNSPECIFIED",
			},
		},
	},
	{
		name: "enum value redefined with package",
		contents: []testFile{
			{"a.proto", "package test; enum A { A_UNSPECIFIED = 0; A_UNSPECIFIED = 1; }"},
		},
		errors: []error{
			&typecheck.FieldNameReusedError{
				File:       "a.proto",
				Line:       1,
				Col:        43,
				ParentName: ".test.A",
				Name:       "A_UNSPECIFIED",
			},
		},
	},
	{
		name: "enum value tag reused",
		contents: []testFile{
			{"a.proto", "enum A { A_UNSPECIFIED = 0; A_FIRST = 0; }"},
		},
		errors: []error{
			&typecheck.FieldTagReusedError{
				File:       "a.proto",
				Line:       1,
				Col:        29,
				ParentName: ".A",
				Tag:        0,
			},
		},
	},
	{
		name: "enum value tag reused with package",
		contents: []testFile{
			{"a.proto", "package test; enum A { A_UNSPECIFIED = 0; A_FIRST = 0; }"},
		},
		errors: []error{
			&typecheck.FieldTagReusedError{
				File:       "a.proto",
				Line:       1,
				Col:        43,
				ParentName: ".test.A",
				Tag:        0,
			},
		},
	},
	{
		name: "enum value tag max",
		contents: []testFile{
			{"a.proto", "enum A { A_UNSPECIFIED = 0; A_ONE = 2147483648; }"},
		},
		errors: []error{
			&typecheck.MaxEnumValueTagError{
				File: "a.proto",
				Line: 1,
				Col:  29,
			},
		},
	},
}
