package cache

import "SocialNetHL/models"

type Cache interface {
	PutFeed(key string, value []models.Post) error
	GetFeed(key string) ([]models.Cache, error)
	GetData(key string, offset int, limit int) models.Feed
	UpdateCacheFromDb(key string) ([]models.Cache, error)
}
