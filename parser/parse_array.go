package parser

import "fmt"

// Parses a PHP serialized array (e.g., a:N:{key;value;...}).
func (p *phpParser) parseArray() (*ASTNode, error) {
	if err := p.expectChar('a'); err != nil {
		return nil, err
	}
	if err := p.expectChar(':'); err != nil {
		return nil, err
	}
	size, err := p.parseNumberString()
	if err != nil {
		return nil, err
	}
	if err := p.expectChar(':'); err != nil {
		return nil, err
	}
	if err := p.expectChar('{'); err != nil {
		return nil, err
	}

	node := p.asignNode("array", make(map[interface{}]interface{}))
	p.references.store(node)
	childrenMap := make(map[interface{}]interface{})

	for i := 0; i < size; i++ {
		keyNode, err := p.parseValue()
		if err != nil {
			return nil, fmt.Errorf("failed to parse array key %d: %w", i, err)
		}
		valNode, err := p.parseValue()
		if err != nil {
			return nil, fmt.Errorf("failed to parse array value %d: %w", i, err)
		}
		p.references.store(valNode)

		key := keyNode.Value
		childrenMap[key] = valNode.Value

		// Add child node for AST representation
		// AST表現のために子ノードを追加
		childNode := *valNode // Make a copy
		childNode.Key = key
		node.Children = append(node.Children, &childNode)
	}

	node.Value = childrenMap // Store the actual Go map

	if err := p.expectChar('}'); err != nil {
		return nil, err
	}
	return node, nil
}
