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

	worker.Run("pages", f)
}
