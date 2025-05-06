// Package datastore provides interfaces and implementations for data persistence
package datastore

import (
	"time"
)

// ID represents a UUID used as an identifier
type ID string

// Blog represents a blog entry in the database
type Blog struct {
	ID        ID        `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Comments  []Comment
}

// Comment represents a comment in the database
type Comment struct {
	ID        ID        `db:"id"`
	BlogID    ID        `db:"blog_id"`
	Content   string    `db:"content"`
	Author    string    `db:"author"`
	CreatedAt time.Time `db:"created_at"`
}

// BlogSummary represents a summary of a blog entry
type BlogSummary struct {
	ID           ID     `db:"id"`
	Title        string `db:"title"`
	CommentCount int32  `db:"comment_count"`
}
