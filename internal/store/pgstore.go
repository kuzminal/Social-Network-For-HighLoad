package store

import (
	"SocialNetHL/models"
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"log"
	"os"
	"sync"
	"time"
)

type postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, connString string) (*postgres, error) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			log.Printf("unable to create connection pool: %v", err)
		}
		err = db.Ping(ctx)
		if err != nil {
			panic(err)
		}
		config := db.Config()
		pgInstance = &postgres{db}
		mdb, err := sql.Open("postgres", config.ConnString())
		err = mdb.Ping()
		if err != nil {
			panic(err)
		}
		migrationsDir := os.Getenv("MIGR_DIR")
		if len(migrationsDir) == 0 {
			migrationsDir = "./internal/migrations"
		}
		err = goose.Up(mdb, migrationsDir)
		if err != nil {
			panic(err)
		}

	})

	return pgInstance, nil
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() error {
	pg.db.Close()
	return nil
}

func (pg *postgres) SaveUser(ctx context.Context, user *models.RegisterUser) (id string, err error) {
	query := `INSERT INTO social.users (id, first_name, second_name, age, birthdate, biography, city, password) VALUES (@id, @firstName, @secondName, @age, @birthDate, @biography, @city, @password) RETURNING id`
	id = uuid.Must(uuid.NewV4()).String()
	bDate, _ := time.Parse("2006-01-02", user.Birthdate)
	age := calculateAge(bDate)
	args := pgx.NamedArgs{
		"id":         id,
		"firstName":  user.FirstName,
		"secondName": user.SecondName,
		"age":        age,
		"birthDate":  bDate,
		"biography":  user.Biography,
		"city":       user.City,
		"password":   fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password))),
	}
	_, err = pg.db.Exec(ctx, query, args)
	if err != nil {
		return "", fmt.Errorf("unable to insert row: %w", err)
	}
	return id, nil
}

func (pg *postgres) LoadUser(ctx context.Context, id string) (usersInfo models.UserInfo, err error) {
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

func (pg *postgres) LoadSession(ctx context.Context, token string) (string, error) {
	query := `SELECT user_id FROM social.session WHERE token = $1`

	row := pg.db.QueryRow(ctx, query, token)

	var userId string
	err := row.Scan(&userId)
	if err != nil {
		return "", fmt.Errorf("unable to scan row: %w", err)
	}
	return userId, nil
}

func (pg *postgres) CreateSession(ctx context.Context, m *models.AuthInfo) (string, error) {
	query := `INSERT INTO social.session (user_id, token, created_at) 
VALUES (@user_id, @token, @created_at)
ON CONFLICT (user_id) DO UPDATE
  SET created_at = now()
returning token;`
	authToken := uuid.Must(uuid.NewV4()).String()
	args := pgx.NamedArgs{
		"user_id":    m.Id,
		"token":      authToken,
		"created_at": time.Now(),
	}
	var token string
	_ = pg.db.QueryRow(ctx, query, args).Scan(&token)

	return token, nil
}

func calculateAge(bDate time.Time) int {
	curDate := time.Now()
	dur := curDate.Sub(bDate)
	return int(dur.Seconds() / 31207680)
}
