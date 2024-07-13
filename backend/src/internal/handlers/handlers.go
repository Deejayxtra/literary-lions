package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
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

	result, err := db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, 'user')", creds.Username, hashedPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"User registered successfully", "user_id":` + strconv.FormatInt(userID, 10) + `}`))
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
	err = db.QueryRow("SELECT username, password, role FROM users WHERE username=?", creds.Username).Scan(&storedCreds.Username, &storedCreds.Password, &role)
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
	w.Write([]byte(`{"message":"User logged in successfully", "token":"` + tokenString + `"}`))
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
	w.Write([]byte(`{"message":"Post created successfully", "post_id":` + strconv.FormatInt(postID, 10) + `}`))
}

func GetPosts(w http.ResponseWriter, r *http.Request, claims *Claims) {
	posts := []map[string]interface{}{}
	rows, err := db.Query("SELECT id, title, content FROM posts WHERE category_id = ?", mux.Vars(r)["category_id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		posts = append(posts, map[string]interface{}{"id": post.ID, "title": post.Title, "content": post.Content})
	}

	json.NewEncoder(w).Encode(posts)
}

func GetPost(w http.ResponseWriter, r *http.Request, claims *Claims) {
	var post struct {
		ID         int    `json:"id"`
		Title      string `json:"title"`
		Content    string `json:"content"`
		UserID     int    `json:"user_id"`
		CategoryID int    `json:"category_id"`
	}

	err := db.QueryRow("SELECT id, title, content, user_id, category_id FROM posts WHERE id = ?", mux.Vars(r)["id"]).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func UpdatePost(w http.ResponseWriter, r *http.Request, claims *Claims) {
	var post struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ? AND user_id = (SELECT id FROM users WHERE username = ?)", post.Title, post.Content, mux.Vars(r)["id"], claims.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Post updated successfully", "post_id":` + mux.Vars(r)["id"] + `}`))
}

func GetAllUsers(w http.ResponseWriter, r *http.Request, claims *Claims) {
	users := []map[string]interface{}{}
	rows, err := db.Query("SELECT id, username, email, role FROM users")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		users = append(users, map[string]interface{}{"id": user.ID, "username": user.Username, "email": user.Email, "role": user.Role})
	}

	json.NewEncoder(w).Encode(users)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, claims *Claims) {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"User deleted successfully", "user_id":` + mux.Vars(r)["id"] + `}`))
}

func DeletePost(w http.ResponseWriter, r *http.Request, claims *Claims) {
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Post deleted successfully", "post_id":` + mux.Vars(r)["id"] + `}`))
}
