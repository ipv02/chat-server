package main

import (
	"context"
	"flag"
	"github.com/ipv02/chat-server/internal/converter"
	"github.com/ipv02/chat-server/internal/service"
	"log"
	"net"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ipv02/chat-server/internal/config"
	"github.com/ipv02/chat-server/internal/config/env"
	chatRepository "github.com/ipv02/chat-server/internal/repository/chat"
	chatService "github.com/ipv02/chat-server/internal/service/chat"
	"github.com/ipv02/chat-server/pkg/chat_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	chat_v1.UnimplementedChatV1Server
	chatService service.ChatService
}

func main() {
	flag.Parse()
	ctx := context.Background()

	dbCtx, dbCancel := context.WithTimeout(ctx, 3*time.Second)
	defer dbCancel()

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

	pool, err := pgxpool.Connect(dbCtx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer pool.Close()

	chatRepo := chatRepository.NewRepository(pool)
	chatServ := chatService.NewService(chatRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	chat_v1.RegisterChatV1Server(s, &server{chatService: chatServ})

	log.Printf("server listening at: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
	}
}

// CreateChat запрос для создания нового чата.
func (s *server) CreateChat(ctx context.Context, req *chat_v1.CreateChatRequest) (*chat_v1.CreateChatResponse, error) {
	id, err := s.chatService.CreateChat(ctx, converter.ToChatCreateFromReq(req))
	if err != nil {
		log.Printf("failed to create chat: %v", err)
		return nil, err
	}

	log.Printf("created chat: %v", id)

	return &chat_v1.CreateChatResponse{
		Id: id,
	}, nil
}

// DeleteChat запрос для удаления чата.
func (s *server) DeleteChat(ctx context.Context, req *chat_v1.DeleteChatRequest) (*emptypb.Empty, error) {
	err := s.chatService.DeleteChat(ctx, req.Id)
	if err != nil {
		log.Printf("failed to delete chat: %v", err)
		return nil, err
	}

	log.Printf("deleted chat: %v", req.Id)

	return &emptypb.Empty{}, nil
}

// SendMessage запрос для отправки сообщения в чат.
func (s *server) SendMessage(ctx context.Context, req *chat_v1.SendMessageRequest) (*emptypb.Empty, error) {
	err := s.chatService.SendMessage(ctx, converter.ToChatSendMessage(req))
	if err != nil {
		log.Printf("failed to send message: %v", err)
		return nil, err
	}

	log.Printf("sent message: %v", req)

	return &emptypb.Empty{}, nil
}
