package kafkaclient

import (
	"match/internal/biz/interfaces"
	"strings"

	"github.com/IBM/sarama"
)

var _ interfaces.KafkaClient = (*SaramaClient)(nil)

// SaramaClient 是 Kafka 的實作，使用 sarama 套件
type SaramaClient struct {
	producer sarama.SyncProducer // 同步 Producer
	consumer sarama.Consumer     // Consumer
	brokers  []string            // Kafka broker 清單
}

// 建立新的 Sarama Kafka client
func NewSaramaClient(brokers []string) (*SaramaClient, error) {
	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}

	return &SaramaClient{
		producer: producer,
		consumer: consumer,
		brokers:  brokers,
	}, nil
}

// CreateTopics 建立 Kafka topic（如果尚未存在）
func (c *SaramaClient) CreateTopics(topic string) error {
	admin, err := sarama.NewClusterAdmin(c.brokers, nil)
	if err != nil {
		return err
	}
	defer admin.Close()

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1, // 預設 1 分區
		ReplicationFactor: 1, // 預設 1 副本
	}

	// 如果 topic 已存在會回傳錯誤，不影響流程
	err = admin.CreateTopic(topic, topicDetail, false)
	if err != nil && !isTopicAlreadyExistsError(err) {
		return err
	}

	return nil
}

// 判斷錯誤是否為 "topic 已存在"
func isTopicAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Topic with this name already exists")
}

// SendMessage 發送訊息到指定 topic
func (c *SaramaClient) SendMessage(topic string, key, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	_, _, err := c.producer.SendMessage(msg)
	return err
}

// ConsumeMessages 開始消費訊息，使用 handler 處理每條訊息
func (c *SaramaClient) ConsumeMessages(topic string, handler func(key, value []byte)) error {
	partitionConsumer, err := c.consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}

	go func() {
		for msg := range partitionConsumer.Messages() {
			handler(msg.Key, msg.Value)
		}
	}()

	return nil
}

// Close 關閉 producer 和 consumer
func (c *SaramaClient) Close() error {
	_ = c.producer.Close()
	return c.consumer.Close()
}
