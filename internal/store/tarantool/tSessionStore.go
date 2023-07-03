package tarantool

import (
	"SocialNetHL/models"
	"context"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"time"
)

type respLoad struct {
	UserId []string
}

type respCreate struct {
	Token []string
}

func (t *TarantoolStore) LoadSession(ctx context.Context, token string) (string, error) {
	var userId respLoad
	err := t.conn.CallTyped("get_session_by_user_id", []interface{}{token}, &userId)
	if err != nil {
		return "", err
	}
	if len(userId.UserId) != 1 {
		return "", errors.Errorf("Cannot find user with id: %s", userId)
	} else {
		return userId.UserId[0], nil
	}
}

func (t *TarantoolStore) CreateSession(ctx context.Context, m *models.AuthInfo) (string, error) {
	var token respCreate
	authToken := uuid.Must(uuid.NewV4()).String()
	sessionId := uuid.Must(uuid.NewV4()).String()
	createdAt := time.Now().Format(time.RFC3339)
	err := t.conn.CallTyped("create_session", []interface{}{sessionId, m.Id, authToken, createdAt}, &token)

	if err != nil {
		return "", err
	}
	if len(token.Token) != 1 {
		return "", errors.Errorf("Cannot create session user with id: %s", m.Id)
	} else {
		return token.Token[0], nil
	}
}
