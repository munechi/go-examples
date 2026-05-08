package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// リクエストデータの構造体 (JSON キー名称あり)
type Request struct {
	Name string `json:"name"`
}

func handler(c *gin.Context) {
	// 変数メモリ確保
	var req Request

	// JSON データを構造体データにバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	// レスポンスを返す
	c.JSON(
		http.StatusOK,
		gin.H{"message": "Hello, " + req.Name + "!"},
	)
}

func main() {
	// 初期化
	r := gin.Default()

	// エンドポイント定義
	r.POST("/hello", handler)

	// サーバ起動
	log.Println("server started at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
