package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	Category  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string // Add this field to store the post author's username
}

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	createTables()
}

func createTables() {
	userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT UNIQUE NOT NULL,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
	_, err := DB.Exec(userTable)
	if err != nil {
		log.Fatalf("Failed to create users table: %v\n", err)
	}

	sessionTable := `
    CREATE TABLE IF NOT EXISTS sessions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        session_token TEXT UNIQUE NOT NULL,
        expires_at DATETIME NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id)
    );`
	_, err = DB.Exec(sessionTable)
	if err != nil {
		log.Fatalf("Failed to create sessions table: %v\n", err)
	}

	postTable := `
    CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        category TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (id)
    );`
	_, err = DB.Exec(postTable)
	if err != nil {
		log.Fatalf("Failed to create posts table: %v\n", err)
	}
}

// CreatePost inserts a new post into the database.
func CreatePost(userID int, title, content, category string) error {
	query := `INSERT INTO posts (user_id, title, content, category) VALUES (?, ?, ?, ?)`
	_, err := DB.Exec(query, userID, title, content, category)
	return err
}

// GetPost retrieves a single post by its ID.
func GetPost(postID int) (*Post, error) {
	query := `SELECT id, user_id, title, content, category, created_at, updated_at FROM posts WHERE id = ?`
	row := DB.QueryRow(query, postID)

	var post Post
	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category, &post.CreatedAt, &post.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &post, nil
}

// GetAllPosts retrieves all posts with optional filtering by category.
func GetAllPosts(category string) ([]Post, error) {
	var rows *sql.Rows
	var err error

	if category != "" {
		query := `SELECT id, user_id, title, content, category, created_at, updated_at FROM posts WHERE category = ? ORDER BY created_at DESC`
		rows, err = DB.Query(query, category)
	} else {
		query := `SELECT id, user_id, title, content, category, created_at, updated_at FROM posts ORDER BY created_at DESC`
		rows, err = DB.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// UpdatePost updates an existing post.
func UpdatePost(postID int, title, content, category string) error {
	query := `UPDATE posts SET title = ?, content = ?, category = ?, updated_at = ? WHERE id = ?`
	_, err := DB.Exec(query, title, content, category, time.Now(), postID)
	return err
}

// DeletePost deletes a post by its ID.
func DeletePost(postID int) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := DB.Exec(query, postID)
	return err
}

// GetCategoryPostCounts retrieves the number of posts for each category.
func GetCategoryPostCounts() (map[string]int, error) {
	// Query to count posts per category
	query := `SELECT category, COUNT(*) FROM posts GROUP BY category`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categoryCounts := make(map[string]int)

	// Loop through the rows and store category counts in the map
	for rows.Next() {
		var category string
		var count int
		err := rows.Scan(&category, &count)
		if err != nil {
			return nil, err
		}
		categoryCounts[category] = count
	}

	return categoryCounts, nil
}

