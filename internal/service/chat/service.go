package chat

import (
	"github.com/ipv02/chat-server/internal/client/db"
	"github.com/ipv02/chat-server/internal/repository"
	"github.com/ipv02/chat-server/internal/service"
)

type serv struct {
	chatRepository repository.ChatRepository
	txManager      db.TxManager
}

// NewService конструктор для создания связи между сервисным слоем и репо слоем
func NewService(chatRepository repository.ChatRepository, txManager db.TxManager) service.ChatService {
	return &serv{
		chatRepository: chatRepository,
		txManager:      txManager,
	}
}
