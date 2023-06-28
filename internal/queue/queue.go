package queue

import (
	"SocialNetHL/models"
	"context"
	"io"
)

type FeedQueue interface {
	io.Closer

	SendPostToFeed(ctx context.Context, post models.Post) error
	GetPostForFeed(ch chan models.Post)
	GetFriendsForUpdateFeed(ch chan models.UpdateFeedRequest)
	SendFriendToUpdateFeed(ctx context.Context, req models.UpdateFeedRequest) error
}
