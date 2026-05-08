package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ハンドラー関数(戻り値の型を `error` で定義するのが echo の流儀)
func handler(c echo.Context) error {
	// JSON 形式でレスポンスを返す
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello world!"})
}

func main() {
	// 初期化
	e := echo.New()

	// エンドポイント定義
	e.GET("/hello", handler)

	// サーバ起動
	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal(err)
	}
}
