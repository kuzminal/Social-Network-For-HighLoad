package handler

import (
	"SocialNetHL/internal/cache"
	"SocialNetHL/internal/queue"
	"SocialNetHL/internal/session"
	"SocialNetHL/internal/store"
	"SocialNetHL/models"
)

type Instance struct {
	userStore        store.UserStore
	sessionStore     store.SessionStore
	postStore        store.PostStore
	friendStore      store.FriendStore
	readStorages     *store.ReadNodes
	Queue            queue.FeedQueue
	cache            cache.Cache
	connectToWs      chan models.ActiveWsUsers
	disconnectFromWs chan models.ActiveWsUsers
	sessionPublisher session.Publisher
}

func NewInstance(
	userStore store.UserStore,
	sessionStore store.SessionStore,
	postStore store.PostStore,
	friendStore store.FriendStore,
	readStorages *store.ReadNodes,
	rabbit queue.FeedQueue,
	cache cache.Cache,
	connectToWs chan models.ActiveWsUsers,
	disconnectFromWs chan models.ActiveWsUsers,
	sessionPublisher session.Publisher,
) *Instance {
	return &Instance{
		userStore:        userStore,
		sessionStore:     sessionStore,
		postStore:        postStore,
		friendStore:      friendStore,
		readStorages:     readStorages,
		Queue:            rabbit,
		cache:            cache,
		connectToWs:      connectToWs,
		disconnectFromWs: disconnectFromWs,
		sessionPublisher: sessionPublisher,
	}
}
