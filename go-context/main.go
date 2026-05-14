package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	fmt.Println("start")

	// slowTask() に ctx を渡しす
	if err := slowTask(ctx); err != nil {
		fmt.Println("err:", err)
		http.Error(w, err.Error(), http.StatusGatewayTimeout)
	}
}

func slowTask(ctx context.Context) error {
	// 複数 channel 待ち受け
	select {
	case <-time.After(5 * time.Second): // 5秒後に発火
		fmt.Println("処理完了")
		return nil

	case <-ctx.Done(): // ctx にキャンセル通知が来たら発火
		fmt.Println("キャンセルされました:", ctx.Err())
		return ctx.Err() // ctx のキャンセル理由を返す
	}
}

func main() {
	http.HandleFunc("/slow", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
