package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// レスポンスの構造体
type Response struct {
	Message string `json:"message"`
	Name    string `json:"name"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// フォームデータ
	name := r.FormValue("name")

	// ファイルサイズが 32MB までメモリ使用、それ以上は一時ファイル
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	// ファイル
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 名前の前後の連続した空白文字を除去
	name = strings.TrimSpace(name)

	// アップロードファイルファイルに日付を付加
	base := filepath.Base(handler.Filename)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s%s", stem, timestamp, ext)

	// ファイルの出力先を開く
	dst, err := os.Create("./" + filename)
	if err != nil {
		http.Error(w, "failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// ファイルデータを出力
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}

	// レスポンスヘッダ設定
	w.Header().Set("Content-Type", "application/json")

	// JSON 形式にエンコードしてレスポンスを返す
	json.NewEncoder(w).Encode(
		Response{
			Message: "upload successful",
			Name:    name,
		},
	)
}

func main() {
	// エンドポイント定義
	http.HandleFunc("/upload", handler)

	// サーバ起動
	log.Println("server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
