package handlers

import (
	"forum/database"
	"net/http"
	"time"
)

type Session struct {
	ID           int
	UserID       int
	SessionToken string
	ExpiresAt    time.Time
	LastActivity time.Time
	UserRole     string
}

func ValidateSession(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}

	var session Session
	err = database.DB.QueryRow(`
		SELECT id, user_id, session_token, expires_at, last_activity, user_role 
		FROM sessions 
		WHERE session_token = ? AND expires_at > ?`,
		cookie.Value, time.Now(),
	).Scan(&session.ID, &session.UserID, &session.SessionToken,
		&session.ExpiresAt, &session.LastActivity, &session.UserRole)

	if err != nil {
		return nil, err
	}

	// Update last activity
	_, err = database.DB.Exec(`
		UPDATE sessions 
		SET last_activity = ? 
		WHERE id = ?`,
		time.Now(), session.ID)

	return &session, err
}
