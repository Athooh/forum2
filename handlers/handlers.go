package handlers

import (
	"context"
	"fmt"
	"net/http"
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
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")

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

    var username string
    query := `SELECT username FROM users WHERE id = ?`
    err := database.DB.QueryRow(query, userID).Scan(&username)
    if err != nil {
        http.Error(w, "User not found", http.StatusInternalServerError)
        return
    }

    data := map[string]interface{}{
        "Title":      "Dashboard - Forum",
        "IsLoggedIn": true,
        "Username":   username,
    }

    err = templates.ExecuteTemplate(w, "dashboard", data)
    if err != nil {
        fmt.Printf("Error executing dashboard template: %v\n", err)
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
        return
    }
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
