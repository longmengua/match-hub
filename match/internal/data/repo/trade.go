package repo

import (
	"database/sql"
	"match/internal/biz/entity"
)

type TradeRepo struct {
	DB *sql.DB
}

func (r *TradeRepo) SaveTrade(t *entity.Trade) error {
	// _, err := r.DB.Exec(`
	// 	INSERT INTO trades (buy_order_id, sell_order_id, price, quantity, timestamp)
	// 	VALUES ($1, $2, $3, $4, $5)
	// 	ON CONFLICT DO NOTHING
	// `, t.ID, t.SellOrderID, t.Price, t.Quantity, t.Timestamp)
	// return err
	return nil // 模擬實現，實際應該執行 SQL 語句
}
