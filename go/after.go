package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kr/beanstalk"
)

type Request struct {
	RequestId uint64                 `json:"request"`
	Action    string                 `json:"action"`
	Data      map[string]interface{} `json:"data"`
}

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func lookup(userId interface{}) *User {
	var user User
	DB.Where("id = ?", userId).First(&user)
	return &user
}

func worker(workerId int, deadWorker chan<- bool) {
	beanConn, err := beanstalk.Dial("tcp", "0.0.0.0:11300")
	if err != nil {
		fmt.Println("BEANSTALK:", err)
		return
	}
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("PANIC:", r)
			deadWorker <- true
		}
		beanConn.Close()
	}()

	var req Request
	var watch = beanstalk.NewTubeSet(beanConn, "after_request")

	for {
		id, msg, err := watch.Reserve(10 * time.Second)
		if err != nil {
			cerr, ok := err.(beanstalk.ConnError)
			if ok && cerr.Err == beanstalk.ErrTimeout {
				continue
			} else {
				fmt.Println("SOMETHING BAD HAPPENED TO BEANSTALK")
				panic("conn err")
			}
		}
		fmt.Printf("[%d:%d] START JOB\n", workerId, time.Now().UnixNano())

		err = json.Unmarshal(msg, &req)
		beanConn.Delete(id)

		if err != nil {
			fmt.Println("JSON:", err)
			panic("json err")
		}
		fmt.Printf("[%d:%d] DATA: %v\n", workerId, time.Now().UnixNano(), req)

		if req.Action != "lookup" {
			continue
		}

		u := lookup(req.Data["user"])
		if (*u).Id > 0 {
			fmt.Printf("[%d:%d] FOUND: %v\n", workerId, time.Now().UnixNano(), *u)
		} else {
			fmt.Printf("[%d:%d] USER NOT FOUND\n", workerId, time.Now().UnixNano())
		}
	}
}

var DB gorm.DB

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	DB, err = gorm.Open("mysql", "root@/go_php")
	if err != nil {
		panic(fmt.Sprintf("Got error when connect database, the error is '%v'", err))
	}

	var i = 0
	deadWorker := make(chan bool)
	for i = 0; i < 5; i++ {
		go worker(i, deadWorker)
	}
	for {
		<-deadWorker
		i++
		go worker(i, deadWorker)
	}
}
