package tarantool

import (
	"SocialNetHL/models"
	"context"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

func (t *TarantoolStore) LoadSession(ctx context.Context, token string) (models.UserSession, error) {
	var session []models.UserSession
	err := t.conn.CallTyped("get_session_by_user_id", []interface{}{token}, &session)
	if err != nil {
		return models.UserSession{}, err
	}
	if len(session) != 1 {
		return models.UserSession{}, errors.Errorf("Cannot find user with id: %s", session)
	} else {
		return session[0], nil
	}
}

func (t *TarantoolStore) CreateSession(ctx context.Context, m *models.AuthInfo) (models.UserSession, error) {
	var session []models.UserSession
	authToken := uuid.Must(uuid.NewV4()).String()
	sessionId := uuid.Must(uuid.NewV4()).String()
	err := t.conn.CallTyped("create_session", []interface{}{sessionId, m.Id, authToken}, &session)
	if err != nil {
		return models.UserSession{}, err
	}
	if len(session) != 1 {
		return models.UserSession{}, errors.Errorf("Cannot create session user with id: %s", m.Id)
	} else {
		return session[0], nil
	}
}
