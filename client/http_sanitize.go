package main

import (
	"encoding/json"
	"fmt"
	"github.com/tmilev/peerstream_servers/backend_request"
	"net/http"
	"net/url"
)

func requestBackend(writer http.ResponseWriter, request *http.Request) {
	result := make (map[string]interface{})
	var commandArray []string
	var ok bool
	commandArray, ok = request.URL.Query()["json"]
	if !ok {
		result["error"] = fmt.Sprint("Command query variable missing. ")
		returnJSON(writer, result)
		return
	}
	if len(commandArray) != 1 {
		result["error"] = fmt.Sprintf("Exactly one query command allowed, you gave me: %v", len(commandArray))
		returnJSON(writer, result)
		return
	}
	command, err := url.QueryUnescape(commandArray[0])
	if err != nil {
		result["error"] = fmt.Sprintf("Failed to unescape your query %v. Error: %v", commandArray[0], err)
		returnJSON(writer, result)
		return
	}
	var parsedRequest backendRequest.BackendMessage
	err = json.Unmarshal([]byte(command), &parsedRequest)
	if err != nil {
		result["error"] = fmt.Sprintf("Failed to parse your JSON query %v. Error: %v", command, err)
		returnJSON(writer, result)
		return
	}
	processCommand(writer, &parsedRequest, result)
	returnJSON(writer, result)
}

func returnJSON (writer http.ResponseWriter, result map[string]interface{}) {
	var resultBytes, err = json.Marshal(result)
	if err != nil {
		fmt.Printf("Bad result object %v.\n", err)
		return
	}
	_, err = fmt.Fprintf(writer, string(resultBytes))
	if err != nil {
		fmt.Printf("Bad result object %v.\n", err)
		return
	}
}
