package service

import (
	"context"

	"github.com/ipv02/chat-server/internal/model"
)

type ChatService interface {
	CreateChat(ctx context.Context, chat *model.ChatCreate) (int64, error)
	DeleteChat(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, chat *model.ChatSendMessage) error
}
