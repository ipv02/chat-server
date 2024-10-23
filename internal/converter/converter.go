package converter

import (
	"github.com/ipv02/chat-server/internal/model"
	"github.com/ipv02/chat-server/pkg/chat_v1"
)

func ToChatCreateFromReq(chat *chat_v1.CreateChatRequest) *model.ChatCreate {
	return &model.ChatCreate{
		UsersId:  chat.UsersId,
		ChatName: chat.ChatName,
	}
}

func ToChatSendMessage(chat *chat_v1.SendMessageRequest) *model.ChatSendMessage {
	return &model.ChatSendMessage{
		From:      chat.From,
		Text:      chat.Text,
		Timestamp: *chat.Timestamp,
	}
}
