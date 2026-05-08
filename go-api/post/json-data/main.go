package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// リクエストの構造体
type Request struct {
	Name string `json:"name"`
}

// レスポンスの構造体
type Response struct {
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// HTTP メソッドの確認
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	}

	// リクエストヘッダの確認
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	// JSON データのデコード
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	// レスポンスデータ生成
	resp := Response{
		Message: "Hello " + req.Name + "!",
	}

	// レスポンスヘッダ設定
	w.Header().Set("Content-Type", "application/json")

	// JSON 形式にエンコードしてレスポンスを返す
	json.NewEncoder(w).Encode(resp)
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
