package ps2

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/ddddddO/ps2/parser"
	"github.com/goccy/go-yaml"
	"github.com/pelletier/go-toml/v2"
)

func Run(input io.Reader, options ...Option) (string, error) {
	cfg := NewConfig(options)

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

	switch cfg.outputType {
	case outputTypeJSON:
		return strings.TrimSuffix(buf.String(), "\n"), nil
	case outputTypeYAML:
		y, err := yaml.JSONToYAML(buf.Bytes())
		if err != nil {
			return "", err
		}
		return strings.TrimSuffix(string(y), "\n"), nil
	case outputTypeTOML:
		d := json.NewDecoder(&buf)
		buf2 := bytes.Buffer{}
		e := toml.NewEncoder(&buf2)
		d.UseNumber()
		e.SetMarshalJsonNumbers(true)
		var v interface{}
		if err := d.Decode(&v); err != nil {
			return "", err
		}
		if err := e.Encode(v); err != nil {
			return "", err
		}
		return strings.TrimSuffix(buf2.String(), "\n"), nil
	default:
		return strings.TrimSuffix(buf.String(), "\n"), nil
	}
}
