package kafkaclient_test

import (
	"context"
	"testing"
	"time"

	"match/pkg/kafkaclient"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func createTestTopic(t *testing.T, broker, topic string) {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		t.Fatalf("failed to connect to Kafka: %v", err)
	}
	defer conn.Close()

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		t.Fatalf("failed to create topic: %v", err)
	}
}

func TestMockKafkaClient(t *testing.T) {
	broker := "kafka.docker-compose-gui.orb.local:9092"
	topic := "topicName"
	groupID := "groupID"

	createTestTopic(t, broker, topic)

	client := kafkaclient.NewKafkaGo(broker, topic, groupID)

	// 測試 Send
	err := client.SendMessage(context.Background(), "k1", "v1")
	assert.NoError(t, err)

	// 測試 Read
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key, val, err := client.ReadMessage(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "k1", key)
	assert.Equal(t, "v1", val)

	// 測試 Read 超出範圍
	_, _, err = client.ReadMessage(ctx)
	assert.Error(t, err)

	_ = client.Close()
}
