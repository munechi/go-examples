package main

import (
	"fmt"
	"log"
	"sync"
)

func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup

	count := 0

	for i := 0; i < 1000; i++ {
		wg.Add(1) // wg に一つ追加

		go func(i int) {
			defer wg.Done()          // 関数を抜けたら wg から一つ削除
			fmt.Printf("wg=%d\n", i) // goroutine の番号を出力

			mu.Lock()
			count++ // 他の goroutine と同時にアクセスしないように mutex で保護
			mu.Unlock()
		}(i) // 無名関数に goroutine 番号を渡す
	}

	wg.Wait() // wg が終了するまで待つ

	log.Print(count)
}
