package parser

import (
	"fmt"
	"strconv"
)

// Parses a PHP serialized float (e.g., d:V;).
func (p *phpParser) parseFloat() (*ASTNode, error) {
	if err := p.expectChar('d'); err != nil {
		return nil, err
	}
	if err := p.expectChar(':'); err != nil {
		return nil, err
	}
	start := p.pos
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if (ch >= '0' && ch <= '9') || ch == '.' || ch == '-' || ch == '+' || ch == 'e' || ch == 'E' {
			p.pos++
		} else {
			break
		}
	}
	numStr := p.input[start:p.pos]
	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid float format: %w", err)
	}
	if err := p.expectChar(';'); err != nil {
		return nil, err
	}
	return p.asignNode("float", val), nil
}
