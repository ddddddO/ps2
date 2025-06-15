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
		"object(include various types)": {
			serialized: `O:13:"MySimpleClass":17:{s:10:"publicProp";s:16:"Top Level Object";s:26:"MySimpleClassprivateProp";i:999;s:16:"*protectedProp";a:3:{s:10:"assoc_key1";s:12:"assoc_value1";s:10:"assoc_key2";i:789;s:17:"deep_nested_array";a:2:{s:9:"sub_key_x";s:11:"sub_value_x";s:9:"sub_key_y";d:12.34;}}s:6:"parent";N;s:15:"nestedArrayData";a:5:{i:0;s:6:"Item A";i:1;s:6:"Item B";i:2;s:6:"Item C";i:3;i:10;i:4;b:0;}s:16:"sharedStringRef1";s:33:"共有される文字列データ";s:16:"sharedStringRef2";s:33:"共有される文字列データ";s:16:"sharedObjectRef1";O:13:"MySimpleClass":6:{s:10:"publicProp";s:24:"共通オブジェクト";s:26:"MySimpleClassprivateProp";i:500;s:16:"*protectedProp";a:1:{s:6:"shared";b:1;}s:6:"parent";N;s:15:"nestedArrayData";a:0:{}s:9:"nullValue";N;}s:16:"sharedObjectRef2";r:19;s:21:"customSerializableObj";O:20:"MyCustomSerializable":2:{s:1:"s";s:36:"オブジェクト内のカスタム";s:1:"n";i:777;}s:12:"userRoleEnum";E:15:"UserRole:Editor";s:10:"statusEnum";E:13:"Status:Active";s:9:"nullValue";N;s:11:"booleanTrue";b:1;s:10:"floatValue";d:45.67;s:12:"integerValue";i:123;s:14:"japaneseString";s:36:"これは日本語の文字列です";}`,
			want: func(json string) (bool, string) {
				wants := []string{
					`"*protectedProp": {`,
					`"assoc_key1": "assoc_value1"`,
					`"assoc_key2": 789`,
					`"deep_nested_array": {`,
					`"sub_key_x": "sub_value_x"`,
					`"sub_key_y": 12.34`,
					`"MySimpleClassprivateProp": 999`,
					`"__class_name": "MySimpleClass"`,
					`"booleanTrue": true`,
					`"customSerializableObj": {`,
					`"__class_name": "MyCustomSerializable"`,
					`"n": 777`,
					`"s": "オブジェクト内のカスタム"`,
					`"floatValue": 45.67`,
					`"integerValue": 123`,
					`"japaneseString": "これは日本語の文字列です"`,
					`"nestedArrayData": [`,
					`"Item A"`,
					`"Item B"`,
					`"Item C"`,
					`10`,
					`false`,
					`"nullValue": null`,
					`"parent": null`,
					`"publicProp": "Top Level Object"`,
					`"sharedObjectRef1": {`,
					`"*protectedProp": {`,
					`"shared": true`,
					`"MySimpleClassprivateProp": 500`,
					`"__class_name": "MySimpleClass"`,
					`"nestedArrayData": []`,
					`"nullValue": null`,
					`"parent": null`,
					`"publicProp": "共通オブジェクト"`,
					`"sharedObjectRef2": "[[PHP_REFERENCE_DATA: map[*protectedProp:map[shared:true] MySimpleClassprivateProp:500 __class_name:MySimpleClass nestedArrayData:map[] nullValue:\u003cnil\u003e parent:\u003cnil\u003e publicProp:共通オブジェクト]]]"`,
					`"sharedStringRef1": "共有される文字列データ"`,
					`"sharedStringRef2": "共有される文字列データ"`,
					`"statusEnum": "Status:Active"`,
					`"userRoleEnum": "UserRole:Editor"`,
				}
				for _, w := range wants {
					if !strings.Contains(json, w) {
						return false, w
					}
				}
				return true, ""
			},
		},
		"array(include various types)": {
			serialized: `a:10:{s:10:"string_val";s:27:"こんにちは、世界！";s:7:"int_val";i:123;s:9:"bool_true";b:1;s:10:"bool_false";b:0;s:8:"null_val";N;s:9:"float_val";d:3.14159;s:18:"nested_assoc_array";a:3:{s:4:"name";s:12:"Go Developer";s:7:"details";a:2:{s:3:"age";i:30;s:6:"status";E:15:"Status:Inactive";}s:7:"hobbies";a:3:{i:0;s:6:"coding";i:1;s:7:"reading";i:2;s:6:"hiking";}}s:13:"indexed_array";a:5:{i:0;s:9:"りんご";i:1;s:9:"バナナ";i:2;s:12:"チェリー";i:3;i:100;i:4;b:1;}s:15:"object_instance";O:8:"MyObject":3:{s:10:"publicProp";s:15:"パブリック";s:16:"*protectedProp";i:456;s:19:"MyObjectprivateProp";a:1:{s:3:"key";s:5:"value";}}s:22:"custom_object_instance";O:9:"CustomObj":1:{s:4:"prop";s:5:"xxxxx";}}`,
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
					`"custom_object_instance": {`, // カスタム
					`"__class_name": "CustomObj"`,
					`"prop": "xxxxx"`,
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

func TestRun_parts(t *testing.T) {
	tests := map[string]struct {
		serialized string
		want       string
	}{
		"array(map)": {
			serialized: `a:2:{s:4:"name";s:12:"Go Developer";s:2:"xx";i:123;}`,
			want: `
{
  "name": "Go Developer",
  "xx": 123
}`,
		},
		"array(list)": {
			serialized: `a:3:{i:0;s:6:"coding";i:1;s:7:"reading";i:2;s:6:"hiking";}`,
			want: `
[
  "coding",
  "reading",
  "hiking"
]`,
		},
		"string": {
			serialized: `s:27:"こんにちは、世界！";`,
			want:       `"こんにちは、世界！"`,
		},
		"enum": {
			serialized: `E:15:"Status:Inactive";`,
			want:       `"Status:Inactive"`,
		},
		"int": {
			serialized: `i:123;`,
			want:       `123`,
		},
		"float": {
			serialized: `d:3.14159;`,
			want:       `3.14159`,
		},
		"bool(false)": {
			serialized: `b:0;`,
			want:       `false`,
		},
		"bool(true)": {
			serialized: `b:1;`,
			want:       `true`,
		},
		"null": {
			serialized: `N;`,
			want:       `null`,
		},
		"object": {
			serialized: `O:12:"SimpleObject":1:{s:4:"name";s:8:"Object A";}`,
			want: `
{
  "__class_name": "SimpleObject",
  "name": "Object A"
}`,
		},
		"custom": {
			serialized: `C:12:"SimpleObject":1:{s:4:"name";s:8:"Object A";}`,
			want: `
{
  "__class_name": "SimpleObject",
  "name": "Object A"
}`,
		},
		"reference(object)": {
			serialized: `a:2:{s:9:"first_obj";O:12:"SimpleObject":1:{s:4:"name";s:8:"Object A";}s:10:"second_obj";r:2;}`,
			want: `
{
  "first_obj": {
    "__class_name": "SimpleObject",
    "name": "Object A"
  },
  "second_obj": "[[PHP_REFERENCE_DATA: map[__class_name:SimpleObject name:Object A]]]"
}`,
		},
		// TODO: "reference(value)" のケース
		// TODO: 以下、自己参照のケースは、ゼロ値になるっぽいからそれ判定してMAYBE_SELF_REFERENCEみたいな文字列出すのいいかも
		"reference(self)": {
			serialized: `O:8:"MyObject":2:{s:4:"name";s:30:"自己参照オブジェクト";s:4:"self";r:1;}`,
			want: `
{
  "__class_name": "MyObject",
  "name": "自己参照オブジェクト",
  "self": "[[PHP_REFERENCE_DATA: map[]]]"
}`,
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

			want := strings.TrimPrefix(tt.want, "\n")
			if got != want {
				t.Errorf("\ngot: \n%s\n\nwant: \n%s", got, want)
			}
		})
	}
}

// TestRun()のobject(include various types)ケースのserializedの元データ
// https://3v4l.org/#live で
/*
<?php

// --- 1. シンプルな独自定義クラス ---
class MySimpleClass
{
    public string $publicProp;
    private int $privateProp; // privateプロパティもシリアライズされる
    protected array $protectedProp;
    public ?MySimpleClass $parent = null; // 自己参照または他のオブジェクト参照用
    public array $nestedArrayData;
    public string $sharedStringRef1;
    public string $sharedStringRef2;
    public object $sharedObjectRef1;
    public object $sharedObjectRef2;
    public object $customSerializableObj;
    public object $userRoleEnum; // Backed Enum
    public object $statusEnum;   // Pure Enum
    public $nullValue;
    public bool $booleanTrue;
    public float $floatValue;
    public int $integerValue;
    public string $japaneseString;

    public function __construct(string $public, int $private, array $protected)
    {
        $this->publicProp = $public;
        $this->privateProp = $private;
        $this->protectedProp = $protected;
        $this->nestedArrayData = []; // 初期化
    }

    public function getPrivateProp(): int
    {
        return $this->privateProp;
    }
}

// --- 2. カスタムシリアライズ可能なクラス (PHP 7.4+ の __serialize/__unserialize) ---
class MyCustomSerializable
{
    public string $dataString;
    private int $dataNumber;

    public function __construct(string $s, int $n)
    {
        $this->dataString = $s;
        $this->dataNumber = $n;
    }

    // オブジェクトがシリアライズされる直前に呼び出される
    // シリアライズしたいプロパティを連想配列で返す
    public function __serialize(): array
    {
        echo "__serialize() called for MyCustomSerializable\n";
        return [
            's' => $this->dataString,
            'n' => $this->dataNumber,
        ];
    }

    // オブジェクトがデシリアライズされた直後に呼び出される
    // __serialize() で返された連想配列を受け取る
    public function __unserialize(array $data): void
    {
        echo "__unserialize() called for MyCustomSerializable\n";
        $this->dataString = $data['s'];
        $this->dataNumber = $data['n'];
    }

    public function getCustomData(): string
    {
        return "Custom: " . $this->dataString . " / " . $this->dataNumber;
    }
}

// --- 3. Enum の定義 (PHP 8.1+) ---
enum UserRole: string // Backed Enum
{
    case Admin = 'admin';
    case Editor = 'editor';
    case Viewer = 'viewer';
}

enum Status // Pure Enum
{
    case Active;
    case Inactive;
}

// --- シリアライズ対象データ構造の構築 (トップレベルがオブジェクト) ---

// オブジェクト参照と値参照のための共通データ
$commonObject = new MySimpleClass('共通オブジェクト', 500, ['shared' => true]);
$commonString = "共有される文字列データ";

// トップレベルのオブジェクトを作成し、多様なデータをプロパティとして持つ
$topLevelObject = new MySimpleClass('Top Level Object', 999, [
    'assoc_key1' => 'assoc_value1',
    'assoc_key2' => 789,
    'deep_nested_array' => [
        'sub_key_x' => 'sub_value_x',
        'sub_key_y' => 12.34,
    ],
]);

// 1. スカラー型 (Scalar Types)
$topLevelObject->japaneseString = "これは日本語の文字列です";
$topLevelObject->integerValue = 123;
$topLevelObject->floatValue = 45.67;
$topLevelObject->booleanTrue = true;
$topLevelObject->nullValue = null;

// 2. リスト (数値キーの配列)
$topLevelObject->nestedArrayData = [
    'Item A',
    'Item B',
    'Item C',
    10,
    false,
];

// 3. 連想配列 (Associative Array)
// $topLevelObject->protectedProp = [ // protectedProp を再利用して連想配列を入れる
//     'assoc_key1' => 'assoc_value1',
//     'assoc_key2' => 789,
//     'deep_nested_array' => [
//         'sub_key_x' => 'sub_value_x',
//         'sub_key_y' => 12.34,
//     ],
// ];

// 4. オブジェクト参照 (r:) - 別のオブジェクトへの参照
$topLevelObject->sharedObjectRef1 = $commonObject;
$topLevelObject->sharedObjectRef2 = $commonObject; // 同じオブジェクトへの参照

// 5. 値参照 (R:) - 同じ文字列への参照
$topLevelObject->sharedStringRef1 = $commonString;
$topLevelObject->sharedStringRef2 = $commonString; // 同じ文字列への参照

// 6. Enum (PHP 8.1+ の場合)
$topLevelObject->userRoleEnum = UserRole::Editor;
$topLevelObject->statusEnum = Status::Active;

// 7. カスタムシリアライズ可能なオブジェクト (C:)
$topLevelObject->customSerializableObj = new MyCustomSerializable('オブジェクト内のカスタム', 777);


// --- シリアライズ実行 ---
$serializedResult = serialize($topLevelObject);

// --- 出力 ---
echo "--- 元データ (var_export) ---\n";
var_export($topLevelObject);
echo "\n\n";

echo "--- シリアライズ結果 ---\n";
echo $serializedResult;
echo "\n\n";


?>
*/
