package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/ipv02/chat-server/grpc/pkg/chat_v1"
)

const grpcPort = 50052

type server struct {
	desc.UnimplementedChatV1Server
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{})

	log.Printf("server listening at: %v", grpcPort)

	if err = s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
	}
}

// CreateChat обрабатывает CreateChatRequest для создания нового чата.
//
// Логирует информацию о запросе (список идентификаторов пользователей и название чата)
// и возвращает CreateChatResponse с сгенерированным идентификатором чата.
//
// Параметры:
//   - _ctx: Контекст для управления временем жизни запроса и дедлайнами.
//   - req: Запрос CreateChatRequest, содержащий список идентификаторов пользователей и название чата.
//
// Возвращает:
//   - *desc.CreateChatResponse: Ответ с идентификатором созданного чата.
//   - error: Возвращает ошибку в случае неудачи, или nil при успешном выполнении.
func (s *server) CreateChat(_ context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	log.Printf("CreateRequest: Users IDs: %v, Chat Name: %s", req.UsersId, req.ChatName)

	return &desc.CreateChatResponse{
		Id: rand.Int63(),
	}, nil
}

// DeleteChat обрабатывает DeleteChatRequest для удаления чата по ID.
//
// Логирует идентификатор удаляемого чата и возвращает пустой ответ
// при успешном выполнении.
//
// Параметры:
//   - ctx: Контекст для управления временем жизни запроса и дедлайнами.
//   - req: Запрос DeleteChatRequest, содержащий ID чата для удаления.
//
// Возвращает:
//   - *emptypb.Empty: Пустой ответ при успешном выполнении операции.
//   - error: Возвращает ошибку в случае неудачи, или nil при успешном выполнении.
func (s *server) DeleteChat(_ context.Context, req *desc.DeleteChatRequest) (*emptypb.Empty, error) {
	log.Printf("Deleting object with ID: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

// SendMessage обрабатывает SendMessageRequest для отправки сообщения в чат.
//
// Логирует информацию о сообщении (отправитель, текст и метка времени) и возвращает пустой ответ
// при успешном выполнении.
//
// Параметры:
//   - _ctx: Контекст для управления временем жизни запроса и дедлайнами.
//   - req: Запрос SendMessageRequest, содержащий данные сообщения (отправитель, текст, метка времени).
//
// Возвращает:
//   - *emptypb.Empty: Пустой ответ при успешном выполнении операции.
//   - error: Возвращает ошибку в случае неудачи, или nil при успешном выполнении.
func (s *server) SendMessage(_ context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("SendMessageRequest - From: %v, Text: %v, Timestamp: %v", req.From, req.Text, req.Timestamp)

	return &emptypb.Empty{}, nil
}
