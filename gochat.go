package gochat

import (
	"syscall"
	"os/signal"
	"os"
	"fmt"
	"net"
)

// Run - starts gochat server
func Run() error {
	listener, err := net.Listen("tcp", ":2900")
	if err != nil {
		fmt.Println("Error creating listener: ", err)
		return err
	}
	room := createChatRoom("Go Chat Room")

	go func() {
		channel := make(chan os.Signal)
		signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)
		<-channel

		listener.Close()
		fmt.Println("Closing TCP Connection")
		close(room.Quit)

		if room.clientCount() > 0 {
			fmt.Println("Client Count is not zero!")
			<-room.messageChannel
		}

		os.Exit(0)
	}()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			break
		}
		go handleNewClientConnection(room, connection)
	}
	
	return nil
}

func handleNewClientConnection(room *chatRoom, connection net.Conn) {
	fmt.Println("New Connection!!! Adding new client.", connection.RemoteAddr())
	room.AddClient(connection)
}