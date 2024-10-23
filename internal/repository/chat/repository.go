package chat

import (
	"context"
	"database/sql"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ipv02/chat-server/internal/model"
	"github.com/ipv02/chat-server/internal/repository"
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

func (r *repo) CreateChat(ctx context.Context, chat *model.ChatCreate) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return 0, err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("tx.Rollback failed: %v", err)
		}
	}()

	builderChatInsert := sq.Insert(tableChatName).
		PlaceholderFormat(sq.Dollar).
		Columns(tableChatNameColumn).
		Values(chat.ChatName).
		Suffix("RETURNING id")

	query, args, err := builderChatInsert.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return 0, err
	}

	var chatID int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return 0, err
	}

	for _, userID := range chat.UsersId {
		builderChatUsersInsert := sq.Insert(tableChatUsersName).
			PlaceholderFormat(sq.Dollar).
			Columns(tableChatUsersChatIdColumn, tableChatUsersUserIdColumn).
			Values(userID, chatID)

		query, args, err := builderChatUsersInsert.ToSql()
		if err != nil {
			log.Printf("failed to build query: %v", err)
			return 0, err
		}

		_, err = r.db.Exec(ctx, query, args...)
		if err != nil {
			log.Printf("failed to execute query: %v", err)
			return 0, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("failed to commit transaction: %v", err)
		return 0, err
	}

	return chatID, nil
}

func (r *repo) DeleteChat(ctx context.Context, id int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("tx.Rollback failed: %v", err)
		}
	}()

	deleteChatBuilder := sq.Delete(tableChatName).
		Where(sq.Eq{tableChatIdColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := deleteChatBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return err
	}

	deleteChatUsersBuilder := sq.Delete(tableChatUsersName).
		Where(sq.Eq{tableChatUsersChatIdColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err = deleteChatUsersBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("failed to commit transaction: %v", err)
		return err
	}

	return err
}

func (r *repo) SendMessage(ctx context.Context, chat *model.ChatSendMessage) error {
	var messageID int64
	insertMessageBuilder := sq.Insert(tableMessagesName).
		PlaceholderFormat(sq.Dollar).
		Columns(tableMessagesUserIdColumn, tableMessagesMessageColumn, tableMessagesCreatedAtColumn).
		Values(chat.From, chat.Text, chat.Timestamp).
		Suffix("RETURNING id")

	query, args, err := insertMessageBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&messageID)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return err
	}

	return err
}
