package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// レスポンスの構造体
type Message struct {
	Message string `json:"message"`
}

// ハンドラ関数
func handler(w http.ResponseWriter, r *http.Request) {
	// レスポンスヘッダ設定
	w.Header().Set("Content-Type", "application/json")

	// レスポンスデータ生成
	resp := Message{
		Message: "Hello world!",
	}

	// JSON 形式にエンコードしてレスポンスを返す
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// エラー処理
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func main() {
	// 初期化
	r := chi.NewRouter()

	// エンドポイント定義
	r.Get("/hello", handler)

	// サーバ起動
	log.Println("server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
