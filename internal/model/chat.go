package model

import "google.golang.org/protobuf/types/known/timestamppb"

type ChatCreate struct {
	UsersId  []string
	ChatName string
}

type ChatSendMessage struct {
	From      string
	Text      string
	Timestamp timestamppb.Timestamp
}
