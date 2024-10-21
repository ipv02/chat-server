package env

import (
	"errors"
	"net"
	"os"

	"github.com/ipv02/chat-server/config"
)

var _ config.GRPCConfig = (*grpcConfig)(nil)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type grpcConfig struct {
	host string
	port string
}

// NewGRPCConfig создает новую конфигурацию для подключения к gRPC серверу.
//
// Получает значения хоста и порта из переменных окружения и использует их
// для создания объекта grpcConfig. Если переменная окружения для хоста или
// порта не найдена или пуста, функция возвращает соответствующую ошибку.
//
// Параметры:
//   - Нет.
//
// Возвращает:
//   - *grpcConfig: Конфигурация подключения к gRPC серверу.
//   - error: Возвращает ошибку, если переменные окружения для хоста или порта не установлены, или nil при успешном создании конфигурации.
func NewGRPCConfig() (*grpcConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
