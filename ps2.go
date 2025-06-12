package ps2

import (
	"bufio"
	"encoding/json"
	"io"
	// UTF-8文字の処理用
)

func Run(input io.Reader) (string, error) {
	scanner := bufio.NewScanner(input)
	phpSerializedString := ""
	for scanner.Scan() {
		phpSerializedString += scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	parser := newPhpParser(phpSerializedString)
	rootNode, err := parser.parseValue()
	if err != nil {
		return "", err
	}

	jsonRootNode := astNodeToJSONNode(rootNode)
	// jsonRootNode.Children = nil
	jsonData, err := json.MarshalIndent(jsonRootNode, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
