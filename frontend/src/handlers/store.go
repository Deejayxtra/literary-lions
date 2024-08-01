package handlers

import (
	"sync"
)

type SessionStore struct {
	sessions map[string]string // map[token]username
	mu       sync.RWMutex
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]string),
	}
}

func (store *SessionStore) Set(token, username string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.sessions[token] = username
}

func (store *SessionStore) Get(token string) (string, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	username, exists := store.sessions[token]
	return username, exists
}

func (store *SessionStore) Delete(token string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.sessions, token)
}

var sessionStore = NewSessionStore()
