package typecheck_test

import "github.com/Clement-Jean/protein/typecheck"

var serviceTests = []typecheckTestCase{
	{
		name: "service redefined",
		contents: []testFile{
			{"a.proto", "service A {} service A {}"},
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
		name: "rpc redefined",
		contents: []testFile{
			{"a.proto", "message C {} service A { rpc B (C) returns (C); rpc B (C) returns (C); }"},
		},
		errors: []error{
			&typecheck.TypeRedefinedError{
				Files: []string{"a.proto", "a.proto"},
				Lines: []int{1, 1},
				Cols:  []int{53, 30},
				Name:  ".A.B",
			},
		},
	},
}
