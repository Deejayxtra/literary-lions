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
	"net/url"
	"strings"
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
    keyword := r.URL.Query().Get("keyword")
    category := r.URL.Query().Get("category")
    startDate := r.URL.Query().Get("start_date")
    endDate := r.URL.Query().Get("end_date")
    filter := r.URL.Query().Get("filter")

    var filterIsSet bool
    var cookie *http.Cookie

    // Construct the API request URL with query parameters
    apiURL := config.BaseApi + "/posts?"
    filteredURL := config.BaseApi + "/filtered-posts?"

    if keyword != "" {
        apiURL += "keyword=" + url.QueryEscape(keyword) + "&"
    }
    if category != "" {
        apiURL += "category=" + url.QueryEscape(category) + "&"
    }
    if startDate != "" {
        apiURL += "start_date=" + url.QueryEscape(startDate) + "&"
    }
    if endDate != "" {
        apiURL += "end_date=" + url.QueryEscape(endDate) + "&"
    }
    if filter != "" {
        filterIsSet = true
        filteredURL += "filter=" + url.QueryEscape(filter) + "&"
        // Extract the session cookie from the header
        var err error
        cookie, err = r.Cookie("session_token")
        if err != nil {
            message := `You are not authorized! Please <a href="/login">login</a> before checking My posts.`
            tmpl := template.Must(template.ParseFiles("templates/index.html"))
            tmpl.Execute(w, map[string]interface{}{
                "Error": template.HTML(message),
            })
            return
        }
    }

    respChan := make(chan models.Data, 1)
    var wg sync.WaitGroup

    wg.Add(1)
    go func() {
        defer wg.Done()
        if filterIsSet {
            SendShowPostsRequestFilter(cookie, filteredURL, respChan)
        } else {
            SendShowPostsRequest(apiURL, respChan)
        }
    }()

    wg.Wait()
    close(respChan)

    response := <-respChan

    if !response.Success {
        handleErrorResponse(w, response)
        return
    }

    posts := response.Posts

    // Handle no posts found
    if len(posts) == 0 {
        currentUser, authenticated := isAuthenticated(r)
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

    // Truncate content if necessary
    for i := range posts {
        posts[i].Content = truncateContent(posts[i].Content, 150)
    }

    currentUser, authenticated := isAuthenticated(r)

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

    RenderTemplate(w, "index.html", data)
}


// Helper function to handle errors and render templates
func handleErrorResponse(w http.ResponseWriter, response models.Data) {
    var tmpl *template.Template
    var err error

    // Load the template
    tmpl, err = template.ParseFiles("templates/index.html")
    if err != nil {
        log.Printf("Template parsing error: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    if response.Status == http.StatusUnauthorized {
        response.Message = `You are not authorized! Please <a href="/login">login</a> before creating a post.`
        tmpl.Execute(w, map[string]interface{}{
            "Error": template.HTML(response.Message),
        })
    } else {
        response.Message = "Oops! Something went wrong. Failed to create post."
        tmpl.Execute(w, map[string]interface{}{
            "Error": response.Message,
        })
    }
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
		return
	}

	// Use an http.Client to make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		message := "Failed to fetch post"
		StatusInternalServerError(w, message)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		message := "Failed to read response"
		StatusInternalServerError(w, message)
		return
	}

	// Parse the JSON response into a PostDetails model
	var response models.PostDetails
	err = json.Unmarshal(body, &response)
	if err != nil {
		message := "Failed to parse response"
		StatusInternalServerError(w, message)
		return
	}

	// Format the created_at date
	formattedDate := response.Post.CreatedAt.Format("January 2, 2006 at 3:04pm")

	// Get the authentication status and the currentUser if any
	currentUser, authenticated := isAuthenticated(r)
	response.Post.FormattedContent = strings.Split(response.Post.Content, "\n")

	data := struct {
		Post          models.Post
		FormattedDate string
		Authenticated bool
		Comments      []models.Comment
		Error         bool
		Username      string
		Likes         int
		Dislikes      int
	}{
		Post:          response.Post,
		FormattedDate: formattedDate,
		Authenticated: authenticated,
		Comments:      response.Comments,
		Error:         false,
		Username:      currentUser,
		Likes:         response.Likes,
		Dislikes:      response.Dislikes,
	}
	// Render the template with posts and authentication status
	RenderTemplate(w, "post.html", data)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	// Get the authentication status and the currentUser if any
	currentUser, authenticated := isAuthenticated(r)

	// Render the page to the user to select from category
	if r.Method == http.MethodGet {
		// Hardcoded categories for now. Might be nice to add feature for creating it and make it dynamic later
		categories := []string{"Random", "News", "Sport", "Technology", "Science", "Health"}
		data := struct {
			Categories    []string
			Error         interface{}
			Authenticated bool
			Username      string
		}{
			Categories:    categories,
			Error:         nil,
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

		// Defines the payload sent to the backend
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
			responseDetails.Message = `You are not authorized! Please <a href="/login">login</a> before creating a post.`
			tmpl := template.Must(template.ParseFiles("templates/create-post.html"))
			tmpl.Execute(w, map[string]interface{}{
				"Error": template.HTML(responseDetails.Message),
			})
		} else {
			responseDetails.Message = "Oops! Something went wrong. Failed to create post."
			tmpl := template.Must(template.ParseFiles("templates/create-post.html"))
			tmpl.Execute(w, map[string]interface{}{
				"Error": responseDetails.Message,
			})
		}
	}
}

func SendCreatePostRequest(cookie *http.Cookie, payload models.Post, waitGroup *sync.WaitGroup, respChan chan models.ResponseDetails) {
	defer waitGroup.Done() // Ensure the channel is closed once this function completes

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

func SendGetPostByIdRequest(id string, w http.ResponseWriter, r *http.Request, waitGroup *sync.WaitGroup, respChan chan models.PostDetails) {
	defer waitGroup.Done()

	// Create a new GET request
	req, err := http.NewRequest("GET", config.BaseApi+"/post/"+id, nil)
	if err != nil {
		message := "Failed to create request"
		StatusInternalServerError(w, message)
		return
	}

	// Use an http.Client to make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		message := "Failed to fetch post"
		StatusInternalServerError(w, message)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		message := "Failed to read response"
		StatusInternalServerError(w, message)
		return
	}

	// Parse the JSON response into a PostDetails model
	var response models.PostDetails
	err = json.Unmarshal(body, &response)
	if err != nil {
		message := "Failed to parse response"
		StatusInternalServerError(w, message)
		return
	}

	respChan <- models.PostDetails{
		Post:     response.Post,
		Comments: response.Comments,
		Likes:    response.Likes,
		Dislikes: response.Dislikes,
		Status:   resp.StatusCode,
	}
}


func SendShowPostsRequest(apiURL string, respChan chan models.Data) {

    req, err := http.NewRequest("GET", apiURL, nil)
    if err != nil {
        respChan <- models.Data{
            Success: false,
            Message: fmt.Sprintf("error creating request: %v", err),
        }
        return
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        respChan <- models.Data{
            Success: false,
            Message: fmt.Sprintf("Failed to fetch post: %v", err),
        }
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        respChan <- models.Data{
            Success: false,
            Message: fmt.Sprintf("Failed to read response: %v", err),
        }
        return
    }

    // Check response status code
    if resp.StatusCode != http.StatusOK {
        respChan <- models.Data{
            Success: false,
            Message: fmt.Sprintf("Unexpected status code: %d", resp.StatusCode),
        }
        return
    }

    // Unmarshal the response directly into a slice of posts
    var posts []models.Post
    if err := json.Unmarshal(body, &posts); err != nil {
        respChan <- models.Data{
            Success: false,
            Message: fmt.Sprintf("error unmarshaling response: %v", err),
        }
        return
    }

    // Successfully return the posts
    respChan <- models.Data{
        Success: true,
        Posts:    posts,
        Message: "Posts fetched successfully",
    }
}


func SendShowPostsRequestFilter(cookie *http.Cookie, apiURL string, respChan chan models.Data) {
    // Create a new GET request
    req, err := http.NewRequest("GET", apiURL, nil)
	// req, err := http.NewRequest("GET", config.BaseApi+"/filtered-posts", nil)
    if err != nil {
        respChan <- models.Data{
            Success: false,
            Message: fmt.Sprintf("Error creating request: %v", err),
        }
        return
    }

    // Add the session cookie to the request
    req.AddCookie(cookie)

    // Use an http.Client to make the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        respChan <- models.Data{
            Success: false,
            Message: fmt.Sprintf("Failed to fetch post: %v", err),
        }
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        respChan <- models.Data{
            Success: false,
            Message: fmt.Sprintf("Failed to read response: %v", err),
        }
        return
    }

    // Check response status code
    if resp.StatusCode != http.StatusOK {
        respChan <- models.Data{
            Success: false,
            Message: fmt.Sprintf("Unexpected status code: %d", resp.StatusCode),
        }
        return
    }

    // Unmarshal the response directly into a slice of posts
    var posts []models.Post
    if err := json.Unmarshal(body, &posts); err != nil {
        respChan <- models.Data{
            Success: false,
            Message: fmt.Sprintf("Error unmarshaling response: %v", err),
        }
        return
    }

    // Successfully return the posts
    respChan <- models.Data{
        Success: true,
        Posts:   posts,
        Message: "Posts fetched successfully",
    }
}
