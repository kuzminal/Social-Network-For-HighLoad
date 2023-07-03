package pg

import (
	"context"
	"log"
	"time"
)

func (pg *Postgres) AddFriend(ctx context.Context, userId string, friendId string) error {
	query := `INSERT INTO social.friends (user_id, friend_id, created_at) 
VALUES ($1, $2, $3);`
	_, err := pg.db.Exec(ctx, query, userId, friendId, time.Now())
	return err
}

func (pg *Postgres) FindFriends(ctx context.Context, userId string) ([]string, error) {
	query := `SELECT user_id FROM social.friends WHERE friend_id=$1;`
	rows, err := pg.db.Query(ctx, query, userId)
	var friends []string
	var friendId string
	for rows.Next() {
		err = rows.Scan(&friendId)
		if err != nil {
			log.Printf("unable to scan row: %v", err)
		}
		friends = append(friends, friendId)
	}
	return friends, err
}

func (pg *Postgres) DeleteFriend(ctx context.Context, userId string, friendId string) error {
	query := `DELETE FROM social.friends WHERE user_id=$1 and friend_id=$2;`
	_, err := pg.db.Exec(ctx, query, userId, friendId)
	return err
}
