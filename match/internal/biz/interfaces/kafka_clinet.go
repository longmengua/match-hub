package interfaces

type KafkaClient interface {
	CreateTopics(topic string) error
	SendMessage(topic string, key, value []byte) error
	ConsumeMessages(topic string, handler func(key, value []byte)) error
	Close() error
}
