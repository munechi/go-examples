package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ハンドラー関数
func handler(c *gin.Context) {
	// JSON 形式でレスポンスを返す
	c.JSON(http.StatusOK, gin.H{"message": "Hello world!"})
}

func main() {
	// 初期化
	r := gin.Default()

	// エンドポイント定義
	r.GET("/hello", handler)

	// サーバ起動
	log.Println("server started at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
