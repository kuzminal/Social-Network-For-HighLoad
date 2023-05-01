package store

import (
	"SocialNetHL/models"
	"context"
	"io"
)

type Store interface {
	io.Closer

	SaveUser(ctx context.Context, registerUser *models.RegisterUser) (id string, err error)
	LoadUser(ctx context.Context, id string) (userInfo models.UserInfo, err error)
	Ping(ctx context.Context) error
}
