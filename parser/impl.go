package parser

import (
	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/codemap"
	"github.com/Clement-Jean/protein/config"
	"github.com/Clement-Jean/protein/internal/bytes"
	"github.com/Clement-Jean/protein/token"
)

type impl struct {
	fm           *codemap.FileMap
	tokens       []token.Token
	prevIdx, idx int
}

func New(tokens []token.Token, fm *codemap.FileMap) Parser {
	return &impl{
		tokens: tokens,
		fm:     fm,
	}
}

func isSpaceOrComment(kind token.Kind) bool {
	return (config.GenerateSourceCodeInfo && kind == token.KindSpace) ||
		(config.KeepComments && kind == token.KindComment)
}

func (p *impl) curr() *token.Token {
	return &p.tokens[p.prevIdx]
}

func (p *impl) peek() *token.Token {
	i := p.idx

	for ; i < len(p.tokens) && isSpaceOrComment(p.tokens[i].Kind); i++ {
	}

	if i >= len(p.tokens) {
		return nil
	}

	return &p.tokens[i]
}

func (p *impl) nextToken() *token.Token {
	for ; p.idx < len(p.tokens) && isSpaceOrComment(p.tokens[p.idx].Kind); p.idx++ {
		p.prevIdx = p.idx
	}

	if p.idx >= len(p.tokens) {
		return nil
	}

	tok := p.tokens[p.idx]
	p.prevIdx = p.idx
	p.idx++
	return &tok
}

var literalToKind = map[string]token.Kind{
	"syntax":   token.KindSyntax,
	"edition":  token.KindEdition,
	"package":  token.KindPackage,
	"import":   token.KindImport,
	"option":   token.KindOption,
	"reserved": token.KindReserved,
	"enum":     token.KindEnum,
	"message":  token.KindMessage,
	"map":      token.KindMap,
}

func (p *impl) Parse() (a ast.Ast, errs []error) {
	for tok := p.nextToken(); tok != nil; tok = p.nextToken() {
		var err error

		if tok.Kind == token.KindSemicolon {
			p.nextToken()
			continue
		}

		kind := token.KindIllegal
		literal := p.fm.Lookup(tok.ID)

		if literal != nil {
			kind = literalToKind[bytes.ToString(literal)]
		}

		switch kind {
		case token.KindSyntax:
			a.Syntax, err = p.parseSyntax()
		case token.KindEdition:
			a.Edition, err = p.parseEdition()
		case token.KindPackage:
			a.Package, err = p.parsePackage()
		case token.KindImport:
			var imp ast.Import

			if imp, err = p.parseImport(); err == nil {
				a.Imports = append(a.Imports, imp)
			}
		case token.KindOption:
			var option ast.Option

			if option, err = p.parseOption(); err == nil {
				a.Options = append(a.Options, option)
			}
		case token.KindEnum:
			var enum ast.Enum

			if enum, err = p.parseEnum(); err == nil {
				a.Enums = append(a.Enums, enum)
			}
		case token.KindMessage:
			var msg ast.Message

			if msg, err = p.parseMessage(1); err == nil {
				a.Messages = append(a.Messages, msg)
			}
		default:
			err = gotUnexpected(
				tok,
				token.KindSyntax, token.KindEdition,
				token.KindPackage, token.KindImport,
			)
		}

		if err != nil {
			errs = append(errs, err)
			// TODO recover error instead of returning
			return a, errs
		}
	}

	return a, nil
}
