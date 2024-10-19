package env

import (
	"errors"
	"os"

	"github.com/ipv02/chat-server/config"
)

var _ config.PGConfig = (*pgConfig)(nil)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

// NewPGConfig создает новую конфигурацию для подключения к PostgreSQL.
//
// Получает значение DSN (Data Source Name) из переменной окружения и
// использует его для создания объекта pgConfig. Если переменная окружения
// не найдена или пуста, возвращает ошибку.
//
// Параметры:
//   - Нет.
//
// Возвращает:
//   - *pgConfig: Конфигурация подключения к PostgreSQL.
//   - error: Возвращает ошибку, если переменная окружения для DSN не установлена, или nil при успешном создании конфигурации.
func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
