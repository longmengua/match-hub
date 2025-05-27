package usecase_test

import (
	"match/internal/biz/entity"
	"match/internal/biz/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMatch_BuyOrder_NoMatch tests the scenario where a limit buy order does not match any sell orders in the order book.
// It verifies that the buy order is added to the order book without any trades being created.
// It also checks that the order book's buy orders contain the new buy order.
func TestMatch_BuyOrder_NoMatch(t *testing.T) {
	ob := &entity.OrderBook{}

	// 無賣單時，限價買單應進入 order book
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

// TestMatch_BuyOrder_MatchesSellOrder tests the scenario where a limit buy order matches an existing sell order in the order book.
// It verifies that a trade is created, the buy order is matched with the sell order, and the remaining quantity of the sell order is updated correctly.
func TestMatch_BuyOrder_MatchesSellOrder(t *testing.T) {
	ob := &entity.OrderBook{}

	// 先加一張賣單進 order book（價格 100，數量 5）
	sellOrder := &entity.Order{
		ID:        "sell-1",
		Price:     100,
		Quantity:  5,
		Type:      entity.OrderTypeLimit,
		Side:      entity.SideSell,
		Timestamp: time.Now(),
	}
	ob.AddOrder(sellOrder)

	// 買單出價 101，要買 3，應該成功撮合
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

	// 檢查剩餘賣單數量（應該剩 2）
	assert.Equal(t, 1, len(ob.SellOrders))
	assert.Equal(t, float64(2), ob.SellOrders[0].Quantity)
}
