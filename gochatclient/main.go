package main

import (
	"os"
	"bufio"
	"net"
	"fmt"
	"time"
	"math/rand"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	name := fmt.Sprintf("Anonymous%d", rand.Intn(1000))

	fmt.Println("Starting, you are a client to gochat!")
	fmt.Println("Your alias: ", name)

	fmt.Println("Connecting to gochat server...")
	connection, err := net.Dial("tcp", "127.0.0.1:2900")
	if err != nil {
		fmt.Println("Could not connect to server, ", err)
	}
	fmt.Println("Connected...")
	name += ": "

	go func() {
		// this scanner reads incoming messages from server
		scanner := bufio.NewScanner(connection)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// this scanner reads messages from terminal which clients want to send
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() && err == nil {
		message := scanner.Text()
		_, err = fmt.Fprintf(connection, name+message+"\n")
	}
}