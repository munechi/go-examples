# log/slog

- log と比べて slog では構造化されたログ出力ができる。
- ログ解析ツールとの相性が良い。
- 速度面では zap/zerolog より遅いが高負荷環境でなければ十分。

## 準備

```bash
go mod init go-slog
```

## プログラムコード

### SetDefault() を使う方法

`SetDefault-echo/main.go`

最初にログの初期化を行う。

```go
func main()
	// ログ出力フォルダを作成
	if err := os.MkdirAll("./log", 0755); err != nil {
		panic("failed create log directory")
	}

	// ログファイルを開く
	f, err := os.OpenFile(
		"./log/app.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// slog 設定(JSON形式)
	logger := slog.New(
		slog.NewJSONHandler(f, nil),
	)
	// slog のデフォルト出力先を設定
	slog.SetDefault(logger)

	// (以下省略)
```

ハンドラー関数内にログ出力処理を記述する。

```go
func handler(c echo.Context) error {
	// ログ出力
	slog.Info(
		"request",
		"path",
		c.Request().URL.Path,
	)
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello world!"})
}
```


`slog.SetDefault()` でアプリ内における slog の出力先を設定できる。

### 依存性注入(DI)な方法

ログ出力フォルダ作成、ログファイルオープン、slog 設定は前述の `SetDefault()` と同じなので説明省略。

ハンドラー生成関数を作る。return には無名関数を渡す（戻り値の型は `echo.HandlerFunc`）

```go
func handler(logger *slog.Logger) echo.HandlerFunc {
	// `echo.HandleFunc` は `func(c echo.Context) error` の別名
	return func(c echo.Context) error { // 無名関数
		// ログ出力
		logger.Info(
			"request",
			"path",
			c.Request().URL.Path,
		)
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello world!"})
	}
}
```

ハンドラー生成関数の引数 slog のオブジェクトを渡す。

```go
	e.GET("/hello", handler(logger))
```

## 実行

```bash ln=false
go run SetDefault-echo/main.go
```

または

```bash ln=false
go run DI-echo/main.go
```

これを何回か実行する。

### ログの出力例

`app/app.log`

```text ln=false
{"time":"2026-05-11T11:51:05.72302+09:00","level":"INFO","msg":"request","path":"/hello"}
{"time":"2026-05-11T11:51:13.424893+09:00","level":"INFO","msg":"request","path":"/hello"}
{"time":"2026-05-11T11:51:30.7356912+09:00","level":"INFO","msg":"request","path":"/hello"}
```
