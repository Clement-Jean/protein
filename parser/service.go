package parser

import (
	"bytes"

	"github.com/Clement-Jean/protein/ast"
	internal_bytes "github.com/Clement-Jean/protein/internal/bytes"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parseRpcType() (isStream bool, id ast.Identifier, err error) {
	if peek := p.peek(); peek.Kind != token.KindLeftParen {
		return false, ast.Identifier{}, gotUnexpected(peek, token.KindLeftParen)
	}
	p.nextToken()

	id, err = p.parseFullyQualifiedIdentifier()

	if err != nil {
		return false, ast.Identifier{}, err
	}

	literal := p.fm.Lookup(id.ID)
	stream := internal_bytes.FromString("stream")

	if bytes.Compare(literal, stream) == 0 {
		isStream = true
		id, err = p.parseFullyQualifiedIdentifier()

		if err != nil {
			return false, ast.Identifier{}, err
		}
	}

	if peek := p.peek(); peek.Kind != token.KindRightParen {
		return false, ast.Identifier{}, gotUnexpected(peek, token.KindRightParen)
	}
	p.nextToken()

	return isStream, id, err
}

func (p *impl) parseRpc() (rpc ast.Rpc, err error) {
	first := p.curr()
	id, err := p.parseIdentifier()

	if err != nil {
		return ast.Rpc{}, err
	}

	isClientStream, inputType, err := p.parseRpcType()

	if err != nil {
		return ast.Rpc{}, err
	}

	if peek := p.peek(); peek.Kind != token.KindIdentifier {
		return ast.Rpc{}, gotUnexpected(peek, token.KindReturns)
	}
	tok := p.nextToken()
	literal := p.fm.Lookup(tok.ID)
	returns := internal_bytes.FromString("returns")

	if bytes.Compare(literal, returns) != 0 {
		return ast.Rpc{}, gotUnexpected(tok, token.KindReturns)
	}

	isServerStream, outputType, err := p.parseRpcType()

	if err != nil {
		return ast.Rpc{}, err
	}

	peek := p.peek()

	if peek.Kind == token.KindLeftBrace {
		p.nextToken()
		peek := p.peek()
		for ; peek.Kind != token.KindRightBrace && peek.Kind != token.KindEOF; peek = p.peek() {
			if peek.Kind == token.KindSemicolon {
				p.nextToken()
				continue
			}

			kind := peek.Kind

			if literal := p.fm.Lookup(peek.ID); literal != nil {
				if k, ok := literalToKind[internal_bytes.ToString(literal)]; ok {
					kind = k
				}
			}

			switch kind {
			case token.KindOption:
				p.nextToken() // point to option keyword
				if option, err := p.parseOption(); err == nil {
					rpc.Options = append(rpc.Options, option)
				}
			default:
				return ast.Rpc{}, gotUnexpected(peek, token.KindOption, token.KindRightBrace)
			}
		}

		if peek.Kind != token.KindRightBrace {
			return ast.Rpc{}, gotUnexpected(peek, token.KindRightBrace)
		}
	} else if peek.Kind != token.KindSemicolon {
		return ast.Rpc{}, gotUnexpected(peek, token.KindSemicolon)
	}

	last := p.nextToken()

	rpc.ID = p.fm.Merge(token.KindRpc, first.ID, last.ID)
	rpc.Name = id
	rpc.InputType = inputType
	rpc.OutputType = outputType
	rpc.IsClientStream = isClientStream
	rpc.IsServerStream = isServerStream
	return rpc, err
}

func (p *impl) parseService() (service ast.Service, err error) {
	first := p.curr()
	id, err := p.parseIdentifier()

	if err != nil {
		return ast.Service{}, err
	}

	if peek := p.peek(); peek.Kind != token.KindLeftBrace {
		return ast.Service{}, gotUnexpected(peek, token.KindLeftBrace)
	}
	p.nextToken()

	peek := p.peek()
	for ; peek.Kind != token.KindRightBrace && peek.Kind != token.KindEOF; peek = p.peek() {
		if peek.Kind == token.KindSemicolon {
			p.nextToken()
			continue
		}

		kind := peek.Kind

		if literal := p.fm.Lookup(peek.ID); literal != nil {
			if k, ok := literalToKind[internal_bytes.ToString(literal)]; ok {
				kind = k
			}
		}

		switch kind {
		case token.KindOption:
			var option ast.Option

			p.nextToken() // point to option keyword
			if option, err = p.parseOption(); err == nil {
				service.Options = append(service.Options, option)
			}
		case token.KindRpc:
			var rpc ast.Rpc

			p.nextToken() // point to rpc keyword
			if rpc, err = p.parseRpc(); err == nil {
				service.Rpcs = append(service.Rpcs, rpc)
			}
		default:
			err = gotUnexpected(peek, token.KindOption, token.KindRpc, token.KindRightBrace)
		}

		if err != nil {
			// TODO report error
			// TODO p.advanceTo(exprEnd)
			return ast.Service{}, err
		}
	}

	if peek.Kind != token.KindRightBrace {
		return ast.Service{}, gotUnexpected(peek, token.KindRightBrace)
	}

	last := p.nextToken()
	service.Name = id
	service.ID = p.fm.Merge(token.KindService, first.ID, last.ID)
	return service, nil
}
