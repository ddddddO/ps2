package parser

// Parses a PHP serialized null (e.g., N;).
func (p *phpParser) parseNull() (*ASTNode, error) {
	if err := p.expectChar('N'); err != nil {
		return nil, err
	}
	if err := p.expectChar(';'); err != nil {
		return nil, err
	}
	return p.asignNode("null", nil), nil
}
