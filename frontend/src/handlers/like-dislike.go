package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"sync"
	"literary-lions/frontend/src/config"
	"literary-lions/frontend/src/models"
)

// Method to like posts
func LikePost(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        r.ParseForm()
        postIDStr := r.URL.Query().Get("postID")

		respChan := make(chan models.ResponseDetails, 1)
        var wg sync.WaitGroup

		wg.Add(1)

		// Extract the session cookie from the header
		cookieToken, err := r.Cookie("session_token")
		if err != nil {
			// User must be logged-in to continue
			message := `You are not authorized! Please <a href="/login">login</a> before liking a post.`
			tmpl := template.Must(template.ParseFiles("templates/post.html"))
			tmpl.Execute(w, map[string]interface{}{
				"Error": template.HTML(message),
			})
			return
		}

		// Calls the function that sends request to the server		
		go SendLikeRequest(postIDStr, cookieToken, &wg, respChan)
        go func() {
            wg.Wait()
            close(respChan)
        }()
        
        responseDetails := <-respChan

        if responseDetails.Status == http.StatusOK {
           http.Redirect(w, r, "/post?id="+postIDStr, http.StatusSeeOther)
        } else {
            tmpl := template.Must(template.ParseFiles("templates/post.html"))
            tmpl.Execute(w, map[string]interface{}{
                "Error": responseDetails.Message,
            })
        }
    }
}

// Method to dislike posts
func DislikePost(w http.ResponseWriter, r *http.Request) {
       if r.Method == http.MethodPost {
        r.ParseForm()
        postIDStr := r.URL.Query().Get("postID")

		respChan := make(chan models.ResponseDetails, 1)
        var wg sync.WaitGroup

		wg.Add(1)

		// Extract the session cookie from the header
		cookieToken, err := r.Cookie("session_token")
		if err != nil {
			// User must be logged-in to continue
			message := `You are not authorized! Please <a href="/login">login</a> before liking a post.`
			tmpl := template.Must(template.ParseFiles("templates/post.html"))
			tmpl.Execute(w, map[string]interface{}{
				"Error": template.HTML(message),
			})
			return
		}

		// Calls the function that sends request to the server		
		go SendDislikeRequest(postIDStr, cookieToken, &wg, respChan)
        go func() {
            wg.Wait()
            close(respChan)
        }()
        
        responseDetails := <-respChan

        if responseDetails.Status == http.StatusOK {
            http.Redirect(w, r, "/post?id="+postIDStr, http.StatusSeeOther)
        } else {
            tmpl := template.Must(template.ParseFiles("templates/post.html"))
            tmpl.Execute(w, map[string]interface{}{
                "Error": responseDetails.Message,
            })
        }
    }
}


// Function to send the like posts request to the backend
func SendLikeRequest(id string, cookie *http.Cookie, waitGroup *sync.WaitGroup, respChan chan models.ResponseDetails) {
    defer waitGroup.Done()
    // Creates request to send the the backend 
    req, err := http.NewRequest("POST", config.BaseApi+"/post/"+id+"/like", nil) 
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Failed to create request"}
        return
    }
	req.AddCookie(cookie)	// adding cookies to the request

    // Sends request to the server
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Request failed"}
        return
    }
    defer resp.Body.Close()

    // Reads the response body 
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error reading response: %v", err),
        }
        return
    }

    // Checking status if OK
    if resp.StatusCode != http.StatusOK {
        var errorResponse map[string]interface{}
        var errorMessage string
        if err := json.Unmarshal(body, &errorResponse); err != nil {
            errorMessage = string(body)
        } else {
            if errMsg, exists := errorResponse["error"]; exists {
                errorMessage = fmt.Sprintf("%v", errMsg)
            } else {
                errorMessage = "unknown error"
            }
        }

        respChan <- models.ResponseDetails{
            Success: false,
            Message: errorMessage,
            Status:  resp.StatusCode,
        }
        return
    }

    // Unmarshals data in JSON format
    var responseMessage map[string]interface{}
    if err := json.Unmarshal(body, &responseMessage); err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error unmarshaling response: %v", err),
            Status:  resp.StatusCode,
        }
        return
    }

    // Checks if response message is OK
    message, ok := responseMessage["message"].(string)
    if !ok {
        message = "Unexpected response format"
    }

    respChan <- models.ResponseDetails{
        Success: true,
        Message: message,
        Status:  resp.StatusCode,
    }
}

// Function to send the dislike posts request to the backend
func SendDislikeRequest(id string, cookie *http.Cookie, waitGroup *sync.WaitGroup, respChan chan models.ResponseDetails) {
    defer waitGroup.Done()
    // Creates request to send the the backend 
	req, err := http.NewRequest("POST", config.BaseApi+"/post/"+id+"/dislike", nil) 
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Failed to create request"}
        return
    }
	req.AddCookie(cookie)	// adding cookies to the request

    // Sends request to the server
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Request failed"}
        return
    }
    defer resp.Body.Close()

    // Reads the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error reading response: %v", err),
        }
        return
    }

    // Checking status if OK
    if resp.StatusCode != http.StatusOK {
        var errorResponse map[string]interface{}
        var errorMessage string
        if err := json.Unmarshal(body, &errorResponse); err != nil {
            errorMessage = string(body)
        } else {
            if errMsg, exists := errorResponse["error"]; exists {
                errorMessage = fmt.Sprintf("%v", errMsg)
            } else {
                errorMessage = "unknown error"
            }
        }

        respChan <- models.ResponseDetails{
            Success: false,
            Message: errorMessage,
            Status:  resp.StatusCode,
        }
        return
    }

    // Unmarshals data in JSON format
    var responseMessage map[string]interface{}
    if err := json.Unmarshal(body, &responseMessage); err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error unmarshaling response: %v", err),
            Status:  resp.StatusCode,
        }
        return
    }

    // Checks the response message if OK
    message, ok := responseMessage["message"].(string)
    if !ok {
        message = "Unexpected response format"
    }

    respChan <- models.ResponseDetails{
        Success: true,
        Message: message,
        Status:  resp.StatusCode,
    }
}

// Method to like comments
func LikeComment(w http.ResponseWriter, r *http.Request) {
    // Method for POST request to like comments
    if r.Method == http.MethodPost {
        r.ParseForm()
        commentIDStr := r.URL.Query().Get("commentID")
        postIDStr := r.URL.Query().Get("postID")

		respChan := make(chan models.ResponseDetails, 1)
        var wg sync.WaitGroup

		wg.Add(1)

		// Extract the session cookie from the header
		cookieToken, err := r.Cookie("session_token")
		if err != nil {
			// User must be logged-in to continue
			message := `You are not authorized! Please <a href="/login">login</a> before liking a post.`
			tmpl := template.Must(template.ParseFiles("templates/post.html"))
			tmpl.Execute(w, map[string]interface{}{
				"Error": template.HTML(message),
			})
			return
		}

		// Calls the function that sends request to the server
		go SendLikeRequestComment(commentIDStr, cookieToken, &wg, respChan)
        go func() {
            wg.Wait()
            close(respChan)
        }()
        
        responseDetails := <-respChan
        // Checks the status if OK
        if responseDetails.Status == http.StatusOK {
           http.Redirect(w, r, "/post?id="+postIDStr+"/#comment-"+commentIDStr, http.StatusSeeOther)
        } else {
            tmpl := template.Must(template.ParseFiles("templates/post.html"))
            tmpl.Execute(w, map[string]interface{}{
                "Error": responseDetails.Message,
            })
        }
    }
}

// Method to dislike comments
func DislikeComment(w http.ResponseWriter, r *http.Request) {  
    // Method for POST request to dislike comments
    if r.Method == http.MethodPost {
        r.ParseForm()
        commentIDStr := r.URL.Query().Get("commentID")
        postIDStr := r.URL.Query().Get("postID")

		respChan := make(chan models.ResponseDetails, 1)
        var wg sync.WaitGroup

		wg.Add(1)

		// Extract the session cookie from the header
		cookieToken, err := r.Cookie("session_token")
		if err != nil {
			// User must be logged-in to continue
			message := `You are not authorized! Please <a href="/login">login</a> before liking a post.`
			tmpl := template.Must(template.ParseFiles("templates/post.html"))
			tmpl.Execute(w, map[string]interface{}{
				"Error": template.HTML(message),
			})
			return
		}

		// Calls the function that sends request to the server		
		go SendDislikeRequestComment(commentIDStr, cookieToken, &wg, respChan)
        go func() {
            wg.Wait()
            close(respChan)
        }()
        
        responseDetails := <-respChan

        // Checks the status if OK
        if responseDetails.Status == http.StatusOK {
            http.Redirect(w, r, "/post?id="+postIDStr+"/#comment-"+commentIDStr, http.StatusSeeOther)
        } else {
            tmpl := template.Must(template.ParseFiles("templates/post.html"))
            tmpl.Execute(w, map[string]interface{}{
                "Error": responseDetails.Message,
            })
        }
    }
}

// Function to send the like comments request to the backend
func SendLikeRequestComment(id string, cookie *http.Cookie, waitGroup *sync.WaitGroup, respChan chan models.ResponseDetails) {
    defer waitGroup.Done()
    // Creates request to send the the backend 
    req, err := http.NewRequest("POST", config.BaseApi+"/comment/"+id+"/like", nil) 
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Failed to create request"}
        return
    }
	req.AddCookie(cookie)	// adding cookies to the request

    // Sends request to the server
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Request failed"}
        return
    }
    defer resp.Body.Close()

    // Reads the response body 
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error reading response: %v", err),
        }
        return
    }

    // Checks the status if OK
    if resp.StatusCode != http.StatusOK {
        var errorResponse map[string]interface{}
        var errorMessage string
        if err := json.Unmarshal(body, &errorResponse); err != nil {
            errorMessage = string(body)
        } else {
            if errMsg, exists := errorResponse["error"]; exists {
                errorMessage = fmt.Sprintf("%v", errMsg)
            } else {
                errorMessage = "unknown error"
            }
        }

        respChan <- models.ResponseDetails{
            Success: false,
            Message: errorMessage,
            Status:  resp.StatusCode,
        }
        return
    }

    var responseMessage map[string]interface{}
    // Unmarshals data in JSON format
    if err := json.Unmarshal(body, &responseMessage); err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error unmarshaling response: %v", err),
            Status:  resp.StatusCode,
        }
        return
    }

    // Check the response message if OK
    message, ok := responseMessage["message"].(string)
    if !ok {
        message = "Unexpected response format"
    }

    respChan <- models.ResponseDetails{
        Success: true,
        Message: message,
        Status:  resp.StatusCode,
    }
}

// Function to send the dislike comments request to the backend
func SendDislikeRequestComment(id string, cookie *http.Cookie, waitGroup *sync.WaitGroup, respChan chan models.ResponseDetails) {
    defer waitGroup.Done()
    // Creates request to send the the backend 
	req, err := http.NewRequest("POST", config.BaseApi+"/comment/"+id+"/dislike", nil) 
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Failed to create request"}
        return
    }
	req.AddCookie(cookie)	// adding cookies to the request

    // Sends request to the server
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Request failed"}
        return
    }
    defer resp.Body.Close()

    // Reads the response body 
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error reading response: %v", err),
        }
        return
    }

    // Checks the status if OK
    if resp.StatusCode != http.StatusOK {
        var errorResponse map[string]interface{}
        var errorMessage string
        if err := json.Unmarshal(body, &errorResponse); err != nil {
            errorMessage = string(body)
        } else {
            if errMsg, exists := errorResponse["error"]; exists {
                errorMessage = fmt.Sprintf("%v", errMsg)
            } else {
                errorMessage = "unknown error"
            }
        }

        respChan <- models.ResponseDetails{
            Success: false,
            Message: errorMessage,
            Status:  resp.StatusCode,
        }
        return
    }

    var responseMessage map[string]interface{}
    // Unmarshals data in JSON format
    if err := json.Unmarshal(body, &responseMessage); err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error unmarshaling response: %v", err),
            Status:  resp.StatusCode,
        }
        return
    }

    // Check the response message if OK
    message, ok := responseMessage["message"].(string)
    if !ok {
        message = "Unexpected response format"
    }

    respChan <- models.ResponseDetails{
        Success: true,
        Message: message,
        Status:  resp.StatusCode,
    }
}
