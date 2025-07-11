package engin

import "match/internal/biz/entity"

// 單一市場撮合引擎
type Engine struct {
	OrderBook *entity.OrderBook
}

func NewEngine() *Engine {
	return &Engine{
		OrderBook: &entity.OrderBook{
			BuyOrders:  []*entity.Order{},
			SellOrders: []*entity.Order{},
		},
	}
}

func (e *Engine) SubmitOrder(order *entity.Order) []*entity.Trade {
	return Match(order, e.OrderBook)
}
