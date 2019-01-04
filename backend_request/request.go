package backendRequest

type ValueWithVersion struct {
	Value   string `json:"value,omitempty"`
	Version int64  `json:"version,omitempty"`
}

type BackendMessage struct {
	// Version: records the last "time" ("version")  this message was known to be sent by the user.
	// If the version is non-positive, it is assumed to be 0.
	// A message with later time takes precedence over one with earlier time.
	// If there are two messages with same keys, same versions
	// and different values (possibly due to malicious behavior)
	// the key with larger value (string-comparison-wise) takes precedence.
	Version                  int64                       `json:"version,omitempty"`
	Command                  string                      `json:"command,omitempty"`
	Key                      string                      `json:"key,omitempty"`
	Value                    string                      `json:"value,omitempty"`
	DestinationServerAddress string                      `json:"serverAddress,omitempty"`
	Result                   string                      `json:"result,omitempty"`
	Error                    string                      `json:"error,omitempty"`
	Comments                 string                      `json:"comments,omitempty"`
	Peer                     string                      `json:"peer,omitempty"`
	Data                     map[string]ValueWithVersion `json:"data,omitempty"`
}

const (
	CommandGetAllKeys = "getAllKeys"
	CommandSetKey     = "setKey"
	CommandGetPeers   = "getPeers"
	CommandAddPeer    = "addPeer"
)
