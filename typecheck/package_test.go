package typecheck_test

import "github.com/Clement-Jean/protein/typecheck"

var packageTests = []typecheckTestCase{
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
	{
		name: "package override",
		contents: []testFile{
			{"a.proto", `package com.google;

import 'google/protobuf/empty.proto';

message A {
  google.protobuf.Empty e = 1;
}`},
			{"google/protobuf/empty.proto", "package google.protobuf; message Empty {}"},
		},
		errors: []error{
			&typecheck.TypeResolvedNotDefinedError{
				File:         "a.proto",
				Line:         6,
				Col:          3,
				Name:         "google.protobuf.Empty",
				ResolvedName: ".com.google.protobuf.Empty",
			},
		},
	},
	{
		name: "package override not exact same path",
		contents: []testFile{
			{"a.proto", `package com.google.notprotobuf;

import 'google/protobuf/empty.proto';

message A {
  google.protobuf.Empty e = 1;
}`},
			{"google/protobuf/empty.proto", "package google.protobuf; message Empty {}"},
		},
		errors: []error{
			&typecheck.TypeResolvedNotDefinedError{
				File:         "a.proto",
				Line:         6,
				Col:          3,
				Name:         "google.protobuf.Empty",
				ResolvedName: ".com.google.protobuf.Empty",
			},
		},
	},
	{
		name: "google protobuf as message names",
		contents: []testFile{
			{"a.proto", `package com.google.notprotobuf;

message google {
  message protobuf {
    message Empty {
    }
  }
}

message A {
  google.protobuf.Empty e = 1;
}`},
		},
	},
	{
		name: "google protobuf as inner message names",
		contents: []testFile{
			{"a.proto", `package com.google.notprotobuf;

message A {
  message google {
    message protobuf {
      message Empty {
      }
    }
  }

  google.protobuf.Empty e = 1;
}`},
		},
	},
	{
		name: "google protobuf as inner message names without package",
		contents: []testFile{
			{"a.proto", `
message A {
  message google {
    message protobuf {
      message Empty {
      }
    }
  }

  google.protobuf.Empty e = 1;
}`},
		},
	},
}
