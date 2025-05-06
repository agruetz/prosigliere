package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/agruetz/prosigliere/internal/datastore"
	"github.com/agruetz/prosigliere/internal/datastore/mocks"
	blogpb "github.com/agruetz/prosigliere/protos/v1/blog"
)

func TestNewBlogService(t *testing.T) {
	mockStore := mocks.NewStore(t)
	service := NewBlogService(mockStore)

	assert.NotNil(t, service)
	assert.Equal(t, mockStore, service.store)
}

func TestBlogService_Create(t *testing.T) {
	tests := []struct {
		name        string
		req         *blogpb.CreateReq
		setupMock   func(mock *mocks.Store)
		expectedID  string
		expectedErr error
	}{
		{
			name: "successful creation",
			req: &blogpb.CreateReq{
				Title:   "Test Blog",
				Content: "This is a test blog content",
			},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("Create", mock.Anything, "Test Blog", "This is a test blog content").
					Return(datastore.ID("123e4567-e89b-12d3-a456-426614174000"), nil)
			},
			expectedID:  "123e4567-e89b-12d3-a456-426614174000",
			expectedErr: nil,
		},
		{
			name: "missing title",
			req: &blogpb.CreateReq{
				Content: "This is a test blog content",
			},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("Create", mock.Anything, "", "This is a test blog content").
					Return(datastore.ID(""), errors.New("missing title"))
			},
			expectedID:  "",
			expectedErr: status.Error(codes.Internal, "failed to create blog: missing title"),
		},
		{
			name: "missing content",
			req: &blogpb.CreateReq{
				Title: "Test Blog",
			},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("Create", mock.Anything, "Test Blog", "").
					Return(datastore.ID(""), errors.New("missing content"))
			},
			expectedID:  "",
			expectedErr: status.Error(codes.Internal, "failed to create blog: missing content"),
		},
		{
			name: "store error",
			req: &blogpb.CreateReq{
				Title:   "Test Blog",
				Content: "This is a test blog content",
			},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("Create", mock.Anything, "Test Blog", "This is a test blog content").
					Return(datastore.ID(""), errors.New("database error"))
			},
			expectedID:  "",
			expectedErr: status.Error(codes.Internal, "failed to create blog: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := mocks.NewStore(t)
			tt.setupMock(mockStore)

			service := NewBlogService(mockStore)
			resp, err := service.Create(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedID, resp.Id.Value)
			}
		})
	}
}

func TestBlogService_Get(t *testing.T) {
	testTime := time.Now().UTC()
	testBlog := &datastore.Blog{
		ID:        datastore.ID("123e4567-e89b-12d3-a456-426614174000"),
		Title:     "Test Blog",
		Content:   "This is a test blog content",
		CreatedAt: testTime,
		UpdatedAt: testTime,
	}

	tests := []struct {
		name        string
		req         *blogpb.GetReq
		setupMock   func(mock *mocks.Store)
		expectedErr error
	}{
		{
			name: "successful get",
			req: &blogpb.GetReq{
				Id: &blogpb.UUID{
					Value: "123e4567-e89b-12d3-a456-426614174000",
				},
			},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("Get", mock.Anything, datastore.ID("123e4567-e89b-12d3-a456-426614174000")).
					Return(testBlog, nil)
			},
			expectedErr: nil,
		},
		{
			name: "missing ID",
			req:  &blogpb.GetReq{},
			setupMock: func(mockStore *mocks.Store) {
				// No mock setup needed
			},
			expectedErr: status.Error(codes.InvalidArgument, "blog ID is required"),
		},
		{
			name: "store error",
			req: &blogpb.GetReq{
				Id: &blogpb.UUID{
					Value: "123e4567-e89b-12d3-a456-426614174000",
				},
			},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("Get", mock.Anything, datastore.ID("123e4567-e89b-12d3-a456-426614174000")).
					Return(nil, errors.New("not found"))
			},
			expectedErr: status.Error(codes.NotFound, "failed to get blog: not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := mocks.NewStore(t)
			tt.setupMock(mockStore)

			service := NewBlogService(mockStore)
			resp, err := service.Get(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, testBlog.ID, datastore.ID(resp.Blog.Id.Value))
				assert.Equal(t, testBlog.Title, resp.Blog.Title)
				assert.Equal(t, testBlog.Content, resp.Blog.Content)
				assert.Equal(t, timestamppb.New(testBlog.CreatedAt).AsTime().Unix(), resp.Blog.CreatedAt.AsTime().Unix())
				assert.Equal(t, timestamppb.New(testBlog.UpdatedAt).AsTime().Unix(), resp.Blog.UpdatedAt.AsTime().Unix())
			}
		})
	}
}

func TestBlogService_Update(t *testing.T) {
	tests := []struct {
		name        string
		req         *blogpb.UpdateReq
		setupMock   func(mock *mocks.Store)
		expectedErr error
	}{
		{
			name: "successful update with both fields",
			req: &blogpb.UpdateReq{
				Id:      &blogpb.UUID{Value: "123e4567-e89b-12d3-a456-426614174000"},
				Title:   stringPtr("Updated Title"),
				Content: stringPtr("Updated Content"),
			},
			setupMock: func(mockStore *mocks.Store) {
				title := "Updated Title"
				content := "Updated Content"
				mockStore.On("Update", mock.Anything, datastore.ID("123e4567-e89b-12d3-a456-426614174000"), &title, &content).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "successful update with title only",
			req: &blogpb.UpdateReq{
				Id:    &blogpb.UUID{Value: "123e4567-e89b-12d3-a456-426614174000"},
				Title: stringPtr("Updated Title"),
			},
			setupMock: func(mockStore *mocks.Store) {
				title := "Updated Title"
				mockStore.On("Update", mock.Anything, datastore.ID("123e4567-e89b-12d3-a456-426614174000"), &title, (*string)(nil)).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "successful update with content only",
			req: &blogpb.UpdateReq{
				Id:      &blogpb.UUID{Value: "123e4567-e89b-12d3-a456-426614174000"},
				Content: stringPtr("Updated Content"),
			},
			setupMock: func(mockStore *mocks.Store) {
				content := "Updated Content"
				mockStore.On("Update", mock.Anything, datastore.ID("123e4567-e89b-12d3-a456-426614174000"), (*string)(nil), &content).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "missing ID",
			req:  &blogpb.UpdateReq{},
			setupMock: func(mockStore *mocks.Store) {
				// No mock setup needed
			},
			expectedErr: status.Error(codes.InvalidArgument, "blog ID is required"),
		},
		{
			name: "store error",
			req: &blogpb.UpdateReq{
				Id:      &blogpb.UUID{Value: "123e4567-e89b-12d3-a456-426614174000"},
				Content: stringPtr("Updated Content"),
			},
			setupMock: func(mockStore *mocks.Store) {
				content := "Updated Content"
				mockStore.On("Update", mock.Anything, datastore.ID("123e4567-e89b-12d3-a456-426614174000"), (*string)(nil), &content).
					Return(errors.New("update error"))
			},
			expectedErr: status.Error(codes.Internal, "failed to update blog: update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := mocks.NewStore(t)
			tt.setupMock(mockStore)

			service := NewBlogService(mockStore)
			resp, err := service.Update(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestBlogService_Delete(t *testing.T) {
	tests := []struct {
		name        string
		req         *blogpb.DeleteReq
		setupMock   func(mock *mocks.Store)
		expectedErr error
	}{
		{
			name: "successful delete",
			req: &blogpb.DeleteReq{
				Id: &blogpb.UUID{Value: "123e4567-e89b-12d3-a456-426614174000"},
			},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("Delete", mock.Anything, datastore.ID("123e4567-e89b-12d3-a456-426614174000")).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "missing ID",
			req:  &blogpb.DeleteReq{},
			setupMock: func(mockStore *mocks.Store) {
				// No mock setup needed
			},
			expectedErr: status.Error(codes.InvalidArgument, "blog ID is required"),
		},
		{
			name: "store error",
			req: &blogpb.DeleteReq{
				Id: &blogpb.UUID{Value: "123e4567-e89b-12d3-a456-426614174000"},
			},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("Delete", mock.Anything, datastore.ID("123e4567-e89b-12d3-a456-426614174000")).
					Return(errors.New("delete error"))
			},
			expectedErr: status.Error(codes.Internal, "failed to delete blog: delete error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := mocks.NewStore(t)
			tt.setupMock(mockStore)

			service := NewBlogService(mockStore)
			resp, err := service.Delete(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestBlogService_List(t *testing.T) {
	testSummaries := []*datastore.BlogSummary{
		{
			ID:           datastore.ID("123e4567-e89b-12d3-a456-426614174000"),
			Title:        "Test Blog 1",
			CommentCount: 5,
		},
		{
			ID:           datastore.ID("223e4567-e89b-12d3-a456-426614174000"),
			Title:        "Test Blog 2",
			CommentCount: 3,
		},
	}

	tests := []struct {
		name           string
		req            *blogpb.ListReq
		setupMock      func(mock *mocks.Store)
		expectedCount  int
		expectedToken  string
		expectedErr    error
		expectedValues func(resp *blogpb.ListResp)
	}{
		{
			name: "successful list with default page size",
			req:  &blogpb.ListReq{},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("List", mock.Anything, int32(10), "").
					Return(testSummaries, "next-token", nil)
			},
			expectedCount: 2,
			expectedToken: "next-token",
			expectedErr:   nil,
			expectedValues: func(resp *blogpb.ListResp) {
				assert.Equal(t, "Test Blog 1", resp.Blogs[0].Title)
				assert.Equal(t, int32(5), resp.Blogs[0].CommentCount)
				assert.Equal(t, "Test Blog 2", resp.Blogs[1].Title)
				assert.Equal(t, int32(3), resp.Blogs[1].CommentCount)
			},
		},
		{
			name: "successful list with custom page size",
			req:  &blogpb.ListReq{PageSize: 20},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("List", mock.Anything, int32(20), "").
					Return(testSummaries, "next-token", nil)
			},
			expectedCount: 2,
			expectedToken: "next-token",
			expectedErr:   nil,
			expectedValues: func(resp *blogpb.ListResp) {
				assert.Equal(t, "Test Blog 1", resp.Blogs[0].Title)
				assert.Equal(t, int32(5), resp.Blogs[0].CommentCount)
			},
		},
		{
			name: "successful list with page token",
			req:  &blogpb.ListReq{PageToken: "token-1"},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("List", mock.Anything, int32(10), "token-1").
					Return(testSummaries, "next-token", nil)
			},
			expectedCount: 2,
			expectedToken: "next-token",
			expectedErr:   nil,
			expectedValues: func(resp *blogpb.ListResp) {
				assert.Equal(t, "Test Blog 1", resp.Blogs[0].Title)
			},
		},
		{
			name: "page size too large",
			req:  &blogpb.ListReq{PageSize: 200},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("List", mock.Anything, int32(100), "").
					Return(testSummaries, "next-token", nil)
			},
			expectedCount: 2,
			expectedToken: "next-token",
			expectedErr:   nil,
		},
		{
			name: "invalid page size",
			req:  &blogpb.ListReq{PageSize: -10},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("List", mock.Anything, int32(10), "").
					Return(testSummaries, "next-token", nil)
			},
			expectedCount: 2,
			expectedToken: "next-token",
			expectedErr:   nil,
		},
		{
			name: "invalid page token",
			req:  &blogpb.ListReq{PageToken: "invalid-token"},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("List", mock.Anything, int32(10), "invalid-token").
					Return(nil, "", errors.New("invalid page token"))
			},
			expectedCount: 0,
			expectedToken: "",
			expectedErr:   status.Error(codes.Internal, "failed to list blogs: invalid page token"),
		},
		{
			name: "store error",
			req:  &blogpb.ListReq{},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("List", mock.Anything, int32(10), "").
					Return(nil, "", errors.New("list error"))
			},
			expectedCount: 0,
			expectedToken: "",
			expectedErr:   status.Error(codes.Internal, "failed to list blogs: list error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := mocks.NewStore(t)
			tt.setupMock(mockStore)

			service := NewBlogService(mockStore)
			resp, err := service.List(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedCount, len(resp.Blogs))
				assert.Equal(t, tt.expectedToken, resp.NextPageToken)
				if tt.expectedValues != nil {
					tt.expectedValues(resp)
				}
			}
		})
	}
}

func TestBlogService_AddComment(t *testing.T) {
	tests := []struct {
		name        string
		req         *blogpb.AddCommentReq
		setupMock   func(mock *mocks.Store)
		expectedErr error
	}{
		{
			name: "successful add comment",
			req: &blogpb.AddCommentReq{
				Id:      &blogpb.UUID{Value: "123e4567-e89b-12d3-a456-426614174000"},
				Content: "Test comment",
				Author:  "Test Author",
			},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("AddComment", mock.Anything, datastore.ID("123e4567-e89b-12d3-a456-426614174000"), "Test comment", "Test Author").
					Return(datastore.ID("comment-id"), nil)
			},
			expectedErr: nil,
		},
		{
			name: "missing ID",
			req: &blogpb.AddCommentReq{
				Content: "Test comment",
				Author:  "Test Author",
			},
			setupMock: func(mockStore *mocks.Store) {
				// No mock setup needed
			},
			expectedErr: status.Error(codes.InvalidArgument, "blog ID is required"),
		},
		{
			name: "missing Content",
			req: &blogpb.AddCommentReq{
				Id:     &blogpb.UUID{Value: "123e4567-e89b-12d3-a456-426614174000"},
				Author: "Test Author",
			},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("AddComment", mock.Anything, datastore.ID("123e4567-e89b-12d3-a456-426614174000"), "", "Test Author").
					Return(datastore.ID("comment-id"), nil)
			},
			expectedErr: nil,
		},
		{
			name: "missing Author",
			req: &blogpb.AddCommentReq{
				Id:      &blogpb.UUID{Value: "123e4567-e89b-12d3-a456-426614174000"},
				Content: "Test comment",
			},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("AddComment", mock.Anything, datastore.ID("123e4567-e89b-12d3-a456-426614174000"), "Test comment", "").
					Return(datastore.ID("comment-id"), nil)
			},
			expectedErr: nil,
		},
		{
			name: "store error",
			req: &blogpb.AddCommentReq{
				Id:      &blogpb.UUID{Value: "123e4567-e89b-12d3-a456-426614174000"},
				Content: "Test comment",
				Author:  "Test Author",
			},
			setupMock: func(mockStore *mocks.Store) {
				mockStore.On("AddComment", mock.Anything, datastore.ID("123e4567-e89b-12d3-a456-426614174000"), "Test comment", "Test Author").
					Return(datastore.ID(""), errors.New("comment error"))
			},
			expectedErr: status.Error(codes.Internal, "failed to add comment: comment error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := mocks.NewStore(t)
			tt.setupMock(mockStore)

			service := NewBlogService(mockStore)
			resp, err := service.AddComment(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
