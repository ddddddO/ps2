package parser

func (p *phpParser) parseReference(ch byte) (*ASTNode, error) {
	if _, err := p.nextChar(); err != nil {
		return nil, err
	}
	if err := p.expectChar(':'); err != nil {
		return nil, err
	}
	referenceID, err := p.parseNumberString() // Reference ID
	if err != nil {
		return nil, err
	}
	if err := p.expectChar(';'); err != nil {
		return nil, err
	}

	t := "reference"
	if ch == 'R' {
		t = "Reference"
	}
	if ref := p.references.getByID(referenceID); ref != nil {
		return p.asignNode(t, ref.Value), nil
	}
	return p.asignNode(t, nil), nil
}
