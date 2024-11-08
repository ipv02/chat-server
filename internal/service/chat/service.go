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

// NewMockService мок конструктор для создания связи между сервисным слоем и репо слоем
func NewMockService(deps ...interface{}) chatService.ChatService {
	service := service{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.ChatRepository:
			service.chatRepository = s
		case db.TxManager:
			service.txManager = s
		}
	}

	if service.chatRepository == nil {
		panic("chatRepository должен быть инициализирован")
	}

	if service.txManager == nil {
		panic("txManager должен быть инициализирован")
	}

	return &service
}
