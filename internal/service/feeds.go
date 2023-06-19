package service

import (
	"SocialNetHL/internal/cache"
	"SocialNetHL/internal/queue"
	"SocialNetHL/internal/store"
	"SocialNetHL/models"
	"context"
	"log"
)

type FeedService struct {
	store           store.Store
	cacheDb         cache.Cache
	queue           queue.FeedQueue
	postChan        chan models.Post
	cacheUpdateChan chan models.UpdateFeedCacheRequest
}

func NewFeedService(
	cacheDb cache.Cache,
	queue queue.FeedQueue,
	postChan chan models.Post,
	cacheUpdateChan chan models.UpdateFeedCacheRequest,
	store store.Store) *FeedService {
	return &FeedService{cacheDb: cacheDb, queue: queue, postChan: postChan, cacheUpdateChan: cacheUpdateChan, store: store}
}

func (f *FeedService) FindFriendsForPost() {
	for {
		select {
		case d := <-f.postChan:
			friends, _ := f.store.FindFriends(context.Background(), d.AuthorUserId)
			for _, fr := range friends {
				_ = f.queue.SendFriendToUpdateFeed(context.Background(), models.UpdateFeedCacheRequest{UserId: fr, Post: d})
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
