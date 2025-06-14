package kafkaclient

import (
	"context"

	"match/internal/biz/interfaces"

	"github.com/segmentio/kafka-go"
)

var _ interfaces.KafkaClient = (*KafkaClient)(nil)

type KafkaClient struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func New(broker, topic, groupID string) *KafkaClient {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupID,
	})

	return &KafkaClient{
		writer: writer,
		reader: reader,
	}
}

func (c *KafkaClient) SendMessage(ctx context.Context, key, value string) error {
	return c.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	})
}

func (c *KafkaClient) ReadMessage(ctx context.Context) (string, string, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return "", "", err
	}
	return string(msg.Key), string(msg.Value), nil
}

func (c *KafkaClient) Close() error {
	err1 := c.writer.Close()
	err2 := c.reader.Close()
	if err1 != nil {
		return err1
	}
	return err2
}
