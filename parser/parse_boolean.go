package parser

import "fmt"

// Parses a PHP serialized boolean (e.g., b:V;).
func (p *phpParser) parseBoolean() (*ASTNode, error) {
	if err := p.expectChar('b'); err != nil {
		return nil, err
	}
	if err := p.expectChar(':'); err != nil {
		return nil, err
	}
	valStr, err := p.nextChar()
	if err != nil {
		return nil, err
	}
	var val bool
	if valStr == '1' {
		val = true
	} else if valStr == '0' {
		val = false
	} else {
		return nil, fmt.Errorf("invalid boolean value '%c' at position %d", valStr, p.pos-1)
	}
	if err := p.expectChar(';'); err != nil {
		return nil, err
	}
	return p.asignNode("bool", val), nil
}
