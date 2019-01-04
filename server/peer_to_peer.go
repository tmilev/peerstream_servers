package main

import (
	"bufio"
	"fmt"
	"github.com/tmilev/peerstream_servers/backend_request"
	"net"
	"os"
	"strings"
)

func main() {
	theServer.initializeAndStart()
}

func startServerOneAttempt(address string) *net.Listener {
	server, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf(
			"Error listening to %v. Error: %v\n",
			address,
			err,
		)
		return nil
	}
	return &server
}

func (server *Server) initializeAndStart() {
	//server.numberOfPeerToPeerMessagesToKeep = 40
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter hostname:port. Example: 127.0.0.1:9000. Enter=auto ->")
	addressRaw, _ := reader.ReadString('\n')
	server.Address = strings.Trim(addressRaw, ",\t\n ")
	if len(server.Address) < 3 {
		server.Address = "localhost:9000"
	}
	fmt.Printf("About to connect to address: %v\n", server.Address)
	server.server = nil
	start := 9000
	for i:= start; i < 9011; i ++ {
		if i > start {
			server.Address = fmt.Sprintf("localhost:%v", i)
			fmt.Printf(
				"Attempting to start with auto-generated address: %v.\n",
				server.Address,
			)
		}
		server.server = startServerOneAttempt(server.Address)
		if server.server != nil {
			break
		}
	}
	if server.server == nil {
		fmt.Printf("Fatal error: could not start peer to peer server.\n")
		return
	}
	server.messageQueueCapacity = 100
	server.messageQueue = make(chan backendRequest.BackendMessage, server.messageQueueCapacity)
	server.start()
}

func (server *Server) start() {
	go server.broadcastMessages()
	defer server.close()
	fmt.Printf("Listening on %v.\n", server.Address)
	numberOfFailedAccepts := 0
	var numberOfSuccessfulConnections int64 = 0

	for {
		// Listen for an incoming connection.
		connection, err := (*server.server).Accept()
		if err != nil {
			numberOfFailedAccepts ++
			fmt.Printf(
				"Error accepting connection %v. So far %v failures encountered.\n",
				err,
				numberOfFailedAccepts,
			)
			continue
		}
		// Handle connections in a new goroutine.
		numberOfSuccessfulConnections ++
		go server.handleConnection(connection, numberOfSuccessfulConnections)
	}
}

func (server *Server) broadcastMessages() {
	var currentMessage MessageWithResponse
	for {
		currentMessage.request = <-server.messageQueue
		fmt.Printf("Message received: %v\n", currentMessage)
		err := backendRequest.ConnectRequestReturn(&currentMessage.request, &currentMessage.result, nil)
		if err != nil {
			fmt.Printf("Error broadcasting message %v. %v\n", currentMessage, err)
		}
		// server.lastPeerToPeerMessagesSent is modified in this goroutine only: no need for locks.
		//
		//if server.numberOfPeerToPeerMessagesToKeep > 0 {
		//	server.lastPeerToPeerMessagesSent = append(server.lastPeerToPeerMessagesSent, currentMessage)
		//}
		//numberOfMessages := len(server.lastPeerToPeerMessagesSent)
		//if numberOfMessages > server.numberOfPeerToPeerMessagesToKeep {
		//	start := numberOfMessages - server.numberOfPeerToPeerMessagesToKeep
		//	server.lastPeerToPeerMessagesSent = server.lastPeerToPeerMessagesSent[start: numberOfMessages]
		//}
	}

}

func (server *Server) close() {
	err := (*server.server).Close()
	if err != nil {
		fmt.Printf("Failed to close server. %v\n", err)
	}
}
