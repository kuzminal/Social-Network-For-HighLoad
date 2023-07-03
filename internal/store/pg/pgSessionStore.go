package pg

import (
	"SocialNetHL/models"
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"time"
)

func (pg *Postgres) LoadSession(ctx context.Context, token string) (string, error) {
	query := `SELECT user_id FROM social.session WHERE token = $1`
	//cont, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	//defer cancel()
	row := pg.db.QueryRow(ctx, query, token)
	var userId string
	err := row.Scan(&userId)
	if err != nil {
		return "", fmt.Errorf("unable to scan row: %w", err)
	}
	return userId, nil
}

func (pg *Postgres) CreateSession(ctx context.Context, m *models.AuthInfo) (string, error) {
	query := `INSERT INTO social.session (user_id, token, created_at) 
VALUES ($1, $2, $3)
ON CONFLICT (user_id) DO UPDATE
  SET created_at = now()
returning token;`
	authToken := uuid.Must(uuid.NewV4()).String()
	var token string
	_ = pg.db.QueryRow(ctx, query, m.Id, authToken, time.Now()).Scan(&token)

	return token, nil
}
