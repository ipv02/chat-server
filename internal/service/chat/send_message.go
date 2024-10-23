package chat

import (
	"context"

	"github.com/ipv02/chat-server/internal/model"
)

func (s *serv) SendMessage(ctx context.Context, chat *model.ChatSendMessage) error {
	err := s.chatRepository.SendMessage(ctx, chat)

	return err
}
