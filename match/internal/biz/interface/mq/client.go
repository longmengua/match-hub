// internal/infra/mq/client.go
package mq

type Client interface {
	Publish(topic string, data []byte) error
	Subscribe(topic string, handler func([]byte)) error
}
