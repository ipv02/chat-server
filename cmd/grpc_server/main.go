package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	chatAPI "github.com/ipv02/chat-server/internal/api/chat"
	"github.com/ipv02/chat-server/internal/config"
	"github.com/ipv02/chat-server/internal/config/env"
	chatRepository "github.com/ipv02/chat-server/internal/repository/chat"
	chatService "github.com/ipv02/chat-server/internal/service/chat"
	"github.com/ipv02/chat-server/pkg/chat_v1"
	"github.com/jackc/pgx/v4/pgxpool"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
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
	chat_v1.RegisterChatV1Server(s, chatAPI.NewImplementation(chatServ))

	log.Printf("server listening at: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
	}
}
