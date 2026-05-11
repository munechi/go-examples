package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func handler(c echo.Context) error {
	// ログ出力
	slog.Info(
		"request",
		"path",
		c.Request().URL.Path,
	)
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello world!"})
}

func main() {
	// ログ出力フォルダを作成
	if err := os.MkdirAll("./log", 0750); err != nil {
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

	e := echo.New()
	e.GET("/hello", handler)
	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal(err)
	}
}
