package main

import (
	"fmt"
	"net/http"

	"forum/database"
	"forum/handlers"
)

func main() {
	// initialize database
	database.InitDB()

	// Serve static files (CSS, JS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Define explicit routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/signup", handlers.SignUpHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	// Wrap the DashboardHandler with AuthMiddleware
	http.Handle("/dashboard", handlers.AuthMiddleware(http.HandlerFunc(handlers.DashboardHandler)))
	http.Handle("/create-post", handlers.AuthMiddleware(http.HandlerFunc(handlers.CreatePostHandler)))
	http.Handle("/posts", http.HandlerFunc(handlers.ListPostsHandler))
	http.Handle("/view-post", http.HandlerFunc(handlers.ViewPostHandler))
	http.Handle("/edit-post", handlers.AuthMiddleware(http.HandlerFunc(handlers.EditPostHandler)))
	http.Handle("/delete-post", handlers.AuthMiddleware(http.HandlerFunc(handlers.DeletePostHandler)))

	// Start the server
	fmt.Println("Server starting at port localhost:8080")
	http.ListenAndServe(":8080", nil)
}
