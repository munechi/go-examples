package main

import (
	"log"
	"time"
)

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
