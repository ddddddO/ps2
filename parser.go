package ps2

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Custom parser struct to manage input string and current position.
// 入力文字列と現在の位置を管理するためのカスタムパーサー構造体
type phpParser struct {
	input string
	pos   int
}

// Creates a new parser instance.
// 新しいパーサーインスタンスを作成
func newPhpParser(input string) *phpParser {
	return &phpParser{input: input, pos: 0}
}

// Reads the next character and advances the position.
// 次の文字を読み込み、位置を進める
func (p *phpParser) nextChar() (rune, error) {
	if p.pos >= len(p.input) {
		return 0, errors.New("unexpected end of input")
	}
	r, size := utf8.DecodeRuneInString(p.input[p.pos:])
	p.pos += size
	return r, nil
}

// Peeks at the next character without advancing the position.
// 位置を進めずに次の文字を覗き見る
func (p *phpParser) peekChar() (rune, error) {
	if p.pos >= len(p.input) {
		return 0, errors.New("unexpected end of input")
	}
	r, _ := utf8.DecodeRuneInString(p.input[p.pos:])
	return r, nil
}

// Expects a specific character at the current position.
// 現在の位置に特定の文字があることを期待
func (p *phpParser) expectChar(expected rune) error {
	ch, err := p.nextChar()
	if err != nil {
		return err
	}
	if ch != expected {
		return fmt.Errorf("expected '%c', but got '%c' at position %d", expected, ch, p.pos-1)
	}
	return nil
}

// Represents a node in the conceptual AST.
// 概念的なASTのノードを表す構造体
type ASTNode struct {
	Type      string      // 例: "array", "object", "string", "int", "bool", "null"
	Value     interface{} // ノードの実際の値 (文字列、数値、マップ、スライスなど)
	Children  []*ASTNode  // 子ノード (配列やオブジェクトの場合)
	Key       interface{} // 親が配列/オブジェクトの場合のキー (string or int)
	PropName  string      // オブジェクトのプロパティ名の場合
	ClassName string      // オブジェクトの場合のクラス名
}

// Parses a PHP serialized value based on its type prefix.
// 型プレフィックスに基づいてPHPのシリアライズされた値を解析
func (p *phpParser) parseValue() (*ASTNode, error) {
	if p.pos >= len(p.input) {
		return nil, errors.New("unexpected end of input when parsing value type")
	}

	ch := p.input[p.pos]
	p.pos++ // Consume the type character
	switch ch {
	case 's':
		p.pos-- // Go back to 's' for parseString
		return p.parseString()
	case 'i':
		p.pos-- // Go back to 'i' for parseInteger
		return p.parseInteger()
	case 'b':
		p.pos-- // Go back to 'b' for parseBoolean
		return p.parseBoolean()
	case 'N':
		p.pos-- // Go back to 'N' for parseNull
		return p.parseNull()
	case 'a':
		p.pos-- // Go back to 'a' for parseArray
		return p.parseArray()
	case 'O':
		p.pos-- // Go back to 'O' for parseObject
		return p.parseObject()
	case 'd':
		p.pos-- // Go back to 'd' for parseFloat
		return p.parseFloat()
	case 'R', 'r': // Reference, currently not fully supported by this parser for deep parsing
		// PHP references (R:N;) point to a previously parsed element.
		// For simplicity, we'll just consume it and return a placeholder.
		// For a full implementation, you'd need to store parsed objects in a map
		// and retrieve them here.
		if err := p.expectChar(':'); err != nil {
			return nil, err
		}
		_, err := p.parseNumberString() // Reference ID
		if err != nil {
			return nil, err
		}
		if err := p.expectChar(';'); err != nil {
			return nil, err
		}
		return &ASTNode{Type: "reference", Value: nil}, nil // Placeholder
	default:
		return nil, fmt.Errorf("unknown PHP serialized type '%c' at position %d", ch, p.pos-1)
	}
}

// Parses an integer value (e.g., "123" from i:123;).
// 整数値（例: i:123; から "123"）を解析
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

// Parses a PHP serialized string (e.g., s:N:"string";).
// PHPのシリアライズされた文字列（例: s:N:"string";）を解析
func (p *phpParser) parseString() (*ASTNode, error) {
	if err := p.expectChar('s'); err != nil {
		return nil, err
	}
	if err := p.expectChar(':'); err != nil {
		return nil, err
	}
	length, err := p.parseNumberString()
	if err != nil {
		return nil, err
	}
	if err := p.expectChar(':'); err != nil {
		return nil, err
	}
	if err := p.expectChar('"'); err != nil {
		return nil, err
	}

	// Read the actual string value based on length
	// 長さに応じて実際の文字列値を読み取る
	// ここで指定された 'length' はバイト数なので、正確にそのバイト数を読み込む
	start := p.pos
	end := start + length

	if end > len(p.input) {
		return nil, fmt.Errorf("string content too short for declared length %d. Current pos %d, End pos %d, Input length %d", length, start, end, len(p.input))
	}
	val := p.input[start:end]

	// 処理追加した。このあたりでバグあるかもしれない
	// 「*」が先頭にある場合、*の前後はnullバイト(ref: https://www.php.net/manual/ja/function.serialize.php#refsect1-function.serialize-parameters の「注意」)
	// ただ、シリアライズされた文字列をコピペしてターミナルに張り付けるとnullバイトが消えるので、その場合はnullバイト分を除くため、end-2する
	if strings.HasPrefix(val, "*") && len(p.input) >= start {
		if start < end-2 {
			end -= 2
			val = p.input[start:end]
		}
	}
	if strings.HasPrefix(val, "�*") {
		// �が3byte分なので、endを、3byte x 2 - 2 する
		end = end + 3*2 - 2
		val = p.input[start:end]
	} else if strings.Contains(val, "�") {
		cnt := strings.Count(val, "�")
		end = end + 3*cnt - cnt
		val = p.input[start:end]
	}

	p.pos = end // posを正確に更新

	if r, err := p.peekChar(); err == nil && r != '"' {
		// private でクラスの変数の場合、[NULLバイト]App\Xxxx[Nullバイト]isFlag みたいに、Nullバイトが2つ分入って進みすぎてしまう
		// よくないと思うけど、ここで次の文字が意図したものでなければ、Nullバイトが含まれていたとみなして、endから-2する
		end -= 2
		val = p.input[start:end]
		p.pos = end
	}

	if err := p.expectChar('"'); err != nil {
		return nil, err
	}
	if err := p.expectChar(';'); err != nil {
		return nil, err
	}

	return &ASTNode{Type: "string", Value: val}, nil
}

// Parses a PHP serialized integer (e.g., i:V;).
// PHPのシリアライズされた整数（例: i:V;）を解析
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
	return &ASTNode{Type: "int", Value: val}, nil
}

// Parses a PHP serialized boolean (e.g., b:V;).
// PHPのシリアライズされた真偽値（例: b:V;）を解析
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
	return &ASTNode{Type: "bool", Value: val}, nil
}

// Parses a PHP serialized null (e.g., N;).
// PHPのシリアライズされたnull（例: N;）を解析
func (p *phpParser) parseNull() (*ASTNode, error) {
	if err := p.expectChar('N'); err != nil {
		return nil, err
	}
	if err := p.expectChar(';'); err != nil {
		return nil, err
	}
	return &ASTNode{Type: "null", Value: nil}, nil
}

// Parses a PHP serialized float (e.g., d:V;).
// PHPのシリアライズされた浮動小数点数（例: d:V;）を解析
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
	return &ASTNode{Type: "float", Value: val}, nil
}

// Parses a PHP serialized array (e.g., a:N:{key;value;...}).
// PHPのシリアライズされた配列（例: a:N:{key;value;...}）を解析
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

	node := &ASTNode{Type: "array", Value: make(map[interface{}]interface{})}
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

// Parses a PHP serialized object (e.g., O:L:"ClassName":N:{prop_name;prop_val;...}).
// PHPのシリアライズされたオブジェクト（例: O:L:"ClassName":N:{prop_name;prop_val;...}）を解析
func (p *phpParser) parseObject() (*ASTNode, error) {
	if err := p.expectChar('O'); err != nil {
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

	node := &ASTNode{Type: "object", ClassName: className, Value: make(map[string]interface{})}
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
