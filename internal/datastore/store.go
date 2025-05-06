// Package datastore provides interfaces and implementations for data persistence
package datastore

import (
	"context"
)

// Store defines the interface for blog data operations
type Store interface {
	// Create creates a new blog entry
	Create(ctx context.Context, title, content string) (ID, error)

	// Get retrieves a blog by ID with its comments
	Get(ctx context.Context, id ID) (*Blog, error)

	// Update updates an existing blog
	Update(ctx context.Context, id ID, title, content *string) error

	// Delete deletes a blog and its comments
	Delete(ctx context.Context, id ID) error

	// List retrieves a paginated list of blog summaries
	List(ctx context.Context, pageSize int32, pageToken string) ([]*BlogSummary, string, error)

	// AddComment adds a comment to a blog
	AddComment(ctx context.Context, blogID ID, content, author string) (ID, error)
}
