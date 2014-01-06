package main

import (
	"fmt"
	"runtime"
	"time"

	"common/go/worker"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

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

var DB gorm.DB

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	DB, err = gorm.Open("mysql", "root@/go_php")
	if err != nil {
		panic(fmt.Sprintf("Got error when connect database, the error is '%v'", err))
	}

	f := func(workerId int, req *worker.Request) (res worker.Response, resTube string, err error) {
		if req.Action != "lookup" {
			return
		}

		d := req.Data.(map[string]interface{})

		u := lookup(d["user"])
		if (*u).Id > 0 {
			fmt.Printf("[%d:%d] FOUND: %v\n", workerId, time.Now().UnixNano(), *u)
		} else {
			fmt.Printf("[%d:%d] USER NOT FOUND\n", workerId, time.Now().UnixNano())
		}
		return
	}

	worker.Run("after_request", f)
}
