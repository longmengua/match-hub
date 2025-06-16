package kafkaclient

import (
	"context"
	"log"
	"match/internal/biz/interfaces"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

// 確保 KafkaGoClient 實作了 KafkaClient 這個 interface
var _ interfaces.KafkaClient = (*KafkaGoClient)(nil)

// KafkaGoClient 是基於 segmentio/kafka-go 套件的 Kafka 客戶端
type KafkaGoClient struct {
	writer *kafka.Writer // 用來發送訊息的 Writer
	reader *kafka.Reader // 用來接收訊息的 Reader
	broker string        // Kafka broker 地址（只取第一個）
	topic  string        // 使用的主題名稱
}

// NewKafkaGoClient 建立一個 Kafka 客戶端，綁定特定主題
func NewKafkaGoClient(brokers []string, topic string) *KafkaGoClient {
	return &KafkaGoClient{
		broker: brokers[0], // 預設選第一個 broker 作為主要連線
		topic:  topic,
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  brokers,             // Kafka broker 列表
			Topic:    topic,               // 預設發送的主題
			Balancer: &kafka.LeastBytes{}, // 使用最少傳輸量分配策略
		}),
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,    // Kafka broker 列表
			Topic:    topic,      // 要讀取的主題
			GroupID:  "my-group", // 消費者群組 ID
			MinBytes: 1,          // 每次讀取最少資料量
			MaxBytes: 10e6,       // 每次讀取最大資料量（10MB）
		}),
	}
}

// CreateTopics 建立 Kafka 主題（如果尚未存在）
func (c *KafkaGoClient) CreateTopics(topic string) error {
	// 建立與 broker 的連線
	conn, err := kafka.Dial("tcp", c.broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 找出 controller broker（只有它有權限建立主題）
	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	// 再建立一個與 controller 的連線
	controllerConn, err := kafka.Dial("tcp", controller.Host+":"+strconv.Itoa(controller.Port))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	// 建立 topic 設定（1 分區、1 副本）
	topicConfigs := []kafka.TopicConfig{{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}}

	// 嘗試建立 topic
	return controllerConn.CreateTopics(topicConfigs...)
}

// SendMessage 發送訊息到 Kafka 主題
func (c *KafkaGoClient) SendMessage(topic string, key, value []byte) error {
	msg := kafka.Message{
		Key:   key,   // 訊息鍵（用於分區）
		Value: value, // 訊息內容
	}
	// 寫入訊息
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.writer.WriteMessages(ctx, msg)
}

// ConsumeMessages 開始消費訊息，收到後透過 handler 處理
func (c *KafkaGoClient) ConsumeMessages(topic string, handler func(key, value []byte)) error {
	go func() {
		for {
			// 讀取一則訊息
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Println("讀取訊息時發生錯誤：", err)
				time.Sleep(time.Second) // 簡單的 retry 機制
				continue
			}
			// 呼叫使用者傳入的 handler 處理訊息
			handler(msg.Key, msg.Value)
		}
	}()
	return nil
}

// Close 關閉 writer 和 reader 的連線
func (c *KafkaGoClient) Close() error {
	_ = c.writer.Close()
	return c.reader.Close()
}
