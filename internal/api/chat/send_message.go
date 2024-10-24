package chat

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ipv02/chat-server/internal/converter"
	"github.com/ipv02/chat-server/pkg/chat_v1"
)

// SendMessage запрос для отправки сообщения в чат.
func (i *Implementation) SendMessage(ctx context.Context, req *chat_v1.SendMessageRequest) (*emptypb.Empty, error) {
	err := i.chatService.SendMessage(ctx, converter.ToChatSendMessage(req))
	if err != nil {
		log.Printf("failed to send message: %v", err)
		return nil, err
	}

	log.Printf("sent message: %v", req)

	return &emptypb.Empty{}, nil
}
