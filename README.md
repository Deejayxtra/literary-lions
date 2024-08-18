# literary-lions

The Literary Lions Application is a web-based platform that allows users to create, read, update, and delete posts, as well as like, dislike, and comment on posts. The application is divided into two main components: the backend, built using Go and Gin, and the frontend, built using Go's net/http package.

## Features

### Backend
- **User Authentication**: Register, login, and logout functionalities.
- **Post Management**: Create, read, update, delete posts.
- **Comment Management**: Add comments to posts.
- **Like/Dislike System**: Users can like or dislike posts and comments.
- **Category Filtering**: View posts filtered by category.
- **Swagger Documentation**: API documentation available via Swagger UI.

### Frontend
- **User Interface**: HTML templates served via the Go net/http package.
- **CRUD Operations**: Interface for creating, viewing, updating, and deleting posts.
- **User Interaction**: Like, dislike, and comment on posts.
- **Profile Management**: Update user profiles.

## Prerequisites

- Go 1.21.5 or later
- Docker
- Docker Compose

## Running the Application from the Terminal

### Backend

1. Navigate to the Backend Directory:
2. Run the Backend Server:
```bash
   cd backend/src/cmd
	go run main.go
```
The backend server will start on http://localhost:8080.

### Frontend

1. Navigate to the frontend Directory:
2. Run the Frontend Server:
```bash
   cd frontend/src
	go run main.go
```
The frontend server will start on http://localhost:8000.

## Building Docker Images

### Backend

1. Navigate to the Backend Directory:
2. Build the Docker Image:
```bash
   cd backend
	docker build -t literary-lions-backend .
```
After building the docker image, to run the docker image, from the same terminal:
```bash
	docker run -d --name literary-lions-backend -p 8080:8080 literary-lions-backend
```
The backend server will start on http://localhost:8080.

### Frontend

1. Navigate to the Backend Directory:
2. Build the Docker Image:
```bash
   cd backend
	docker build -t literary-lions-frontend .
```
After building the docker image, to run the docker image, from the same terminal:
```bash
	ddocker run -d --name literary-lions-frontend -p 8000:8000 literary-lions-frontend
```
The frontend server will start on http://localhost:8000.

## Running the Application with Docker Compose

1. Navigate to the Project Root Directory:
2. Run the Docker Compose Command:

```bash
   cd ~/literary-lions
	docker-compose up --build .
```
This will build and start both the frontend and backend containers.
3. Access the Application:

- Frontend: Open http://localhost:8000 in your browser.
- Backend API: API endpoints are available at http://localhost:8080/api/v1.0. Swagger documentation is available at http://localhost:8080/swagger/index.html.

## Application Structure

## Backend

- cmd/main.go: Entry point for the backend application.
- internal/: Contains the core logic for handlers, models, and middleware.
- config/: Configuration files.
- docs/: Swagger API documentation files.
- literary_lions.db: SQLite database file.

## Frontend

- src/main.go: Entry point for the frontend application.
- src/handlers/: Contains handlers for different routes.
- src/templates/: HTML templates for the frontend.
- src/static/: Static files like CSS and images.


## Explanation of the Sections

1. **Features**: Provides a high-level overview of what the application does.
2. **Prerequisites**: Lists the tools needed before running or building the application.
3. **Running the Application from the Terminal**: Provides step-by-step instructions to run both the frontend and backend components without Docker.
4. **Building Docker Images**: Guides on how to build Docker images for both the frontend and backend individually.
5. **Running the Application with Docker Compose**: Shows how to use Docker Compose to build and run the entire application.
6. **Application Structure**: Provides a brief overview of the directory structure and the purpose of the main files and directories.
7. **Contributing**: An open invitation for others to contribute to the project.
8. **License**: Information about the project's licensing.

This README should serve as a comprehensive guide for both developers and users of the Literary Lions Application.


