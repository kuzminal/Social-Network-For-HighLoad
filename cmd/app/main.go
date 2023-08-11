package main

import (
	"SocialNetHL/internal/cache"
	"SocialNetHL/internal/handler"
	"SocialNetHL/internal/helper"
	_ "SocialNetHL/internal/migrations/go"
	"SocialNetHL/internal/queue"
	"SocialNetHL/internal/router"
	"SocialNetHL/internal/service"
	"SocialNetHL/internal/session"
	"SocialNetHL/internal/store"
	"SocialNetHL/internal/store/pg"
	"SocialNetHL/internal/store/tarantool"
	"SocialNetHL/internal/tracing"
	"SocialNetHL/models"
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

var (
	master             *pg.Postgres
	slave              *pg.Postgres
	tarantoolMaster    *tarantool.TarantoolStore
	tarantoolReadNodes store.ReadNodes[store.UserStore]
	queues             *queue.Rabbit
)

func main() {
	port := helper.GetEnvValue("PORT", "8080")

	traceServer := helper.GetEnvValue("TRACE_SERVER", "trace")
	tracePort := helper.GetEnvValue("TRACE_PORT", "14268")
	tracer, err := tracing.TracerProvider(fmt.Sprintf("http://%s:%s/api/traces", traceServer, tracePort))
	if err != nil {
		log.Fatal(err)
	}
	defer tracer.Shutdown(context.Background())

	initDb()
	initQueue()
	initTarantoolDb()
	tarant, _ := cache.NewTarantool()

	connectToWsChan := make(chan models.ActiveWsUsers, 10)
	disconnectToWsChan := make(chan models.ActiveWsUsers, 10)
	sessionPublisher := session.NewSessionPublisher()
	app := handler.NewInstance(
		tarantoolMaster,
		tarantoolMaster,
		master,
		slave,
		master,
		&tarantoolReadNodes,
		queues,
		tarant,
		connectToWsChan,
		disconnectToWsChan,
		sessionPublisher,
	)
	r := router.NewRouter(app)
	go service.NewTokenServiceServer(tarantoolMaster, tracer)

	go tarantoolReadNodes.HealthCheck()

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

	log.Printf("Starting http serer on port: %v", port)
	log.Fatalln(http.ListenAndServe(":"+port, r))
}

func initDb() {
	pgHost := helper.GetEnvValue("PGHOST", "localhost")
	pgPort := helper.GetEnvValue("PGPORT", "5432")
	pgSlaveHost := helper.GetEnvValue("PG_SLAVE_HOST", "localhost")
	pgSlavePort := helper.GetEnvValue("PG_SLAVE_PORT", "5000")
	pgUser := helper.GetEnvValue("PGUSER", "user")
	pgPassword := helper.GetEnvValue("PGPASSWORD", "password")
	pgDbName := helper.GetEnvValue("PGDBNAME", "social")
	master, _ = pg.NewMaster(
		context.Background(),
		fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPassword, pgHost, pgPort, pgDbName))
	slave, _ = pg.NewSlave(
		context.Background(),
		fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPassword, pgSlaveHost, pgSlavePort, pgDbName))
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
	var nodes []store.Backend[store.UserStore]
	var storage store.UserStore
	for _, str := range readConnStr {
		hosts := strings.Split(str, ":")
		storage, _ = tarantool.NewTarantoolSlave(hosts[0], hosts[1], "user", "password")
		nodes = append(nodes, store.Backend[store.UserStore]{
			Id:     uuid.Must(uuid.NewV4()).String(),
			IsDead: false,
			Store:  storage,
		})
	}
	tarantoolReadNodes = store.NewReadNode[store.UserStore]()
	tarantoolReadNodes.Current = 0
	tarantoolReadNodes.Nodes = nodes
}
