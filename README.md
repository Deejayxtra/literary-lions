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





Short guide on how to build and run docker containers

Build the Docker Images
You'll need to open your terminal and navigate to the directory where each Dockerfile is located to build the Docker images.

1. Navigate to Backend Directory and Build Backend Image
Open your terminal and navigate to the backend directory:

cd ~/literary-lions/backend/

Now, build the Docker image for the backend:

docker build -t literary-lions-backend .

-t literary-lions-backend: This tags your image with the name literary-lions-backend.
.: This represents the current directory, which contains the Dockerfile.

2. Navigate to Frontend Directory and Build Frontend Image
Next, navigate to the frontend directory:

cd ~/literary-lions/frontend

Build the Docker image for the frontend:

docker build -t literary-lions-frontend .


-t literary-lions-frontend: This tags your image with the name literary-lions-frontend.
.: Again, this represents the current directory.



Step 2: Run the Docker Containers
Now that you've built the images, you can run the containers.

1. Run the Backend Container
In the terminal, run the backend container:

docker run -d --name literary-lions-backend -p 8080:8080 literary-lions-backend


-d: Runs the container in detached mode (in the background).
--name literary-lions-backend: Names the container literary-lions-backend

(../../literary-lions$ docker run -d --name literary-lions-backend -p 8080:8080 literary-lions-backend)
-p 8080:8080: Maps port 8080 of the container to port 8080 on your local machine.
literary-lions-backend: This is the name of the Docker image you built.


2. Run the Frontend Container
In another terminal, run the frontend container:

docker run -d --name literary-lions-frontend -p 8000:8000 literary-lions-frontend


-d: Runs the container in detached mode.
--name literary-lions-frontend: Names the container literary-lions-frontend.
-p 8000:8000: Maps port 8000 of the container to port 8000 on your local machine.
literary-lions-frontend: This is the name of the Docker image you built.


Step 3: Verify that the Containers are Running
To check if your containers are running, use:

docker ps


This will list all running containers. You should see both literary-lions-backend and literary-lions-frontend in the list.

Step 4: Access Your Application
Now you can access your application using your web browser:

Backend API Documentation (Swagger UI):
Open http://localhost:8080/swagger/index.html to see the Swagger UI for your backend API.


Frontend Application:
Open http://localhost:8000 to see the frontend of your application.
Step 5: Managing Containers


Stopping Containers
To stop the containers, use:

docker stop literary-lions-backend literary-lions-frontend


Removing Containers
If you want to remove the containers after stopping them, use:

docker rm literary-lions-backend literary-lions-frontend


Summary of Commands
Build Backend Image: docker build -t literary-lions-backend .
Build Frontend Image: docker build -t literary-lions-frontend .
Run Backend Container: docker run -d --name literary-lions-backend -p 8080:8080 literary-lions-backend
Run Frontend Container: docker run -d --name literary-lions-frontend -p 8000:8000 literary-lions-frontend
Check Running Containers: docker ps
Stop Containers: docker stop literary-lions-backend literary-lions-frontend
Remove Containers: docker rm literary-lions-backend literary-lions-frontend