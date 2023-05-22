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
}
