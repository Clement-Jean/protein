package typecheck_test

import "github.com/Clement-Jean/protein/typecheck"

var messageTests = []typecheckTestCase{
	{
		name: "message nested type",
		contents: []testFile{
			{"a.proto", "message A { message B {} B b = 1; }"},
		},
	},
	{
		name: "message not defined",
		contents: []testFile{
			{"a.proto", "message A { C c = 1; }"},
		},
		errors: []error{
			&typecheck.TypeNotDefinedError{
				File: "a.proto",
				Line: 1,
				Col:  13,
				Name: "C",
			},
		},
	},
	{
		name: "message redefined",
		contents: []testFile{
			{"a.proto", "message A {} message A {}"},
		},
		errors: []error{
			&typecheck.TypeRedefinedError{
				Files: []string{"a.proto", "a.proto"},
				Lines: []int{1, 1},
				Cols:  []int{22, 9},
				Name:  ".A",
			},
		},
	},
	{
		name: "message redefined across files",
		contents: []testFile{
			{"a.proto", "message A {}"},
			{"b.proto", "message A {}"},
		},
		errors: []error{
			&typecheck.TypeRedefinedError{
				Files: []string{"b.proto", "a.proto"},
				Lines: []int{1, 1},
				Cols:  []int{9, 9},
				Name:  ".A",
			},
		},
	},
	{
		name: "message field redefined",
		contents: []testFile{
			{"a.proto", "message A { message B {} int32 b = 1; B b = 2; }"},
		},
		errors: []error{
			&typecheck.FieldNameReusedError{
				File:       "a.proto",
				Line:       1,
				Col:        39,
				ParentName: ".A",
				Name:       "b",
			},
		},
	},
	{
		name: "message field redefined with package",
		contents: []testFile{
			{"a.proto", "package test; message A { message B {} int32 b = 1; B b = 2; }"},
		},
		errors: []error{
			&typecheck.FieldNameReusedError{
				File:       "a.proto",
				Line:       1,
				Col:        53,
				ParentName: ".test.A",
				Name:       "b",
			},
		},
	},
	{
		name: "message field tag reused",
		contents: []testFile{
			{"a.proto", "message A { int32 a = 1; int32 b = 1; }"},
		},
		errors: []error{
			&typecheck.FieldTagReusedError{
				File:       "a.proto",
				Line:       1,
				Col:        26,
				ParentName: ".A",
				Tag:        1,
			},
		},
	},
	{
		name: "message field tag reused with package",
		contents: []testFile{
			{"a.proto", "package test; message A { int32 a = 1; int32 b = 1; }"},
		},
		errors: []error{
			&typecheck.FieldTagReusedError{
				File:       "a.proto",
				Line:       1,
				Col:        40,
				ParentName: ".test.A",
				Tag:        1,
			},
		},
	},
	{
		name: "field tag max",
		contents: []testFile{
			{"a.proto", "message A { string a = 536870912; }"},
		},
		errors: []error{
			&typecheck.MaxFieldTagError{
				File: "a.proto",
				Line: 1,
				Col:  13,
			},
		},
	},
}

var mapTests = []typecheckTestCase{
	{
		name: "map",
		contents: []testFile{
			{"a.proto", "message B {} message A { map<string, B> names = 1; }"},
		},
	},
	{
		name: "map value",
		contents: []testFile{
			{"a.proto", "message A { map<int32, B> b = 1; } message B {}"},
		},
	},
	{
		name: "map value across files",
		contents: []testFile{
			{"a.proto", "import 'b.proto'; message A { map<int32, B> b = 1; }"},
			{"b.proto", "message B {}"},
		},
	},
	{
		name: "map value type not defined",
		contents: []testFile{
			{"a.proto", "message A { map<string, B> names = 1; }"},
		},
		errors: []error{
			&typecheck.TypeNotDefinedError{
				File: "a.proto",
				Line: 1,
				Col:  25,
				Name: "B",
			},
		},
	},
	{
		name: "map with preceding dot",
		contents: []testFile{
			{"a.proto", "package a; import 'b.proto'; message A { map<string, .b.B> names = 1; }"},
			{"b.proto", "package b; message B {}"},
		},
	},
	{
		name: "unknown import map value across files",
		contents: []testFile{
			{"a.proto", "import 'b.proto'; message A { map<int32, B> b = 1; }"},
		},
		unknown: []testFile{
			{"b.proto", "message B {}"},
		},
	},
	{
		name: "package map value across files",
		contents: []testFile{
			{"a.proto", "import 'b.proto'; message A { map<int32, google.protobuf.B> b = 1; }"},
			{"b.proto", "package google.protobuf; message B {}"},
		},
	},
	{
		name: "unknown import package map value across files",
		contents: []testFile{
			{"a.proto", "import 'b.proto'; message A { map<int32, google.protobuf.B> b = 1; }"},
		},
		unknown: []testFile{
			{"b.proto", "package google.protobuf; message B {}"},
		},
	},
	{
		name: "map name redefined",
		contents: []testFile{
			{"a.proto", "message B {} message A { map<string, string> a = 1; map<int32, string> a = 2; }"},
		},
		errors: []error{
			&typecheck.FieldNameReusedError{
				File:       "a.proto",
				Line:       1,
				Col:        72,
				Name:       "a",
				ParentName: ".A",
			},
		},
	},
	{
		name: "map tag redefined",
		contents: []testFile{
			{"a.proto", "message B {} message A { map<string, string> a = 1; map<int32, string> b = 1; }"},
		},
		errors: []error{
			&typecheck.FieldTagReusedError{
				File:       "a.proto",
				Line:       1,
				Col:        72,
				Tag:        1,
				ParentName: ".A",
			},
		},
	},
	{
		name: "map field tag max",
		contents: []testFile{
			{"a.proto", "message A { map<string, string> a = 536870912; }"},
		},
		errors: []error{
			&typecheck.MaxFieldTagError{
				File: "a.proto",
				Line: 1,
				Col:  33,
			},
		},
	},
}

var oneofTests = []typecheckTestCase{
	{
		name: "oneof",
		contents: []testFile{
			{"a.proto", "message A { oneof B {} B b = 1; }"},
		},
		errors: []error{
			&typecheck.NotTypeError{
				File: "a.proto",
				Line: 1,
				Col:  24,
				Name: "B",
			},
		},
	},
	{
		name: "oneof redefined",
		contents: []testFile{
			{"a.proto", "message A { oneof B {} oneof B {} }"},
		},
		errors: []error{
			&typecheck.TypeRedefinedError{
				Files: []string{"a.proto", "a.proto"},
				Lines: []int{1, 1},
				Cols:  []int{30, 19},
				Name:  ".A.B",
			},
		},
	},
	{
		name: "oneof field name reused",
		contents: []testFile{
			{"a.proto", "message A { oneof B { int32 a = 1; } int64 a = 2; }"},
		},
		errors: []error{
			&typecheck.FieldNameReusedError{
				File:       "a.proto",
				Line:       1,
				Col:        38,
				Name:       "a",
				ParentName: ".A",
			},
		},
	},
	{
		name: "oneof field tag reused",
		contents: []testFile{
			{"a.proto", "message A { oneof B { int32 a = 1; } int64 b = 1; }"},
		},
		errors: []error{
			&typecheck.FieldTagReusedError{
				File:       "a.proto",
				Line:       1,
				Col:        38,
				Tag:        1,
				ParentName: ".A",
			},
		},
	},
}
