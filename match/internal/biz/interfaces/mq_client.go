package interfaces

type MQClient interface {
	Publish(data []byte, topic string) error
	Subscribe(handler func([]byte), topic string, groupId string) error
}
