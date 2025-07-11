# 🎮 PS2

[![GitHub release](https://img.shields.io/github/release/ddddddO/ps2.svg?label=Release&color=darkcyan)](https://github.com/ddddddO/ps2/releases)
[![codecov](https://codecov.io/gh/ddddddO/ps2/graph/badge.svg?token=6E0G81K2H0)](https://codecov.io/gh/ddddddO/ps2)
[![ci](https://github.com/ddddddO/ps2/actions/workflows/ci.yaml/badge.svg)](https://github.com/ddddddO/ps2/actions/workflows/ci.yaml)


**phP Serialize() To {JSON|YAML|TOML}**

Web and CLI tool are available!</br>

> [!WARNING]
> Use it for easy human checking.

Converts data serialized by the PHP serialize function into some format (JSON or YAML or TOML for now).</br>
This may be useful if you want to convert the payload string stored in the Laravel job queue 👍 </br>
If you find a bug or have a request, please create issue or Pull Request!

Web 👉 https://ddddddo.github.io/ps2/ </br>
Conversion can be done offline as it is processed in Wasm.

## Supported Types
|-|Type|Note|
|--|--|--|
|⭕|`a:5;`||
|⭕|`s:4:"xxxx";`||
|⭕|`E:4:"xxxx";`||
|⭕|`i:555`||
|⭕|`d:3.14`||
|⭕|`b:0` / `b:1`||
|⭕|`N`||
|⭕|`r:1` / `R:1`| filled with `[[PHP_REFERENCE_DATA: <actual referenced unserialized data>]]`<br>or `[[MAYBE_PHP_SELF_REFERENCE]]` |
|⭕|`O:3:"Obj":2:{...}`||
|⭕|`C:3:"Csm":2:{...}`||

## Install

### Homebrew

```console
brew install ddddddO/tap/ps2
```

### Go
```console
go install github.com/ddddddO/ps2/cmd/ps2@latest
```

## Usage

Format:
```console
$ ps2 <<< '< Data serialized by PHP serialize function >'
```

- Use `<<<` or `printf` to pass serialized data to stdin, since using echo as in `echo 'O:13:"App\UpdateJob....' | ps2` may cut off the `\U` part.

Example serialized data:
```console
O:13:"MySimpleClass":17:{s:10:"publicProp";s:16:"Top Level Object";s:26:"MySimpleClassprivateProp";i:999;s:16:"*protectedProp";a:3:{s:10:"assoc_key1";s:12:"assoc_value1";s:10:"assoc_key2";i:789;s:17:"deep_nested_array";a:2:{s:9:"sub_key_x";s:11:"sub_value_x";s:9:"sub_key_y";d:12.34;}}s:6:"parent";N;s:15:"nestedArrayData";a:5:{i:0;s:6:"Item A";i:1;s:6:"Item B";i:2;s:6:"Item C";i:3;i:10;i:4;b:0;}s:16:"sharedStringRef1";s:33:"共有される文字列データ";s:16:"sharedStringRef2";R:17;s:16:"sharedObjectRef1";O:13:"MySimpleClass":6:{s:10:"publicProp";s:24:"共通オブジェクト";s:26:"MySimpleClassprivateProp";i:500;s:16:"*protectedProp";a:1:{s:6:"shared";b:1;}s:6:"parent";N;s:15:"nestedArrayData";a:0:{}s:9:"nullValue";N;}s:16:"sharedObjectRef2";r:18;s:21:"customSerializableObj";O:20:"MyCustomSerializable":2:{s:1:"s";s:36:"オブジェクト内のカスタム";s:1:"n";i:777;}s:12:"userRoleEnum";E:15:"UserRole:Editor";s:10:"statusEnum";E:13:"Status:Active";s:9:"nullValue";N;s:11:"booleanTrue";b:1;s:10:"floatValue";d:45.67;s:12:"integerValue";i:123;s:14:"japaneseString";s:36:"これは日本語の文字列です";}
```

Example execution:

```console
$ ps2 <<< 'O:13:"MySimpleClass":17:{s:10:"publicProp";s:16:"Top Level Object";s:26:"MySimpleClassprivateProp";i:999;s:16:"*protectedProp";a:3:{s:10:"assoc_key1";s:12:"assoc_value1";s:10:"assoc_key2";i:789;s:17:"deep_nested_array";a:2:{s:9:"sub_key_x";s:11:"sub_value_x";s:9:"sub_key_y";d:12.34;}}s:6:"parent";N;s:15:"nestedArrayData";a:5:{i:0;s:6:"Item A";i:1;s:6:"Item B";i:2;s:6:"Item C";i:3;i:10;i:4;b:0;}s:16:"sharedStringRef1";s:33:"共有される文字列データ";s:16:"sharedStringRef2";R:17;s:16:"sharedObjectRef1";O:13:"MySimpleClass":6:{s:10:"publicProp";s:24:"共通オブジェクト";s:26:"MySimpleClassprivateProp";i:500;s:16:"*protectedProp";a:1:{s:6:"shared";b:1;}s:6:"parent";N;s:15:"nestedArrayData";a:0:{}s:9:"nullValue";N;}s:16:"sharedObjectRef2";r:18;s:21:"customSerializableObj";O:20:"MyCustomSerializable":2:{s:1:"s";s:36:"オブジェクト内のカスタム";s:1:"n";i:777;}s:12:"userRoleEnum";E:15:"UserRole:Editor";s:10:"statusEnum";E:13:"Status:Active";s:9:"nullValue";N;s:11:"booleanTrue";b:1;s:10:"floatValue";d:45.67;s:12:"integerValue";i:123;s:14:"japaneseString";s:36:"これは日本語の文字列です";}'
{
  "*protectedProp": {
    "assoc_key1": "assoc_value1",
    "assoc_key2": 789,
    "deep_nested_array": {
      "sub_key_x": "sub_value_x",
      "sub_key_y": 12.34
    }
  },
  "MySimpleClassprivateProp": 999,
  "__class_name": "MySimpleClass",
  "booleanTrue": true,
  "customSerializableObj": {
    "__class_name": "MyCustomSerializable",
    "n": 777,
    "s": "オブジェクト内のカスタム"
  },
  "floatValue": 45.67,
  "integerValue": 123,
  "japaneseString": "これは日本語の文字列です",
  "nestedArrayData": [
    "Item A",
    "Item B",
    "Item C",
    10,
    false
  ],
  "nullValue": null,
  "parent": null,
  "publicProp": "Top Level Object",
  "sharedObjectRef1": {
    "*protectedProp": {
      "shared": true
    },
    "MySimpleClassprivateProp": 500,
    "__class_name": "MySimpleClass",
    "nestedArrayData": [],
    "nullValue": null,
    "parent": null,
    "publicProp": "共通オブジェクト"
  },
  "sharedObjectRef2": "[[PHP_REFERENCE_DATA: map[*protectedProp:map[shared:true] MySimpleClassprivateProp:500 __class_name:MySimpleClass nestedArrayData:map[] nullValue:<nil> parent:<nil> publicProp:共通オブジェクト]]]",
  "sharedStringRef1": "共有される文字列データ",
  "sharedStringRef2": "[[PHP_REFERENCE_DATA: 共有される文字列データ]]",
  "statusEnum": "Status:Active",
  "userRoleEnum": "UserRole:Editor"
}
$
```

### Output Types
- `--json`
  - Default. Convert to JSON.
- `--yaml`
  - Convert to YAML.
- `--toml`
  - Convert to TOML.

## Related tool
- **[PHP Serialized Object viewer](https://github.com/haradakunihiko/intellij-plugin-php-serialized-object-viewer)**
    - If you are using IntelliJ IDEA / PhpStorm, please try this plugin!

## Memo

- The first version of this tool was generated by `Gemini 2.5 Flash`.
    ```
    phpのserialize関数でシリアライズされた文字列（元のデータは、ネストした連想配列と、配列と、適当なクラスのオブジェクトを持つ）から、astを生成するgoプログラムを教えてください。
    ```

    - 一発で動かない & 微調整の必要があったため手を入れている
      - 継続して手を入れ続けている
    - 生成されたコードは「Go言語でつくるインタプリタ」の内容にかなり似ている

- PHP serialize()
    - https://www.php.net/manual/ja/function.serialize.php

- This tool has been verified to work with serialized data generated by PHP 7.4.33.

- [Reddit: r/PHP](https://www.reddit.com/r/PHP/comments/1l61qw7/github_ddddddops2_tool_to_convert_from_serialized/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button)