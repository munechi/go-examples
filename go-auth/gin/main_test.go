package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var ginSecret = []byte("test-secret")

type GinClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func createGinToken(userID int) (string, error) {
	claims := GinClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ginSecret)
}

// 認証のミドルウェア
func ginAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		if tokenString == auth {
			c.String(http.StatusUnauthorized, "missing bearer token")
			c.Abort()
			return
		}

		claims := &GinClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return ginSecret, nil
		})
		if err != nil || !token.Valid {
			c.String(http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

func TestGinJWTMiddlewareOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()             // 初期化
	r.Use(ginAuthMiddleware()) // ミドルウェア登録

	// エンドポイント
	r.GET("/private", func(c *gin.Context) {
		userID := c.GetInt("user_id")
		c.String(http.StatusOK, fmt.Sprintf("user_id: %d", userID))
	})

	// userID = 1 のトークン作成
	token, err := createGinToken(1)
	assert.NoError(t, err)

	// テスト用リクエスト作成
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// テスト用レスポンス格納エリア作成
	rec := httptest.NewRecorder()

	// ハンドラ chain
	r.ServeHTTP(rec, req)

	// リターンコード比較
	assert.Equal(t, http.StatusOK, rec.Code)
	// レスポンス比較
	assert.Equal(t, "user_id: 1", rec.Body.String())
}

func TestGinJWTMiddlewareNG(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(ginAuthMiddleware())

	r.GET("/private", func(c *gin.Context) {
		c.String(http.StatusOK, "secret")
	})

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	// リクエストヘッダにトークンが無いので 401
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
