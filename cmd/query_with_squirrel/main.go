package main

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
)

const dbDSN = "host=localhost port=54322 dbname=chat user=chat-user password=chat-password sslmode=disable"

func main() {
	ctx := context.Background()

	dbCtx, dbCancel := context.WithTimeout(ctx, 3*time.Second)
	defer dbCancel()

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(dbCtx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Делаем запрос на вставку записи в таблицу chat
	builderInsert := sq.Insert("chat").
		PlaceholderFormat(sq.Dollar).
		Columns("name").
		Values(gofakeit.Name()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var chatID int
	err = pool.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		log.Fatalf("failed to insert note: %v", err)
	}

	log.Printf("inserted note with id: %d", chatID)

	// Делаем запрос на выборку записей из таблицы chat
	builderSelect := sq.Select("id", "name").
		From("chat").
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Limit(10)

	query, args, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}

	var id int
	var name string

	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatalf("failed to scan note: %v", err)
		}

		log.Printf("id: %d, name: %s", id, name)
	}

	// Делаем запрос на обновление записи в таблице chat
	builderUpdate := sq.Update("chat").
		PlaceholderFormat(sq.Dollar).
		Set("name", gofakeit.Name()).
		Where(sq.Eq{"id": chatID})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update note: %v", err)
	}

	log.Printf("updated %d rows", res.RowsAffected())

	// Делаем запрос на получение измененной записи из таблицы chat
	builderSelectOne := sq.Select("id", "name").
		From("chat").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": chatID}).
		Limit(1)

	query, args, err = builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&id, &name)
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}

	log.Printf("id: %d, name: %s", id, name)
}
