package kafkaclient

import (
	"github.com/IBM/sarama"
)

type SaramaClient struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
}

func NewSaramaClient(brokers []string) (*SaramaClient, error) {
	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}

	return &SaramaClient{producer: producer, consumer: consumer}, nil
}

func (c *SaramaClient) SendMessage(topic string, key, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	_, _, err := c.producer.SendMessage(msg)
	return err
}

func (c *SaramaClient) ConsumeMessages(topic string, handler func(key, value []byte)) error {
	pc, err := c.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	go func() {
		for msg := range pc.Messages() {
			handler(msg.Key, msg.Value)
		}
	}()
	return nil
}

func (c *SaramaClient) Close() error {
	_ = c.producer.Close()
	return c.consumer.Close()
}
