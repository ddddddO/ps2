package parser

func (p *phpParser) parseReference(ch byte) (*ASTNode, error) {
	if err := p.expectChars('R', 'r'); err != nil {
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

	t := ASTNodeTypeReference1
	if ch == 'R' {
		t = ASTNodeTypeReference2
	}
	if ref := p.references.getByID(referenceID); ref != nil {
		return p.asignNode(t, ref.Value), nil
	}
	return p.asignNode(t, nil), nil
}
