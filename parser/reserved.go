package parser

import (
	"bytes"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parseRange() (ast.Range, error) {
	first, err := p.parseInt()

	if err != nil {
		return ast.Range{}, err
	}

	last := first
	if peek := p.peek(); peek.Kind == token.KindIdentifier && bytes.Compare(p.fm.Lookup(peek.ID), []byte("to")) == 0 {
		p.nextToken()

		switch peek = p.peek(); peek.Kind {
		case token.KindInt:
			last, _ = p.parseInt()
		case token.KindIdentifier:
			if bytes.Compare(p.fm.Lookup(peek.ID), []byte("max")) == 0 {
				id, _ := p.parseIdentifier()
				last = ast.Integer{ID: id.ID}
				break
			}
			fallthrough
		default:
			return ast.Range{}, gotUnexpected(peek, token.KindMax, token.KindInt)
		}
	}

	var id token.UniqueID

	if last != first {
		id = p.fm.Merge(token.KindRange, first.ID, last.ID)
	} else {
		id = first.ID
	}
	return ast.Range{ID: id, Start: first, End: last}, err
}

func (p *impl) parseReservedTags() (reserved ast.ReservedTags, err error) {
	var r ast.ReservedTags

	first := p.peek()
	if first.Kind != token.KindInt {
		return r, gotUnexpected(first, token.KindInt)
	}
	item, err := p.parseRange()

	if err != nil {
		return r, err
	}

	items := []ast.Range{item}
	peek := p.peek()
	for ; peek.Kind == token.KindComma; peek = p.peek() {
		p.nextToken()

		if item, err = p.parseRange(); err != nil {
			return r, err
		}

		items = append(items, item)
	}

	last := items[len(items)-1]
	if peek.Kind != token.KindSemicolon {
		return r, gotUnexpected(peek, token.KindSemicolon)
	}

	if len(items) > 1 {
		r.ID = p.fm.Merge(token.KindReserved, first.ID, last.ID)
	} else {
		r.ID = last.ID
	}
	r.Items = items
	return r, nil
}

func (p *impl) parseReservedNames() (reserved ast.ReservedNames, err error) {
	var r ast.ReservedNames

	if peek := p.peek(); peek.Kind != token.KindStr {
		return r, gotUnexpected(peek, token.KindStr)
	}
	first := p.tokens[p.idx]
	item, _ := p.parseString()
	items := []ast.String{item}

	peek := p.peek()
	for ; peek.Kind == token.KindComma; peek = p.peek() {
		p.nextToken()

		if item, err = p.parseString(); err != nil {
			return r, err
		}

		items = append(items, item)
	}

	last := items[len(items)-1]
	if peek.Kind != token.KindSemicolon {
		return r, gotUnexpected(peek, token.KindSemicolon)
	}

	if len(items) > 1 {
		r.ID = p.fm.Merge(token.KindReserved, first.ID, last.ID)
	} else {
		r.ID = last.ID
	}
	r.Items = items
	return r, nil
}
