package main

import (
	"fmt"
	"github.com/tmilev/peerstream_servers/backend_request"
	"net"
	"sync"
)

type MessageWithResponse struct {
	request backendRequest.BackendMessage
	result  backendRequest.BackendMessage
}

type Server struct {
	Address              string
	keyValueStore        map[string]backendRequest.ValueWithVersion
	theLock              sync.Mutex
	server               *net.Listener
	peers                map[string]string
	messageQueue         chan backendRequest.BackendMessage
	messageQueueCapacity int64
	//lastPeerToPeerMessagesSent       []MessageWithResponse
	//numberOfPeerToPeerMessagesToKeep int
}

var theServer Server

func init() {
	theServer.keyValueStore = make(map[string]backendRequest.ValueWithVersion)
	theServer.peers = make(map[string]string)
}

func closeConnection(connection net.Conn, numberOfSuccessfulConnections int64) {
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

func (server *Server) handleConnection(connection net.Conn, numberOfSuccessfulConnections int64) {
	defer closeConnection(connection, numberOfSuccessfulConnections)

	var request backendRequest.BackendMessage

	err := backendRequest.ReadInterface(&request, connection)
	var result backendRequest.BackendMessage
	result.DestinationServerAddress = server.Address
	if err != nil {
		fmt.Printf("Error: first read on connection failed. %v\n", err)
		result.Error = fmt.Sprintf("Failed to parse your json bytes to a backend request. %v", err)
		backendRequest.WriteJSON(connection, result)
		return
	}
	switch request.Command {
	case backendRequest.CommandSetKey:
		server.setKey(&request, &result)
	case backendRequest.CommandGetAllKeys:
		server.getAllKeys(&result)
	case backendRequest.CommandGetPeers:
		server.getPeers(&result)
	case backendRequest.CommandAddPeer:
		server.addPeer(&request, &result)
	default:
		result.Error = fmt.Sprintf("Unrecognized command: %v", request.Command)
	}
	backendRequest.WriteJSON(connection, &result)
}

func (server *Server) getAllKeys(result *backendRequest.BackendMessage) {
	server.theLock.Lock()
	defer server.theLock.Unlock()
	if result.Data == nil {
		result.Data = make(map[string]backendRequest.ValueWithVersion)
	}
	for key, value := range server.keyValueStore {
		result.Data[key] = value
	}
}

func (server *Server) setKey(request *backendRequest.BackendMessage, result *backendRequest.BackendMessage) {
	server.theLock.Lock()
	defer server.theLock.Unlock()
	incomingValue := backendRequest.ValueWithVersion{
		Value:   request.Value,
		Version: request.Version,
	}
	currentValue, alreadyKnown := server.keyValueStore[request.Key]
	if !alreadyKnown {
		result.Comments = fmt.Sprintf(
			"Key %v added for the first time with value: %v",
			request.Key,
			request.Value,
		)
		// key is new to us, broadcast it
		server.setKeyAndBroadcast(request, result, &incomingValue)
		return
	}
	if incomingValue.Version < currentValue.Version {
		result.Comments = fmt.Sprintf("Key %v already known with a more recent version. ", request.Key)
		return
	}
	if incomingValue.Version == currentValue.Version {
		if incomingValue.Value < currentValue.Value {
			result.Comments = fmt.Sprintf(
				"Key %v has conflicting versions. Not updated (incoming value is smaller). ",
				request.Key,
			)
			return
		}
		if incomingValue.Value == currentValue.Value {
			result.Command = fmt.Sprintf(
				"Key %v already known. ",
				request.Key,
			)
			return
		}
		result.Comments = fmt.Sprintf(
			"Key %v has conflicting versions. Larger value takes precedence. ",
			request.Key,
		)
	}
	// key is new to us, broadcast it
	server.setKeyAndBroadcast(request, result, &incomingValue)
}

func (server *Server) setKeyAndBroadcast(
	request *backendRequest.BackendMessage,
	result *backendRequest.BackendMessage,
	incomingValue *backendRequest.ValueWithVersion,
) {
	// We broadcast a key only when it is new to us.
	// This guarantees that we broadcast every key only once (to all peers), which in turn guarantees that messages
	// do not bounce infinitely around.
	//Broadcasting is lock-less, so we must copy all info that requires locking
	server.keyValueStore[request.Key] = *incomingValue
	newKeyMap := make (map[string]backendRequest.ValueWithVersion)
	newKeyMap[request.Key] = *incomingValue
	peersToInform := server.getPeerNames()
	go server.broadcastToPeers_NoLocks(peersToInform, newKeyMap)
}

func (server *Server) getPeers(result *backendRequest.BackendMessage) {
	server.theLock.Lock()
	defer server.theLock.Unlock()
	if result.Data == nil {
		result.Data = make(map[string]backendRequest.ValueWithVersion)
	}
	for key, _ := range server.peers {
		result.Data[key] = backendRequest.ValueWithVersion{Value: key}
	}
}

func (server *Server) getPeerNames() []string {
	//server lock must be held
	result := make([]string, len(server.peers))
	index := 0
	for key := range server.peers {
		result[index] = key
		index ++
	}
	return result
}

func (server *Server) addPeer(request *backendRequest.BackendMessage, result *backendRequest.BackendMessage) {
	server.theLock.Lock()
	defer server.theLock.Unlock()
	peer := request.Peer
	_, alreadyKnown := server.peers[peer]
	if alreadyKnown {
		// no double-introductions.
		result.Comments = fmt.Sprintf("Peer %v already known. ", peer)
		return
	}
	server.peers[peer] = ""
	// Attention: introduction is tricky: we need to make sure that we don't block the server
	// when the key-value store is very large.
	keyValueStoreCopy := make(map[string]backendRequest.ValueWithVersion)
	for key, value := range server.keyValueStore {
		keyValueStoreCopy[key] = value
	}
	go server.introduceMyself_NoLocks(peer, keyValueStoreCopy)
	peersSnapShot := server.getPeerNames()
	go server.introducePeer_NoLocks(peersSnapShot, peer)
}

func (server *Server) introducePeer_NoLocks(peers []string, newComer string) {
	// Warning: locks not allowed here! This goroutine may block if the message queue is full.
	// In turn, the message queue can block if the lock is held
	// Warning: accessing the keyvalue store not allowed here: as already explained, locks are not allowed here.
	var introduceNewComer backendRequest.BackendMessage
	for i:= 0; i < len(peers); i ++ {
		introduceNewComer.DestinationServerAddress = peers[i]
		introduceNewComer.Peer = newComer
		introduceNewComer.Command = backendRequest.CommandAddPeer
		server.messageQueue <- introduceNewComer
	}
}

func (server *Server) introduceMyself_NoLocks(peer string, snapShot map[string]backendRequest.ValueWithVersion) {
	// Warning: locks not allowed here! This goroutine may block if the message queue is full.
	// In turn, the message queue can block if the lock is held.
	// Warning: accessing the keyvalue store not allowed here: as already explained, locks are not allowed here.
	var introduceMyself backendRequest.BackendMessage
	introduceMyself.DestinationServerAddress = peer
	introduceMyself.Peer = server.Address
	introduceMyself.Command = backendRequest.CommandAddPeer
	server.messageQueue <- introduceMyself
	server.broadcastToPeers_NoLocks([]string{peer}, snapShot)
}

func (server *Server) broadcastToPeers_NoLocks(peers []string, snapShot map[string]backendRequest.ValueWithVersion) {
	// Warning: locks not allowed here! This goroutine may block if the message queue is full.
	// In turn, the message queue can block if the lock is held.
	// Warning: accessing the keyvalue store not allowed here: as already explained, locks are not allowed here.
	for i := 0; i < len(peers); i ++ {
		for key, value := range snapShot {
			broadcastNextKey := backendRequest.BackendMessage{
				Key:     key,
				Version: value.Version,
				Value:   value.Value,
				DestinationServerAddress: peers[i],
				Command:                  backendRequest.CommandSetKey,
			}
			server.messageQueue <- broadcastNextKey
		}
	}
}
