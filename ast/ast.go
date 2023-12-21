package ast

import "github.com/Clement-Jean/protein/token"

type Ast struct {
	Imports    []Import
	Options    []Option
	Enums      []Enum
	Messages   []Message
	Services   []Service
	Extensions []Extend
	Package    Package
	Syntax     Syntax
	Edition    Edition
}

type Expression interface {
	expressionNode()
	GetID() token.UniqueID
}
