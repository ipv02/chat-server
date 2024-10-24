package chat

import (
	"context"
	"log"

	"github.com/ipv02/chat-server/internal/converter"
	"github.com/ipv02/chat-server/pkg/chat_v1"
)

// CreateChat запрос для создания нового чата.
func (i *Implementation) CreateChat(ctx context.Context, req *chat_v1.CreateChatRequest) (*chat_v1.CreateChatResponse, error) {
	id, err := i.chatService.CreateChat(ctx, converter.ToChatCreateFromReq(req))
	if err != nil {
		log.Printf("failed to create chat: %v", err)
		return nil, err
	}

	log.Printf("created chat: %v", id)

	return &chat_v1.CreateChatResponse{
		Id: id,
	}, nil
}
