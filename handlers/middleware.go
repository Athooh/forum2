package handlers

import (
	"context"
	"forum/database"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var userID int
		err = database.DB.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		context := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(context))
	})
}

func OwnershipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(userIDKey).(int)
		resourceID := r.URL.Query().Get("id")

		isOwner, err := database.VerifyResourceOwnership(userID, resourceID, r.URL.Path)
		if err != nil || !isOwner {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GuestMiddleware allows both authenticated and guest users
func GuestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			// Guest user - continue without user context
			next.ServeHTTP(w, r)
			return
		}

		// Attempt to get user context if authenticated
		var userID int
		err = database.DB.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID)
		if err != nil {
			// Invalid session - continue as guest
			next.ServeHTTP(w, r)
			return
		}

		// Add user context if authenticated
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
