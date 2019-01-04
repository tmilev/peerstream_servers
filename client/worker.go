package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/tmilev/peerstream_servers/backend_request"
	"net"
	"net/http"
	"net/url"
	"time"
)

func processCommand(writer http.ResponseWriter, request *backendRequest.BackendRequest, result map[string]interface{}) {
	result["command"] = request.Command
	if request.Command == "" {
		request.Command = "[empty]"
	}
	connection, err := net.DialTimeout("tcp", request.ServerAddress, time.Duration(2) * time.Second)
	if err!= nil {
		result["error"] = fmt.Sprintf("Failed to dial %v. Error: %v. ", request.ServerAddress, err)
		result["comments"] = fmt.Sprintf(
			"Please don't forget to build and start the server: cd server && go build && ./server",
		)
		return
	}
	defer closeConnection(connection, result)
	result["comments"] = "Please add timeouts for reading/writing. "

	bytesToSend, _ := json.Marshal(request)
	bytesToSendString := url.QueryEscape(string(bytesToSend))
	_, err = fmt.Fprint(connection, bytesToSendString + "\n")
	if err != nil {
		result["error"] = fmt.Sprintf(
			"Failed to write bytes: %v to address: %v",
			bytesToSendString,
			request.ServerAddress,
		)
		return
	}
	resultMessage, err := bufio.NewReader(connection).ReadString('\n')
	resultUnescaped, err := url.QueryUnescape(resultMessage)
	if err != nil {
		result["error"] = fmt.Sprintf("Failed to unescape the server's response. ")
		return
	}
	result["result"] = resultUnescaped
}

func closeConnection (connection net.Conn, result map[string]interface{}) {
	err := connection.Close()
	if err != nil {
		result["error"] = err
	}
}