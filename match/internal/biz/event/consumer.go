// internal/event/consumer.go
package event

import (
	"encoding/json"
	"log"
	"match/internal/biz/entity"
	"match/internal/data/repo"
)

type Consumer struct {
	SQLRepo *repo.TradeRepo
}

func (c *Consumer) HandleTradeEvent(data []byte) {
	var trade entity.Trade
	if err := json.Unmarshal(data, &trade); err != nil {
		log.Println("Invalid trade event:", err)
		return
	}

	if err := c.SQLRepo.SaveTrade(&trade); err != nil {
		log.Println("Failed to write trade to SQL:", err)
	}
}
