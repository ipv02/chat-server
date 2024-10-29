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

// CreateChat выполняет создание нового чата в базе данных
func (r *repo) CreateChat(ctx context.Context, chat *model.ChatCreate) (int64, error) {
	chatID, err := r.insertChat(ctx, chat.ChatName)
	if err != nil {
		return 0, err
	}

	if err := r.insertChatUsers(ctx, chatID, chat.UsersID); err != nil {
		return 0, err
	}

	return chatID, nil
}

// insertChat Вставка записи чата в таблицу chat
func (r *repo) insertChat(ctx context.Context, chatName string) (int64, error) {
	builderChatInsert := sq.Insert(tableChatName).
		Columns(tableChatNameColumn).
		Values(chatName).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := builderChatInsert.ToSql()
	if err != nil {
		log.Printf("failed to build chat insert query: %v", err)
		return 0, err
	}

	q := db.Query{
		Name:     "chat_repository.Create",
		QueryRaw: query,
	}

	var chatID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)
	if err != nil {
		log.Printf("failed to execute chat insert query: %v", err)
		return 0, err
	}

	return chatID, nil
}

// insertChatUsers Вставка пользователей в таблицу chat_users за одно обращение к БД
func (r *repo) insertChatUsers(ctx context.Context, chatID int64, userIDs []string) error {
	builderChatUsersInsert := sq.Insert(tableChatUsersName).
		Columns(tableChatUsersChatIDColumn, tableChatUsersUserIDColumn).
		PlaceholderFormat(sq.Dollar)

	for _, userID := range userIDs {
		builderChatUsersInsert = builderChatUsersInsert.Values(chatID, userID)
	}

	query, args, err := builderChatUsersInsert.ToSql()
	if err != nil {
		log.Printf("failed to build chat_users insert query: %v", err)
		return err
	}

	q := db.Query{
		Name:     "chat_users_repository.Create",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("failed to execute chat_users insert query: %v", err)
		return err
	}

	return nil
}

// DeleteChat удаление чата в базе данных
func (r *repo) DeleteChat(ctx context.Context, id int64) error {
	if err := r.deleteChatByID(ctx, id); err != nil {
		log.Printf("failed to delete chat: %v", err)
		return err
	}

	if err := r.deleteChatUsersByChatID(ctx, id); err != nil {
		log.Printf("failed to delete chat users: %v", err)
		return err
	}

	return nil
}

// deleteChatByID удаляет чат из таблицы чатов по его ID.
func (r *repo) deleteChatByID(ctx context.Context, id int64) error {
	deleteChatBuilder := sq.Delete(tableChatName).
		Where(sq.Eq{tableChatIDColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := deleteChatBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build delete chat query: %v", err)
		return err
	}

	q := db.Query{
		Name:     "chat_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("failed to execute delete chat query: %v", err)
		return err
	}

	return nil
}

// deleteChatUsersByChatID удаляет пользователей из таблицы chat_users, связанных с указанным ID чата.
func (r *repo) deleteChatUsersByChatID(ctx context.Context, chatID int64) error {
	deleteChatUsersBuilder := sq.Delete(tableChatUsersName).
		Where(sq.Eq{tableChatUsersChatIDColumn: chatID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := deleteChatUsersBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build delete chat users query: %v", err)
		return err
	}

	q := db.Query{
		Name:     "chat_users_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("failed to execute delete chat users query: %v", err)
		return err
	}

	return nil
}

// SendMessage запись в базу данных отправленных сообщений
func (r *repo) SendMessage(ctx context.Context, chat *model.ChatSendMessage) error {
	var messageID int64
	insertMessageBuilder := sq.Insert(tableMessagesName).
		Columns(tableMessagesUserIDColumn, tableMessagesMessageColumn, tableMessagesCreatedAtColumn).
		Values(chat.From, chat.Text, chat.Timestamp).
		PlaceholderFormat(sq.Dollar).
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
