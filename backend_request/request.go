package backendRequest

type BackendRequest struct {
	Command       string `json:"command"`
	Key           string `json:"key"`
	Value         string `json:"value"`
	ServerAddress string `json:"serverAddress"`
}

const (
	CommandGetAllKeys = "getAllKeys"
	CommandSetKey     = "setAllKeys"
)
