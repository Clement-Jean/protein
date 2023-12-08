package parser

import "github.com/Clement-Jean/protein/ast"

type Parser interface {
	Parse() (ast.Ast, []error)
}
