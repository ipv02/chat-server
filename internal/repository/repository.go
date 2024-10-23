package repository

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ipv02/chat-server/pkg/chat_v1"
)

type ChatRepository interface {
	CreateChat(ctx context.Context, req *chat_v1.CreateChatRequest) (*chat_v1.CreateChatResponse, error)
	DeleteChat(ctx context.Context, req *chat_v1.DeleteChatRequest) (*emptypb.Empty, error)
	SendMessage(ctx context.Context, req *chat_v1.SendMessageRequest) (*emptypb.Empty, error)
}
