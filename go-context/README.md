# context

context は「この処理をまだ続けていいかどうか」を伝えるものです。

例：

```text ln=false
ユーザーがブラウザを閉じた
 ↓
HTTP リクエストの context がキャンセルされる
 ↓
DB 問い合わせも止めたい
 ↓
外部 API 呼び出しも止めたい
 ↓
goroutine も止めたい
```

具体的には、

```text ln=false
`context.Context` を DB 問い合わせ関数の引数に渡す。
 ↓
ユーザーがブラウザを閉じたら、その情報が伝搬して DB 問い合わせ関数に届く。
 ↓
DB 問い合わせ関数は実行中の処理を終了する。
```

当然ながら、DB 問い合わせ関数が `context.Context` に対応(引数に取る)していなければなりません。

`go-db` のコード例で `context.Background()` を渡していましたが、
これは `context.Context` が引数として必須の関数に対してとりあえず渡すもので、
キャンセルや Timeout が発生しないすーっと生きている context です。

## 起点

HTTP リクエストなら `r.Context()`

```go ln=false
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    user, err := getUser(ctx)
}
```

重要なのは context は Immutable(変更不可) であるということです。
context に TimeOut の設定を加えたい場合、親から受け取った context にタイムアウトの設定を付加して、
子の context を作ります。そしてそれを子の関数に渡していきます。

なお、`r.Context() にはブラウザを閉じたなどの情報は伝搬しますが、タイムアウトに関する設定は入っていません。
そのため BD 問い合わせ時などには自分で `context.WithTimeout()` を指定する必要があります。

```go ln=false
func handler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(
        r.Context(),
        3*time.Second,
    )
    defer cancel()

    err := heavyDB(ctx)
}
```

## 動作確認

HTTP アクセスで、

- タイムアウト
- 強制終了（Web ブラウザ閉じ）

が発生したときの動作を確認する go のコードが以下になります。

コード例：

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	fmt.Println("start")

	// slowTask() に ctx を渡しす
	if err := slowTask(ctx); err != nil {
		fmt.Println("err:", err)
		http.Error(w, err.Error(), http.StatusGatewayTimeout)
	}
}

func slowTask(ctx context.Context) error {
	// 複数 channel 待ち受け
	select {
	case <-time.After(5 * time.Second): // 5秒後に発火
		fmt.Println("処理完了")
		return nil

	case <-ctx.Done(): // ctx にキャンセル通知が来たら発火
		fmt.Println("キャンセルされました:", ctx.Err())
		return ctx.Err() // ctx のキャンセル理由を返す
	}
}

func main() {
	http.HandleFunc("/slow", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
```

サーバ起動：

```bash ln=false
go run main.go
```

### タイムアウト

下記コマンドを実行して 2 秒待ちます。

すると、タイムアウトのメッセージが出力されます。

実行：

```bash ln=false
curl http://localhost:8080/slow
```

出力例：

```text ln=false
start
キャンセルされました: context deadline exceeded
err: context deadline exceeded
```

### ブラウザ閉じ

下記コマンドを実行して 2 秒以内に `Ctrl` + `c` を押します （Web ブラウザを閉じたのと同じ）

すると、キャンセルのメッセージが出力されます。

実行：

```bash ln=false
curl http://localhost:8080/slow
```

出力例：

```text ln=false
start
キャンセルされました: context canceled
err: context canceled
```
