package parser

import (
	"fmt"
	"strings"
)

// Parses a PHP serialized custom (e.g., C:L:"ClassName":N:{prop_name;prop_val;...}).
func (p *phpParser) parseCustom() (*ASTNode, error) {
	if err := p.expectChar('C'); err != nil {
		return nil, err
	}
	if err := p.expectChar(':'); err != nil {
		return nil, err
	}
	classNameLen, err := p.parseNumberString()
	if err != nil {
		return nil, err
	}
	if err := p.expectChar(':'); err != nil {
		return nil, err
	}
	if err := p.expectChar('"'); err != nil {
		return nil, err
	}

	classNameStart := p.pos
	if classNameStart+classNameLen > len(p.input) {
		return nil, fmt.Errorf("class name length mismatch: expected %d, available %d", classNameLen, len(p.input)-classNameStart)
	}
	className := p.input[classNameStart : classNameStart+classNameLen]
	p.pos += classNameLen

	if err := p.expectChar('"'); err != nil {
		return nil, err
	}
	if err := p.expectChar(':'); err != nil {
		return nil, err
	}

	numProps, err := p.parseNumberString()
	if err != nil {
		return nil, err
	}
	if err := p.expectChar(':'); err != nil {
		return nil, err
	}
	if err := p.expectChar('{'); err != nil {
		return nil, err
	}

	node := p.asignNodeWithClassname("custom", className, make(map[string]interface{}))
	p.references.store(node)
	propertiesMap := make(map[string]interface{})
	propertiesMap["__class_name"] = className

	for i := 0; i < numProps; i++ {
		// Property name is a string (s:N:"prop_name";)
		// プロパティ名は文字列 (s:N:"prop_name";)
		propNameNode, err := p.parseString()
		if err != nil {
			return nil, fmt.Errorf("failed to parse object property name %d: %w", i, err)
		}
		propName := propNameNode.Value.(string)

		// 厳密には else if ブロックのコメント
		// PHP object properties can be public, protected, or private.
		// Protected properties start with a null byte (0x00), then '*' then null byte.
		// Private properties start with a null byte, then class name, then null byte.
		// For simplicity, we just extract the name after null bytes if present.
		// Public properties have no prefix.
		cleanPropName := propName
		if strings.HasPrefix(propName, "�") {
			parts := strings.Split(propName, "�")
			if len(parts) >= 3 {
				cleanPropName = parts[2] // Private: �ClassName�propName, Protected: �*�propName
			}
		} else if strings.HasPrefix(propName, "\x00") {
			parts := strings.Split(propName, "\x00")
			if len(parts) >= 3 {
				cleanPropName = parts[2] // Private: \x00ClassName\x00propName, Protected: \x00*\x00propName
			}
		}

		propValNode, err := p.parseValue()
		if err != nil {
			return nil, fmt.Errorf("failed to parse object property value %d: %w", i, err)
		}
		p.references.store(propValNode)

		propertiesMap[cleanPropName] = propValNode.Value

		childNode := *propValNode // Make a copy
		childNode.PropName = cleanPropName
		node.Children = append(node.Children, &childNode)
	}

	node.Value = propertiesMap // Store the actual Go map

	if err := p.expectChar('}'); err != nil {
		return nil, err
	}
	return node, nil
}
