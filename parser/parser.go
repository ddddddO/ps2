package parser

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

// Custom parser struct to manage input string and current position.
type phpParser struct {
	input       string
	pos         int
	nodeIndexer uint // 各ノードを一意にする識別するための数
	references  *references
}

func New(input string) *phpParser {
	return &phpParser{input: input, pos: 0, references: newReferences()}
}

func (p *phpParser) Parse() (*ASTNode, error) {
	return p.parseValue()
}

// Parses a PHP serialized value based on its type prefix.
func (p *phpParser) parseValue() (*ASTNode, error) {
	if p.pos >= len(p.input) {
		return nil, errors.New("unexpected end of input when parsing value type")
	}

	ch := p.input[p.pos]
	switch ch {
	case 's':
		return p.parseString()
	case 'i':
		return p.parseInteger()
	case 'b':
		return p.parseBoolean()
	case 'N':
		return p.parseNull()
	case 'a':
		return p.parseArray()
	case 'O':
		return p.parseObject()
	case 'C':
		return p.parseCustom()
	case 'd':
		return p.parseFloat()
	case 'E':
		return p.parseEnum()
	case 'R', 'r':
		return p.parseReference(ch)
	default:
		return nil, fmt.Errorf("unknown PHP serialized type '%c' at position %d", ch, p.pos-1)
	}
}

// Reads the next character and advances the position.
func (p *phpParser) nextChar() (rune, error) {
	if p.pos >= len(p.input) {
		return 0, errors.New("unexpected end of input")
	}
	r, size := utf8.DecodeRuneInString(p.input[p.pos:])
	p.pos += size
	return r, nil
}

// Peeks at the next character without advancing the position.
func (p *phpParser) peekChar() (rune, error) {
	if p.pos >= len(p.input) {
		return 0, errors.New("unexpected end of input")
	}
	r, _ := utf8.DecodeRuneInString(p.input[p.pos:])
	return r, nil
}

// Expects a specific character at the current position.
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
type ASTNode struct {
	Index     uint        // ノードの通し番号
	Type      string      // 例: "array", "object", "string", "int", "bool", "null"
	Value     interface{} // ノードの実際の値 (文字列、数値、マップ、スライスなど)
	Children  []*ASTNode  // 子ノード (配列やオブジェクトの場合)
	Key       interface{} // 親が配列/オブジェクトの場合のキー (string or int)
	PropName  string      // オブジェクトのプロパティ名の場合
	ClassName string      // オブジェクトの場合のクラス名
}

func (p *phpParser) asignNode(typ string, value interface{}) *ASTNode {
	index := p.nodeIndexer
	p.nodeIndexer++
	return &ASTNode{
		Index: index,
		Type:  typ,
		Value: value,
	}
}

func (p *phpParser) asignNodeWithClassname(typ string, classname string, value interface{}) *ASTNode {
	node := p.asignNode(typ, value)
	node.ClassName = classname
	return node
}

const (
	ASTNodeTypeString     = "string"
	ASTNodeTypeInt        = "int"
	ASTNodeTypeBool       = "bool"
	ASTNodeTypeNull       = "null"
	ASTNodeTypeFloat      = "float"
	ASTNodeTypeEnum       = "enum"
	ASTNodeTypeReference1 = "reference"
	ASTNodeTypeReference2 = "Reference"
	ASTNodeTypeArray      = "array"
	ASTNodeTypeObject     = "object"
	ASTNodeTypeCustom     = "custom"
)
