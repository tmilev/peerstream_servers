package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tmilev/peerstream_servers/backend_request"
	"net"
	"net/url"
	"sync"
)


type Server struct {
	keyValueStore map[string]string
	theLock sync.Mutex
}
var theServer Server
func init () {
	theServer.keyValueStore = make (map[string]string)
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
	if requestLength < 500 {
		bufferToShow := buffer[:requestLength]
		fmt.Printf("Received bytes: %v\n", hex.EncodeToString(bufferToShow))
		fmt.Printf("UTF8 encoding: %v\n", string(bufferToShow))
	} else {
		bufferToShow := buffer[:500]
		fmt.Printf("Received bytes, utf8: %v ...\n", hex.EncodeToString(bufferToShow))
		fmt.Printf("UTF8 encoding: %v ...\n", string(bufferToShow))
	}
	if len(buffer) <= 5 {
		fmt.Printf("Message too short. ")
		return
	}
	if buffer[len(buffer) - 1] != '\n' {
		fmt.Printf(
			"Invalid message encoding: message does not end with a new line but rather with: %v. \n",
			buffer[len(buffer) - 1],
		)

		return
	}
	buffer = buffer[: len(buffer) - 1]
	var result = make(map[string]interface{})
	unescapedString, err := url.QueryUnescape(string(buffer))
	if err != nil {
		result["error"] = fmt.Sprintf(
			"Failed to url-decode your input. Error: %v",
			err,
		)
		writeJSON(connection, result)
		return
	}
	fmt.Printf("Unescaped input: %v\n", unescapedString)
	var request backendRequest.BackendRequest
	err = json.Unmarshal([]byte(unescapedString), &request)
	if err != nil {
		fmt.Printf("Failed to unmarshal your json bytes from: %v. %v\n", unescapedString, err)
		result["error"] = fmt.Sprintf("Failed to parse your json bytes to a backend request. %v", err)
		writeJSON(connection, result)
		return
	}
	fmt.Printf("DEBUG: go to command handling.")
	switch request.Command {
	case backendRequest.CommandSetKey:
		setKey(&request, result)
	case backendRequest.CommandGetAllKeys:
		getAllKeys(result)
	}
	writeJSON(connection, result)
}

func getAllKeys(result map[string]interface{}) {
	theServer.theLock.Lock()
	defer theServer.theLock.Unlock()

	for key, value := range theServer.keyValueStore {
		result[key] = value
	}
}

func setKey(request* backendRequest.BackendRequest, result map[string]interface{}) {
	theServer.theLock.Lock()
	defer theServer.theLock.Unlock()

	theServer.keyValueStore[request.Key] = request.Value
	result["comments"] = fmt.Sprintf("Key %v modified to %v", request.Key, request.Value)
}

func writeJSON(connection net.Conn, result map[string]interface{}) {
	resultBytes, _ := json.Marshal(result)
	resultEscapedBytes := url.QueryEscape(string(resultBytes))
	_, err := fmt.Fprintf(connection, resultEscapedBytes + "\n")
	if err != nil {
		fmt.Printf("Failed to return bytes to backend. ")
	}
}

