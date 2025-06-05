// internal/event/publisher.go
package event

import (
	"encoding/json"
	"match/internal/biz/entity"
	"match/internal/biz/interface/mq"
)

type Producer struct {
	MQ mq.Client
}

func (p *Producer) PublishTrade(trade *entity.Trade) error {
	payload, _ := json.Marshal(trade)
	return p.MQ.Publish("trade_events", payload)
}
