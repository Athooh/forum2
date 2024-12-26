package database

import (
	"database/sql"
	"fmt"
	"forum/utils"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	ID             int
	UserID         int
	Title          string
	Content        string
	Preview        string // Truncated content for preview
	Category       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Username       string // Add this field to store the post author's username
	Likes          int    // Number of likes
	Dislikes       int    // Number of dislikes
	CommentsCount  int    // New field to store comment count
	CreatedAtHuman string // Human-readable time difference
	ImageURL       string // New field to store the image URL
}

type Comment struct {
	ID             int
	PostID         int
	UserID         int
	Content        string
	CreatedAt      time.Time
	Username       string // Add this field to store the comment author's username
	CreatedAtHuman string // Human-readable time difference
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
		image_url TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (id)
    );`
	_, err = DB.Exec(postTable)
	if err != nil {
		log.Fatalf("Failed to create posts table: %v\n", err)
	}

	commentTable := `
    CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES posts (id),
        FOREIGN KEY (user_id) REFERENCES users (id)
    );`
	_, tableErr := DB.Exec(commentTable)
	if tableErr != nil {
		log.Fatalf("Failed to create comments table: %v\n", tableErr)
	}

	categoryTable := `
    CREATE TABLE IF NOT EXISTS categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT UNIQUE NOT NULL
    );`
	_, err = DB.Exec(categoryTable)
	if err != nil {
		log.Fatalf("Failed to create categories table: %v\n", err)
	}

	// Insert predefined categories
	categorySeed := `
    INSERT OR IGNORE INTO categories (name) VALUES 
    ('Technology'),
    ('Design'),
    ('Marketing'),
    ('Development'),
    ('Science'),
    ('Health'),
    ('Education'),
    ('Business'),
    ('Lifestyle'),
    ('Entertainment');`
	_, err = DB.Exec(categorySeed)
	if err != nil {
		log.Fatalf("Failed to seed categories: %v\n", err)
	}

	reactionTable := `
	CREATE TABLE IF NOT EXISTS post_reactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		reaction TEXT CHECK(reaction IN ('like', 'dislike')) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id),
		FOREIGN KEY (post_id) REFERENCES posts (id),
		UNIQUE(user_id, post_id)
	);`
	_, err = DB.Exec(reactionTable)
	if err != nil {
		log.Fatalf("Failed to create post_reactions table: %v\n", err)
	}

}

// AddComment inserts a new comment into the database.
func AddComment(postID, userID int, content string) error {
	query := `INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, postID, userID, content)
	return err
}

// GetCommentsByPostID retrieves all comments for a specific post ID.
func GetCommentsByPostID(postID int) ([]Comment, error) {
	query := `
    SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username
    FROM comments c
    INNER JOIN users u ON c.user_id = u.id
    WHERE c.post_id = ?
    ORDER BY c.created_at ASC`

	rows, err := DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Username)
		if err != nil {
			return nil, err
		}
		comment.CreatedAtHuman = utils.TimeAgo(comment.CreatedAt) // Populate human-readable time
		comments = append(comments, comment)
	}
	return comments, nil
}

// CreatePost inserts a new post into the database.
func CreatePost(userID int, title, content, category string) error {
	query := `INSERT INTO posts (user_id, title, content, category) VALUES (?, ?, ?, ?)`
	_, err := DB.Exec(query, userID, title, content, category)
	return err
}

// GetPost retrieves a single post by its ID.
func GetPost(postID int) (*Post, error) {
	query := `SELECT id, user_id, title, content, category, image_url, created_at, updated_at FROM posts WHERE id = ?`
	row := DB.QueryRow(query, postID)

	var post Post
	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	post.CreatedAtHuman = utils.TimeAgo(post.CreatedAt)
	return &post, nil
}

// GetAllPosts retrieves all posts with optional filtering by category, including like, dislike, and comment counts.
func GetAllPosts(category string) ([]Post, error) {
	var rows *sql.Rows
	var err error

	if category != "" {
		query := `
		SELECT p.id, p.user_id, p.title, p.content, p.category, p.image_url, p.created_at, p.updated_at,
		       IFNULL(likes.count, 0) AS likes, 
		       IFNULL(dislikes.count, 0) AS dislikes,
		       IFNULL(comments.count, 0) AS comments_count
		FROM posts p
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS count
			FROM post_reactions
			WHERE reaction = 'like'
			GROUP BY post_id
		) AS likes ON p.id = likes.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS count
			FROM post_reactions
			WHERE reaction = 'dislike'
			GROUP BY post_id
		) AS dislikes ON p.id = dislikes.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS count
			FROM comments
			GROUP BY post_id
		) AS comments ON p.id = comments.post_id
		WHERE p.category = ?
		ORDER BY p.created_at DESC`

		rows, err = DB.Query(query, category)
	} else {
		query := `
		SELECT p.id, p.user_id, p.title, p.content, p.category, p.image_url, p.created_at, p.updated_at,
		       IFNULL(likes.count, 0) AS likes, 
		       IFNULL(dislikes.count, 0) AS dislikes,
		       IFNULL(comments.count, 0) AS comments_count
		FROM posts p
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS count
			FROM post_reactions
			WHERE reaction = 'like'
			GROUP BY post_id
		) AS likes ON p.id = likes.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS count
			FROM post_reactions
			WHERE reaction = 'dislike'
			GROUP BY post_id
		) AS dislikes ON p.id = dislikes.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS count
			FROM comments
			GROUP BY post_id
		) AS comments ON p.id = comments.post_id
		ORDER BY p.created_at DESC`

		rows, err = DB.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var imageURL sql.NullString // Handle nullable image_url

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.Category,
			&imageURL, // Added to retrieve the image URL
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Likes,
			&post.Dislikes,
			&post.CommentsCount,
		)
		if err != nil {
			return nil, err
		}

		// Handle NULL image_url
		if imageURL.Valid {
			post.ImageURL = imageURL.String
		} else {
			post.ImageURL = "/static/images/default-post.jpg"
		}

		// Truncate content for preview and format time
		post.Preview = utils.TruncateContent(post.Content, 30) // Limit to 30 words
		post.CreatedAtHuman = utils.TimeAgo(post.CreatedAt)    // Populate human-readable time

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

// GetCategories retrieves all categories from the database.
func GetCategories() ([]string, error) {
	query := `SELECT name FROM categories ORDER BY name ASC`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		categories = append(categories, name)
	}
	return categories, nil
}

// AddReaction toggles a like or dislike reaction on a post.
func ToggleReaction(userID, postID int, reaction string) error {
	var existingReaction string

	// Check if a reaction already exists
	query := `SELECT reaction FROM post_reactions WHERE user_id = ? AND post_id = ?`
	err := DB.QueryRow(query, userID, postID).Scan(&existingReaction)

	if err == sql.ErrNoRows {
		// No existing reaction, insert a new one
		insertQuery := `INSERT INTO post_reactions (user_id, post_id, reaction) VALUES (?, ?, ?)`
		_, err := DB.Exec(insertQuery, userID, postID, reaction)
		return err
	} else if err != nil {
		return err // Return any other unexpected errors
	}

	if existingReaction == reaction {
		// If the reaction is the same, remove it (toggle off)
		deleteQuery := `DELETE FROM post_reactions WHERE user_id = ? AND post_id = ?`
		_, err := DB.Exec(deleteQuery, userID, postID)
		return err
	}

	// If the reaction is different, update it
	updateQuery := `UPDATE post_reactions SET reaction = ? WHERE user_id = ? AND post_id = ?`
	_, err = DB.Exec(updateQuery, reaction, userID, postID)
	return err
}

// GetReactionCounts returns the number of likes and dislikes for a post.
func GetReactionCounts(postID int) (int, int, error) {
	var likes, dislikes int

	likeQuery := `SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND reaction = 'like'`
	err := DB.QueryRow(likeQuery, postID).Scan(&likes)
	if err != nil {
		return 0, 0, err
	}

	dislikeQuery := `SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND reaction = 'dislike'`
	err = DB.QueryRow(dislikeQuery, postID).Scan(&dislikes)
	if err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}

// CreatePostWithImage creates a new post with an optional image.
func CreatePostWithImage(userID int, title, content, category, imageURL string) error {
	query := `INSERT INTO posts (user_id, title, content, category, image_url) VALUES (?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, userID, title, content, category, imageURL)
	return err
}

func VerifyResourceOwnership(userID int, resourceID string, resourceType string) (bool, error) {
	var ownerID int
	var query string

	switch {
	case strings.Contains(resourceType, "post"):
		query = "SELECT user_id FROM posts WHERE id = ?"
	case strings.Contains(resourceType, "comment"):
		query = "SELECT user_id FROM comments WHERE id = ?"
	default:
		return false, fmt.Errorf("invalid resource type")
	}

	err := DB.QueryRow(query, resourceID).Scan(&ownerID)
	if err != nil {
		return false, err
	}

	return ownerID == userID, nil
}
