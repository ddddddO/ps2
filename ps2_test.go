package ps2_test

import (
	"strings"
	"testing"

	"github.com/ddddddO/ps2"
)

func TestRun(t *testing.T) {
	tests := map[string]struct {
		serialized string
		want       func(json string) (bool, string)
	}{
		// TODO: object のパターン(ジョブキューのpayloadに載ってるやつ)を追加したい

		"array(include various types)": {
			serialized: `a:9:{s:10:"string_val";s:27:"こんにちは、世界！";s:7:"int_val";i:123;s:9:"bool_true";b:1;s:10:"bool_false";b:0;s:8:"null_val";N;s:9:"float_val";d:3.14159;s:18:"nested_assoc_array";a:3:{s:4:"name";s:12:"Go Developer";s:7:"details";a:2:{s:3:"age";i:30;s:6:"status";E:15:"Status:Inactive";}s:7:"hobbies";a:3:{i:0;s:6:"coding";i:1;s:7:"reading";i:2;s:6:"hiking";}}s:13:"indexed_array";a:5:{i:0;s:9:"りんご";i:1;s:9:"バナナ";i:2;s:12:"チェリー";i:3;i:100;i:4;b:1;}s:15:"object_instance";O:8:"MyObject":3:{s:10:"publicProp";s:15:"パブリック";s:16:"*protectedProp";i:456;s:19:"MyObjectprivateProp";a:1:{s:3:"key";s:5:"value";}}}`,
			want: func(json string) (bool, string) {
				wants := []string{
					`"string_val": "こんにちは、世界！"`,
					`"int_val": 123`,
					`"bool_true": true`,
					`"bool_false": false`,
					`"null_val": null`,
					`"float_val": 3.14159`,
					`"nested_assoc_array": {`,
					`"name": "Go Developer"`,
					`"details": {`, // 連想配列
					`"age": 30`,
					`"status": "Status:Inactive"`, // Enum
					`"hobbies": [`,                // 配列
					`"coding"`,
					`"reading"`,
					`"hiking"`,
					`"indexed_array": [`,
					`"りんご"`,
					`"バナナ"`,
					`"チェリー"`,
					`100`,
					`true`,
					`"object_instance": {`, // クラス
					`"*protectedProp": 456`,
					`"MyObjectprivateProp": {`,
					`"key": "value"`,
					`"__class_name": "MyObject"`,
					`"publicProp": "パブリック"`,
				}
				for _, w := range wants {
					if !strings.Contains(json, w) {
						return false, w
					}
				}
				return true, ""
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			input := strings.NewReader(tt.serialized)
			got, err := ps2.Run(input)
			if err != nil {
				t.Fatal(err)
			}
			if ok, want := tt.want(got); !ok {
				t.Errorf("\ngot: \n%s\n\nwant including: \n%s", got, want)
			}
		})

	}
}
