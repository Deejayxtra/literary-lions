package handlers

import (
    "fmt"
    "net/http"
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

        respChan := make(chan models.Data)
        go ShowPosts(w, r, respChan)

        // Retrieve posts from the channel
        responseDetails := <-respChan
        // fmt.Println("responseDetails:", responseDetails)
        // fmt.Println("Response:", responseDetails.Posts)

        data := struct {
            RoomID   string
            RoomName string
            Messages []models.Message
        }{
            RoomID:   roomID,
            RoomName: getRoomName(roomID), // Function to get the room name based on roomID
            Messages: responseDetails.Posts,
        }
        RenderTemplate(w, "conversation-room.html", data)
        return
    } else if r.Method == http.MethodPost {
        // Extract Message from form values
        content := r.FormValue("content")
        title := r.FormValue("title")

        // Sample Message
        comment := models.Message{
            Title:   title,
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
        "category1": "Recent Posts",
        "category2": "News",
        "category3": "Entertainment",
        "category4": "Music",
        "category5": "Sports",
        "category6": "General",
    }
    if name, ok := roomNames[roomID]; ok {
        return name
    }
    return "Unknown Room"
}
