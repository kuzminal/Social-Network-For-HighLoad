package store

import (
	"SocialNetHL/models"
	"context"
)

type PostStore interface {
	Store

	AddPost(ctx context.Context, post models.Post) (string, error)
	DeletePost(ctx context.Context, userId string, postId string) error
	UpdatePost(ctx context.Context, post models.Post) error
	GetPost(ctx context.Context, postId string) (models.Post, error)
	FeedPost(ctx context.Context, offset int, limit int, userId string) ([]models.Post, error)
}
