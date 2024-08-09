package handlers

import (
    "sync"
)

// UserSession stores both username and email
type UserSession struct {
    Username string `json:"username"`
    Email    string `json:"email"`
}

type SessionStore struct {
    sessions map[string]UserSession // map[token]UserSession
    mu       sync.RWMutex
}

func NewSessionStore() *SessionStore {
    return &SessionStore{
        sessions: make(map[string]UserSession),
    }
}

func (store *SessionStore) Set(token, username, email string) {
    store.mu.Lock()
    defer store.mu.Unlock()
    store.sessions[token] = UserSession{Username: username, Email: email}
}

func (store *SessionStore) Get(token string) (UserSession, bool) {
    store.mu.RLock()
    defer store.mu.RUnlock()
    session, exists := store.sessions[token]
    return session, exists
}

func (store *SessionStore) Delete(token string) {
    store.mu.Lock()
    defer store.mu.Unlock()
    delete(store.sessions, token)
}

var sessionStore = NewSessionStore()