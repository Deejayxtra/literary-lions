package handlers

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	// "html/template"
	// "log"
	"net/http"
	// "sync"
	// "time"
	// "io/ioutil"

	"literary-lions/frontend/src/models"
)

// ConversationRoom handles the conversation room.
func ConversationRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		roomID := r.URL.Query().Get("room_id")
		if roomID == "" {
			http.Error(w, "Missing room_id", http.StatusBadRequest)
			return
		}

		conversationsLock.Lock()
		messages := conversations[roomID]
		conversationsLock.Unlock()

		data := struct {
			RoomID   string
			RoomName string
			Messages []models.Message
		}{
			RoomID:   roomID,
			RoomName: getRoomName(roomID), // Function to get the room name based on roomID
			Messages: messages,
		}
		RenderTemplate(w, "conversation-room.html", data)
		return
	} else if r.Method == http.MethodPost {
		// Extract Message from form values
		content := r.FormValue("content")

		// Sample Message
		comment := models.Message{
			Content: content,
		}

		roomID := r.URL.Query().Get("room_id")
		if roomID == "" {
			http.Error(w, "Missing room_id", http.StatusBadRequest)
			return
		}

		conversationsLock.Lock()
		conversations[roomID] = append(conversations[roomID], comment)
		conversationsLock.Unlock()

		// Redirect to the same conversation room to display the updated conversation
		redirectURL := fmt.Sprintf("/conversation-room?room_id=%s", roomID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// Helper function to get room name based on roomID
func getRoomName(roomID string) string {
	roomNames := map[string]string{
		"channel1": "General",
		"channel2": "News",
		"channel3": "Entertainment",
		"channel4": "Music",
		"channel5": "Sports",
		"channel6": "Random",
	}
	if name, ok := roomNames[roomID]; ok {
		return name
	}
	return "Unknown Room"
}
