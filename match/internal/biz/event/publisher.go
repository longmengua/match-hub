// internal/event/publisher.go
package event

import (
	"encoding/json"
	"match/internal/biz/entity"
	"match/internal/biz/interfaces"
)

type Producer struct {
	MQ interfaces.MQClient
}

func (p *Producer) PublishTrade(trade *entity.Trade) error {
	payload, _ := json.Marshal(trade)
	return p.MQ.Publish(payload, "")
}
