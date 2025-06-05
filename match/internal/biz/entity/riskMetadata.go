package entity

type RiskMetadata struct {
	MaxSlippagePercent float64 // 最大滑點百分比
	LiquidationPrice   float64 // 強平價格
	TriggerStopLoss    float64 // 觸發止損價格
	TriggerTakeProfit  float64 // 觸發止盈價格
}
