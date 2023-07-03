package store

import (
	"SocialNetHL/models"
	"context"
)

type SessionStore interface {
	Store

	LoadSession(ctx context.Context, id string) (string, error)
	CreateSession(ctx context.Context, m *models.AuthInfo) (string, error)
}
