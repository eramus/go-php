package main

import (
	"fmt"
	"time"

	"common/go/worker"
)

func main() {
	f := func(workerId int, req *worker.Request) (res worker.Response, resTube string, err error) {
		fmt.Printf("[%d:%d] PAGE DATA: %v\n", workerId, time.Now().UnixNano(), req)
		return
	}

	var i = 0
	deadWorker := make(chan bool)
	for i = 0; i < 5; i++ {
		go worker.Run(i, "pages", f, deadWorker)
	}
	for {
		<-deadWorker
		i++
		go worker.Run(i, "pages", f, deadWorker)
	}
}
