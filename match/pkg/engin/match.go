package engin

import (
	"match/internal/biz/entity"
	"time"
)

// Match 負責進行撮合邏輯：根據訂單方向、類型，與對手方訂單進行成交比對。
func Match(order *entity.Order, ob *entity.OrderBook) []*entity.Trade {
	trades := []*entity.Trade{}

	// 若訂單方向、型別不合法或數量為 0，則直接返回空成交結果
	if !isValidSide(order.Side) || !isValidType(order.Type) || order.LeavesQuantity <= 0 {
		return trades
	}

	var (
		isBuy      = order.Side == entity.SideBuy           // 判斷是否為買單
		oppOrders  *[]*entity.Order                         // 對手方訂單列表（買單對賣單，賣單對買單）
		priceMatch func(orderPrice, bookPrice float64) bool // 價格匹配條件
	)

	if isBuy {
		// 買單撮合：對手方為賣單
		oppOrders = &ob.SellOrders
		priceMatch = func(orderPrice, bookPrice float64) bool {
			// 市價單總是可以成交；限價單需價格 >= 對手方
			return order.Type == entity.TypeMarket || orderPrice >= bookPrice
		}
	} else {
		// 賣單撮合：對手方為買單
		oppOrders = &ob.BuyOrders
		priceMatch = func(orderPrice, bookPrice float64) bool {
			// 市價單總是可以成交；限價單需價格 <= 對手方
			return order.Type == entity.TypeMarket || orderPrice <= bookPrice
		}
	}

	// 與對手方訂單逐筆撮合，直到訂單數量為 0 或無法繼續撮合
	for i := 0; i < len(*oppOrders) && order.LeavesQuantity > 0; {
		opp := (*oppOrders)[i]

		// 價格無法成交時跳出（限價單才有這個限制）
		if !priceMatch(order.Price, opp.Price) {
			break
		}

		// 計算撮合數量：取兩方剩餘數量的最小值
		matchQty := min(order.LeavesQuantity, opp.LeavesQuantity)

		// 判斷哪一方是買單／賣單（便於建立 Trade 結構）
		var buyOrder, sellOrder *entity.Order
		if isBuy {
			buyOrder = order
			sellOrder = opp
		} else {
			buyOrder = opp
			sellOrder = order
		}

		// 建立成交紀錄
		trade := &entity.Trade{
			BuyOrder:  buyOrder,
			SellOrder: sellOrder,
			OrderSide: opp.Side,  // 使用被動方作為成交方向
			Price:     opp.Price, // 採用被動方價格
			Quantity:  matchQty,
			Timestamp: time.Now(),
		}
		trades = append(trades, trade)

		// 扣除撮合數量
		order.LeavesQuantity -= matchQty
		opp.LeavesQuantity -= matchQty

		// 若對手方訂單撮合完畢，從簿子中移除，否則前進下一筆
		if opp.LeavesQuantity == 0 {
			ob.RemoveOrder(oppOrders, i)
		} else {
			i++
		}
	}

	// 若尚有剩餘數量且為限價單，則掛單到訂單簿上
	if order.LeavesQuantity > 0 && order.Type == entity.TypeLimit {
		ob.AddOrder(order)
	}

	return trades
}

// 取兩數中較小值
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// 檢查訂單方向是否合法（只接受 BUY/SELL）
func isValidSide(side entity.OrderSide) bool {
	return side == entity.SideBuy || side == entity.SideSell
}

// 檢查訂單型別是否合法（只接受 LIMIT/MARKET）
func isValidType(typ entity.OrderType) bool {
	return typ == entity.TypeLimit || typ == entity.TypeMarket
}
