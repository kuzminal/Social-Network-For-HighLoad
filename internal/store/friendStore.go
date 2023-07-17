package store

import (
	"context"
)

type FriendStore interface {
	Store

	AddFriend(ctx context.Context, userId string, friendId string) error
	FindFriends(ctx context.Context, userId string) ([]string, error)
	DeleteFriend(ctx context.Context, userId string, friendId string) error
}
