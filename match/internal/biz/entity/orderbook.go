package entity

import (
	"sort"
)

type OrderBook struct {
	BuyOrders  []*Order
	SellOrders []*Order
}

// NewOrderBook 創建一個新的訂單簿
func NewOrderBook() *OrderBook {
	return &OrderBook{
		BuyOrders:  []*Order{},
		SellOrders: []*Order{},
	}
}

// 價格由高到低，時間早的在前
func (ob *OrderBook) sortBuy() {
	sort.Slice(ob.BuyOrders, func(i, j int) bool {
		if ob.BuyOrders[i].Price == ob.BuyOrders[j].Price {
			return ob.BuyOrders[i].Timestamp.Before(ob.BuyOrders[j].Timestamp)
		}
		return ob.BuyOrders[i].Price > ob.BuyOrders[j].Price
	})
}

// 價格由低到高，時間早的在前
func (ob *OrderBook) sortSell() {
	sort.Slice(ob.SellOrders, func(i, j int) bool {
		if ob.SellOrders[i].Price == ob.SellOrders[j].Price {
			return ob.SellOrders[i].Timestamp.Before(ob.SellOrders[j].Timestamp)
		}
		return ob.SellOrders[i].Price < ob.SellOrders[j].Price
	})
}

func (ob *OrderBook) AddOrder(order *Order) {
	if order.Side == SideBuy {
		ob.BuyOrders = append(ob.BuyOrders, order)
		ob.sortBuy()
	} else {
		ob.SellOrders = append(ob.SellOrders, order)
		ob.sortSell()
	}
}

func (ob *OrderBook) RemoveOrder(orderList *[]*Order, index int) {
	*orderList = append((*orderList)[:index], (*orderList)[index+1:]...)
}
