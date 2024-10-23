package chat

import (
	"github.com/ipv02/chat-server/internal/repository"
	"github.com/ipv02/chat-server/internal/service"
)

type serv struct {
	chatRepository repository.ChatRepository
}

func NewService(chatRepository repository.ChatRepository) service.ChatService {
	return &serv{chatRepository: chatRepository}
}
