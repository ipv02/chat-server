package chat

import (
	"context"

	"github.com/ipv02/chat-server/internal/model"
)

// CreateChat выполняет создание нового чата в сервисном слое
func (s *service) CreateChat(ctx context.Context, chat *model.ChatCreate) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.chatRepository.CreateChat(ctx, chat)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
