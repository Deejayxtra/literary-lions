package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"literary-lions/frontend/src/config"
	"literary-lions/frontend/src/models"
	"log"
	"net/http"
	"sync"
	"unicode/utf8"
)

// Helper function to truncate post content to 150 characters
func truncateContent(content string, limit int) string {
	if utf8.RuneCountInString(content) > limit {
		runes := []rune(content)
		return string(runes[:limit]) + "..."
	}
	return content
}

func ShowPosts(w http.ResponseWriter, r *http.Request) {
	// Get the search query parameters
	keyword := r.URL.Query().Get("query")
	category := r.URL.Query().Get("category")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// Construct the API request URL with query parameters
	apiURL := config.BaseApi + "/posts?"
	if keyword != "" {
		apiURL += "keyword=" + keyword + "&"
	}
	if category != "" {
		apiURL += "category=" + category + "&"
	}
	if startDate != "" {
		apiURL += "start_date=" + startDate + "&"
	}
	if endDate != "" {
		apiURL += "end_date=" + endDate + "&"
	}

	// Make an HTTP GET request to the /api/posts endpoint
	resp, err := http.Get(apiURL)
	if err != nil {
		message := "Failed to fetch posts"
		StatusInternalServerError(w, message)
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
	var posts []models.Post
	err = json.Unmarshal(body, &posts)
	if err != nil {
		message := "Error rendering template"
		StatusInternalServerError(w, message)
		log.Fatalf("Error rendering template: %v", err)
		return
	}

	// Handle the case where no posts match the search criteria
	if len(posts) == 0 {
		// Get the authentication status and the currentUser if any
		currentUser, authenticated := isAuthenticated(r)
		// No posts found, so display a friendly message
		data := struct {
			Posts         []models.Post
			Authenticated bool
			Categories    []string
			Username      string
			NoPostsFound  bool
			SearchMessage string
		}{
			Posts:         posts,
			Authenticated: authenticated,
			Categories:    []string{"Random", "News", "Sport", "Technology", "Science", "Health"},
			Username:      currentUser,
			NoPostsFound:  true,
			SearchMessage: "No posts found for the selected criteria.",
		}

		RenderTemplate(w, "index.html", data)
		return
	}

	// Truncate the content of each post
	for i := range posts {
		posts[i].Content = truncateContent(posts[i].Content, 150)
	}

	// Get the authentication status and the currentUser if any
	currentUser, authenticated := isAuthenticated(r)

	// Hardcoded categories for now. Might be nice to add feature for creating it and make it dynamic later
	categories := []string{"Random", "News", "Sport", "Technology", "Science", "Health"}

	data := struct {
		Posts         []models.Post
		Authenticated bool
		Categories    []string
		Username      string
		NoPostsFound  bool
	}{
		Posts:         posts,
		Authenticated: authenticated,
		Categories:    categories,
		Username:      currentUser,
		NoPostsFound:  false,
	}

	// Render the template with posts and authentication status
	RenderTemplate(w, "index.html", data)
}

// Display posts by category.
func ShowPostsByCategory(w http.ResponseWriter, r *http.Request) {
	// Extract the category query parameter from the URL
	category := r.URL.Query().Get("category")

	var url string
	if category == "" {
		// If category is empty, fetch all posts
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		// If category is specified, fetch posts by category
		url = config.BaseApi + "/posts/category/" + category
	}

	// Make an HTTP GET request to the /api/posts endpoint
	resp, err := http.Get(url)
	if err != nil {
		message := "Failed to fetch posts"
		StatusInternalServerError(w, message)
		// http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Get the authentication status and the currentUser if any
	currentUser, authenticated := isAuthenticated(r)

	// Hardcoded categories for now. Might be nice to add feature for creating it and make it dynamic later
	categories := []string{"Random", "News", "Sport", "Technology", "Science", "Health"}

	// Check response status code
	if resp.StatusCode == http.StatusNotFound {
		var emptyPosts []models.Post
		data := struct {
			Posts         []models.Post
			Authenticated bool
			Categories    []string
			Username      string
		}{
			Posts:         emptyPosts,
			Authenticated: authenticated,
			Categories:    categories,
			Username:      currentUser,
		}

		// Render the template with posts and authentication status
		RenderTemplate(w, "index.html", data)
		return
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		message := "Failed to read response"
		StatusInternalServerError(w, message)
		// http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Parse the JSON response into a slice of Post models
	var posts []models.Post
	err = json.Unmarshal(body, &posts)
	if err != nil {
		message := "Failed to parse response"
		StatusInternalServerError(w, message)
		// http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	// Truncate the content of each post
	for i := range posts {
		posts[i].Content = truncateContent(posts[i].Content, 150)
	}

	data := struct {
		Posts         []models.Post
		Authenticated bool
		Categories    []string
		Username      string
	}{
		Posts:         posts,
		Authenticated: authenticated,
		Categories:    categories,
		Username:      currentUser,
	}
	// Render the template with posts and authentication status
	RenderTemplate(w, "index.html", data)
}

// ShowPostByID handles displaying a post by its ID
func ShowPostByID(w http.ResponseWriter, r *http.Request) {
	// Extract the id query parameter from the URL
	id := r.URL.Query().Get("id")

	// Create a new GET request
	req, err := http.NewRequest("GET", config.BaseApi+"/post/"+id, nil)
	if err != nil {
		message := "Failed to create request"
		StatusInternalServerError(w, message)
		// http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Use an http.Client to make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		message := "Failed to fetch post"
		StatusInternalServerError(w, message)
		// http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		message := "Failed to read response"
		StatusInternalServerError(w, message)
		// http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Parse the JSON response into a PostDetails model
	var response models.PostDetails
	err = json.Unmarshal(body, &response)
	if err != nil {
		message := "Failed to parse response"
		StatusInternalServerError(w, message)
		// http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	// Format the created_at date
	formattedDate := response.Post.CreatedAt.Format("January 2, 2006 at 3:04pm")

	// Get the authentication status and the currentUser if any
	currentUser, authenticated := isAuthenticated(r)

	data := struct {
		Post          models.Post
		FormattedDate string
		Authenticated bool
		Comments      []models.Comment
		Error         bool
		Username      string
		Likes		  int
		Dislikes      int
	}{
		Post:          response.Post,
		FormattedDate: formattedDate,
		Authenticated: authenticated,
		Comments:      response.Comments,
		Error:         false,
		Username:      currentUser,
		Likes:		   response.Likes,
		Dislikes:      response.Dislikes,
	}
	// Render the template with posts and authentication status
	RenderTemplate(w, "post.html", data)
}


func CreatePost(w http.ResponseWriter, r *http.Request) {
	// Get the authentication status and the currentUser if any
	currentUser, authenticated := isAuthenticated(r)

	// Renders page to user to select from category
	if r.Method == http.MethodGet {
		// Hardcoded categories for now. Might be nice to add feature for creating it and make it dynamic later
		categories := []string{"Random", "News", "Sport", "Technology", "Science", "Health"}
		data := struct {
			Categories    []string
			Error         bool
			Authenticated bool
			Username      string
		}{
			Categories:    categories,
			Error:         false,
			Authenticated: authenticated,
			Username:      currentUser,
		}
		tmpl := template.Must(template.ParseFiles("templates/create-post.html"))
		tmpl.Execute(w, data)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		category := r.FormValue("category")
		title := r.FormValue("title")
		content := r.FormValue("content")

		respChan := make(chan models.ResponseDetails, 1)
		var wg sync.WaitGroup

		// Defines the paylod sent to the backend
		payload := models.Post{
			Category: category,
			Title:    title,
			Content:  content,
		}

		// Extract the session cookie from the header
		cookieToken, err := r.Cookie("session_token")
		if err != nil {
			// User must be logged-in to continue
			message := `You are not authorized! Please <a href="/login">login</a> before creating a post.`
			tmpl := template.Must(template.ParseFiles("templates/create-post.html"))
			tmpl.Execute(w, map[string]interface{}{
				"Error": template.HTML(message),
			})
			return
		}

		// Calls the function that sends request to the server
		wg.Add(1)
		go func() {
			SendCreatePostRequest(cookieToken, payload, &wg, respChan)
		}()

		go func() {
			wg.Wait()
			close(respChan)
		}()

		responseDetails := <-respChan

		if responseDetails.Status == http.StatusCreated {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else if responseDetails.Status == http.StatusUnauthorized {
			// responseDetails.Status = http.StatusUnauthorized
			responseDetails.Message = `You are not authorized! Please <a href="/login">login</a> before creating a post.`
			tmpl := template.Must(template.ParseFiles("templates/create-post.html"))
			tmpl.Execute(w, map[string]interface{}{
				"Error": template.HTML(responseDetails.Message),
			})
		} else {
			// responseDetails.Status = resp.StatusCode
			responseDetails.Message = "Oops! Something went wrong. Failed to create post."
			tmpl := template.Must(template.ParseFiles("templates/create-post.html"))
			tmpl.Execute(w, map[string]interface{}{
				"Error": responseDetails.Message,
			})
		}
	}
}

func SendCreatePostRequest(cookie *http.Cookie, payload models.Post, waitGroup *sync.WaitGroup, respChan chan models.ResponseDetails) {
	defer waitGroup.Done()	// Ensure the channel is closed once this function completes

	// Convert payload to JSON
	postData, err := json.Marshal(payload)
	if err != nil {
		respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Failed to marshal payload"}
		return
	}

	// Create a POST request
	req, err := http.NewRequest("POST", config.BaseApi+"/post", bytes.NewBuffer(postData))
	if err != nil {
		respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Failed to create request"}
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Set the session cookie in the request
	req.AddCookie(cookie)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Request failed"}
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintf("error reading response: %v", err),
		}
		return
	}

	// Check response status code
	if resp.StatusCode != http.StatusCreated {
		// Attempt to parse the error message from the response
		var errorResponse map[string]interface{}
		var errorMessage string
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			errorMessage = string(body) // Use raw body as fallback
		} else {
			if errMsg, exists := errorResponse["error"]; exists {
				errorMessage = fmt.Sprintf("%v", errMsg)
			} else {
				errorMessage = "unknown error"
			}
		}

		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintln(errorMessage),
			Status:  resp.StatusCode,
		}
		return
	}

	// Optionally, you can further process the response body if needed
	var responseMessage map[string]interface{}
	if err := json.Unmarshal(body, &responseMessage); err != nil {
		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintf("error unmarshaling response: %v", err),
			Status:  resp.StatusCode,
		}
		return
	}

	// Extracting the message from the response map
	message, ok := responseMessage["message"].(string)
	if !ok {
		message = "Unexpected response format"
	}

	respChan <- models.ResponseDetails{
		Success: true,
		Message: fmt.Sprintln(message), // displays server response to the user
		Status:  resp.StatusCode,
	}
}
