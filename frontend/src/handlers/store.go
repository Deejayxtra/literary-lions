package handlers

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "sync"
)

type UserSession struct {
    Username string `json:"username"`
    Email    string `json:"email"`
}

type SessionStore struct {
    sessions map[string]UserSession
    mu       sync.RWMutex
    filename string
}

func NewSessionStore(filename string) *SessionStore {
    absPath, err := filepath.Abs(filename)
    if err != nil {
        fmt.Printf("Error getting absolute path: %v\n", err)
        return nil
    }
    store := &SessionStore{
        sessions: make(map[string]UserSession),
        filename: absPath,
    }
    store.LoadFromFile()
    return store
}

func (store *SessionStore) Set(token, username, email string) {
    store.mu.Lock()
    store.sessions[token] = UserSession{Username: username, Email: email}
    store.mu.Unlock()

    // Perform file save asynchronously to avoid blocking
    go func() {
        err := store.SaveToFile()
        if err != nil {
            fmt.Printf("Error saving to file: %v\n", err)
        }
    }()
}

func (store *SessionStore) Get(token string) (UserSession, bool) {
    store.mu.RLock()
    defer store.mu.RUnlock()
    session, exists := store.sessions[token]
    return session, exists
}

func (store *SessionStore) Delete(token string) {
    store.mu.Lock()
    delete(store.sessions, token)
    store.mu.Unlock()

    // Perform file save asynchronously to avoid blocking
    go func() {
        err := store.SaveToFile()
        if err != nil {
            fmt.Printf("Error saving to file: %v\n", err)
        }
    }()
}

func (store *SessionStore) SaveToFile() error {
    store.mu.RLock()
    defer store.mu.RUnlock()
    data, err := json.Marshal(store.sessions)
    if err != nil {
        return fmt.Errorf("error marshaling sessions: %v", err)
    }
    err = ioutil.WriteFile(store.filename, data, 0644)
    if err != nil {
        return fmt.Errorf("error writing to file %s: %v", store.filename, err)
    }
    return nil
}

func (store *SessionStore) LoadFromFile() {
    store.mu.Lock()
    defer store.mu.Unlock()
    if _, err := os.Stat(store.filename); os.IsNotExist(err) {
        fmt.Printf("File %s does not exist; will create new file upon saving.\n", store.filename)
        return
    }

    data, err := ioutil.ReadFile(store.filename)
    if err != nil {
        fmt.Printf("Error reading file %s: %v\n", store.filename, err)
        return
    }
    err = json.Unmarshal(data, &store.sessions)
    if err != nil {
        fmt.Printf("Error unmarshaling sessions: %v\n", err)
    }
}

var sessionStore = NewSessionStore("sessions.json")
