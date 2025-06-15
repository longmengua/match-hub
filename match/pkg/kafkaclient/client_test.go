package kafkaclient_test

import (
	"bytes"
	"testing"
	"time"

	"match/internal/biz/interfaces"
	"match/pkg/kafkaclient"
)

// Helper：測試流程封裝
func testKafkaClient(t *testing.T, client interfaces.KafkaClient) {
	// 先確認是否有topic存在，若不存在則創建
	err := client.CreateTopics("test-topic")
	if err != nil {
		t.Fatalf("CreateTopics failed: %v", err)
	}
	// ✅ 測試發送和接收消息
	topic := "test-topic"
	msgKey := []byte("test-key")
	msgVal := []byte("test-value")

	err = client.SendMessage(topic, msgKey, msgVal)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	received := make(chan []byte, 1)

	err = client.ConsumeMessages(topic, func(key, value []byte) {
		if bytes.Equal(key, msgKey) {
			received <- value
		}
	})
	if err != nil {
		t.Fatalf("ConsumeMessages failed: %v", err)
	}

	select {
	case val := <-received:
		if !bytes.Equal(val, msgVal) {
			t.Errorf("Expected %s, got %s", msgVal, val)
		}
	case <-time.After(5 * time.Second):
		t.Error("Timed out waiting for message")
	}

	client.Close()
}

var address = []string{"localhost:19092", "localhost:19093", "localhost:19094"}

// ✅ 測試 Sarama 實作
func TestSaramaClient(t *testing.T) {
	client, err := kafkaclient.NewSaramaClient(address)
	if err != nil {
		t.Fatalf("Failed to create SaramaClient: %v", err)
	}
	testKafkaClient(t, client)
}

// ✅ 測試 kafka-go 實作
func TestKafkaGoClient(t *testing.T) {
	client := kafkaclient.NewKafkaGoClient(address, "test-topic")
	testKafkaClient(t, client)
}
