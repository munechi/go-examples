# auth

通常、API エンドポイントとハンドラの間に middleware として auth プログラムを挟んで運用します。

```text ln=false
request
  ↓
middleware
  ↓
handler
  ↓
middleware
  ↓
response
```

middleware として組み込まれるものとしては以下の様なものがります。

- 認証チェック
- ログ出力
- エラーハンドリング
- CORS
- リクエスト ID 付与
- タイムアウト

認証チェックであれば、認証後、ログイン済みであることのデータに JWT を利用します。

## JWT

JWT は JSON Web Token の略です。

JWT 方式のやりとり：

```text ln=false
ログイン成功
↓
JWT 発行
↓
クライアントが JWT 保持
↓
毎回ヘッダに JWT 送る
↓
サーバは署名検証だけ
```

※ https 通信ではヘッダも暗号化されているのでクライアントは JWT を header 付けて送ります。

### JWT に含まれる情報

- ヘッダ（暗号化方式、タイプ）
- Payload（ユーザ ID、名前、JWT の有効期限など。Base64 でエンコードされている）
- 署名

## middleware での検証

```go
func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		// JWT検証
		// OKなら次へ
		// NGなら401

		next.ServeHTTP(w, r)
	})
}
```

auto middleware はリクエストのヘッダから JWT を取り出してログイン済みか検証します。
NG だったら 401 を返します。

そして、JWT からユーザ ID などを取り出して、このモジュールへ context に載せて情報を伝搬します。

```go ln=false
ctx := context.WithValue(r.Context(), "user_id", 123)
```

そして `user_id` を読みときは、

```go ln=false
userID := r.Context().Value("user_id")
```

このようにします。

### middleware でハンドラーを包む

複数の middleware を利用する場合、その分だけ handler を包む必要があります。

```go ln=false
http.Handle(
	"/users",
	auth(logging(recovery(http.HandlerFunc(usersHandler)))),
)
```

> `http.HandleFunc()` ではなく `http.Handle()` で合っています。

このようにハンドラの定義が何層にも包まれ記述が長くなって地獄です。

なのでフレームワークではもっとコーディングしやすくなるように工夫されています。

#### Echo の例

```go ln=false
e.Use(middleware.Logger())
e.Use(middleware.Recover())
e.Use(authMiddleware)

e.GET("/users", usersHandler)
```

#### Gin の例

```go ln=false
r.Use(gin.Logger())
r.Use(gin.Recovery())
r.Use(authMiddleware)
```

また、router をグループ化して認証の middleware を付けることもできます。

このように、middleware を使用するシステムではフレームワークを使った方がすっきりとしたコードを書けます。

## コード例

### net/http

コードの場所：

`net_http/main_test.go`

実行：

```bash ln=false
go test -v net_http/main_test.go
```

出力例：

```text ln=false
=== RUN   TestNetHTTPJWTMiddlewareOK
--- PASS: TestNetHTTPJWTMiddlewareOK (0.00s)
=== RUN   TestNetHTTPJWTMiddlewareNG
--- PASS: TestNetHTTPJWTMiddlewareNG (0.00s)
PASS
ok      command-line-arguments  0.227s
```

### Gin

コードの場所：

`gin/main_test.go`

実行：

```bash ln=false
go test -v gin/main_test.go
```

出力例：

```text ln=false
=== RUN   TestGinJWTMiddlewareOK
--- PASS: TestGinJWTMiddlewareOK (0.00s)
=== RUN   TestGinJWTMiddlewareNG
--- PASS: TestGinJWTMiddlewareNG (0.00s)
PASS
ok      command-line-arguments  0.313s
```
