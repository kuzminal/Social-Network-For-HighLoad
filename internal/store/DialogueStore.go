package store

import (
	"SocialNetHL/models"
	"context"
)

type DialogueStore interface {
	Store

	GetMessages(ctx context.Context, userFrom string, userTo string) ([]models.Message, error)
	SaveMessage(ctx context.Context, message models.Message) error
	GetChatId(ctx context.Context, userFrom string, userTo string) (string, error)
	CreateChat(ctx context.Context, fromUser string, toUser string) (string, error)
}
