package kafkaclient

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaGoClient struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaGoClient(brokers []string, topic string) *KafkaGoClient {
	return &KafkaGoClient{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic,
		}),
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: "my-group",
		}),
	}
}

func (c *KafkaGoClient) SendMessage(topic string, key, value []byte) error {
	return c.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   key,
		Value: value,
	})
}

func (c *KafkaGoClient) ConsumeMessages(topic string, handler func(key, value []byte)) error {
	go func() {
		for {
			msg, err := c.reader.ReadMessage(context.Background())
			if err != nil {
				log.Println("error reading message:", err)
				continue
			}
			handler(msg.Key, msg.Value)
		}
	}()
	return nil
}

func (c *KafkaGoClient) Close() error {
	_ = c.writer.Close()
	return c.reader.Close()
}
