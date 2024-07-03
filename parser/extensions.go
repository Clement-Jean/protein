package parser

func (p *Parser) parseExtensions() {
	p.pushState(stateReservedFinish)
	p.pushState(stateReservedRange)
}
