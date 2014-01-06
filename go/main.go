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

func signup(workerId int, req *worker.Request) (res worker.Response, resTube string, err error) {
	fmt.Println("SIGNUP")
	d := req.Data.(map[string]interface{})

	var u User
	u.Name = d["name"].(string)
	u.Email = d["email"].(string)

	fmt.Println("save it", u, &u)
	DB.Save(&u)
	fmt.Println("saved")

	if (time.Now().UnixNano() % 2) == 1 {
		res.Success = true
	} else {
		res.Success = false
	}

	res.Data = u
	resTube = fmt.Sprintf("response_%d", req.RequestId)

	return
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
		go worker.Run(i, "default", signup, deadWorker)
	}
	for {
		<-deadWorker
		i++
		go worker.Run(i, "default", signup, deadWorker)
	}
}
