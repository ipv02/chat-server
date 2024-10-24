package chat

import (
	"context"
	"log"

	"github.com/ipv02/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteChat запрос для удаления чата.
func (i *Implementation) DeleteChat(ctx context.Context, req *chat_v1.DeleteChatRequest) (*emptypb.Empty, error) {
	err := i.chatService.DeleteChat(ctx, req.Id)
	if err != nil {
		log.Printf("failed to delete chat: %v", err)
		return nil, err
	}

	log.Printf("deleted chat: %v", req.Id)

	return &emptypb.Empty{}, nil
}
