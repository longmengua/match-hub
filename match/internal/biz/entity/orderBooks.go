package entity

type OrderBooks struct {
	Books map[string]OrderBook
}

// NewOrderBook 創建一個新的訂單簿
func NewOrderBooks() *OrderBooks {
	return &OrderBooks{
		Books: make(map[string]OrderBook),
	}
}

var DefaultOrderBooks = NewOrderBooks()

// GetAllSymbols 獲取所有訂單簿的交易對符號
// 返回一個字符串切片，包含所有訂單簿的交易對符號
func (ob *OrderBooks) GetAllSymbols() []string {
	symbols := make([]string, 0, len(ob.Books))
	for symbol := range ob.Books {
		symbols = append(symbols, symbol)
	}
	return symbols
}

// GetOrderBook 根據交易對符號獲取訂單簿
// 如果不存在則返回 nil
func (ob *OrderBooks) GetOrderBook(symbol string) *OrderBook {
	if book, exists := ob.Books[symbol]; exists {
		return &book
	}
	return nil
}

// AddOrderBook 添加新的訂單簿
// 如果已存在則不做任何操作
func (ob *OrderBooks) AddOrderBook(symbol string, book OrderBook) {
	if _, exists := ob.Books[symbol]; !exists {
		ob.Books[symbol] = book
	}
}

// RemoveOrderBook 移除指定的訂單簿
// 如果不存在則不做任何操作
func (ob *OrderBooks) RemoveOrderBook(symbol string) {
	delete(ob.Books, symbol)
}
