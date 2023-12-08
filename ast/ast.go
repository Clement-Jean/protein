package ast

type Ast struct {
	Imports []Import
	Syntax  Syntax
	Edition Edition
	Package Package
}
