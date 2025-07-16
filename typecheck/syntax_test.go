package typecheck_test

import "github.com/Clement-Jean/protein/typecheck"

var syntaxTests = []typecheckTestCase{
	{
		name: "syntax proto2",
		contents: []testFile{
			{"a.proto", "syntax = 'proto2';"},
		},
	},
	{
		name: "syntax proto3",
		contents: []testFile{
			{"a.proto", "syntax = 'proto3';"},
		},
	},
	{
		name: "syntax unknown",
		contents: []testFile{
			{"a.proto", "syntax = 'unknown';"},
		},
		errors: []error{
			&typecheck.UnknownSyntaxError{
				File:  "a.proto",
				Line:  1,
				Col:   10,
				Value: "unknown",
			},
		},
	},
}
