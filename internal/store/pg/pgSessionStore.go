package pg

import (
	"SocialNetHL/models"
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"time"
)

func (pg *Postgres) LoadSession(ctx context.Context, token string) (models.UserSession, error) {
	query := `SELECT id, user_id, token, created_at FROM social.session WHERE token = $1`
	//cont, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	//defer cancel()
	row := pg.db.QueryRow(ctx, query, token)
	var userSession models.UserSession
	err := row.Scan(&userSession)
	if err != nil {
		return userSession, fmt.Errorf("unable to scan row: %w", err)
	}
	return userSession, nil
}

func (pg *Postgres) CreateSession(ctx context.Context, m *models.AuthInfo) (models.UserSession, error) {
	query := `INSERT INTO social.session (user_id, token, created_at) 
VALUES ($1, $2, $3)
ON CONFLICT (user_id) DO UPDATE
  SET created_at = now()
returning id, user_id, token, created_at;`
	authToken := uuid.Must(uuid.NewV4()).String()
	var userSession models.UserSession
	_ = pg.db.QueryRow(ctx, query, m.Id, authToken, time.Now()).Scan(&userSession)

	return userSession, nil
}
