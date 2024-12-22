package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"forum/database"
	"forum/utils"
)

var templates *template.Template

func init() {
	// Parse all templates in one go
	templates = template.Must(template.ParseGlob("templates/*.html"))

	// Print actual template names
	for _, t := range templates.Templates() {
		fmt.Println("Parsed template name:", t.Name())
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect to login if accessing root directly
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginHandler invoked")

	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"Title": "Login - Forum",
		}
		templates.ExecuteTemplate(w, "login", data)
		return
	}
	// List all available templates
	// fmt.Println("Available templates:", templates.Templates())
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var id int
		var hashedPassword string

		query := `SELECT id, password FROM users WHERE username = ?`
		err := database.DB.QueryRow(query, username).Scan(&id, &hashedPassword)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		err = utils.CheckPassword(hashedPassword, password)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		sessionToken, err := utils.GenerateSessionToken()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		expiresAt := time.Now().Add(24 * time.Hour)
		query = `INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)`
		_, err = database.DB.Exec(query, id, sessionToken, expiresAt)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Expires:  expiresAt,
			HttpOnly: true,
		})

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SignupHandler invoked")

	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"Title": "Signup - Forum",
		}
		err := templates.ExecuteTemplate(w, "signup", data)
		if err != nil {
			fmt.Printf("Error executing signup template: %v\n", err)
			http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Validate password strength
		err := utils.ValidatePasswordStrength(password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		query := `INSERT INTO users (email, username, password) VALUES (?, ?, ?)`
		_, err = database.DB.Exec(query, email, username, hashedPassword)
		if err != nil {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SignupHandler invoked")

	// List all available templates
	// fmt.Println("Available templates:", templates.Templates())

	userID := r.Context().Value("user_id").(int) // Retrieve the logged-in user's ID from the context

	// var username string
	// query := `SELECT username FROM users WHERE id = ?`
	// err := database.DB.QueryRow(query, userID).Scan(&username)
	// if err != nil {
	//     http.Error(w, "User not found", http.StatusInternalServerError)
	//     return
	// }
	posts, err := database.GetAllPosts("")
	if err != nil {
		http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
		return
	}

	// Retrieve the username for each post's userID
	for i := range posts {
		// Query the users table for the username using the userID of each post
		var username string
		query := `SELECT username FROM users WHERE id = ?`
		err := database.DB.QueryRow(query, posts[i].UserID).Scan(&username)
		if err != nil {
			// If there's an error retrieving the username, skip adding the username
			posts[i].Username = "Unknown"
		} else {
			posts[i].Username = username
		}
	}

	// Retrieve category post counts
	categoryCounts, err := database.GetCategoryPostCounts()
	if err != nil {
		http.Error(w, "Error retrieving category counts", http.StatusInternalServerError)
		return
	}

	// Pass the posts and logged-in user's username to the template
	data := map[string]interface{}{
		"Title":      "Dashboard - Forum",
		"IsLoggedIn": true,
		"Username":   getUsername(userID), // Utility function to fetch username
		"Posts":      posts,
		"Categories": categoryCounts, // Make sure categoryCounts is used here
		"UserID":     userID,
	}

	err = templates.ExecuteTemplate(w, "dashboard", data)
	if err != nil {
		fmt.Printf("Error executing dashboard template: %v\n", err)
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getUsername(userID int) string {
	var username string
	query := `SELECT username FROM users WHERE id = ?`
	err := database.DB.QueryRow(query, userID).Scan(&username)
	if err != nil {
		return ""
	}
	return username
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var userID int
		query := `SELECT user_id FROM sessions WHERE session_token = ? AND expires_at > ?`
		err = database.DB.QueryRow(query, cookie.Value, time.Now()).Scan(&userID)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		context := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(context))
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	query := `DELETE FROM sessions 	WHERE session_token = ?`
	_, err = database.DB.Exec(query, cookie.Value)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// CreatePostHandler handles creating new posts
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"Title":      "Create Post - Forum",
			"IsLoggedIn": true,
		}
		err := templates.ExecuteTemplate(w, "create-post", data)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		userID := r.Context().Value("user_id").(int)
		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")

		err := database.CreatePost(userID, title, content, category)
		if err != nil {
			http.Error(w, "Error creating post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// ListPostsHandler displays all posts
func ListPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := database.GetAllPosts("")
	if err != nil {
		http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":      "All Posts - Forum",
		"IsLoggedIn": true,
		"Posts":      posts,
	}

	err = templates.ExecuteTemplate(w, "post-list", data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// ViewPostHandler displays a specific post
func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id") // Example: /view-post?id=123

	// Convert postID from string to int
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	post, err := database.GetPost(postID)
	if err != nil || post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	comments, err := database.GetCommentsByPostID(postID)
	if err != nil {
		http.Error(w, "Error retrieving comments", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":      post.Title + " - Forum",
		"IsLoggedIn": true,
		"Post":       post,
		"Comments":   comments,
	}

	err = templates.ExecuteTemplate(w, "view-post", data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// EditPostHandler handles editing an existing post
func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id")

	// Convert postID from string to int
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		post, err := database.GetPost(postID)
		if err != nil || post == nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		// Debugging: Log the retrieved post
		fmt.Printf("Post data: %+v\n", post)

		data := map[string]interface{}{
			"Title":      "Edit Post - Forum",
			"IsLoggedIn": true,
			"Post":       post,
		}

		// Debugging: Log the template data
		fmt.Printf("Template data: %+v\n", data)

		err = templates.ExecuteTemplate(w, "edit-post", data)
		if err != nil {
			fmt.Printf("Error executing edit-post template: %v\n", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
		return
	}

	// Handle post updates (if using POST)
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")

		err := database.UpdatePost(postID, title, content, category)
		if err != nil {
			http.Error(w, "Error updating post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// DeletePostHandler handles deleting a post
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id")

	// Convert postID from string to int
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	err = database.DeletePost(postID)
	if err != nil {
		http.Error(w, "Error deleting post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// LikePostHandler handles liking a post.
func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	postIDStr := r.URL.Query().Get("id")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	err = database.AddReaction(userID, postID, "like")
	if err != nil {
		http.Error(w, "Failed to like post", http.StatusInternalServerError)
		return
	}

	// Return updated reaction counts
	likes, dislikes, _ := database.GetReactionCounts(postID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"likes": %d, "dislikes": %d}`, likes, dislikes)))
}

// DislikePostHandler handles disliking a post.
func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	postIDStr := r.URL.Query().Get("id")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	err = database.AddReaction(userID, postID, "dislike")
	if err != nil {
		http.Error(w, "Failed to dislike post", http.StatusInternalServerError)
		return
	}

	// Return updated reaction counts
	likes, dislikes, _ := database.GetReactionCounts(postID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"likes": %d, "dislikes": %d}`, likes, dislikes)))
}

// AddCommentHandler handles adding a comment to a post.
func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)
	postIDStr := r.URL.Query().Get("post_id")
	content := r.FormValue("comment")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	if content == "" {
		http.Error(w, "Comment content cannot be empty", http.StatusBadRequest)
		return
	}

	err = database.AddComment(postID, userID, content)
	if err != nil {
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/view-post?id=%d", postID), http.StatusSeeOther)
}
