# バリデーション

### 注目ポイント

リクエストデータの構造体のメンバー `Age` の型を `int` ではなく `*int` になっている。

```go ln=false
type Request struct {
	Name string `json:"name" binding:"required"`
	Age  *int   `json:"age" binding:"required,gte=0"`
}
```

【理由】

もし `age=0` の場合、送られてきたデータ `0` は `int` の初期値 `0` と同じ値になる。
そのため、`required` バリデーションではゼロ値として扱われ、未入力と判定されてバリデーションエラーになる。

`Age` を `*int` にすると、未入力時は `nil` となり、`age=0` が送信された場合は `0` を指すポインタになる。
そのため、`0` を未入力と区別でき、意図したとおり `age=0` を有効な値として扱える。

これは `Gin` が使っている `validator` の `required` の挙動によるものである。

## サーバー起動

```bash ln=false
go run validation-gin/main.go
```

## テスト

### (1) `name` が無い

```bash ln=false
curl --json '{"age":20}' http://localhost:8080/hello
```

#### 出力

```text ln=false
{"errors":[{"field":"Name","message":"必須項目です"}]}
```

### (2) `age` が無い

```bash ln=false
curl --json '{"name":"Taro"}' http://localhost:8080/hello
```

#### 出力

```text ln=false
{"errors":[{"field":"Age","message":"必須項目です"}]}
```

### (3) `name` も `age` も無い

```bash ln=false
curl --json '{}' http://localhost:8080/hello
```

#### 出力

```text ln=false
{"errors":[{"field":"Name","message":"必須項目です"},{"field":"Age","message":"必須項目です"}]}
```

### (4) `age` が負の数

```bash ln=false
curl --json '{"name":"Taro","age":-1}' http://localhost:8080/hello
```

#### 出力

```text ln=false
{"errors":[{"field":"Age","message":"0以上で入力してください"}]}
```
