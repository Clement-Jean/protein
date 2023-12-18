package parser

import (
	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/token"
)

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
