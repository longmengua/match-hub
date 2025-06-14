package interfaces

import "context"

type KafkaClient interface {
	// Producer 功能
	SendMessage(ctx context.Context, key, value string) error
	// Consumer 功能
	ReadMessage(ctx context.Context) (key string, value string, err error)
	// 通用資源釋放
	Close() error
}
