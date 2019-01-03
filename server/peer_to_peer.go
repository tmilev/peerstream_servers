package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
)


func main() {
	startServer()
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

func startServer() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter hostname:port. Example: 127.0.0.1:9000. Enter=auto ->")
	addressRaw, _ := reader.ReadString('\n')
	address := strings.Trim(addressRaw, ",\t\n ")
	if len(address) < 3 {
		address = "localhost:9000"

	}
	fmt.Printf("About to connect to address: %v\n", address)
	var server *net.Listener = nil
	start := 8999
	for i:= start; i < 9011; i ++ {
		if i > start {
			address = fmt.Sprintf("localhost:%v", i)
			fmt.Printf(
				"Attempting to start with auto-generated address: %v.\n",
				address,
			)
		}
		server = startServerOneAttempt(address)
		if server != nil {
			break
		}
	}
	if server == nil {
		fmt.Printf("Fatal error: could not start peer to peer server.\n")
		return
	}
	defer closeServer(server)
	fmt.Printf("Listening on %v.\n", address)
	numberOfFailedAccepts := 0
	var numberOfSuccessfulConnections int64 = 0
	for {
		// Listen for an incoming connection.
		connection, err := (*server).Accept()
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
		go handleConnection(connection, numberOfSuccessfulConnections)
	}
}

func closeServer(server *net.Listener) {
	err := (*server).Close()
	if err != nil {
		fmt.Printf("Failed to close server. %v\n", err)
	}
}

func closeConnection (connection net.Conn, numberOfSuccessfulConnections int64) {
	err := connection.Close()
	if err == nil {
		fmt.Printf(
			"Successfully closed connection %v\n",
			numberOfSuccessfulConnections,
		)
	} else {
		fmt.Printf(
			"Failed to properly close connection %v. Error: %v\n",
			numberOfSuccessfulConnections, err,
		)
	}
}

func handleConnection(connection net.Conn, numberOfSuccessfulConnections int64) {
	defer closeConnection(connection, numberOfSuccessfulConnections)
	buffer := make([]byte, 20000)
	// Read the incoming connection into the buffer.
	requestLength, err := connection.Read(buffer)
	if err != nil {
		fmt.Printf(
			"Error reading connections %v. Error message: %v\n",
			numberOfSuccessfulConnections,
			err,
		)
		return
	}
	fmt.Printf(
		"Received %v bytes on connection %v. ",
		requestLength, numberOfSuccessfulConnections,
	)
	if requestLength < 100 {
		bufferToShow := buffer[:requestLength]
		fmt.Printf("Received bytes: %v\n", hex.EncodeToString(bufferToShow))
		fmt.Printf("UTF8 encoding: %v\n", string(bufferToShow))
	} else {
		bufferToShow := buffer[:100]
		fmt.Printf("Received bytes, utf8: %v ...\n", hex.EncodeToString(bufferToShow))
		fmt.Printf("UTF8 encoding: %v\n", string(bufferToShow))
	}
}
