package main

import (
	"log"
	"time"
)

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
