package usecase_test

import (
	"match/internal/biz/entity"
	"match/internal/biz/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 工具方法：快速建立一張 Order
func newOrder(id string, price, qty float64, orderType entity.OrderType, side entity.OrderSide, tsOffset int) *entity.Order {
	return &entity.Order{
		ID:        id,
		Price:     price,
		Quantity:  qty,
		Type:      orderType,
		Side:      side,
		Timestamp: time.Now().Add(time.Duration(tsOffset) * time.Second),
	}
}

// 測試：限價買單遇到空的 order book，應該進入買單簿而不是成交
func TestMatch_BuyOrder_NoMatch(t *testing.T) {
	ob := &entity.OrderBook{}

	buyOrder := &entity.Order{
		ID:        "buy-2",
		Price:     100,
		Quantity:  5,
		Type:      entity.OrderTypeLimit,
		Side:      entity.SideBuy,
		Timestamp: time.Now(),
	}

	trades := usecase.Match(buyOrder, ob)

	assert.Len(t, trades, 0)
	assert.Equal(t, 1, len(ob.BuyOrders))
	assert.Equal(t, "buy-2", ob.BuyOrders[0].ID)
}

// 測試：限價買單價格高於現有賣單時，應該進行撮合成交
func TestMatch_BuyOrder_MatchesSellOrder(t *testing.T) {
	ob := &entity.OrderBook{}

	sellOrder := &entity.Order{
		ID:        "sell-1",
		Price:     100,
		Quantity:  5,
		Type:      entity.OrderTypeLimit,
		Side:      entity.SideSell,
		Timestamp: time.Now(),
	}
	ob.AddOrder(sellOrder)

	buyOrder := &entity.Order{
		ID:        "buy-1",
		Price:     101,
		Quantity:  3,
		Type:      entity.OrderTypeLimit,
		Side:      entity.SideBuy,
		Timestamp: time.Now(),
	}
	trades := usecase.Match(buyOrder, ob)

	assert.Len(t, trades, 1)
	assert.Equal(t, "buy-1", trades[0].BuyOrderID)
	assert.Equal(t, "sell-1", trades[0].SellOrderID)
	assert.Equal(t, float64(100), trades[0].Price)
	assert.Equal(t, float64(3), trades[0].Quantity)
	assert.Equal(t, float64(2), ob.SellOrders[0].Quantity)
}

// 測試：限價買單價格低於賣單，應該無法撮合並掛入 order book
func TestMatch_LimitBuy_PriceTooLow_NoMatch(t *testing.T) {
	ob := &entity.OrderBook{}
	ob.AddOrder(newOrder("sell-1", 105, 5, entity.OrderTypeLimit, entity.SideSell, 0))

	buy := newOrder("buy-3", 100, 5, entity.OrderTypeLimit, entity.SideBuy, 1)
	trades := usecase.Match(buy, ob)

	assert.Len(t, trades, 0)
	assert.Equal(t, 1, len(ob.BuyOrders))
}

// 測試：市價買單應優先吃價格最低的賣單，直到數量滿足或賣單耗盡
func TestMatch_MarketBuy_MatchWithBestPrice(t *testing.T) {
	ob := &entity.OrderBook{}
	ob.AddOrder(newOrder("sell-1", 100, 5, entity.OrderTypeLimit, entity.SideSell, 0))
	ob.AddOrder(newOrder("sell-2", 99, 3, entity.OrderTypeLimit, entity.SideSell, 1))

	buy := newOrder("buy-4", 0, 6, entity.OrderTypeMarket, entity.SideBuy, 2)
	trades := usecase.Match(buy, ob)

	assert.Len(t, trades, 2)
	assert.Equal(t, float64(3), trades[0].Quantity)
	assert.Equal(t, float64(99), trades[0].Price)
	assert.Equal(t, float64(3), trades[1].Quantity)
	assert.Equal(t, float64(100), trades[1].Price)
}

// 測試：限價賣單價格低於現有買單時，應該進行成交
func TestMatch_LimitSell_Match_Success(t *testing.T) {
	ob := &entity.OrderBook{}
	ob.AddOrder(newOrder("buy-1", 101, 5, entity.OrderTypeLimit, entity.SideBuy, 0))

	sell := newOrder("sell-1", 100, 3, entity.OrderTypeLimit, entity.SideSell, 1)
	trades := usecase.Match(sell, ob)

	assert.Len(t, trades, 1)
	assert.Equal(t, float64(3), trades[0].Quantity)
	assert.Equal(t, float64(101), trades[0].Price)
}

// 測試：市價賣單應優先吃價格最高的買單
func TestMatch_MarketSell_MatchWithBestBuy(t *testing.T) {
	ob := &entity.OrderBook{}
	ob.AddOrder(newOrder("buy-1", 101, 2, entity.OrderTypeLimit, entity.SideBuy, 0))
	ob.AddOrder(newOrder("buy-2", 100, 3, entity.OrderTypeLimit, entity.SideBuy, 1))

	sell := newOrder("sell-1", 0, 4, entity.OrderTypeMarket, entity.SideSell, 2)
	trades := usecase.Match(sell, ob)

	assert.Len(t, trades, 2)
	assert.Equal(t, float64(2), trades[0].Quantity)
	assert.Equal(t, float64(101), trades[0].Price)
	assert.Equal(t, float64(2), trades[1].Quantity)
	assert.Equal(t, float64(100), trades[1].Price)
}

// 測試：一張單可以吃掉多張對手單（連續撮合），若不夠則掛入簿中
func TestMatch_MultiMatch_PartialRemaining(t *testing.T) {
	ob := &entity.OrderBook{}
	ob.AddOrder(newOrder("sell-1", 99, 2, entity.OrderTypeLimit, entity.SideSell, 0))
	ob.AddOrder(newOrder("sell-2", 100, 2, entity.OrderTypeLimit, entity.SideSell, 1))

	buy := newOrder("buy-1", 100, 5, entity.OrderTypeLimit, entity.SideBuy, 2)
	trades := usecase.Match(buy, ob)

	assert.Len(t, trades, 2)
	assert.Equal(t, float64(4), trades[0].Quantity+trades[1].Quantity)
	assert.Equal(t, 1, len(ob.BuyOrders))
	assert.Equal(t, float64(1), ob.BuyOrders[0].Quantity)
}

// 測試：價格相同時，應以時間先後順序進行撮合（時間優先原則）
func TestMatch_TimePriority_WhenPriceEqual(t *testing.T) {
	ob := &entity.OrderBook{}
	ob.AddOrder(newOrder("sell-1", 100, 2, entity.OrderTypeLimit, entity.SideSell, 0)) // 較早
	ob.AddOrder(newOrder("sell-2", 100, 2, entity.OrderTypeLimit, entity.SideSell, 1))

	buy := newOrder("buy-1", 100, 3, entity.OrderTypeLimit, entity.SideBuy, 2)
	trades := usecase.Match(buy, ob)

	assert.Len(t, trades, 2)
	assert.Equal(t, "sell-1", trades[0].SellOrderID)
	assert.Equal(t, "sell-2", trades[1].SellOrderID)
}
