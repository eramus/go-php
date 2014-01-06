package main

import (
	"fmt"
)

type TestFunc func(string) string

func test(t string) string {
	fmt.Println(t)
	return t
}

func acceptTest(f TestFunc, param string) string {
	return f(param)
}

func main() {
	t := "Hello World"
	r := acceptTest(test, t)

	fmt.Println("RETURN:", r)
}
