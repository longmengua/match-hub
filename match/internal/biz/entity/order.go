package entity

import "time"

// / 訂單類型（OrderType）定義
type OrderType string

const (
	// 限價單：指定價格成交
	TypeLimit OrderType = "LIMIT"

	// 市價單：不指定價格，依市場最優價格立即成交
	TypeMarket OrderType = "MARKET"
)

// / 訂單方向（OrderSide）定義
type OrderSide string

const (
	// 買單：買入資產
	SideBuy OrderSide = "BUY"

	// 賣單：賣出資產
	SideSell OrderSide = "SELL"
)

// / 訂單結構（Order）
// / 描述一筆尚未完全成交的委託單
type Order struct {
	ID             string    // 訂單唯一識別碼（UUID 或自定字串）
	Price          float64   // 委託價格（市價單可為 0）
	Quantity       float64   // 委託數量
	LeavesQuantity float64   // 委託剩餘數量（撮合過程中會遞減）
	Type           OrderType // 訂單類型（限價 / 市價）
	Side           OrderSide // 訂單方向（買 / 賣）
	Timestamp      time.Time // 下單時間（用於排序撮合）
}

// / 成交紀錄結構（Trade）
// / 描述一筆實際撮合成交的紀錄
type Trade struct {
	BuyOrder  *Order    // 買方訂單（主動單或被動單皆可）
	SellOrder *Order    // 賣方訂單（主動單或被動單皆可）
	OrderSide OrderSide // 被動方方向（用來區分是哪邊的掛單被吃）
	Price     float64   // 實際成交價格（通常以被動方價格為準）
	Quantity  float64   // 成交數量
	Timestamp time.Time // 成交時間
}
