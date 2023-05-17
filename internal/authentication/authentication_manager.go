package authentication

import "sync"

var authentication = &Authentication{}

type Authentication struct {
	username   string
	credential string
	mu         sync.RWMutex
}

func GetCurrentUser() string {
	return authentication.username
}

func GetCurrentUserCredential() string {
	return authentication.credential
}

func SetCurrentUser(username, token string) {
	authentication.mu.Lock()
	defer authentication.mu.Unlock()
	authentication.username = username
	authentication.credential = token
}

func ClearCurrentUser() {
	authentication.mu.Lock()
	defer authentication.mu.Unlock()
	authentication.username = ""
	authentication.credential = ""
}
