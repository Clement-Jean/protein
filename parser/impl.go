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

	syncPos int // last synchronization position
	syncCnt int // number of parser.advance calls without progress
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
		eofID := p.tokens[len(p.tokens)-1].ID
		return &token.Token{ID: eofID, Kind: token.KindEOF}
	}

	return &p.tokens[i]
}

func (p *impl) nextToken() *token.Token {
	for ; p.idx < len(p.tokens) && isSpaceOrComment(p.tokens[p.idx].Kind); p.idx++ {
		p.prevIdx = p.idx
	}

	if p.idx >= len(p.tokens) {
		eofID := p.tokens[len(p.tokens)-1].ID
		return &token.Token{ID: eofID, Kind: token.KindEOF}
	}

	tok := p.tokens[p.idx]
	p.prevIdx = p.idx
	p.idx++
	return &tok
}

var literalToKind = map[string]token.Kind{
	"syntax":     token.KindSyntax,
	"edition":    token.KindEdition,
	"package":    token.KindPackage,
	"import":     token.KindImport,
	"option":     token.KindOption,
	"reserved":   token.KindReserved,
	"enum":       token.KindEnum,
	"message":    token.KindMessage,
	"map":        token.KindMap,
	"oneof":      token.KindOneOf,
	"extensions": token.KindExtensions,
	"service":    token.KindService,
	"rpc":        token.KindRpc,
	"extend":     token.KindExtend,
}

func (p *impl) Parse() (a ast.Ast, errs []error) {
	for tok := p.nextToken(); tok.Kind != token.KindEOF; tok = p.nextToken() {
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
			opt, innerErrs := p.parseOption()
			if len(innerErrs) == 0 {
				a.Options = append(a.Options, opt)
			}
			errs = append(errs, innerErrs...)
		case token.KindEnum:
			enum, innerErrs := p.parseEnum()
			if len(innerErrs) == 0 {
				a.Enums = append(a.Enums, enum)
			}
			errs = append(errs, innerErrs...)
		case token.KindMessage:
			msg, innerErrs := p.parseMessage(1)
			if len(innerErrs) == 0 {
				a.Messages = append(a.Messages, msg)
			}
			errs = append(errs, innerErrs...)
		case token.KindService:
			svc, innerErrs := p.parseService()
			if len(innerErrs) == 0 {
				a.Services = append(a.Services, svc)
			}
			errs = append(errs, innerErrs...)
		case token.KindExtend:
			extend, innerErrs := p.parseExtend()
			if len(innerErrs) == 0 {
				a.Extensions = append(a.Extensions, extend)
			}
			errs = append(errs, innerErrs...)
		default:
			err = gotUnexpected(
				tok,
				token.KindSyntax, token.KindEdition,
				token.KindPackage, token.KindImport, token.KindOption,
				token.KindMessage, token.KindEnum, token.KindService, token.KindExtend,
			)
		}

		if err != nil {
			errs = append(errs, err)
			p.advanceTo(protoTopLevelStart)
		}

		err = nil
	}

	return a, errs
}
