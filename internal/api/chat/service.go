package chat

import (
	"github.com/ipv02/chat-server/internal/service"
	"github.com/ipv02/chat-server/pkg/chat_v1"
)

// Implementation структура описывающая сервер
type Implementation struct {
	chat_v1.UnimplementedChatV1Server
	chatService service.ChatService
}

// NewImplementation конструктор создает реализацию сервера и связывает ее с бизнес-логиклй
func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{
		chatService: chatService,
	}
}
