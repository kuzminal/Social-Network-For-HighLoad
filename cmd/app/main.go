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
	"SocialNetHL/internal/store/pg"
	"SocialNetHL/internal/store/tarantool"
	"SocialNetHL/models"
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

var (
	master             *pg.Postgres
	tarantoolMaster    *tarantool.TarantoolStore
	tarantoolReadNodes store.ReadNodes
	queues             *queue.Rabbit
)

func main() {
	initDb()
	initQueue()
	initTarantoolDb()
	tarant, _ := cache.NewTarantool()

	connectToWsChan := make(chan models.ActiveWsUsers, 10)
	disconnectToWsChan := make(chan models.ActiveWsUsers, 10)

	app := handler.NewInstance(
		tarantoolMaster,
		tarantoolMaster,
		master,
		master,
		master,
		&tarantoolReadNodes,
		queues,
		tarant,
		connectToWsChan,
		disconnectToWsChan,
	)
	r := router.NewRouter(app)

	go store.HealthCheck(&tarantoolReadNodes)

	postChan := make(chan models.Post, 10)
	cacheCh := make(chan models.UpdateFeedRequest, 10)

	go app.Queue.GetPostForFeed(postChan)
	go app.Queue.GetFriendsForUpdateFeed(cacheCh)

	feedService := service.NewFeedService(
		tarant,
		queues,
		postChan,
		cacheCh,
		master,
		connectToWsChan,
		disconnectToWsChan,
	)
	go feedService.FindFriendsForPost()
	go feedService.UpdateCacheForFriends()
	go feedService.AddActiveClient(connectToWsChan)
	go feedService.DeleteActiveClient(disconnectToWsChan)
	log.Fatalln(http.ListenAndServe(":8080", r))
}

func initDb() {
	pghost := helper.GetEnvValue("PGHOST", "localhost")
	pgport := helper.GetEnvValue("PGPORT", "5432")
	master, _ = pg.NewMaster(context.Background(), fmt.Sprintf("postgresql://postgres:postgres@%s:%s/postgres?sslmode=disable", pghost, pgport))
}

func initQueue() {
	rhost := helper.GetEnvValue("RABBIT_HOST", "localhost")
	ruser := helper.GetEnvValue("RABBIT_USER", "user")
	rpassword := helper.GetEnvValue("RABBIT_PASSWORD", "password")
	connStr := fmt.Sprintf("amqp://%s:%s@%s:5672/", ruser, rpassword, rhost)
	queues, _ = queue.NewFeedQueue(connStr, "posts", "friends")
}

func initTarantoolDb() {
	thost := helper.GetEnvValue("TARANTOOL_HOST", "localhost")
	tport := "3301" //пока так
	tuser := helper.GetEnvValue("TARANTOOL_USER_NAME", "user")
	tpassword := helper.GetEnvValue("TARANTOOL_USER_PASSWORD", "password")
	readConnStr := strings.Split(
		helper.GetEnvValue("SLAVE_HOST_PORT",
			"localhost:3302,localhost:3303"),
		",")

	tarantoolMaster, _ = tarantool.NewTarantoolMaster(thost, tport, tuser, tpassword)
	var nodes []store.Backend
	var storage store.UserStore
	for _, str := range readConnStr {
		hosts := strings.Split(str, ":")
		storage, _ = tarantool.NewTarantoolSlave(hosts[0], hosts[1], "user", "password")
		nodes = append(nodes, store.Backend{IsDead: false, Store: storage})
	}
	tarantoolReadNodes = store.ReadNodes{
		Current: 0,
		Nodes:   nodes,
	}
}
