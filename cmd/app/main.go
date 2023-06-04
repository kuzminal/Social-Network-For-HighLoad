package main

import (
	"SocialNetHL/internal/handler"
	"SocialNetHL/internal/helper"
	_ "SocialNetHL/internal/migrations/go"
	"SocialNetHL/internal/router"
	"SocialNetHL/internal/store"
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

func main() {
	pghost := helper.GetEnvValue("PGHOST", "localhost")
	pgport := helper.GetEnvValue("PGPORT", "5432")
	readConnStr := strings.Split(
		helper.GetEnvValue("SLAVE_HOST_PORT",
			"localhost:5433,localhost:5434"),
		",")

	master, _ := store.NewMaster(context.Background(), fmt.Sprintf("postgresql://postgres:postgres@%s:%s/postgres?sslmode=disable", pghost, pgport))
	var nodes []store.Backend
	var storage store.Store
	for _, str := range readConnStr {
		hosts := strings.Split(str, ":")
		storage, _ = store.NewSlave(context.Background(), fmt.Sprintf("postgresql://postgres:postgres@%s:%s/postgres?sslmode=disable", hosts[0], hosts[1]))
		nodes = append(nodes, store.Backend{IsDead: false, Store: storage})
	}
	readNodes := store.ReadNodes{
		Current: 0,
		Nodes:   nodes,
	}

	app := handler.NewInstance(master, &readNodes)
	r := router.NewRouter(app)
	go store.HealthCheck(&readNodes)
	http.ListenAndServe(":8080", r)

}
