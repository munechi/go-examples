package main

import (
	"log"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "done"
	}()

	msg := <-ch
	log.Print(msg)
}
