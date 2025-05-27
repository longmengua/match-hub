package entity

import "time"

type OrderType string
type OrderSide string

const (
	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeMarket OrderType = "MARKET"

	SideBuy  OrderSide = "BUY"
	SideSell OrderSide = "SELL"
)

type Order struct {
	ID        string
	Price     float64
	Quantity  float64
	Type      OrderType
	Side      OrderSide
	Timestamp time.Time
}

type Trade struct {
	BuyOrderID  string
	SellOrderID string
	Price       float64
	Quantity    float64
	Timestamp   time.Time
}
