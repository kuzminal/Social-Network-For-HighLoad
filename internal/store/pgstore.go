package store

import (
	"SocialNetHL/internal/helper"
	"SocialNetHL/models"
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"log"
	"os"
	"sync"
	"time"
)

type Postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

func NewMaster(ctx context.Context, connString string) (*Postgres, error) {
	pgOnce.Do(func() {
		conf, err := pgxpool.ParseConfig(connString)
		if err != nil {
			log.Printf("unable to create connection pool: %v", err)
		}
		conf.MaxConns = 990
		//conf.MinConns = 10

		db, err := pgxpool.ConnectConfig(ctx, conf)
		//db, err := pgxpool.New(ctx, connString)
		if err != nil {
			log.Printf("unable to create connection pool: %v", err)
		}
		err = db.Ping(ctx)
		if err != nil {
			panic(err)
		}
		config := db.Config()
		pgInstance = &Postgres{db}
		mdb, err := sql.Open("postgres", config.ConnString())
		err = mdb.Ping()
		if err != nil {
			panic(err)
		}
		migrationsDir := helper.GetEnvValue("MIGR_DIR", "./internal/migrations")
		err = goose.Up(mdb, migrationsDir)
		if err != nil {
			panic(err)
		}

	})

	return pgInstance, nil
}

func NewSlave(ctx context.Context, connString string) (*Postgres, error) {
	conf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Printf("unable to create connection pool: %v", err)
	}
	conf.MaxConns = 990

	db, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		log.Printf("unable to create connection pool: %v", err)
	}
	err = db.Ping(ctx)
	if err != nil {
		panic(err)
	}
	pgInstance = &Postgres{db}

	return pgInstance, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *Postgres) Close() error {
	pg.db.Close()
	return nil
}

func (pg *Postgres) SaveUser(ctx context.Context, user models.RegisterUser) (id string, err error) {
	query := `INSERT INTO social.users (id, first_name, second_name, age, birthdate, biography, city, password) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	id = uuid.Must(uuid.NewV4()).String()
	bDate, _ := time.Parse("2006-01-02", user.Birthdate)
	age := calculateAge(bDate)
	_, err = pg.db.Exec(ctx, query, id, user.FirstName, user.SecondName, age, bDate, user.Biography, user.City, fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password))))
	if err != nil {
		return "", fmt.Errorf("unable to insert row: %w", err)
	}
	os.WriteFile("filename", []byte(id+"\n"), 0644)
	return id, nil
}

func (pg *Postgres) LoadUser(ctx context.Context, id string) (usersInfo models.UserInfo, err error) {
	query := `SELECT id, first_name, second_name, age, birthdate, biography, city, password FROM social.users WHERE id = $1`

	row := pg.db.QueryRow(ctx, query, id)
	if err != nil {
		return models.UserInfo{}, fmt.Errorf("unable to query users: %w", err)
	}
	var bDate time.Time
	user := models.UserInfo{}
	err = row.Scan(&user.Id, &user.FirstName, &user.SecondName, &user.Age, &bDate, &user.Biography, &user.City, &user.Password)
	if err != nil {
		return models.UserInfo{}, fmt.Errorf("unable to scan row: %w", err)
	}
	user.Birthdate = bDate.Format("2006-01-02")
	return user, nil
}

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

func calculateAge(bDate time.Time) int {
	curDate := time.Now()
	dur := curDate.Sub(bDate)
	return int(dur.Seconds() / 31207680)
}

func (pg *Postgres) SearchUser(ctx context.Context, request models.UserSearchRequest) (users []models.UserInfo, err error) {
	query := `SELECT id, first_name, second_name, age, birthdate, biography, city, password FROM social.users WHERE first_name LIKE $1 AND second_name LIKE $2 ORDER BY id`
	//cont, cancel := context.WithTimeout(ctx, 2*time.Second)
	//defer cancel()
	rows, err := pg.db.Query(ctx, query, request.FirstName+"%", request.LastName+"%")
	defer rows.Close()
	if err != nil {
		return []models.UserInfo{}, fmt.Errorf("unable to query users: %w", err)
	}
	user := models.UserInfo{}
	var bDate time.Time
	for rows.Next() {
		//user := models.UserInfo{}
		err = rows.Scan(&user.Id, &user.FirstName, &user.SecondName, &user.Age, &bDate, &user.Biography, &user.City, &user.Password)
		if err != nil {
			log.Printf("unable to scan row: %v", err)
		}
		user.Birthdate = bDate.Format("2006-01-02")
		users = append(users, user)
	}
	return users, nil
}

func (pg *Postgres) CheckIfExistsUser(ctx context.Context, userId string) (bool, error) {
	query := `SELECT id FROM social.users WHERE id=$1`
	row := pg.db.QueryRow(ctx, query, userId)

	var user string
	err := row.Scan(&user)
	if err != nil {
		return false, err
	}
	if len(user) == 0 {
		return false, nil
	}
	return true, nil
}

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

func (pg *Postgres) GetMessages(ctx context.Context, userFrom string, userTo string) ([]models.Message, error) {
	var messages []models.Message
	query := `select id, from_user, to_user, text from social.messages where chat_id =$1 ORDER BY created_at DESC;`
	chatId, err := pg.GetChatId(ctx, userFrom, userTo)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rows, err := pg.db.Query(ctx, query, chatId)
	defer rows.Close()
	if err != nil {
		return []models.Message{}, fmt.Errorf("unable to query posts: %w", err)
	}
	message := models.Message{}
	for rows.Next() {
		err = rows.Scan(&message.Id, &message.FromUser, &message.ToUser, &message.Text)
		if err != nil {
			log.Printf("unable to scan row: %v", err)
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (pg *Postgres) GetChatId(ctx context.Context, userFrom string, userTo string) (string, error) {
	query := `select chat_id from social.chats where (user_from=$1 and user_to=$2) or (user_from=$2 and user_to=$1);`
	row := pg.db.QueryRow(ctx, query, userFrom, userTo)
	var chatId string
	err := row.Scan(&chatId)
	if err != nil && err.Error() == "no rows in result set" {
		id, _ := pg.CreateChat(ctx, userFrom, userTo)
		return id, nil
	} else if err != nil {
		return "", err
	}
	return chatId, nil
}

func (pg *Postgres) SaveMessage(ctx context.Context, msg models.Message) error {
	query := `INSERT INTO social.messages (id, "text", to_user, from_user, created_at, chat_id) 
VALUES ($1, $2, $3, $4, $5, $6);`
	_, err := pg.db.Exec(ctx, query, msg.Id, msg.Text, msg.ToUser, msg.FromUser, msg.CreatedAt, msg.ChatId)
	if err != nil {
		return err
	}
	return nil
}

func (pg *Postgres) CreateChat(ctx context.Context, fromUser string, toUser string) (string, error) {
	query := `INSERT INTO social.chats (chat_id, user_to, user_from) 
VALUES ($1, $2, $3);`
	chatId := uuid.Must(uuid.NewV4()).String()
	_, err := pg.db.Exec(ctx, query, chatId, fromUser, toUser)
	if err != nil {
		return "", err
	}
	return chatId, nil
}
