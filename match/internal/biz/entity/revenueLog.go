package entity

import "time"

type RevenueLog struct {
	ID        string
	Source    string    // e.g., "trading", "staking", "referral"
	Symbol    string    // e.g., "BTC/USD", "ETH/USD"
	Amount    float64   // 收益金額
	RelatedID string    // 相關ID，例如交易ID、質押ID等
	Timestamp time.Time // 收益時間
}
