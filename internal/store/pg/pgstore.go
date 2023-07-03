package pg

import (
	"SocialNetHL/internal/helper"
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"log"
	"sync"
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
		conf.MaxConns = 90
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
	conf.MaxConns = 90

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
