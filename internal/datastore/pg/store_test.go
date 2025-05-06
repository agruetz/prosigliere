package pg_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/agruetz/prosigliere/internal/datastore"
	"github.com/agruetz/prosigliere/internal/datastore/pg"
)

func TestCreate(t *testing.T) {
	// Define test cases
	tests := []struct {
		name        string
		title       string
		content     string
		mockSetup   func(mock sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name:    "successful creation",
			title:   "Test Title",
			content: "Test Content",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO blogs").
					WithArgs(sqlmock.AnyArg(), "Test Title", "Test Content").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: false,
		},
		{
			name:    "database error",
			title:   "Test Title",
			content: "Test Content",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO blogs").
					WithArgs(sqlmock.AnyArg(), "Test Title", "Test Content").
					WillReturnError(errors.New("database error"))
			},
			expectError: true,
			errorMsg:    "failed to create blog",
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock database
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Create a new store with the mock database
			store := pg.NewWithDB(db)

			// Set up expectations
			tc.mockSetup(mock)

			// Call the method
			id, err := store.Create(context.Background(), tc.title, tc.content)

			// Assert expectations
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, id)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGet(t *testing.T) {
	// Define test cases
	tests := []struct {
		name        string
		id          datastore.ID
		mockSetup   func(mock sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
		expected    *datastore.Blog
	}{
		{
			name: "successful retrieval",
			id:   datastore.ID("test-id"),
			mockSetup: func(mock sqlmock.Sqlmock) {
				testID := datastore.ID("test-id")
				testTitle := "Test Title"
				testContent := "Test Content"
				testCreatedAt := time.Now()
				testUpdatedAt := time.Now()

				rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
					AddRow(testID, testTitle, testContent, testCreatedAt, testUpdatedAt)

				mock.ExpectQuery(`SELECT id, title, content, created_at, updated_at FROM blogs WHERE id = \$1`).
					WithArgs(string(testID)).
					WillReturnRows(rows)
			},
			expectError: false,
			expected: &datastore.Blog{
				ID:      datastore.ID("test-id"),
				Title:   "Test Title",
				Content: "Test Content",
				// CreatedAt and UpdatedAt will be set by the database
			},
		},
		{
			name: "blog not found",
			id:   datastore.ID("non-existent-id"),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, title, content, created_at, updated_at FROM blogs WHERE id = ?").
					WithArgs("non-existent-id").
					WillReturnError(sql.ErrNoRows)
			},
			expectError: true,
			errorMsg:    "blog not found",
		},
		{
			name: "database error",
			id:   datastore.ID("test-id"),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, title, content, created_at, updated_at FROM blogs WHERE id = ?").
					WithArgs("test-id").
					WillReturnError(errors.New("database error"))
			},
			expectError: true,
			errorMsg:    "failed to get blog",
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock database
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Create a new store with the mock database
			store := pg.NewWithDB(db)

			// Set up expectations
			tc.mockSetup(mock)

			// Call the method
			blog, err := store.Get(context.Background(), tc.id)

			// Assert expectations
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, blog)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.ID, blog.ID)
				assert.Equal(t, tc.expected.Title, blog.Title)
				assert.Equal(t, tc.expected.Content, blog.Content)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdate(t *testing.T) {
	// Define test cases
	testTitle := "Updated Title"
	testContent := "Updated Content"

	tests := []struct {
		name        string
		id          datastore.ID
		title       *string
		content     *string
		mockSetup   func(mock sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name:    "successful update with both fields",
			id:      datastore.ID("test-id"),
			title:   &testTitle,
			content: &testContent,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE blogs SET").
					WithArgs(testTitle, testContent, string(datastore.ID("test-id"))).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectError: false,
		},
		{
			name:    "successful update with title only",
			id:      datastore.ID("test-id"),
			title:   &testTitle,
			content: nil,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE blogs SET").
					WithArgs(testTitle, string(datastore.ID("test-id"))).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectError: false,
		},
		{
			name:    "successful update with content only",
			id:      datastore.ID("test-id"),
			title:   nil,
			content: &testContent,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE blogs SET").
					WithArgs(testContent, string(datastore.ID("test-id"))).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectError: false,
		},
		{
			name:    "blog not found",
			id:      datastore.ID("non-existent-id"),
			title:   &testTitle,
			content: &testContent,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE blogs SET").
					WithArgs(testTitle, testContent, string(datastore.ID("non-existent-id"))).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectError: true,
			errorMsg:    "blog not found",
		},
		{
			name:    "database error",
			id:      datastore.ID("test-id"),
			title:   &testTitle,
			content: &testContent,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE blogs SET").
					WithArgs(testTitle, testContent, string(datastore.ID("test-id"))).
					WillReturnError(errors.New("database error"))
			},
			expectError: true,
			errorMsg:    "failed to update blog",
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock database
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Create a new store with the mock database
			store := pg.NewWithDB(db)

			// Set up expectations
			tc.mockSetup(mock)

			// Call the method
			err = store.Update(context.Background(), tc.id, tc.title, tc.content)

			// Assert expectations
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDelete(t *testing.T) {
	// Define test cases
	tests := []struct {
		name        string
		id          datastore.ID
		mockSetup   func(mock sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful deletion",
			id:   datastore.ID("test-id"),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM blogs WHERE id = ?").
					WithArgs(string(datastore.ID("test-id"))).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectError: false,
		},
		{
			name: "blog not found",
			id:   datastore.ID("non-existent-id"),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM blogs WHERE id = ?").
					WithArgs(string(datastore.ID("non-existent-id"))).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectError: true,
			errorMsg:    "blog not found",
		},
		{
			name: "database error",
			id:   datastore.ID("test-id"),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM blogs WHERE id = ?").
					WithArgs(string(datastore.ID("test-id"))).
					WillReturnError(errors.New("database error"))
			},
			expectError: true,
			errorMsg:    "failed to delete blog",
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock database
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Create a new store with the mock database
			store := pg.NewWithDB(db)

			// Set up expectations
			tc.mockSetup(mock)

			// Call the method
			err = store.Delete(context.Background(), tc.id)

			// Assert expectations
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestList(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		pageSize      int32
		pageToken     string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectError   bool
		errorMsg      string
		expectedBlogs []*datastore.BlogSummary
		nextPageToken string
	}{
		{
			name:      "successful list without pagination",
			pageSize:  10,
			pageToken: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				testID1 := datastore.ID("test-id-1")
				testID2 := datastore.ID("test-id-2")
				testTitle1 := "Test Title 1"
				testTitle2 := "Test Title 2"
				commentCount1 := int32(5)
				commentCount2 := int32(10)

				rows := sqlmock.NewRows([]string{"id", "title", "comment_count"}).
					AddRow(testID1, testTitle1, commentCount1).
					AddRow(testID2, testTitle2, commentCount2)

				mock.ExpectQuery("SELECT b.id, b.title, COUNT\\(c.id\\) as comment_count FROM blogs b LEFT JOIN comments c ON b.id = c.blog_id").
					WillReturnRows(rows)
			},
			expectError: false,
			expectedBlogs: []*datastore.BlogSummary{
				{
					ID:           datastore.ID("test-id-1"),
					Title:        "Test Title 1",
					CommentCount: 5,
				},
				{
					ID:           datastore.ID("test-id-2"),
					Title:        "Test Title 2",
					CommentCount: 10,
				},
			},
			nextPageToken: "",
		},
		{
			name:      "successful list with pagination",
			pageSize:  1,
			pageToken: "test-id-1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				testID2 := datastore.ID("test-id-2")
				testTitle2 := "Test Title 2"
				commentCount2 := int32(10)

				rows := sqlmock.NewRows([]string{"id", "title", "comment_count"}).
					AddRow(testID2, testTitle2, commentCount2)

				mock.ExpectQuery("SELECT b.id, b.title, COUNT\\(c.id\\) as comment_count FROM blogs b LEFT JOIN comments c ON b.id = c.blog_id WHERE b.id > \\$1 GROUP BY b.id, b.title ORDER BY b.id LIMIT \\$2").
					WithArgs("test-id-1", int32(2)). // pageSize + 1 = 1 + 1 = 2
					WillReturnRows(rows)
			},
			expectError: false,
			expectedBlogs: []*datastore.BlogSummary{
				{
					ID:           datastore.ID("test-id-2"),
					Title:        "Test Title 2",
					CommentCount: 10,
				},
			},
			nextPageToken: "",
		},
		{
			name:      "database error",
			pageSize:  10,
			pageToken: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT b.id, b.title, COUNT\\(c.id\\) as comment_count FROM blogs b LEFT JOIN comments c ON b.id = c.blog_id").
					WillReturnError(errors.New("database error"))
			},
			expectError: true,
			errorMsg:    "failed to list blogs",
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock database
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Create a new store with the mock database
			store := pg.NewWithDB(db)

			// Set up expectations
			tc.mockSetup(mock)

			// Call the method
			summaries, nextPageToken, err := store.List(context.Background(), tc.pageSize, tc.pageToken)

			// Assert expectations
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, summaries)
				assert.Empty(t, nextPageToken)
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(tc.expectedBlogs), len(summaries))
				for i, expected := range tc.expectedBlogs {
					assert.Equal(t, expected.ID, summaries[i].ID)
					assert.Equal(t, expected.Title, summaries[i].Title)
					assert.Equal(t, expected.CommentCount, summaries[i].CommentCount)
				}
				assert.Equal(t, tc.nextPageToken, nextPageToken)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAddComment(t *testing.T) {
	// Define test cases
	tests := []struct {
		name        string
		blogID      datastore.ID
		content     string
		author      string
		mockSetup   func(mock sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name:    "successful comment addition",
			blogID:  datastore.ID("test-blog-id"),
			content: "Test Comment",
			author:  "Test Author",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Set up expectations for checking if blog exists
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(1)
				mock.ExpectQuery("SELECT 1 FROM blogs WHERE id = ?").
					WithArgs(string(datastore.ID("test-blog-id"))).
					WillReturnRows(rows)

				// Set up expectations for inserting comment
				mock.ExpectExec("INSERT INTO comments").
					WithArgs(sqlmock.AnyArg(), string(datastore.ID("test-blog-id")), "Test Comment", "Test Author").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: false,
		},
		{
			name:    "blog not found",
			blogID:  datastore.ID("non-existent-blog-id"),
			content: "Test Comment",
			author:  "Test Author",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT 1 FROM blogs WHERE id = ?").
					WithArgs(string(datastore.ID("non-existent-blog-id"))).
					WillReturnError(sql.ErrNoRows)
			},
			expectError: true,
			errorMsg:    "blog not found",
		},
		{
			name:    "database error on check",
			blogID:  datastore.ID("test-blog-id"),
			content: "Test Comment",
			author:  "Test Author",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT 1 FROM blogs WHERE id = ?").
					WithArgs(string(datastore.ID("test-blog-id"))).
					WillReturnError(errors.New("database error"))
			},
			expectError: true,
			errorMsg:    "failed to check blog existence",
		},
		{
			name:    "database error on insert",
			blogID:  datastore.ID("test-blog-id"),
			content: "Test Comment",
			author:  "Test Author",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Set up expectations for checking if blog exists
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(1)
				mock.ExpectQuery("SELECT 1 FROM blogs WHERE id = ?").
					WithArgs(string(datastore.ID("test-blog-id"))).
					WillReturnRows(rows)

				// Set up expectations for inserting comment with error
				mock.ExpectExec("INSERT INTO comments").
					WithArgs(sqlmock.AnyArg(), string(datastore.ID("test-blog-id")), "Test Comment", "Test Author").
					WillReturnError(errors.New("database error"))
			},
			expectError: true,
			errorMsg:    "failed to add comment",
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock database
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Create a new store with the mock database
			store := pg.NewWithDB(db)

			// Set up expectations
			tc.mockSetup(mock)

			// Call the method
			commentID, err := store.AddComment(context.Background(), tc.blogID, tc.content, tc.author)

			// Assert expectations
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Empty(t, commentID)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, commentID)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
