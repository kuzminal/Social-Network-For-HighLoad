package store

import (
	"SocialNetHL/models"
	"context"
)

type SessionStore interface {
	Store

	LoadSession(ctx context.Context, id string) (models.UserSession, error)
	CreateSession(ctx context.Context, m *models.AuthInfo) (models.UserSession, error)
}
