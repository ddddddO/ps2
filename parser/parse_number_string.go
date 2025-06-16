package parser

import (
	"fmt"
	"strconv"
)

// Parses an integer value (e.g., "123" from i:123;).
func (p *phpParser) parseNumberString() (int, error) {
	start := p.pos
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch >= '0' && ch <= '9' || ch == '-' {
			p.pos++
		} else {
			break
		}
	}
	numStr := p.input[start:p.pos]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("invalid number format: %w", err)
	}
	return num, nil
}
