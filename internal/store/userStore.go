package store

import (
	"SocialNetHL/models"
	"context"
)

type UserStore interface {
	Store

	SaveUser(ctx context.Context, registerUser models.RegisterUser) (id string, err error)
	LoadUser(ctx context.Context, id string) (userInfo models.UserInfo, err error)
	SearchUser(ctx context.Context, request models.UserSearchRequest) (users []models.UserInfo, err error)
	CheckIfExistsUser(ctx context.Context, userId string) (bool, error)
}
