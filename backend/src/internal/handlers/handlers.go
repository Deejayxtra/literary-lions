package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func InitHandlers(database *sql.DB) {
	db = database
}

func Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (email, username, password, role) VALUES (?, ?, ?, 'user')", creds.Email, creds.Username, hashedPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered successfully"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var storedCreds Credentials
	var role string
	err = db.QueryRow("SELECT email, username, password, role FROM users WHERE username=?", creds.Username).Scan(&storedCreds.Email, &storedCreds.Username, &storedCreds.Password, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User logged in successfully"))
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request, *Claims), role string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if role != "" && claims.Role != role {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		endpoint(w, r, claims)
	})
}

func CreatePost(w http.ResponseWriter, r *http.Request, claims *Claims) {
	var post struct {
		CategoryID int    `json:"category_id"`
		Title      string `json:"title"`
		Content    string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO posts (user_id, category_id, title, content) VALUES ((SELECT id FROM users WHERE username=?), ?, ?, ?)", claims.Username, post.CategoryID, post.Title, post.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	postID, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message":  "Post created successfully",
		"post_id":  postID,
		"username": claims.Username,
	}
	json.NewEncoder(w).Encode(response)
}

func UpdatePost(w http.ResponseWriter, r *http.Request, claims *Claims) {
	var post struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	postID := mux.Vars(r)["id"]

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ? AND user_id = (SELECT id FROM users WHERE username = ?)", post.Title, post.Content, postID, claims.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":  "Post updated successfully",
		"post_id":  postID,
		"username": claims.Username,
	}
	json.NewEncoder(w).Encode(response)
}

func DeletePost(w http.ResponseWriter, r *http.Request, claims *Claims) {
	postID := mux.Vars(r)["id"]

	_, err := db.Exec("DELETE FROM posts WHERE id = ? AND user_id = (SELECT id FROM users WHERE username = ?)", postID, claims.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":  "Post deleted successfully",
		"post_id":  postID,
		"username": claims.Username,
	}
	json.NewEncoder(w).Encode(response)
}

func GetPost(w http.ResponseWriter, r *http.Request, claims *Claims) {
	postID := mux.Vars(r)["id"]

	var post struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	err := db.QueryRow("SELECT id, title, content FROM posts WHERE id = ?", postID).Scan(&post.ID, &post.Title, &post.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"id":      post.ID,
		"title":   post.Title,
		"content": post.Content,
	}
	json.NewEncoder(w).Encode(response)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request, claims *Claims) {
	rows, err := db.Query("SELECT id, username, email, role FROM users")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []map[string]interface{}{}
	for rows.Next() {
		var user struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		users = append(users, map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, claims *Claims) {
	userID := mux.Vars(r)["id"]

	_, err := db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "User deleted successfully",
		"user_id": userID,
	}
	json.NewEncoder(w).Encode(response)
}

func UpdateUserRole(w http.ResponseWriter, r *http.Request, claims *Claims) {
	var roleUpdate struct {
		Role string `json:"role"`
	}

	userID := mux.Vars(r)["id"]

	err := json.NewDecoder(r.Body).Decode(&roleUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE users SET role = ? WHERE id = ?", roleUpdate.Role, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "User role updated successfully",
		"user_id": userID,
		"role":    roleUpdate.Role,
	}
	json.NewEncoder(w).Encode(response)
}

func InitAdminUser() error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = db.Exec("INSERT INTO users (email, username, password, role) VALUES ('admin@example.com', 'admin', ?, 'admin')", hashedPassword)
		if err != nil {
			return err
		}
	}
	return nil
}
