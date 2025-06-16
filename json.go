package ps2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/ddddddO/ps2/parser"
)

// JSON出力用の構造体。ASTNodeの情報をJSONにマッピングする。
// Represents a JSON-friendly version of ASTNode for output.
type JSONNode struct {
	Type      string      `json:"type"`                   // ノードの型
	Value     interface{} `json:"value,omitempty"`        // プリミティブな値、または配列/オブジェクトの実際の値（マップやスライス）
	ClassName string      `json:"__class_name,omitempty"` // オブジェクトの場合のクラス名
	Key       interface{} `json:"key,omitempty"`          // 親が配列/オブジェクトの場合のキー (このノードが子ノードの場合)
	PropName  string      `json:"prop_name,omitempty"`    // 親がオブジェクトの場合のプロパティ名 (このノードがプロパティの場合)
	Children  []*JSONNode `json:"children,omitempty"`     // 子ノードのリスト (AST構造を維持するためのもの)
}

func (j *JSONNode) MarshalJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(j.Value); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Converts an ASTNode tree to a JSONNode tree.
// この関数は、ASTNodeの構造をJSONNodeに変換し、特に配列の'Value'フィールドを
// PHPのjson_encodeの挙動に合わせてJSON配列またはJSONオブジェクトに変換します。
func astNodeToJSONNode(astNode *parser.ASTNode) *JSONNode {
	if astNode == nil {
		return nil
	}

	jsonNode := &JSONNode{
		Type:      astNode.Type,
		ClassName: astNode.ClassName,
		Key:       astNode.Key,
		PropName:  astNode.PropName,
	}

	switch astNode.Type {
	case parser.ASTNodeTypeString, parser.ASTNodeTypeInt, parser.ASTNodeTypeBool, parser.ASTNodeTypeNull, parser.ASTNodeTypeFloat, parser.ASTNodeTypeEnum:
		// プリミティブ型の場合、Valueを直接設定
		jsonNode.Value = astNode.Value
	case parser.ASTNodeTypeReference1, parser.ASTNodeTypeReference2:
		jsonNode.Value = fmt.Sprintf("[[PHP_REFERENCE_DATA: %+v]]", astNode.Value)

		// FIXME: ちょっとここ大変で、↑みたいにREFERENCEされてたらそのデータを文字列として出すようにした

		// 自己参照型は、それ自身のパースが終わってなくてValueがゼロ値となるよう?で、その場合は自己参照型と見做してプレースホルダを出力
		// 自己参照型でなく、そもそもValueがゼロ値なreferenceの可能性はあるかもしれない
		if astNode.Value != nil {
			switch m := astNode.Value.(type) {
			case map[string]interface{}:
				if len(m) == 0 {
					jsonNode.Value = "[[MAYBE_PHP_SELF_REFERENCE]]"
				}
			}
		}
		// } else {
		// 	switch m := astNode.Value.(type) {
		// 	case map[string]interface{}:
		// 		if len(m) == 0 {
		// 			jsonNode.Value = fmt.Sprintf("[[PHP_SELF_REFERENCE or zero value: %+v]]", m)
		// 		}
		// 	case string:
		// 		if len(m) == 0 {
		// 			jsonNode.Value = fmt.Sprintf("[[PHP_SELF_REFERENCE or zero value: %s]]", m)
		// 		}
		// 	case int, float64:
		// 		if m == 0 {
		// 			jsonNode.Value = fmt.Sprintf("[[PHP_SELF_REFERENCE or zero value: %d]]", m)
		// 		}
		// 	case bool:
		// 		if !m {
		// 			jsonNode.Value = fmt.Sprintf("[[PHP_SELF_REFERENCE or zero value: %b]]", m)
		// 		}
		// 	}
		// }
	case parser.ASTNodeTypeArray:
		phpMap := astNode.Value.(map[interface{}]interface{})

		// PHP配列が純粋な数値インデックスの連続した配列であるかを判定
		isSequentialArray := true
		numKeys := len(phpMap)
		if numKeys > 0 {
			intKeys := make([]int, 0, numKeys)
			for k := range phpMap {
				if intKey, ok := k.(int); ok {
					intKeys = append(intKeys, intKey)
				} else {
					isSequentialArray = false // 整数以外のキーが存在する
					break
				}
			}

			if isSequentialArray { // 全てのキーが整数である場合のみ、連続性をチェック
				sort.Ints(intKeys) // キーをソート
				for i := 0; i < numKeys; i++ {
					if intKeys[i] != i {
						isSequentialArray = false // キーが0から連続していない
						break
					}
				}
			}
		} else {
			// 空の配列はJSON配列として扱う
			isSequentialArray = true
		}

		if isSequentialArray {
			// 数値インデックスの連続した配列の場合、JSON配列（Goのスライス）に変換
			jsonArray := make([]interface{}, numKeys)
			for i := 0; i < numKeys; i++ {
				// 該当する子ASTNodeを見つけて、その値を再帰的にJSONValueに変換
				var childAST *parser.ASTNode
				for _, child := range astNode.Children {
					if child.Key != nil {
						if k, ok := child.Key.(int); ok && k == i {
							childAST = child
							break
						}
					}
				}
				if childAST != nil {
					// Recursively convert the child's value to its appropriate JSON representation.
					jsonArray[i] = astNodeToJSONNode(childAST).Value
				} else {
					// Fallback: child ASTNodeが見つからない場合は、生の値をそのまま使用（ただし、複雑な型の場合再帰的にならない）
					jsonArray[i] = phpMap[i]
				}
			}
			jsonNode.Value = jsonArray
		} else {
			// 連想配列、または非連続な数値キーを持つ配列の場合、JSONオブジェクト（Goのmap[string]interface{}）に変換
			jsonMap := make(map[string]interface{})
			for k, v := range phpMap {
				var jsonKey string
				if keyStr, ok := k.(string); ok {
					jsonKey = keyStr
				} else if keyInt, ok := k.(int); ok {
					jsonKey = fmt.Sprintf("%d", keyInt) // 整数キーはJSONオブジェクトのキーとして文字列に変換
				} else {
					jsonKey = fmt.Sprintf("%v", k) // その他の型のキーは文字列に変換
				}

				// 該当する子ASTNodeを見つけて、その値を再帰的にJSONValueに変換
				var childAST *parser.ASTNode
				for _, child := range astNode.Children {
					if child.Key == k {
						childAST = child
						break
					}
				}
				if childAST != nil {
					jsonMap[jsonKey] = astNodeToJSONNode(childAST).Value
				} else {
					jsonMap[jsonKey] = v
				}
			}
			jsonNode.Value = jsonMap
		}
	case parser.ASTNodeTypeObject, parser.ASTNodeTypeCustom:
		// オブジェクトの場合、Goのmap[string]interface{}に変換（プロパティ名は既に文字列）
		phpObjectMap := astNode.Value.(map[string]interface{})
		jsonMap := make(map[string]interface{})
		for k := range phpObjectMap { // 元のキーをイテレート
			jsonKey := k // プロパティ名は既に文字列

			// 該当する子ASTNodeを見つけて、その値を再帰的にJSONValueに変換
			var childAST *parser.ASTNode
			for _, child := range astNode.Children {
				if child.PropName == k {
					childAST = child
					break
				}
			}
			if childAST != nil {
				jsonMap[jsonKey] = astNodeToJSONNode(childAST).Value
			} else {
				jsonMap[jsonKey] = phpObjectMap[k]
			}
		}
		jsonNode.Value = jsonMap
	}

	// JSONNodeの'Children'フィールドはAST構造そのものを表すため、常に再帰的に構築する
	if len(astNode.Children) > 0 {
		jsonNode.Children = make([]*JSONNode, len(astNode.Children))
		for i, child := range astNode.Children {
			jsonNode.Children[i] = astNodeToJSONNode(child) // 子ノードのJSONNodeを再帰的に構築
		}
	}

	return jsonNode
}
