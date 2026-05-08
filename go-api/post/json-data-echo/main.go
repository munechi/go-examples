package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// リクエストデータの構造体 (JSON キー名称あり)
type Request struct {
	Name string `json:"name"`
}

// レスポンスデータの構造体 (JSON キー名称あり)
type Message struct {
	Message string `json:"name"`
}

func handler(c echo.Context) error {
	// 変数メモリ確保
	var req Request

	// JSON データをデコード(`Bind()` には `req` のアドレスを渡す)
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "Bad Request"},
		)
	}

	// レスポンスデータ生成
	resp := Message{
		Message: "Hello " + req.Name + "!",
	}

	// レスポンスを返す
	return c.JSON(http.StatusOK, resp)
}

func main() {
	// 初期化
	e := echo.New()

	// エンドポイント定義
	e.POST("/hello", handler)

	// サーバ起動
	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal(err)
	}
}
