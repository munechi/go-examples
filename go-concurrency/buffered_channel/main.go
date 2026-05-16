package main

import (
	"log"
)

func main() {
	// `make()` にバッファ数を記述すると buffered channel になる。
	ch := make(chan string, 1)

	ch <- "done-1"

	msg := <-ch
	log.Print(msg)
}
