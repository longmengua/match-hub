package mq

import (
	"context"
	"log"
	"match/internal/biz/interfaces"
	"time"

	"github.com/segmentio/kafka-go"
)

const defaultKafkaGroupID = "match-engine-default-group"
const defaultKafkaTopic = "match-engine-default-topic"

type KafkaMQClient struct {
	Brokers []string
}

// 靜態保證 KafkaMQClient 實作了 interfaces.MQClient
var _ interfaces.MQClient = (*KafkaMQClient)(nil)

func NewKafkaMQClient(brokers []string) *KafkaMQClient {
	return &KafkaMQClient{
		Brokers: brokers,
	}
}

func (k *KafkaMQClient) Publish(data []byte, topic string) error {
	topicToUse := topic
	if topicToUse == "" {
		topicToUse = defaultKafkaTopic
	}
	writer := &kafka.Writer{
		Addr:         kafka.TCP(k.Brokers...),
		Topic:        topicToUse,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}
	defer writer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   nil,
		Value: data,
	})
	if err != nil {
		log.Printf("failed to publish message to topic %s: %v", topic, err)
		return err
	}

	return nil
}

func (k *KafkaMQClient) Subscribe(handler func([]byte), topic string, groupId string) error {
	group := groupId
	if group == "" {
		group = defaultKafkaTopic
	}
	topicToUse := topic
	if topicToUse == "" {
		topicToUse = defaultKafkaTopic
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: k.Brokers,
		Topic:   topicToUse,
		GroupID: group,
	})
	go func() {
		for {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("failed to read message from topic %s: %v", topic, err)
				continue
			}
			handler(msg.Value)
		}
	}()
	return nil
}
