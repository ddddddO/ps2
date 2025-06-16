package parser

// Parses a PHP serialized integer (e.g., i:V;).
func (p *phpParser) parseInteger() (*ASTNode, error) {
	if err := p.expectChar('i'); err != nil {
		return nil, err
	}
	if err := p.expectChar(':'); err != nil {
		return nil, err
	}
	val, err := p.parseNumberString()
	if err != nil {
		return nil, err
	}
	if err := p.expectChar(';'); err != nil {
		return nil, err
	}
	return p.asignNode("int", val), nil
}
