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
}
