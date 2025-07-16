package typecheck_test

import "github.com/Clement-Jean/protein/typecheck"

var optionTests = []typecheckTestCase{
	{
		name: "option file",
		contents: []testFile{
			{"a.proto", "import 'google/protobuf/descriptor.proto'; option deprecated = true;"},
			{"google/protobuf/descriptor.proto", "package google.protobuf; message FileOptions { bool deprecated = 1; }"},
		},
	},
	{
		name: "option file not defined",
		contents: []testFile{
			{"a.proto", "import 'google/protobuf/descriptor.proto'; option nope = true;"},
			{"google/protobuf/descriptor.proto", "package google.protobuf; message FileOptions { bool deprecated = 1; }"},
		},
		errors: []error{
			&typecheck.OptionUnknownError{
				File: "a.proto",
				Line: 1,
				Col:  51,
				Name: ".nope",
			},
		},
	},
	{
		name: "option file custom not defined",
		contents: []testFile{
			{"a.proto", "import 'b.proto'; option (my.custom.opt) = true;"},
			{"b.proto", "package my.custom; import 'google/protobuf/descriptor.proto'; extend google.protobuf.FileOptions { bool option = 1; }"},
			{"google/protobuf/descriptor.proto", "package google.protobuf; message FileOptions {}"},
		},
		errors: []error{
			&typecheck.OptionUnknownError{
				File: "a.proto",
				Line: 1,
				Col:  26,
				Name: ".my.custom.opt",
			},
		},
	},
	{
		name: "option file custom",
		contents: []testFile{
			{"a.proto", "import 'b.proto'; option (my.custom.option) = true;"},
			{"b.proto", "package my.custom; import 'google/protobuf/descriptor.proto'; extend google.protobuf.FileOptions { bool option = 1; }"},
			{"google/protobuf/descriptor.proto", "package google.protobuf; message FileOptions {}"},
		},
		errors: []error{},
	},
	{
		name: "option file custom full ident",
		contents: []testFile{
			{"a.proto", "import 'b.proto'; option (.my.custom.opt) = true;"},
			{"b.proto", "package my.custom; import 'google/protobuf/descriptor.proto'; extend google.protobuf.FileOptions { bool option = 1; }"},
			{"google/protobuf/descriptor.proto", "package google.protobuf; message FileOptions {}"},
		},
		errors: []error{
			&typecheck.OptionUnknownError{
				File: "a.proto",
				Line: 1,
				Col:  26,
				Name: ".my.custom.opt",
			},
		},
	},
	{
		name: "option file custom full ident with field undefined",
		contents: []testFile{
			{"a.proto", "import 'b.proto'; option (.my.custom.opt).a = true;"},
			{"b.proto", "package my.custom; import 'google/protobuf/descriptor.proto'; message A { } extend google.protobuf.FileOptions { A opt = 1; }"},
			{"google/protobuf/descriptor.proto", "package google.protobuf; message FileOptions {}"},
		},
		errors: []error{
			&typecheck.OptionUnknownError{
				File: "a.proto",
				Line: 1,
				Col:  26,
				Name: "(.my.custom.opt).a",
			},
		},
	},
	{
		name: "option file custom full ident with field",
		contents: []testFile{
			{"a.proto", "import 'b.proto'; option (.my.custom.opt).a = true;"},
			{"b.proto", "package my.custom; import 'google/protobuf/descriptor.proto'; message A { bool a = 1; } extend google.protobuf.FileOptions { A opt = 1; }"},
			{"google/protobuf/descriptor.proto", "package google.protobuf; message FileOptions {}"},
		},
	},
	{
		name: "option file custom full ident with fields",
		contents: []testFile{
			{"a.proto", "import 'b.proto'; option (.my.custom.opt).a.b = true;"},
			{"b.proto", "package my.custom; import 'google/protobuf/descriptor.proto'; message B { bool b = 1; } message A { B a = 1; } extend google.protobuf.FileOptions { A opt = 1; }"},
			{"google/protobuf/descriptor.proto", "package google.protobuf; message FileOptions {}"},
		},
	},
	{
		name: "option file custom without parens",
		contents: []testFile{
			{"a.proto", "import 'b.proto'; option my.a.b = true;"},
			{"b.proto", "package my.custom; import 'google/protobuf/descriptor.proto'; message B { bool b = 1; } message A { B a = 1; } extend google.protobuf.FileOptions { A my = 1; }"},
			{"google/protobuf/descriptor.proto", "package google.protobuf; message FileOptions {}"},
		},
		errors: []error{
			&typecheck.OptionUnknownError{
				File: "a.proto",
				Line: 1,
				Col:  26,
				Name: ".my.a.b",
			},
		},
	},
}
