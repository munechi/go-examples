# json

内容：

- Marshal
- Unmarshal
- tag
- pointer
- RawMessage

## Marshal

Go の struct → JSON 文字列の変換

### 基本形

```go ln=false
data, err := json.Marshal(v)
```

### コード例

```go
package main

import (
	"encoding/json"
	"fmt"
)

// JSON に対応する構造体
type User struct {
	Name string
	Age  int
}

func main() {
	// 変換データの定義
	user := User{
		Name: "Taro",
		Age:  20,
	}

	// 変換
	data, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	// 変換結果出力
	fmt.Println(string(data)) // Marshal の戻り値は []byte
}
```

実行：

```bash ln=false
go run marshal/main.go
```

出力例：

```text ln=false
{"Name":"Taro","Age":20}
```

### strict のフィールドが大文字始まり

Go では

- 大文字始まり ... export (外部公開)
- 小文字始まり ... private

というルールがあります。`encoding/json` は外部ライブラリなので大文字始まりにしないと認識されません。

### マップ

```go ln=false
	user := map[string]string{
		"Name": "Jiro",
		"Age":  "19",
	}
	data, err := json.Marshal(user)
```

### スライス

```go ln=false
	user := []User{
		{Name: "Saburo", Age: 18},
		{Name: "Shiro", Age: 17},
	}
	data, err := json.MarshalIndent(user, "", "  ")
```

インデント整形された JSON が欲しい場合は `MarshalIndent()` を使います。

### 実務では

`Marshal()` よりも `json.NewEncoder(w).Encode(user)` を使うケースが多いようです。

理由は、

- 直接 io.Writer に書ける
- メモリ効率の良さ
- HTTP レスポンスとの相性の良さ

## Unmarshal

Go の JSON 文字列 → struct の変換

### 基本形

```go ln=false
err := json.Unmarshal(data, &user)
```

### コード例

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// 変換用の構造体
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
    // 変換する JSON
	data := []byte(`{"name":"Taro","age":20}`)

    // メモリ確保
	var user User

    // 変換
	err := json.Unmarshal(data, &user) // user はアドレスを渡す！
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(user.Name)
	fmt.Println(user.Age)
}
```

実行：

```bash ln=false
go run unmarshal/main.go
```

出力例：

```text ln=false
Taro
20
```

### スライス＋マップ

```go ln=false
	var slice_user []map[string]any
	err = json.Unmarshal(
		[]byte(`[{"name":"Saburo","age":18},{"name":"Shiro","age":17}]`),
		&slice_user,
	)
	fmt.Println(slice_user[0]["name"], slice_user[0]["age"])
	fmt.Println(slice_user[1]["name"], slice_user[1]["age"])
```

### any

どうしても構造が不明な JSON を受けるときに `any` 使用する場合があります。

## タグ

既に struct の定義で登場している、

```go ln=false
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
```

この

```text ln=false
`json:"name"`
`json:"age"`
```

がタグです。

これは単なる文字列であり struct のフィールドにあるメタ情報です。
Go の `reflect` が読み取って `encoding/json` が解釈しています。

### omitempty

空値だった場合は出力しないようにする定義です。

```go ln=false
type User struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty`
}
```

```go ln=false
u := User{Name: "Taro"}
data, err := json.Marshal(u) // data は `{"name":"Taro"}`
```

### 無視

キー名を書かず `"-"` と書きます。

```go ln=false
type User struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"-"`
}
```

```go ln=false
u := User{Name: "Taro", Age: 20}
data, err := json.Marshal(u) // data は `{"name":"Taro"}`
```

## ポインター

**未入力**と**ゼロ値**の区別をしたいときに使用します（入力データのバリデーションなど）

例：

年齢データの int 型とした場合、ゼロ値が 0 になります。

送られてきたデータ JSON データが `{}` でも `{"age":0}` でもプログラム内では `age=0` となり、
「未入力」なのか「年齢 0 歳のデータ」かの判定ができなくなります。


### 通常

```go ln=false
type User0 struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
```

この場合、未入力の判定を、

```go ln=false
	if user1.Age == 0 {
        // エラー処理
    }
```

このように書くと `age=0` のデータもエラーになります（誤判定）

### age を ポインターにする

```go ln=false
type User0 struct {
	Name string `json:"name"`
	Age  *int   `json:"age"`
}
```

こうすれば、未入力の判定を、

```go ln=false
	if user1.Age == nil {
        // エラー処理
    }
```

とすることができて、本当に age データが無いときだけエラーにすることができます。

### Go での JSON 設計

必須項目：

```go ln=false
Name string `json:"name"
```

オプション項目：

```go ln=false
Age *int `json:"age,omitempty"
```

この使い分けが多いようです。

## RawMessage

受信した JSON の一部を生データとして保持しておくための型です。

### 使いどころ

例えば `type` キーの値によってデータの構造が変わるようなケースで使用します。

```json
{
    "type": ＜タイプ＞,
    "data": ＜タイプによってデータ構造が変わる＞
}
```

普通に struct で受けることができない！

### 使用例

イベントの struct

```go ln=false
type Hoge struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
```

`Type` の違いによるに struct

```go ln=false
// type: "foo"
type Foo struct {
	Name       string `json:"name"`
	Age        *int   `json:"age"`
}

// type: "bar"
type Bar struct {
	Address    string `json:"address"`
	Email      string `json:"email"`
}
```

分岐処理

```go ln=false
var d Hoge
err := json.Unmarshal(src, &d)

switch d.Type {
case "foo":
	var f Foo
	err := json.Unmarshal(d.Data, &f)
case "bar":
	var b Bar
	err := json.Unmarshal(d.Data, &b)
default:
	// 未知の Type が来た時の処理
}
```
