package ps2

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/ddddddO/ps2/parser"
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

	parser := parser.New(phpSerializedString)
	rootNode, err := parser.Parse()
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	jsonRootNode := astNodeToJSONNode(rootNode)
	// jsonRootNode.Children = nil
	if err := encoder.Encode(jsonRootNode); err != nil {
		return "", err
	}

	return strings.TrimSuffix(buf.String(), "\n"), nil
}
