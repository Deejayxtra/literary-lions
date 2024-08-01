# literary-lions

To generate swagger UI
<!-- swag init --dir /home/femsworld/literary-lions/backend/src --output /home/femsworld/literary-lions/backend/src/docs -->
swag init -g cmd/main.go

...././.:~/literary-lions/backend$ swag init --dir ./src/cmd --output ./docs

http://localhost:8080/swagger/index.html



// Package docs contains auto-generated Swagger API documentation.
// To generate or update the documentation, run `swag init` in the project root .
// (~/literary-lions/backend$ swag init -g src/cmd/main.go)
//Then go the backend from ~/literary-lions/backend/src/cmd$ go run .

You can see the swagger UI on: http://localhost:8080/swagger/index.html


**Note:**
Handler => login func ====DO SAME TO LOGIN====  *** Define error message explicitly***
if err := auth.RegisterUser(db, creds.Email, creds.Username, creds.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})           # Define error message explicitly
		return
	}

This is login request from the frontend:
email: admin@mail.com, password: admin123
row: &{<nil> 0xc00012a240}
[GIN] 2024/08/01 - 07:04:13 | 200 |   81.618298ms |       127.0.0.1 | POST     "/login"
[GIN] 2024/08/01 - 07:04:15 | 500 |     166.878µs |       127.0.0.1 | GET      "/api/posts"

Why is the user being redirected to the 