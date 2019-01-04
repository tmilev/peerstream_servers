package main

import (
	"fmt"
	"github.com/tmilev/peerstream_servers/backend_request"
	"net/http"
)

func processCommand(writer http.ResponseWriter, request *backendRequest.BackendMessage, result map[string]interface{}) {
	result["command"] = request.Command
	if request.Command == "" {
		request.Command = "[empty]"
	}
	var finalResult backendRequest.BackendMessage
	var comments string
	err := backendRequest.ConnectRequestReturn(request, &finalResult, &comments)
	if comments != "" {
		result["comments"] = comments
	}
	if err != nil {
		result["error"] = fmt.Sprintf("Error connecting to server. %v", err)
	} else {
		result["result"] = finalResult
	}
}