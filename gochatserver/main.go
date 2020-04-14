package main

import (
	"fmt"
	"github.com/krnblni/GoChat"
)

func main() {
	err := gochat.Run()
	fmt.Println(err)
}