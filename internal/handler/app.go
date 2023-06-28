package handler

import (
	"SocialNetHL/internal/cache"
	"SocialNetHL/internal/queue"
	"SocialNetHL/internal/store"
	"SocialNetHL/models"
)

type Instance struct {
	store            store.Store
	readStorages     *store.ReadNodes
	Queue            queue.FeedQueue
	cache            cache.Cache
	connectToWs      chan models.ActiveWsUsers
	disconnectFromWs chan models.ActiveWsUsers
}

func NewInstance(
	writeStorage store.Store,
	readStorages *store.ReadNodes,
	rabbit queue.FeedQueue,
	cache cache.Cache,
	connectToWs chan models.ActiveWsUsers,
	disconnectFromWs chan models.ActiveWsUsers,
) *Instance {
	return &Instance{
		store:            writeStorage,
		readStorages:     readStorages,
		Queue:            rabbit,
		cache:            cache,
		connectToWs:      connectToWs,
		disconnectFromWs: disconnectFromWs,
	}
}
