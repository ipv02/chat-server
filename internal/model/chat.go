package model

import "google.golang.org/protobuf/types/known/timestamppb"

// ChatCreate модель для конвертации из протомодели в модель бизнес-логики
type ChatCreate struct {
	UsersID  []string
	ChatName string
}

// ChatSendMessage модель для конвертации из протомодели в модель бизнес-логики
type ChatSendMessage struct {
	From      string
	Text      string
	Timestamp *timestamppb.Timestamp
}
