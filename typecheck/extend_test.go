package typecheck_test

import "github.com/Clement-Jean/protein/typecheck"

var extendTests = []typecheckTestCase{
	{
		name: "extend not defined",
		contents: []testFile{
			{"a.proto", "extend google.protobuf.FileOptions { bool option = 1; }"},
		},
		errors: []error{
			&typecheck.TypeNotDefinedError{
				File: "a.proto",
				Line: 1,
				Col:  8,
				Name: ".google.protobuf.FileOptions",
			},
		},
	},
	{
		name: "extend full ident not defined",
		contents: []testFile{
			{"a.proto", "extend .google.protobuf.FileOptions { bool option = 1; }"},
		},
		errors: []error{
			&typecheck.TypeNotDefinedError{
				File: "a.proto",
				Line: 1,
				Col:  9,
				Name: ".google.protobuf.FileOptions",
			},
		},
	},
	{
		name: "extend",
		contents: []testFile{
			{"a.proto", "import 'b.proto'; extend google.protobuf.FileOptions { bool option = 1; }"},
			{"b.proto", "package google.protobuf; message FileOptions {}"},
		},
	},
	{
		name: "extend full ident",
		contents: []testFile{
			{"a.proto", "import 'b.proto'; extend .google.protobuf.FileOptions { bool option = 1; }"},
			{"b.proto", "package google.protobuf; message FileOptions {}"},
		},
	},
	{
		name: "extend message field",
		contents: []testFile{
			{"a.proto", "import 'b.proto'; message A {} extend google.protobuf.FileOptions { A option = 1; }"},
			{"b.proto", "package google.protobuf; message FileOptions {}"},
		},
	},
	{
		name: "extend already defined",
		contents: []testFile{
			{"a.proto", "message A {} message B {} message C {} extend B { A a = 1; } extend C { A a = 1; }"},
		},
		errors: []error{
			&typecheck.TypeRedefinedError{
				Files: []string{"a.proto", "a.proto"},
				Lines: []int{1, 1},
				Cols:  []int{73, 51},
				Name:  ".a",
			},
		},
	},
	{
		name: "extend map not allowed",
		contents: []testFile{
			{"a.proto", "message B {} extend B { map<string, string> a = 1; }"},
		},
		errors: []error{
			&typecheck.ExtendMapNotAllowedError{
				File: "a.proto",
				Line: 1,
				Col:  45,
			},
		},
	},
	{
		name: "extend message name clash",
		contents: []testFile{
			{"a.proto", "message B {} message a {} extend B { bool a = 1; }"},
		},
		errors: []error{
			&typecheck.TypeRedefinedError{
				Files: []string{"a.proto", "a.proto"},
				Lines: []int{1, 1},
				Cols:  []int{38, 22},
				Name:  ".a",
			},
		},
	},
}
