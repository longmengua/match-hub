package entity

type Position struct {
	ID               string  // 持倉ID，唯一標識
	UserID           string  // 用戶ID
	Symbol           string  // 交易對，例如 "BTC/USD"
	EntryPrice       float64 // 開倉價格
	Size             float64 // 持倉大小，正數表示多頭，負數表示空頭
	PositionStatus   string  // 持倉狀態，例如 "OPEN", "CLOSED", "LIQUIDATED"
	IsCross          bool    // 是否為全倉模式
	IsIsolated       bool    // 是否為逐倉模式
	PositionType     string  // 持倉類型，例如 "SPOT", "FUTURE", "OPTIONS"
	Margin           float64 // 保證金
	Leverage         float64 // 杠杆倍數
	LiquidationPrice float64 // 強平價格
	ClosedSize       float64 // 已平倉大小
	UnrealizedPnL    float64 // 未實現盈虧
	ClosedPnL        float64 // 已實現盈虧
	CreatedAt        int64   // 創建時間戳
	UpdatedAt        int64   // 更新時間戳
	PositionTag      string  // 持倉標籤，用於標識不同的交易策略或用途
}
