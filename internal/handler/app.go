package handler

import (
	"SocialNetHL/internal/cache"
	"SocialNetHL/internal/queue"
	"SocialNetHL/internal/store"
)

type Instance struct {
	store        store.Store
	readStorages *store.ReadNodes
	Queue        queue.FeedQueue
	cache        cache.Cache
}

func NewInstance(writeStorage store.Store, readStorages *store.ReadNodes, rabbit queue.FeedQueue, cache cache.Cache) *Instance {
	return &Instance{
		store:        writeStorage,
		readStorages: readStorages,
		Queue:        rabbit,
		cache:        cache,
	}
}
