package gochat

import (
	"io"
	"fmt"
	"bufio"
)

type chatClient struct {
	// reads the message from terminal which client wants to send to other users
	*bufio.Reader
	*bufio.Writer
	writeChannel chan string
}

// StartClient - creates and starts a new client
func StartClient(messageChannel chan<- string, connection io.ReadWriteCloser, quit chan struct{}) (chan<- string, <-chan struct{}) {
	client := new(chatClient)
	client.Reader = bufio.NewReader(connection)
	client.Writer = bufio.NewWriter(connection)
	client.writeChannel = make(chan string)
	done := make(chan struct{})

	// setting up reader
	// this reads incoming messages from the client to the server
	// CLIENT -------> SERVER
	go func() {
		scanner := bufio.NewScanner(client.Reader)
		for scanner.Scan() {
			fmt.Println("This is READER!!!")
			fmt.Println(scanner.Text())
			// send message into chat room message channel
			// CLIENT MESSAGE -------> ROOM MESSAGE CHANNEL
			messageChannel <- scanner.Text()
		}
		done <- struct{}{}
	}()

	// setting up writer
	// this checks if the client has recieved any new message
	// it then writes the new message from room on the connection 
	// and then flushes it to send it back to client
	// CLIENT <------ CONNECTION <------ SERVER (ROOM)
	go func() {
		for s := range client.writeChannel {
			fmt.Println("Sending... This is WRITER!!!")
			client.WriteString(s + "\n")
			client.Flush()
		}
	}()

	go func() {
		select {
		case <- quit:
			fmt.Println("Server closed connection!")
			connection.Close()
		case <- done:
		}
	}()

	return client.writeChannel, done
}