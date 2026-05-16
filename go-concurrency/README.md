# concurrency

並行処理に関する部分（並列処理ではない）

Go には以下の様な仕組みが備わっていいます。

```text ln=false
goroutine   = 処理を並行に動かす
channel     = goroutine 間で値を渡す
buffered    = channel に貯められる
select      = 複数 channel を待つ
worker pool = 並行処理数を制御する
mutex       = 共有データを守る
```

## goroutine

こういう書式になります。

```go ln=false
go func() {
    // 別 goroutine で動く処理
}()
```

無名関数の定義の後に `()` が付いているとこれは即時に実行という意味になります。

先頭の `go` が goroutine が並行処理として動きます。

以下のコード例では、

- 1秒後に goroutine-2 が終了
- 2秒後に goroutine-1 が終了
- 3秒後に main が終了

という動きになります。

`goroutine/main.go`

```go
func main() {
	log.Print("main start")

	go func() {
		log.Print("goroutine-1 start")
		time.Sleep(2 * time.Second)
		log.Print("goroutine-1 end")
	}()

	go func() {
		log.Print("goroutine-2 start")
		time.Sleep(1 * time.Second)
		log.Print("goroutine-2 end")
	}()

	time.Sleep(3 * time.Second)
	log.Print("main end")
}
```

## channel

goroutine の無名関数は戻り値を返せません。

goroutine 内から結果を返すには channel という仕組みを使います。

`channel/main.go`

```go
func main() {
	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "done"
	}()

    msg := <-ch
    log.Print(msg)
}
```

`make(chan string)` で string 型の channel を作ります。

goroutine は処理が終わったら channel に結果を渡します。

`msg := <-ch` は channel にデータが来るのをずーっと待ち続けます。

## buffered channel

`buffered_channel/main.go`

ます下記コードについて、

```go
func main() {
	ch := make(chan string)

    ch <- "done"

    msg := <-ch
    log.Print(msg)
}
```

これを実行すると、

```text ln=false
fatal error: all goroutines are asleep - deadlock!
```

になります。理由は `ch <- "done"` 実行時に受け先 `msg := <-ch` がまだ存在しないので送信処理がブロックします。
その結果、main goroutine が次の行へ進めず、全 goroutine が停止状態になって deadlock になります。

そこで `make()` にバッファ数を記述して buffered channel にします。

```go
func main() {
	ch := make(chan string, 1)

    ch <- "done"

    msg := <-ch
    log.Print(msg)
}
```

そうすると動作します。 理由は buffered channel にしたことで、`ch<-"done"` 実行時にまだ `msg := <-ch` に到達していなくても、
いったんバッファに溜まって `msg := <-ch` に到達したときにバッファから取り出されて実行されるためです。

ちなみに、

```go
func main() {
	ch := make(chan string)

    go func(){
        ch <- "done"
    }()

    msg := <-ch
    log.Print(msg)
}
```

unbuffered channel でも動作するのは `ch <- "done"` を 別の goroutine の中に入れたことで、送信待ちで止まるのは別の goroutine ということにになります。
main goroutine はそのまま `msg := <-ch` まで進むため、送信側と受信側が揃って値の受け渡しが行われます。

## select

複数の channel を待つのが select 文です。

`select/main.go`

```go
func main() {
	ch := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch <- "done"
	}()

	select {
	case msg := <-ch:
		log.Print(msg)

	case <-time.After(2 * time.Second):
		log.Print("timeout")
	}
}
```

select 内の case で channel データが来るのをまちます。

コード例では 1 秒後に goroutine の処理が終了して channel に "done" を返します。

すると select で ch にデータが来たのに反応して ch のデータを msg に代入し "done" と表示されます。

もし goroutine のスリープ時間を 3 秒にすると `time.After(2 * time.Second)` が先にデータを返して "timeout" と表示されます。

## worker pool

worker pool とは Go での並行処理の設計方法における同時実行数を制限することです。

goroutine と buffered channel を使い buffer channel に溜まったジョブを空いている worker が処理していきます。

コード例：

```go
func worker(id int, jobs <-chan int) {
	for job := range jobs {
		log.Printf("worker %d start job %d\n", id, job)
		time.Sleep(1 * time.Second)
		log.Printf("worker %d finish job %d\n", id, job)
	}

}

func main() {
	jobs := make(chan int, 10)

	// worker を 3 つ起動する
	for i := 1; i <= 3; i++ {
		go worker(i, jobs)
	}

	// job 1～5 投入
	for j := 1; j <= 5; j++ {
		jobs <- j
	}

	// buffer が空になったら終了(close)する
	close(jobs)

	// worker が終わるまでの待ち（example向けの簡易手段）
	time.Sleep(3 * time.Second)
}
```

`worker` 関数の for ループ文は、一見、`jobs` 内のデータをループ処理しているようにみえるけど、
goroutine の待機と実行を繰り返すといった複雑な処理を裏で行っています（どちらかというとイベントループっぽい）

ジョブの投入も `jobs <- j` とあまりにもシンプルに書けるので実は裏で複雑な処理をしている感じがしないですね。

このコード例ではシンプルにするために最後に `time.Sleep()` で時間を指定して main goroutine の終了を遅らせていますが、
実務で worker pool の終了を待つときは `sync.WaitGroup` を使用します。

## mutex

複数の goroutine で一つの変数を更新したいとき、同時に更新して値が壊れないようにするために変数を保護する仕組みです。

```go
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
```

このコードでは 1000 個の goroutine を起動しつつ `count` をインクリメントしています。

出力される goroutine 番号を見ると、必ずしも 1 から順に実行されていないことがわかります。
ここは Go runtime の scheduler 次第になります。

なお、無名関数は外側の変数を参照できます（便利だが危うい面もある）
