package gochat

import (
	"io"
	"fmt"
	"sync"
)

type chatRoom struct {
	name string
	messageChannel chan string
	clientChannels map[chan <- string]struct{}
	Quit chan struct{}
	*sync.RWMutex
}

func createChatRoom(name string) *chatRoom {
	room := &chatRoom{
		name: name,
		messageChannel: make(chan string),
		clientChannels: make(map[chan <- string]struct{}),
		Quit: make(chan struct{}),
		RWMutex: new(sync.RWMutex),
	}

	// run the channel after creation
	room.Run()
	return room
}

func (room *chatRoom) Run() {
	fmt.Println("Starting the chat room...")
	fmt.Println("Name: ", room.name)
	// simply listen to msg channel of room in separate go routine
	go func() {
		for message := range room.messageChannel {
			// whenever any message is recieved to room message channel 
			// simply broadcase it
			room.broadcastMessage(message)
		}
	}()
}

func (room *chatRoom) broadcastMessage(message string) {
	// apply read lock on clients of room
	room.RLock()
	defer room.RUnlock()

	fmt.Println("Recieved message: ", message)
	for writeChannel := range room.clientChannels {
		go func(writeChannel chan<- string) {
			fmt.Println("sending recieved message")
			writeChannel <- message
		}(writeChannel)
	}
}

func (room *chatRoom) AddClient(connection io.ReadWriteCloser) {
	fmt.Println("Adding new Client: ", connection)
	room.Lock()
	writeChannel, done := StartClient(room.messageChannel,connection, room.Quit)
	room.clientChannels[writeChannel] = struct{}{}
	room.Unlock()
	
	go func() {
		<-done
		room.RemoveClient(writeChannel)
	}()
}

func (room *chatRoom) RemoveClient(writeChannel chan<- string) {
	fmt.Println("Removing client: ", writeChannel)
	room.Lock()
	close(writeChannel)
	delete(room.clientChannels, writeChannel)
	room.Unlock()

	select {
	case <-room.Quit:
		if len(room.clientChannels) == 0 {
			close(room.messageChannel)
		}
	default:
	}
}

func (room *chatRoom) clientCount() int {
	return len(room.clientChannels)
}