package main

import (
	"SocialNetHL/internal/handler"
	_ "SocialNetHL/internal/migrations/go"
	"SocialNetHL/internal/router"
	"SocialNetHL/internal/store"
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	pghost := "localhost"
	pgport := "5432"
	if len(os.Getenv("PGHOST")) > 0 {
		pghost = os.Getenv("PGHOST")
	}
	if len(os.Getenv("PGPORT")) > 0 {
		pgport = os.Getenv("PGPORT")
	}
	storage, _ := store.NewPG(context.Background(), fmt.Sprintf("postgresql://postgres:postgres@%s:%s/postgres?sslmode=disable", pghost, pgport))
	app := handler.NewInstance(storage)
	r := router.NewRouter(app)

	http.ListenAndServe(":8080", r)

}
