package ast

type Ast struct {
	Imports []Import
	Package Package
	Syntax  Syntax
	Edition Edition
}
