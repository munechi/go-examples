package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func handler(c *gin.Context) {
	// フォームデータ
	name := c.PostForm("name")

	// ファイル
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	// 名前の前後の連続した空白文字を除去
	name = strings.TrimSpace(name)

	// アップロードファイルファイルに日付を付加
	base := filepath.Base(file.Filename)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s%s", stem, timestamp, ext)

	// ファイルを保存する
	err = c.SaveUploadedFile(file, "./"+filename)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "save failed"},
		)
		return
	}

	// JSON 形式のレスポンスデータを返す
	c.JSON(
		http.StatusOK,
		gin.H{"message": "upload successful", "name": name},
	)
}

func main() {
	// 初期化
	r := gin.Default()

	// エンドポイント定義
	r.POST("/upload", handler)

	// サーバ起動
	log.Println("server started at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
