package engin_test

import (
	"match/internal/biz/entity"
	"match/pkg/engin"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 測試：限價買單遇到空的 order book，應該進入買單簿而不是成交
func TestMatch_BuyOrder_NoMatch(t *testing.T) {
	ob := &entity.OrderBook{}
	buyOrder := entity.NewOrder("buy-2", 100, 5, entity.TypeLimit, entity.SideBuy, 0)

	trades := engin.Match(buyOrder, ob)

	assert.Len(t, trades, 0)
	assert.Equal(t, 1, len(ob.BuyOrders))
	assert.Equal(t, "buy-2", ob.BuyOrders[0].ID)
}

// 測試：限價買單價格高於現有賣單時，應該撮合成交
func TestMatch_BuyOrder_MatchesSellOrder(t *testing.T) {
	ob := &entity.OrderBook{}
	ob.AddOrder(entity.NewOrder("sell-1", 100, 5, entity.TypeLimit, entity.SideSell, 0))
	buyOrder := entity.NewOrder("buy-1", 101, 3, entity.TypeLimit, entity.SideBuy, 1)

	trades := engin.Match(buyOrder, ob)

	assert.Len(t, trades, 1)
	assert.Equal(t, "buy-1", trades[0].BuyOrder.ID)
	assert.Equal(t, "sell-1", trades[0].SellOrder.ID)
	assert.Equal(t, float64(100), trades[0].Price)               // 成交價格應為對手賣單價格
	assert.Equal(t, float64(3), trades[0].Quantity)              // 成交數量
	assert.Equal(t, float64(2), ob.SellOrders[0].LeavesQuantity) // 賣單剩餘 2
}

// 測試：限價買單價格低於賣單，無法撮合，應掛入買單簿
func TestMatch_LimitBuy_PriceTooLow_NoMatch(t *testing.T) {
	ob := &entity.OrderBook{}
	ob.AddOrder(entity.NewOrder("sell-1", 105, 5, entity.TypeLimit, entity.SideSell, 0))
	buy := entity.NewOrder("buy-3", 100, 5, entity.TypeLimit, entity.SideBuy, 1)

	trades := engin.Match(buy, ob)

	assert.Len(t, trades, 0)
	assert.Equal(t, 1, len(ob.BuyOrders))
}

// 測試：市價買單應優先吃價格最低的賣單，直到數量滿足或賣單耗盡
func TestMatch_MarketBuy_MatchWithBestPrice(t *testing.T) {
	ob := &entity.OrderBook{}
	ob.AddOrder(entity.NewOrder("sell-1", 100, 5, entity.TypeLimit, entity.SideSell, 0))
	ob.AddOrder(entity.NewOrder("sell-2", 99, 3, entity.TypeLimit, entity.SideSell, 1))

	buy := entity.NewOrder("buy-4", 0, 6, entity.TypeMarket, entity.SideBuy, 2)
	trades := engin.Match(buy, ob)

	assert.Len(t, trades, 2)
	assert.Equal(t, float64(3), trades[0].Quantity)
	assert.Equal(t, float64(99), trades[0].Price)
	assert.Equal(t, float64(3), trades[1].Quantity)
	assert.Equal(t, float64(100), trades[1].Price)
}

// 測試：限價賣單價格低於現有買單時，應進行撮合
func TestMatch_LimitSell_Match_Success(t *testing.T) {
	ob := &entity.OrderBook{}
	ob.AddOrder(entity.NewOrder("buy-1", 101, 5, entity.TypeLimit, entity.SideBuy, 0))
	sell := entity.NewOrder("sell-1", 100, 3, entity.TypeLimit, entity.SideSell, 1)

	trades := engin.Match(sell, ob)

	assert.Len(t, trades, 1)
	assert.Equal(t, float64(3), trades[0].Quantity)
	assert.Equal(t, float64(101), trades[0].Price)
}

// 測試：市價賣單應優先吃價格最高的買單
func TestMatch_MarketSell_MatchWithBestBuy(t *testing.T) {
	ob := &entity.OrderBook{}
	ob.AddOrder(entity.NewOrder("buy-1", 101, 2, entity.TypeLimit, entity.SideBuy, 0))
	ob.AddOrder(entity.NewOrder("buy-2", 100, 3, entity.TypeLimit, entity.SideBuy, 1))

	sell := entity.NewOrder("sell-1", 0, 4, entity.TypeMarket, entity.SideSell, 2)
	trades := engin.Match(sell, ob)

	assert.Len(t, trades, 2)
	assert.Equal(t, float64(2), trades[0].Quantity)
	assert.Equal(t, float64(101), trades[0].Price)
	assert.Equal(t, float64(2), trades[1].Quantity)
	assert.Equal(t, float64(100), trades[1].Price)
}

// 測試：限價單可連續吃多筆對手單，若未吃完則剩餘掛簿
func TestMatch_MultiMatch_PartialRemaining(t *testing.T) {
	ob := &entity.OrderBook{}
	ob.AddOrder(entity.NewOrder("sell-1", 99, 2, entity.TypeLimit, entity.SideSell, 0))
	ob.AddOrder(entity.NewOrder("sell-2", 100, 2, entity.TypeLimit, entity.SideSell, 1))
	buy := entity.NewOrder("buy-1", 100, 5, entity.TypeLimit, entity.SideBuy, 2)

	trades := engin.Match(buy, ob)

	assert.Len(t, trades, 2)
	assert.Equal(t, float64(4), trades[0].Quantity+trades[1].Quantity)
	assert.Equal(t, 1, len(ob.BuyOrders))
	assert.Equal(t, float64(1), ob.BuyOrders[0].LeavesQuantity)
}

// 測試：相同價格時，應根據時間先後順序（時間優先原則）撮合
func TestMatch_TimePriority_WhenPriceEqual(t *testing.T) {
	ob := &entity.OrderBook{}
	ob.AddOrder(entity.NewOrder("sell-1", 100, 2, entity.TypeLimit, entity.SideSell, 0))
	ob.AddOrder(entity.NewOrder("sell-2", 100, 2, entity.TypeLimit, entity.SideSell, 1))
	buy := entity.NewOrder("buy-1", 100, 3, entity.TypeLimit, entity.SideBuy, 2)

	trades := engin.Match(buy, ob)

	assert.Len(t, trades, 2)
	assert.Equal(t, "sell-1", trades[0].SellOrder.ID)
	assert.Equal(t, "sell-2", trades[1].SellOrder.ID)
}

// 測試：0 數量訂單應被忽略，不可進入 order book
func TestMatch_ZeroQuantityOrder(t *testing.T) {
	ob := &entity.OrderBook{}
	order := entity.NewOrder("zero-qty", 100, 0, entity.TypeLimit, entity.SideBuy, 0)

	trades := engin.Match(order, ob)
	assert.Len(t, trades, 0)
	assert.Len(t, ob.BuyOrders, 0)
	assert.Len(t, ob.SellOrders, 0)
}

// 測試：負數價格的限價單（不合理）預期仍進入簿中（視實作處理）
func TestMatch_NegativePriceOrder(t *testing.T) {
	ob := &entity.OrderBook{}
	order := entity.NewOrder("negative-price", -50, 10, entity.TypeLimit, entity.SideSell, 0)

	trades := engin.Match(order, ob)
	assert.Len(t, trades, 0)
	assert.Equal(t, 1, len(ob.SellOrders))
	assert.Equal(t, float64(-50), ob.SellOrders[0].Price)
}

// 測試：市價單在對手簿為空時，不會撮合也不會掛單
func TestMatch_MarketOrder_NoCounterOrders(t *testing.T) {
	ob := &entity.OrderBook{}
	order := entity.NewOrder("market-alone", 0, 5, entity.TypeMarket, entity.SideBuy, 0)

	trades := engin.Match(order, ob)
	assert.Len(t, trades, 0)
	assert.Len(t, ob.BuyOrders, 0)
	assert.Len(t, ob.SellOrders, 0)
}

// 測試：限價買單完全撮合後應不殘留在 order book 中
func TestMatch_LimitBuy_FullyMatched(t *testing.T) {
	ob := &entity.OrderBook{}
	ob.AddOrder(entity.NewOrder("sell-1", 100, 3, entity.TypeLimit, entity.SideSell, 0))
	buy := entity.NewOrder("buy-1", 101, 3, entity.TypeLimit, entity.SideBuy, 1)

	trades := engin.Match(buy, ob)
	assert.Len(t, trades, 1)
	assert.Equal(t, float64(3), trades[0].Quantity)
	assert.Empty(t, ob.BuyOrders)
	assert.Empty(t, ob.SellOrders)
}

// 測試：大數量撮合壓力測試
func TestMatch_LargeVolumeOrder(t *testing.T) {
	ob := &entity.OrderBook{}
	ob.AddOrder(entity.NewOrder("sell-huge", 100, 1e6, entity.TypeLimit, entity.SideSell, 0))
	buy := entity.NewOrder("buy-huge", 100, 1e6, entity.TypeLimit, entity.SideBuy, 1)

	trades := engin.Match(buy, ob)
	assert.Len(t, trades, 1)
	assert.Equal(t, float64(1e6), trades[0].Quantity)
	assert.Empty(t, ob.BuyOrders)
	assert.Empty(t, ob.SellOrders)
}
