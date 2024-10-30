package chat

import (
	"github.com/ipv02/chat-server/internal/client/db"
	"github.com/ipv02/chat-server/internal/repository"
	chatService "github.com/ipv02/chat-server/internal/service"
)

type service struct {
	chatRepository repository.ChatRepository
	txManager      db.TxManager
}

// NewService конструктор для создания связи между сервисным слоем и репо слоем
func NewService(chatRepository repository.ChatRepository, txManager db.TxManager) chatService.ChatService {
	return &service{
		chatRepository: chatRepository,
		txManager:      txManager,
	}
}
