package main

import (
	"log"
	"time"
)

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
