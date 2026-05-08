package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// リクエストデータの構造体 (バリデーション定義あり)
type Request struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required,gte=0"`
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "必須項目です"
	case "gte":
		return "0以上で入力してください"
	default:
		return "不正な値です"
	}
}

func handler(c *gin.Context) {
	var req Request

	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors

		if errors.As(err, &ve) {
			var errs []gin.H

			for _, fe := range ve {
				errs = append(
					errs,
					gin.H{
						"field":   fe.Field(),
						"message": validationMessage(fe),
					},
				)
			}

			c.JSON(
				http.StatusBadRequest,
				gin.H{"errors": errs},
			)
			return
		}

		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid request"},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "Hello " + req.Name + "!",
			"age":     req.Age,
		},
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
