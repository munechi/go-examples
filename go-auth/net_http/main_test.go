package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// 秘密鍵(トークンをパースするときに必要)
var secret = []byte("test-secret")

// context のキーは string で直書き非推奨なので新しい型を作る
type ctxKey string

// 変更しないので定数にする
const userIDKey ctxKey = "user_id"

// ユーザID と JWT 標準の claims
type MyClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// レスポンス用 struct
type Response struct {
	UserID int `json:"user_id"`
}

// JWT トークン作成
func createToken(userID int) (string, error) {
	claims := MyClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			// 有効期限
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			// 発行日時
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	// HS256 アルゴリズムでトークン作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 秘密鍵で署名して返す
	// (トークンのパースでも秘密鍵を使うので「共通秘密鍵方式」になります)
	return token.SignedString(secret)
}

// 認証のミドルウェア
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ヘッダから生トークンを取り出す
		auth := r.Header.Get("Authorization")

		// トークン本体を取り出す
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		if tokenString == auth {
			// トークン先頭に "Bearer " が無かったらエラー
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		// クレームのメモリ確保
		claims := MyClaims{}

		// トークンのパース
		// ParseWithClaim() の第3引数には秘密キーを渡すが型は func である。
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
			return secret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// context に userID を加える
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		// 次へ
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ハンドラー関数
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	// context から userID を取り出す
	userID := r.Context().Value(userIDKey).(int)
	// plain text のレスポンスを返す
	w.Write([]byte(fmt.Sprintf("user_id: %d", userID)))
}

func TestNetHTTPJWTMiddlewareOK(t *testing.T) {
	// userID = 1 のトークン作成
	token, err := createToken(1)
	assert.NoError(t, err)

	// auth ミドルウェアを動かす
	h := authMiddleware(http.HandlerFunc(protectedHandler))

	// テスト用リクエスト作成
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// テスト用レスポンス格納エリア作成
	rec := httptest.NewRecorder()

	// ハンドラ chain
	h.ServeHTTP(rec, req)

	// リターンコード比較
	assert.Equal(t, http.StatusOK, rec.Code)
	// レスポンス比較
	assert.Equal(t, "user_id: 1", rec.Body.String())
}

func TestNetHTTPJWTMiddlewareNG(t *testing.T) {
	h := authMiddleware(http.HandlerFunc(protectedHandler))

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	// リクエストヘッダにトークンが無いので 401
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
