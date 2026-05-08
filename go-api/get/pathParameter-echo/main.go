package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func handler(c echo.Context) error {
	// パスパラメータ取得
	name := c.Param("name")

	// JSON 形式でレスポンスを返す
	return c.JSON(
		http.StatusOK,
		map[string]string{"message": "Hello " + name + "!"},
	)
}

func main() {
	// 初期化
	e := echo.New()

	// エンドポイント定義
	e.GET("/hello/:name", handler)

	// サーバ起動
	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal(err)
	}
}
