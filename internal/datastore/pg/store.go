// Package pg provides a PostgreSQL implementation of the datastore.Store interface
package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"prosigliere/internal/datastore"
)

// Create creates a new blog entry
func (s *Store) Create(ctx context.Context, title, content string) (datastore.ID, error) {
	id := uuid.New().String()
	query := `
		INSERT INTO blogs (id, title, content)
		VALUES ($1, $2, $3)
	`
	_, err := s.db.ExecContext(ctx, query, id, title, content)
	if err != nil {
		return "", fmt.Errorf("failed to create blog: %w", err)
	}
	return datastore.ID(id), nil
}

// Get retrieves a blog by ID with its comments
func (s *Store) Get(ctx context.Context, id datastore.ID) (*datastore.Blog, error) {
	// First get the blog
	query := `
		SELECT id, title, content, created_at, updated_at
		FROM blogs
		WHERE id = $1
	`
	var blog datastore.Blog
	var createdAt, updatedAt time.Time

	err := s.db.QueryRowContext(ctx, query, string(id)).Scan(
		&blog.ID, &blog.Title, &blog.Content, &createdAt, &updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("blog not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get blog: %w", err)
	}

	blog.CreatedAt = createdAt
	blog.UpdatedAt = updatedAt

	// Note: In a real implementation with sqlx, we could use StructScan to directly scan into the struct
	// This would require adding github.com/jmoiron/sqlx as a dependency

	// We're not returning comments as the interface doesn't specify it
	// But we could fetch them if needed

	return &blog, nil
}

// Update updates an existing blog
func (s *Store) Update(ctx context.Context, id datastore.ID, title, content *string) error {
	// Build the query dynamically based on which fields are provided
	query := "UPDATE blogs SET"
	args := []interface{}{}
	paramCount := 1
	updateParts := []string{}

	if title != nil {
		updateParts = append(updateParts, fmt.Sprintf(" title = $%d", paramCount))
		args = append(args, *title)
		paramCount++
	}

	if content != nil {
		updateParts = append(updateParts, fmt.Sprintf(" content = $%d", paramCount))
		args = append(args, *content)
		paramCount++
	}

	if len(updateParts) == 0 {
		return nil // Nothing to update
	}

	query += strings.Join(updateParts, ",")

	// Add WHERE clause
	query += fmt.Sprintf(" WHERE id = $%d", paramCount)
	args = append(args, string(id))

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update blog: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("blog not found")
	}

	return nil
}

// Delete deletes a blog and its comments
func (s *Store) Delete(ctx context.Context, id datastore.ID) error {
	// Comments will be deleted automatically due to ON DELETE CASCADE
	query := `DELETE FROM blogs WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, string(id))
	if err != nil {
		return fmt.Errorf("failed to delete blog: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("blog not found")
	}

	return nil
}

// List retrieves a paginated list of blog summaries
func (s *Store) List(ctx context.Context, pageSize int32, pageToken string) ([]*datastore.BlogSummary, string, error) {
	query := `
		SELECT b.id, b.title, COUNT(c.id) as comment_count
		FROM blogs b
		LEFT JOIN comments c ON b.id = c.blog_id
	`
	args := []interface{}{}
	paramCount := 1

	// Add pagination if pageToken is provided
	if pageToken != "" {
		query += fmt.Sprintf(" WHERE b.id > $%d", paramCount)
		args = append(args, pageToken)
		paramCount++
	}

	query += `
		GROUP BY b.id, b.title
		ORDER BY b.id
		LIMIT $` + fmt.Sprintf("%d", paramCount)

	args = append(args, pageSize+1) // Fetch one extra to determine if there are more results

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list blogs: %w", err)
	}
	defer rows.Close()

	var summaries []*datastore.BlogSummary
	for rows.Next() {
		var summary datastore.BlogSummary
		err := rows.Scan(&summary.ID, &summary.Title, &summary.CommentCount)
		if err != nil {
			return nil, "", fmt.Errorf("failed to scan blog summary: %w", err)
		}
		summaries = append(summaries, &summary)
	}

	if err := rows.Err(); err != nil {
		return nil, "", fmt.Errorf("error iterating blog summaries: %w", err)
	}

	// Handle pagination
	var nextPageToken string
	if len(summaries) > int(pageSize) {
		// We fetched one extra result, so there are more pages
		nextPageToken = string(summaries[len(summaries)-1].ID)
		summaries = summaries[:len(summaries)-1] // Remove the extra result
	}

	return summaries, nextPageToken, nil
}

// AddComment adds a comment to a blog
func (s *Store) AddComment(ctx context.Context, blogID datastore.ID, content, author string) (datastore.ID, error) {
	// First check if the blog exists
	checkQuery := `SELECT 1 FROM blogs WHERE id = $1`
	var exists int
	err := s.db.QueryRowContext(ctx, checkQuery, string(blogID)).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("blog not found")
		}
		return "", fmt.Errorf("failed to check blog existence: %w", err)
	}

	// Insert the comment
	id := uuid.New().String()
	query := `
		INSERT INTO comments (id, blog_id, content, author)
		VALUES ($1, $2, $3, $4)
	`
	_, err = s.db.ExecContext(ctx, query, id, string(blogID), content, author)
	if err != nil {
		return "", fmt.Errorf("failed to add comment: %w", err)
	}

	return datastore.ID(id), nil
}
