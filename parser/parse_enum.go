package parser

import (
	"fmt"
	"strings"
)

// Parses a PHP serialized enum (e.g., E:15:"Status:Inactive";).
func (p *phpParser) parseEnum() (*ASTNode, error) {
	if err := p.expectChar('E'); err != nil {
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

	return p.asignNode("enum", val), nil
}
