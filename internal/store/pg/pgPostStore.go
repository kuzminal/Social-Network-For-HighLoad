package pg

import (
	"SocialNetHL/models"
	"context"
	"fmt"
	"log"
	"time"
)

func (pg *Postgres) AddPost(ctx context.Context, post models.Post) (string, error) {
	query := `INSERT INTO social.posts (id, "text", author_user_id, created_at) 
VALUES ($1, $2, $3, $4) RETURNING id;`
	var postId string
	row := pg.db.QueryRow(ctx, query, post.Id, post.Text, post.AuthorUserId, time.Now())
	err := row.Scan(&postId)
	if err != nil {
		return "", err
	}
	return postId, nil
}

func (pg *Postgres) DeletePost(ctx context.Context, userId string, postId string) error {
	query := `DELETE FROM social.posts WHERE author_user_id=$1 and id=$2;`
	_, err := pg.db.Exec(ctx, query, userId, postId)
	return err
}

func (pg *Postgres) UpdatePost(ctx context.Context, post models.Post) error {
	query := `UPDATE social.posts SET "text"=$1, created_at=$2 WHERE id=$3;`
	_, err := pg.db.Exec(ctx, query, post.Text, time.Now(), post.Id)
	if err != nil {
		return fmt.Errorf("unable to scan row: %w", err)
	}
	return nil
}

func (pg *Postgres) GetPost(ctx context.Context, postId string) (models.Post, error) {
	query := `SELECT id, "text", author_user_id, created_at FROM social.posts WHERE id = $1;`

	row := pg.db.QueryRow(ctx, query, postId)

	post := models.Post{}
	var created time.Time
	err := row.Scan(&post.Id, &post.Text, &post.AuthorUserId, &created)
	if err != nil {
		return models.Post{}, fmt.Errorf("unable to scan row: %w", err)
	}
	post.CreatedAt = created.Format("2006-01-02")
	return post, nil
}

func (pg *Postgres) FeedPost(ctx context.Context, offset int, limit int, userId string) (posts []models.Post, err error) {
	query := `SELECT p.id, p."text", p.author_user_id, p.created_at
FROM social.friends f
RIGHT JOIN social.posts p ON p.author_user_id=f.friend_id
WHERE f.user_id=$1
OFFSET $2
LIMIT $3;`
	rows, err := pg.db.Query(ctx, query, userId, offset, limit)
	defer rows.Close()
	if err != nil {
		return []models.Post{}, fmt.Errorf("unable to query posts: %w", err)
	}
	post := models.Post{}
	var created time.Time
	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Text, &post.AuthorUserId, &created)
		if err != nil {
			log.Printf("unable to scan row: %v", err)
		}
		post.CreatedAt = created.Format("2006-01-02")
		posts = append(posts, post)
	}
	return posts, nil
}
