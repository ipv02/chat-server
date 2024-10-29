package converter

import (
	"log"

	"github.com/ipv02/chat-server/internal/model"
	"github.com/ipv02/chat-server/pkg/chat_v1"
)

// ToChatCreateFromReq конвертер протомодели в модель бизнес-логики
func ToChatCreateFromReq(chat *chat_v1.CreateChatRequest) *model.ChatCreate {
	if chat == nil {
		log.Println("ToChatCreateFromReq: nil chat, returning nil response")
		return nil
	}

	return &model.ChatCreate{
		UsersID:  chat.UsersId,
		ChatName: chat.ChatName,
	}
}

// ToChatSendMessage конвертер протомодели в модель бизнес-логики
func ToChatSendMessage(chat *chat_v1.SendMessageRequest) *model.ChatSendMessage {
	if chat == nil {
		log.Println("ToChatSendMessage: nil chat, returning nil response")
		return nil
	}

	return &model.ChatSendMessage{
		From:      chat.From,
		Text:      chat.Text,
		Timestamp: chat.Timestamp,
	}
}
