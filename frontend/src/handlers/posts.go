package handlers

import (
	//"html/template"
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"literary-lions/frontend/src/models"
	//"log"

	"net/http"
)

// ShowPosts handles the fetching and displaying of all posts in Recent Posts.
func ShowPosts(w http.ResponseWriter, r *http.Request, respChan chan<- models.Data) {
    // Make an HTTP GET request to the /api/posts endpoint
    resp, err := http.Get("http://localhost:8080/api/posts")
    if err != nil {
        http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, "Failed to read response", http.StatusInternalServerError)
        return
    }

    // Parse the JSON response into a slice of Post models
    var posts []models.Message
    err = json.Unmarshal(body, &posts)
    if err != nil {
        http.Error(w, "Failed to parse response", http.StatusInternalServerError)
        return
    }
    //log.Println("posts:", posts)

    // Send the data through the channel
    respChan <- models.Data{
        Posts: posts,
    }
    close(respChan)
}

// CreatePost renders the template for creating a post
func CreatePost(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
		// Extract form values
		category := r.FormValue("category")
		title := r.FormValue("title")
		content := r.FormValue("content")

		// Create a new post
		comment := models.Message{
			Title: title,
			Content: content,
		}

		// Determine the roomID based on category
		roomID := "category1" // Default to Recent Posts
		if category != "Recent Posts" {
			roomID = mapCategoryToRoomID(category)
		}

		conversationsLock.Lock()
		conversations[roomID] = append(conversations[roomID], comment)
		conversationsLock.Unlock()

		// Store the roomID in a cookie or session to use in the confirmation page
		http.SetCookie(w, &http.Cookie{
			Name:  "lastRoomID",
			Value: roomID,
			Path:  "/",
		})

		// Redirect to a confirmation page
		http.Redirect(w, r, "/conversation-room?room_id=category1", http.StatusSeeOther)
	} else {
		// Show the form for GET requests
		categories := []models.Category{
			{Name: "News"},
			{Name: "Entertainment"},
			{Name: "Music"},
			{Name: "Sports"},
			{Name: "Random"},
		}

		data := struct {
			Category []models.Category
		}{
			Category: categories,
		}

		RenderTemplate(w, "create-post.html", data)
	}
}

func PostConfirmation(w http.ResponseWriter, r *http.Request) {
	// Get the roomID from the cookie
	cookie, err := r.Cookie("lastRoomID")
	if err != nil {
		http.Error(w, "Room ID not found", http.StatusBadRequest)
		return
	}

	roomID := cookie.Value

	data := struct {
		RoomID string
	}{
		RoomID: roomID,
	}

	RenderTemplate(w, "post-confirmation.html", data)
		// After displaying the confirmation, redirect to Recent Posts
	http.Redirect(w, r, "/conversation-room?room_id=category1", http.StatusSeeOther)
}


// Helper function to map category to room ID
func mapCategoryToRoomID(category string) string {
	roomIDs := map[string]string{
		"Recent Posts":  "category1",
		"News":          "category2",
		"Entertainment": "category3",
		"Music":         "category4",
		"Sports":        "category5",
		"Random":        "category6",
	}
	if roomID, ok := roomIDs[category]; ok {
		return roomID
	}
	return "category1" // Default to Recent Posts
}