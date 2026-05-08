package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func handler(c echo.Context) error {
	// フォームデータ
	name := c.FormValue("name")

	// ファイル
	fileHandler, err := c.FormFile("file")
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "File is required"},
		)
	}

	// ファイルを開く
	src, err := fileHandler.Open()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "failed to open upload file"},
		)
	}
	defer src.Close()

	// 名前の前後の連続した空白文字を除去
	name = strings.TrimSpace(name)

	// アップロードファイルファイルに日付を付加
	base := filepath.Base(fileHandler.Filename)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s%s", stem, timestamp, ext)

	// ファイルの出力先を開く
	dst, err := os.Create("./" + filename)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "failed to create file"},
		)
	}
	defer dst.Close()

	// ファイルデータを出力
	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "failed to save file"},
		)
	}

	// JSON 形式のレスポンスデータを返す
	return c.JSON(
		http.StatusOK,
		map[string]string{"message": "upload successful", "name": name},
	)
}

func main() {
	// 初期化
	e := echo.New()

	// エンドポイント定義
	e.POST("/upload", handler)

	// サーバ起動
	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal(err)
	}
}
