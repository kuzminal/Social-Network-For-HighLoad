package service

import (
	"SocialNetHL/internal/cache"
	"SocialNetHL/internal/queue"
	"SocialNetHL/internal/store"
	"SocialNetHL/models"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type ActiveClients struct {
	sync.Mutex
	Clients map[string]*websocket.Conn
}

var clients = ActiveClients{Clients: map[string]*websocket.Conn{}}

type FeedService struct {
	store            store.Store
	cacheDb          cache.Cache
	queue            queue.FeedQueue
	postChan         chan models.Post
	cacheUpdateChan  chan models.UpdateFeedRequest
	connectToWs      chan models.ActiveWsUsers
	disconnectFromWs chan models.ActiveWsUsers
}

func NewFeedService(
	cacheDb cache.Cache,
	queue queue.FeedQueue,
	postChan chan models.Post,
	cacheUpdateChan chan models.UpdateFeedRequest,
	store store.Store,
	connectToWs chan models.ActiveWsUsers,
	disconnectFromWs chan models.ActiveWsUsers,
) *FeedService {
	return &FeedService{
		cacheDb:          cacheDb,
		queue:            queue,
		postChan:         postChan,
		cacheUpdateChan:  cacheUpdateChan,
		store:            store,
		connectToWs:      connectToWs,
		disconnectFromWs: disconnectFromWs,
	}
}

func (f *FeedService) FindFriendsForPost() {
	for {
		select {
		case d := <-f.postChan:
			friends, _ := f.store.FindFriends(context.Background(), d.AuthorUserId)
			for _, fr := range friends {
				_ = f.queue.SendFriendToUpdateFeed(context.Background(), models.UpdateFeedRequest{UserId: fr, Post: d})
				go sendToWs(fr, d)
			}
		}
	}
}

func (f *FeedService) UpdateCacheForFriends() {
	for {
		select {
		case d := <-f.cacheUpdateChan:
			db, err := f.cacheDb.UpdateCacheFromDb(d.UserId)
			if err != nil {
				return
			}
			log.Printf("Update data : %v", db)
		}
	}
}

func (f *FeedService) AddActiveClient(ch chan models.ActiveWsUsers) {
	for {
		select {
		case d := <-ch:
			clients.Mutex.Lock()
			clients.Clients[d.User] = d.Conn
			clients.Mutex.Unlock()
		}
	}
}

func (f *FeedService) DeleteActiveClient(ch chan models.ActiveWsUsers) {
	for {
		select {
		case d := <-ch:
			clients.Mutex.Lock()
			delete(clients.Clients, d.User)
			clients.Mutex.Unlock()
		}
	}
}
func sendToWs(userId string, post models.Post) {
	clients.Mutex.Lock()
	if conn, ok := clients.Clients[userId]; ok {
		p, err := json.Marshal(post)
		if err != nil {
			log.Println("Could not marshal post")
		}
		err = conn.WriteMessage(websocket.TextMessage, p)
		if err != nil {
			log.Println("Could not sent message to ws")
		}
	}
	clients.Mutex.Unlock()
}
