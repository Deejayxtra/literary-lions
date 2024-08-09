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


// Function to send the like/dislike request to the backend
func SendLikeRequest(id string, cookie *http.Cookie, waitGroup *sync.WaitGroup, respChan chan models.ResponseDetails) {
    defer waitGroup.Done()
    req, err := http.NewRequest("POST", config.BaseApi+"/post/"+id+"/like", nil) 
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Failed to create request"}
        return
    }
	req.AddCookie(cookie)	// adding cookies to the request

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Request failed"}
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error reading response: %v", err),
        }
        return
    }

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
    if err := json.Unmarshal(body, &responseMessage); err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error unmarshaling response: %v", err),
            Status:  resp.StatusCode,
        }
        return
    }

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


func SendDislikeRequest(id string, cookie *http.Cookie, waitGroup *sync.WaitGroup, respChan chan models.ResponseDetails) {
    defer waitGroup.Done()

	    req, err := http.NewRequest("POST", config.BaseApi+"/post/"+id+"/dislike", nil) 
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Failed to create request"}
        return
    }
	req.AddCookie(cookie)	// adding cookies to the request

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Request failed"}
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error reading response: %v", err),
        }
        return
    }

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
    if err := json.Unmarshal(body, &responseMessage); err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error unmarshaling response: %v", err),
            Status:  resp.StatusCode,
        }
        return
    }

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


func LikeComment(w http.ResponseWriter, r *http.Request) {
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
		
		go SendLikeRequestComment(commentIDStr, cookieToken, &wg, respChan)
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

func DislikeComment(w http.ResponseWriter, r *http.Request) {
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
		
		go SendDislikeRequestComment(commentIDStr, cookieToken, &wg, respChan)
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

// Function to send the like/dislike request to the backend
func SendLikeRequestComment(id string, cookie *http.Cookie, waitGroup *sync.WaitGroup, respChan chan models.ResponseDetails) {
    defer waitGroup.Done()
    req, err := http.NewRequest("POST", config.BaseApi+"/comment/"+id+"/like", nil) 
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Failed to create request"}
        return
    }
	req.AddCookie(cookie)	// adding cookies to the request

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Request failed"}
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error reading response: %v", err),
        }
        return
    }

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
    if err := json.Unmarshal(body, &responseMessage); err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error unmarshaling response: %v", err),
            Status:  resp.StatusCode,
        }
        return
    }

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

func SendDislikeRequestComment(id string, cookie *http.Cookie, waitGroup *sync.WaitGroup, respChan chan models.ResponseDetails) {
    defer waitGroup.Done()

	    req, err := http.NewRequest("POST", config.BaseApi+"/comment/"+id+"/dislike", nil) 
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Failed to create request"}
        return
    }
	req.AddCookie(cookie)	// adding cookies to the request

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Request failed"}
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error reading response: %v", err),
        }
        return
    }

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
    if err := json.Unmarshal(body, &responseMessage); err != nil {
        respChan <- models.ResponseDetails{
            Success: false,
            Message: fmt.Sprintf("error unmarshaling response: %v", err),
            Status:  resp.StatusCode,
        }
        return
    }

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
