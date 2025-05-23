// Package service provides implementations of the gRPC services
package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/agruetz/prosigliere/internal/datastore"
	blogpb "github.com/agruetz/prosigliere/protos/v1/blog"
)

// BlogService implements the blog.v1.BlogsServer interface
type BlogService struct {
	blogpb.UnimplementedBlogsServer
	store datastore.Store
}

// NewBlogService creates a new BlogService with the given datastore
func NewBlogService(store datastore.Store) *BlogService {
	return &BlogService{
		store: store,
	}
}

// Create creates a new blog
func (s *BlogService) Create(ctx context.Context, req *blogpb.CreateReq) (*blogpb.CreateResp, error) {
	// Validate inputs
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if req.GetContent() == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}

	id, err := s.store.Create(ctx, req.GetTitle(), req.GetContent())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create blog: %v", err)
	}

	return &blogpb.CreateResp{
		Id: &blogpb.UUID{
			Value: string(id),
		},
	}, nil
}

// Get retrieves a blog by ID
func (s *BlogService) Get(ctx context.Context, req *blogpb.GetReq) (*blogpb.GetResp, error) {
	if req.GetId() == nil {
		return nil, status.Error(codes.InvalidArgument, "blog ID is required")
	}

	id := datastore.ID(req.GetId().GetValue())
	blog, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get blog: %v", err)
	}

	// Create a slice to hold valid comments
	comments := make([]*blogpb.Comment, 0, len(blog.Comments))
	for _, comment := range blog.Comments {
		// Only add comments with valid IDs
		if comment.ID != "" {
			comments = append(comments, &blogpb.Comment{
				Id:        &blogpb.UUID{Value: string(comment.ID)},
				Content:   comment.Content,
				Author:    comment.Author,
				CreatedAt: timestamppb.New(comment.CreatedAt),
			})
		}
	}

	return &blogpb.GetResp{
		Blog: &blogpb.Blog{
			Id: &blogpb.UUID{
				Value: string(blog.ID),
			},
			Title:     blog.Title,
			Content:   blog.Content,
			CreatedAt: timestamppb.New(blog.CreatedAt),
			UpdatedAt: timestamppb.New(blog.UpdatedAt),
			Comments:  comments,
		},
	}, nil
}

// Update updates an existing blog
func (s *BlogService) Update(ctx context.Context, req *blogpb.UpdateReq) (*emptypb.Empty, error) {
	if req.GetId() == nil {
		return nil, status.Error(codes.InvalidArgument, "blog ID is required")
	}

	id := datastore.ID(req.GetId().GetValue())
	var title, content *string

	// Handle optional fields
	if req.Title != nil {
		titleVal := req.GetTitle()
		title = &titleVal
	}
	if req.Content != nil {
		contentVal := req.GetContent()
		content = &contentVal
	}

	err := s.store.Update(ctx, id, title, content)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update blog: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// Delete deletes a blog
func (s *BlogService) Delete(ctx context.Context, req *blogpb.DeleteReq) (*emptypb.Empty, error) {
	if req.GetId() == nil {
		return nil, status.Error(codes.InvalidArgument, "blog ID is required")
	}

	id := datastore.ID(req.GetId().GetValue())
	err := s.store.Delete(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete blog: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// List lists blogs with pagination
func (s *BlogService) List(ctx context.Context, req *blogpb.ListReq) (*blogpb.ListResp, error) {
	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}
	if pageSize > 100 {
		pageSize = 100 // Maximum page size
	}

	summaries, nextPageToken, err := s.store.List(ctx, pageSize, req.GetPageToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list blogs: %v", err)
	}

	pbSummaries := make([]*blogpb.BlogSummary, len(summaries))
	for i, summary := range summaries {
		pbSummaries[i] = &blogpb.BlogSummary{
			Id: &blogpb.UUID{
				Value: string(summary.ID),
			},
			Title:        summary.Title,
			CommentCount: summary.CommentCount,
		}
	}

	return &blogpb.ListResp{
		Blogs:         pbSummaries,
		NextPageToken: nextPageToken,
	}, nil
}

// AddComment adds a comment to a blog
func (s *BlogService) AddComment(ctx context.Context, req *blogpb.AddCommentReq) (*emptypb.Empty, error) {
	if req.GetId() == nil {
		return nil, status.Error(codes.InvalidArgument, "blog ID is required")
	}

	// Validate inputs
	if req.GetContent() == "" {
		return nil, status.Error(codes.InvalidArgument, "comment content is required")
	}
	if req.GetAuthor() == "" {
		return nil, status.Error(codes.InvalidArgument, "comment author is required")
	}

	id := datastore.ID(req.GetId().GetValue())
	_, err := s.store.AddComment(ctx, id, req.GetContent(), req.GetAuthor())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add comment: %v", err)
	}

	return &emptypb.Empty{}, nil
}
