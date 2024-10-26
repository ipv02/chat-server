package chat

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/ipv02/chat-server/internal/client/db"
	"github.com/ipv02/chat-server/internal/model"
	"github.com/ipv02/chat-server/internal/repository"
)

const (
	tableChatName       = "chat"
	tableChatIDColumn   = "id"
	tableChatNameColumn = "name"

	tableChatUsersName         = "chat_users"
	tableChatUsersChatIDColumn = "chat_id"
	tableChatUsersUserIDColumn = "user_id"

	tableMessagesName            = "messages"
	tableMessagesUserIDColumn    = "user_id"
	tableMessagesMessageColumn   = "message"
	tableMessagesCreatedAtColumn = "created_at"
)

type repo struct {
	db db.Client
}

// NewRepository создает новый экземпляр UserRepository с подключением к базе данных
func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

func (r *repo) CreateChat(ctx context.Context, chat *model.ChatCreate) (int64, error) {
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

	q := db.Query{
		Name:     "chat_repository.Create",
		QueryRaw: query,
	}

	var chatID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return 0, err
	}

	for _, userID := range chat.UsersID {
		builderChatUsersInsert := sq.Insert(tableChatUsersName).
			PlaceholderFormat(sq.Dollar).
			Columns(tableChatUsersChatIDColumn, tableChatUsersUserIDColumn).
			Values(userID, chatID)

		query, args, err := builderChatUsersInsert.ToSql()
		if err != nil {
			log.Printf("failed to build query: %v", err)
			return 0, err
		}

		q := db.Query{
			Name:     "chat_users_repository.Create",
			QueryRaw: query,
		}

		_, err = r.db.DB().ExecContext(ctx, q, args...)
		if err != nil {
			log.Printf("failed to execute query: %v", err)
			return 0, err
		}
	}

	return chatID, nil
}

func (r *repo) DeleteChat(ctx context.Context, id int64) error {
	deleteChatBuilder := sq.Delete(tableChatName).
		Where(sq.Eq{tableChatIDColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := deleteChatBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return err
	}

	q := db.Query{
		Name:     "chat_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return err
	}

	deleteChatUsersBuilder := sq.Delete(tableChatUsersName).
		Where(sq.Eq{tableChatUsersChatIDColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err = deleteChatUsersBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return err
	}

	q = db.Query{
		Name:     "chat_users_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return err
	}

	return err
}

func (r *repo) SendMessage(ctx context.Context, chat *model.ChatSendMessage) error {
	var messageID int64
	insertMessageBuilder := sq.Insert(tableMessagesName).
		PlaceholderFormat(sq.Dollar).
		Columns(tableMessagesUserIDColumn, tableMessagesMessageColumn, tableMessagesCreatedAtColumn).
		Values(chat.From, chat.Text, chat.Timestamp).
		Suffix("RETURNING id")

	query, args, err := insertMessageBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return err
	}

	q := db.Query{
		Name:     "chat_repository.SendMessage",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&messageID)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return err
	}

	return err
}
