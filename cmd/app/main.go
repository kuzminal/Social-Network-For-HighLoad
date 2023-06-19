package main

import (
	"SocialNetHL/internal/cache"
	"SocialNetHL/internal/handler"
	"SocialNetHL/internal/helper"
	_ "SocialNetHL/internal/migrations/go"
	"SocialNetHL/internal/queue"
	"SocialNetHL/internal/router"
	"SocialNetHL/internal/service"
	"SocialNetHL/internal/store"
	"SocialNetHL/models"
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

var (
	master    *store.Postgres
	readNodes store.ReadNodes
	queues    *queue.Rabbit
)

func main() {
	initDb()
	initQueue()
	tarant, _ := cache.NewTarantool()

	app := handler.NewInstance(master, &readNodes, queues, tarant)
	r := router.NewRouter(app)

	go store.HealthCheck(&readNodes)

	postChan := make(chan models.Post, 10)
	cacheCh := make(chan models.UpdateFeedCacheRequest, 10)
	go app.Queue.GetPostForFeed(postChan)
	go app.Queue.GetFriendsForUpdateFeed(cacheCh)

	feedService := service.NewFeedService(tarant, queues, postChan, cacheCh, master)
	go feedService.FindFriendsForPost()
	go feedService.UpdateCacheForFriends()

	http.ListenAndServe(":8080", r)
}

func initDb() {
	pghost := helper.GetEnvValue("PGHOST", "localhost")
	pgport := helper.GetEnvValue("PGPORT", "5432")
	readConnStr := strings.Split(
		helper.GetEnvValue("SLAVE_HOST_PORT",
			"localhost:5433,localhost:5434"),
		",")

	master, _ = store.NewMaster(context.Background(), fmt.Sprintf("postgresql://postgres:postgres@%s:%s/postgres?sslmode=disable", pghost, pgport))
	var nodes []store.Backend
	var storage store.Store
	for _, str := range readConnStr {
		hosts := strings.Split(str, ":")
		storage, _ = store.NewSlave(context.Background(), fmt.Sprintf("postgresql://postgres:postgres@%s:%s/postgres?sslmode=disable", hosts[0], hosts[1]))
		nodes = append(nodes, store.Backend{IsDead: false, Store: storage})
	}
	readNodes = store.ReadNodes{
		Current: 0,
		Nodes:   nodes,
	}
}

func initQueue() {
	queues, _ = queue.NewFeedQueue("amqp://user:password@localhost:5672/", "posts", "friends")
}
