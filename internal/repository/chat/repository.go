package chat

import (
	"context"
	"database/sql"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ipv02/chat-server/internal/repository"
	"github.com/ipv02/chat-server/pkg/chat_v1"
)

const (
	tableChatName       = "chat"
	tableChatIdColumn   = "id"
	tableChatNameColumn = "name"

	tableChatUsersName         = "chat_users"
	tableChatUsersChatIdColumn = "chat_id"
	tableChatUsersUserIdColumn = "user_id"

	tableMessagesName            = "messages"
	tableMessagesUserIdColumn    = "user_id"
	tableMessagesMessageColumn   = "message"
	tableMessagesCreatedAtColumn = "created_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.ChatRepository {
	return &repo{db: db}
}

func (r *repo) CreateChat(ctx context.Context, req *chat_v1.CreateChatRequest) (*chat_v1.CreateChatResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("tx.Rollback failed: %v", err)
		}
	}()

	builderChatInsert := sq.Insert(tableChatName).
		PlaceholderFormat(sq.Dollar).
		Columns(tableChatNameColumn).
		Values(req.ChatName).
		Suffix("RETURNING id")

	query, args, err := builderChatInsert.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	var chatID int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return nil, err
	}

	for _, userID := range req.UsersId {
		builderChatUsersInsert := sq.Insert(tableChatUsersName).
			PlaceholderFormat(sq.Dollar).
			Columns(tableChatUsersChatIdColumn, tableChatUsersUserIdColumn).
			Values(userID, chatID)

		query, args, err := builderChatUsersInsert.ToSql()
		if err != nil {
			log.Printf("failed to build query: %v", err)
			return nil, err
		}

		_, err = r.db.Exec(ctx, query, args...)
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

func (r *repo) DeleteChat(ctx context.Context, req *chat_v1.DeleteChatRequest) (*emptypb.Empty, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("tx.Rollback failed: %v", err)
		}
	}()

	deleteChatBuilder := sq.Delete(tableChatName).
		Where(sq.Eq{tableChatIdColumn: req.Id}).
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

	deleteChatUsersBuilder := sq.Delete(tableChatUsersName).
		Where(sq.Eq{tableChatUsersChatIdColumn: req.Id}).
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

func (r *repo) SendMessage(ctx context.Context, req *chat_v1.SendMessageRequest) (*emptypb.Empty, error) {
	var messageID int64
	insertMessageBuilder := sq.Insert(tableMessagesName).
		PlaceholderFormat(sq.Dollar).
		Columns(tableMessagesUserIdColumn, tableMessagesMessageColumn, tableMessagesCreatedAtColumn).
		Values(req.From, req.Text, req.Timestamp).
		Suffix("RETURNING id")

	query, args, err := insertMessageBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&messageID)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
