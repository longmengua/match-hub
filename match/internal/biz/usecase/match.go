package usecase

import (
	"match/internal/biz/entity"
	"time"
)

func Match(order *entity.Order, ob *entity.OrderBook) []*entity.Trade {
	trades := []*entity.Trade{}

	// 拒絕不合法的 order，return 空撮合
	if !isValidSide(order.Side) || !isValidType(order.Type) || order.Quantity <= 0 {
		return trades
	}

	if order.Side == entity.SideBuy {
		// 與最便宜的賣單撮合
		for i := 0; i < len(ob.SellOrders) && order.Quantity > 0; {
			sell := ob.SellOrders[i]

			// 處理限價單不能接受高價
			if order.Type == entity.OrderTypeLimit && order.Price < sell.Price {
				break
			}

			matchQty := min(order.Quantity, sell.Quantity)

			trade := &entity.Trade{
				BuyOrderID:  order.ID,
				SellOrderID: sell.ID,
				Price:       sell.Price,
				Quantity:    matchQty,
				Timestamp:   time.Now(),
			}
			trades = append(trades, trade)

			order.Quantity -= matchQty
			sell.Quantity -= matchQty

			if sell.Quantity == 0 {
				ob.RemoveOrder(&ob.SellOrders, i)
			} else {
				i++
			}
		}
	}

	if order.Side == entity.SideSell {
		for i := 0; i < len(ob.BuyOrders) && order.Quantity > 0; {
			buy := ob.BuyOrders[i]

			if order.Type == entity.OrderTypeLimit && order.Price > buy.Price {
				break
			}

			matchQty := min(order.Quantity, buy.Quantity)

			trade := &entity.Trade{
				BuyOrderID:  buy.ID,
				SellOrderID: order.ID,
				Price:       buy.Price,
				Quantity:    matchQty,
				Timestamp:   time.Now(),
			}
			trades = append(trades, trade)

			order.Quantity -= matchQty
			buy.Quantity -= matchQty

			if buy.Quantity == 0 {
				ob.RemoveOrder(&ob.BuyOrders, i)
			} else {
				i++
			}
		}
	}

	// 掛回簿子（僅限價單）
	if order.Quantity > 0 && order.Type == entity.OrderTypeLimit {
		ob.AddOrder(order)
	}

	return trades
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func isValidSide(side entity.OrderSide) bool {
	return side == entity.SideBuy || side == entity.SideSell
}

func isValidType(typ entity.OrderType) bool {
	return typ == entity.OrderTypeLimit || typ == entity.OrderTypeMarket
}
