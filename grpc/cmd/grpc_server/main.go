package main

import (
	"context"
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ipv02/chat-server/config"
	"github.com/ipv02/chat-server/config/env"
	desc "github.com/ipv02/chat-server/grpc/pkg/chat_v1"
	"github.com/jackc/pgx/v4/pgxpool"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedChatV1Server
	pool *pgxpool.Pool
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{pool: pool})

	log.Printf("server listening at: %v", lis.Addr())

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
		Id: 1,
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
