package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// レスポンスの構造体
type Message struct {
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// レスポンスヘッダ設定
	w.Header().Set("Content-Type", "application/json")

	// クエリデータ取得
	name := r.URL.Query().Get("name")

	// レスポンスデータ生成
	resp := Message{
		Message: fmt.Sprintf("Hello %s!", name), //  ケツカンマ必須よ！
	}

	// JSON 形式にエンコードしてレスポンスを返す
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		//　エラー処理
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func main() {
	// エンドポイント定義
	http.HandleFunc("/hello", handler)

	// サーバ起動
	log.Println("server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
