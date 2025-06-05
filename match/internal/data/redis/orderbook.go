package redis

import (
	"context"
	"encoding/json"
	"match/internal/biz/entity"

	"github.com/redis/go-redis/v9"
)

type OrderBookRepo struct {
	client *redis.Client
}

func NewOrderBookRepo(rdb *redis.Client) *OrderBookRepo {
	return &OrderBookRepo{client: rdb}
}

func zsetKey(side entity.OrderSide) string {
	if side == entity.SideBuy {
		return "orderbook:buy"
	}
	return "orderbook:sell"
}

func hashKey() string {
	return "orderbook:orders"
}

func (r *OrderBookRepo) AddOrder(order *entity.Order) error {
	ctx := context.TODO()
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	score := order.Price
	if order.Side == entity.SideBuy {
		score = -order.Price
	}

	pipe := r.client.TxPipeline()
	pipe.HSet(ctx, hashKey(), order.ID, data)
	pipe.ZAdd(ctx, zsetKey(order.Side), redis.Z{
		Score:  score,
		Member: order.ID,
	})
	_, err = pipe.Exec(ctx)
	return err
}

func (r *OrderBookRepo) RemoveOrder(orderID string, side entity.OrderSide) error {
	ctx := context.TODO()
	pipe := r.client.TxPipeline()
	pipe.ZRem(ctx, zsetKey(side), orderID)
	pipe.HDel(ctx, hashKey(), orderID)
	_, err := pipe.Exec(ctx)
	return err
}

// 讀取累積數量 >= targetQty 為止的賣單（價格低到高）
func (r *OrderBookRepo) GetSellOrders(targetQty float64) ([]*entity.Order, error) {
	ctx := context.TODO()
	var result []*entity.Order
	var cursor int64 = 0
	var accumulatedQty float64 = 0

	for {
		orderIDs, err := r.client.ZRange(ctx, zsetKey(entity.SideSell), cursor, cursor+9).Result()
		if err != nil || len(orderIDs) == 0 {
			break
		}
		cursor += int64(len(orderIDs))

		orders, err := r.getOrdersByIDs(orderIDs)
		if err != nil {
			return nil, err
		}

		for _, order := range orders {
			result = append(result, order)
			accumulatedQty += order.Quantity
			if accumulatedQty >= targetQty {
				return result, nil
			}
		}
	}

	return result, nil
}

// 讀取累積數量 >= targetQty 為止的買單（價格高到低）
func (r *OrderBookRepo) GetBuyOrders(targetQty float64) ([]*entity.Order, error) {
	ctx := context.TODO()
	var result []*entity.Order
	var cursor int64 = 0
	var accumulatedQty float64 = 0

	for {
		orderIDs, err := r.client.ZRange(ctx, zsetKey(entity.SideBuy), cursor, cursor+9).Result()
		if err != nil || len(orderIDs) == 0 {
			break
		}
		cursor += int64(len(orderIDs))

		orders, err := r.getOrdersByIDs(orderIDs)
		if err != nil {
			return nil, err
		}

		for _, order := range orders {
			result = append(result, order)
			accumulatedQty += order.Quantity
			if accumulatedQty >= targetQty {
				return result, nil
			}
		}
	}

	return result, nil
}

func (r *OrderBookRepo) getOrdersByIDs(orderIDs []string) ([]*entity.Order, error) {
	if len(orderIDs) == 0 {
		return nil, nil
	}
	ctx := context.TODO()
	datas, err := r.client.HMGet(ctx, hashKey(), orderIDs...).Result()
	if err != nil {
		return nil, err
	}

	var result []*entity.Order
	for i, data := range datas {
		str, ok := data.(string)
		if !ok {
			continue
		}
		var order entity.Order
		if err := json.Unmarshal([]byte(str), &order); err == nil {
			order.ID = orderIDs[i]
			result = append(result, &order)
		}
	}
	return result, nil
}

func (r *OrderBookRepo) GetBestMatchOrder(side entity.OrderSide) (*entity.Order, error) {
	ctx := context.TODO()
	matchSide := entity.SideBuy
	if side == entity.SideBuy {
		matchSide = entity.SideSell
	}

	orderIDs, err := r.client.ZRange(ctx, zsetKey(matchSide), 0, 0).Result()
	if err != nil || len(orderIDs) == 0 {
		return nil, nil
	}

	data, err := r.client.HGet(ctx, hashKey(), orderIDs[0]).Result()
	if err != nil {
		return nil, err
	}

	var order entity.Order
	if err := json.Unmarshal([]byte(data), &order); err != nil {
		return nil, err
	}
	order.ID = orderIDs[0]
	return &order, nil
}
