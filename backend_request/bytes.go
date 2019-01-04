package backendRequest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"time"
)

func ReadUnescapedString(buffer []byte, connection net.Conn) (string, error) {
	if len(buffer) < 100 {
		panic("Buffer too small or not allocated. ")
	}
	requestLength, err := connection.Read(buffer)
	if err != nil {
		return "", err
	}
	fmt.Printf("Received %v bytes:\n%v\n", requestLength, string(buffer))
	if requestLength <= 2 {
		return "", errors.New("Message too short. ")
	}
	if buffer[requestLength - 1] != '\n' {
		return "", errors.New(fmt.Sprintf(
			"Invalid message encoding: message does not end with a new line but rather with: %v.\n",
			buffer[requestLength - 1],
		))
	}
	slicedInput := buffer[: requestLength - 1] //requestLength ensured to be > 5.
	unescapedString, err := hex.DecodeString(string(slicedInput))
	fmt.Printf("Unescaped string has %v bytes: %v\n", len(unescapedString), string(unescapedString))
	if err != nil {
		return "", err
	}
	return string(unescapedString), nil
}

func ReadInterface(output interface{}, connection net.Conn) error {
	buffer := make([]byte, 20000)
	jsonString, err := ReadUnescapedString(buffer, connection)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(jsonString), output)
	if err != nil {
		return err
	}
	return nil
}

func WriteJSON(connection net.Conn, data interface{}) bool {
	resultBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("While writing bytes to connection, failed to convert data to JSON.")
		return false
	}
	fmt.Printf(
		"About to write %v bytes:\n%v\n",
		len(resultBytes),
		string(resultBytes),
	)
	resultEscapedBytes := hex.EncodeToString(resultBytes) + "\n"
	fmt.Printf(
		"The data encodes to %v bytes:\n%v",
		len(resultEscapedBytes),
		resultEscapedBytes,
	)
	_, err = fmt.Fprintf(connection, resultEscapedBytes)
	if err != nil {
		fmt.Println("Failed to return bytes to backend.")
		return false
	}
	return true
}

func closeConnection (connection net.Conn, comments *string) {
	err := connection.Close()
	if err != nil && comments != nil{
		*comments = fmt.Sprintf("Error closing connection %v. ", err)
	}
}

func ConnectRequestReturn(request *BackendMessage, resultMessage *BackendMessage, comments *string) error {
	connection, err := net.DialTimeout("tcp", request.DestinationServerAddress, time.Duration(2) * time.Second)
	if err!= nil {
		if comments != nil {
			*comments = fmt.Sprintf(
				"Please don't forget to build and start the server: cd server && go build && ./server",
			)
		}
		return err
	}
	defer closeConnection(connection, comments)
	if comments != nil {
		*comments = "Please add timeouts for reading/writing. "
	}
	success := WriteJSON(connection, request)
	if ! success {
		return errors.New(fmt.Sprintf(
			"Failed to write JSON. Server address: %v. ",
			request.DestinationServerAddress,
		))
	}
	err = ReadInterface(resultMessage, connection)
	if err != nil {
		fmt.Printf("Failed read the server's response. %v", err)
	}
	return err
}