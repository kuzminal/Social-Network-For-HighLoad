package cache

import (
	"SocialNetHL/internal/helper"
	"SocialNetHL/models"
	"fmt"
	"github.com/tarantool/go-tarantool"
	"sync"
)

type Tarantool struct {
	conn *tarantool.Connection
}

var (
	tarantoolInstance *Tarantool
	tarantoolOnce     sync.Once
)

func NewTarantool() (*Tarantool, error) {
	thost := helper.GetEnvValue("TARANTOOL_HOST", "localhost")
	tuser := helper.GetEnvValue("TARANTOOL_USER", "user")
	tpassword := helper.GetEnvValue("TARANTOOL_PASSWORD", "password")
	tarantoolOnce.Do(func() {
		opts := tarantool.Opts{
			User: tuser,
			Pass: tpassword,
			/*Timeout:       2500 * time.Millisecond,
			Reconnect:     1 * time.Second,
			MaxReconnects: 3,*/
		}
		connStr := fmt.Sprintf("%s:3301", thost)
		conn, err := tarantool.Connect(connStr, opts)
		if err != nil {
			fmt.Println("Connection refused:", err)
		}
		tarantoolInstance = &Tarantool{conn: conn}
	})
	return tarantoolInstance, nil
}

func (t *Tarantool) GetFeed(key string) ([]models.Cache, error) {
	var posts []models.Cache
	err := t.conn.SelectTyped("posts", "primary", 0, 1, tarantool.IterEq, tarantool.StringKey{S: key}, &posts)
	if err != nil {
		fmt.Println(err)
	}

	return posts, nil
}

func (t *Tarantool) GetData(key string, offset int, limit int) models.Feed {
	var feed models.Feed
	err := t.conn.CallTyped("get_data", []interface{}{key, offset, limit}, &feed)
	if err != nil {
		return models.Feed{Posts: []models.Post{}}
	}
	return feed
}

func (t *Tarantool) PutFeed(key string, value []models.Post) error {
	_, err := t.conn.Replace("posts", []interface{}{key, &value})
	if err != nil {
		return err
	}
	return nil
}

func (t *Tarantool) UpdateCacheFromDb(key string) ([]models.Cache, error) {
	var cache []models.Cache
	err := t.conn.CallTyped("update_cache_from_db", []interface{}{key}, &cache)
	if err != nil {
		return []models.Cache{}, err
	}
	return cache, nil
}
