package interfaces

type KafkaClient interface {
	SendMessage(topic string, key, value []byte) error
	ConsumeMessages(topic string, handler func(key, value []byte)) error
	Close() error
}
