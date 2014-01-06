package main

import (
	"fmt"
	"runtime"
	"time"

	worker "go-php/go/beanstalk"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Request struct {
	RequestId uint64      `json:"request"`
	Action    string      `json:"action"`
	Data      interface{} `json:"data"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

/*func signup(data interface{}) (User, error) {
	fmt.Println("SIGNUP")
	d := data.(map[string]interface{})

	var u User
	u.Name = d["name"].(string)
	u.Email = d["email"].(string)

	fmt.Println("save it", u, &u)
	DB.Save(&u)
	fmt.Println("saved")
	return u, nil
}*/

func signup(workerId int, req *Request) (res Response, resTube string, err error) {
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

/*func worker(workerId int, deadWorker chan<- bool) {
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
	var res Response

	for {
		beanConn.Tube.Name = "default"
		id, msg, err := beanConn.Reserve(10 * time.Second)
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

		fmt.Println("SIGNUP")
		d := req.Data.(map[string]interface{})

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

		jsonRes, err := json.Marshal(res)
		if err != nil {
			fmt.Println("JSON:", err)
			panic("json err")
		}
		fmt.Printf("[%d:%d] %s\n", workerId, time.Now().UnixNano(), string(jsonRes))

		resTube := fmt.Sprintf("response_%d", req.RequestId)
		beanConn.Tube.Name = resTube

		_, err = beanConn.Put(jsonRes, 0, 0, (3600 * time.Second))
		if err != nil {
			fmt.Println("BEANSTALK WRITE:", err)
			panic("write err")
		}

		fmt.Printf("[%d:%d] FINISHED REQUEST\n", workerId, time.Now().UnixNano())
	}
}*/

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
		go worker.Run(i, signup, deadWorker)
		//		go worker(i, deadWorker)
	}
	for {
		<-deadWorker
		i++
		go worker.Run(i, signup, deadWorker)
		//		go worker(i, deadWorker)
	}
}
