package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ipv02/chat-server/config"
	"github.com/ipv02/chat-server/config/env"
	"github.com/ipv02/chat-server/pkg/chat_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	chat_v1.UnimplementedChatV1Server
	pool *pgxpool.Pool
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

	s := grpc.NewServer()
	reflection.Register(s)
	chat_v1.RegisterChatV1Server(s, &server{pool: pool})

	log.Printf("server listening at: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
	}
}

// CreateChat запрос для создания нового чата.
func (s *server) CreateChat(ctx context.Context, req *chat_v1.CreateChatRequest) (*chat_v1.CreateChatResponse, error) {
	log.Printf("CreateRequest: Users IDs: %v, Chat Name: %s", req.UsersId, req.ChatName)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return nil, err
	}

	errRollback := tx.Rollback(ctx)
	if errRollback != nil {
		log.Printf("failed to rollback transaction: %v", errRollback)
		return nil, errRollback
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("tx.Rollback failed: %v", err)
		}
	}()

	builderChatInsert := sq.Insert("chat").
		PlaceholderFormat(sq.Dollar).
		Columns("name").
		Values(req.ChatName).
		Suffix("RETURNING id")

	query, args, err := builderChatInsert.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	var chatID int64
	err = tx.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return nil, err
	}

	for _, userID := range req.UsersId {
		builderChatUsersInsert := sq.Insert("chat_users").
			PlaceholderFormat(sq.Dollar).
			Columns("chat_id", "user_id").
			Values(userID, chatID)

		query, args, err := builderChatUsersInsert.ToSql()
		if err != nil {
			log.Printf("failed to build query: %v", err)
			return nil, err
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			log.Printf("failed to execute query: %v", err)
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("failed to commit transaction: %v", err)
		return nil, err
	}

	return &chat_v1.CreateChatResponse{
		Id: chatID,
	}, nil
}

// DeleteChat запрос для удаления чата.
func (s *server) DeleteChat(ctx context.Context, req *chat_v1.DeleteChatRequest) (*emptypb.Empty, error) {
	log.Printf("Deleting object with ID: %d", req.GetId())

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("tx.Rollback failed: %v", err)
		}
	}()

	deleteChatBuilder := sq.Delete("chat").
		Where(sq.Eq{"id": req.Id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := deleteChatBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return nil, err
	}

	deleteChatUsersBuilder := sq.Delete("chat_users").
		Where(sq.Eq{"chat_id": req.Id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err = deleteChatUsersBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("failed to commit transaction: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// SendMessage запрос для отправки сообщения в чат.
func (s *server) SendMessage(ctx context.Context, req *chat_v1.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("SendMessageRequest - From: %v, Text: %v, Timestamp: %v", req.From, req.Text, req.Timestamp)

	var messageID int64
	insertMessageBuilder := sq.Insert("messages").
		PlaceholderFormat(sq.Dollar).
		Columns("user_id", "message", "created_at").
		Values(req.From, req.Text, req.Timestamp).
		Suffix("RETURNING id")

	query, args, err := insertMessageBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	err = s.pool.QueryRow(ctx, query, args...).Scan(&messageID)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
