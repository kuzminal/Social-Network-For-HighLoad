package store

import (
	"SocialNetHL/models"
	"context"
	"io"
)

type Store interface {
	io.Closer

	SaveUser(ctx context.Context, registerUser models.RegisterUser) (id string, err error)
	LoadUser(ctx context.Context, id string) (userInfo models.UserInfo, err error)
	LoadSession(ctx context.Context, id string) (string, error)
	Ping(ctx context.Context) error
	CreateSession(ctx context.Context, m *models.AuthInfo) (string, error)
	SearchUser(ctx context.Context, request models.UserSearchRequest) (users []models.UserInfo, err error)
	CheckIfExistsUser(ctx context.Context, userId string) (bool, error)
	AddFriend(ctx context.Context, userId string, friendId string) error
	FindFriends(ctx context.Context, userId string) ([]string, error)
	DeleteFriend(ctx context.Context, userId string, friendId string) error
	AddPost(ctx context.Context, post models.Post) (string, error)
	DeletePost(ctx context.Context, userId string, postId string) error
	UpdatePost(ctx context.Context, post models.Post) error
	GetPost(ctx context.Context, postId string) (models.Post, error)
	FeedPost(ctx context.Context, offset int, limit int, userId string) ([]models.Post, error)
}
