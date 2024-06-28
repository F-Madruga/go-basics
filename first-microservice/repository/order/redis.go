package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/F-Madruga/go-basics/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func orderIdKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := orderIdKey(order.OrderId)

	trx := r.Client.TxPipeline()

	res := trx.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		trx.Discard()
		return fmt.Errorf("failed to set: %w", err)
	}

	err = trx.SAdd(ctx, "orders", key).Err()
	if err != nil {
		trx.Discard()
		return fmt.Errorf("failed to add orders set: %w", err)
	}

	_, err = trx.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

var ErrNotExist = errors.New("order does not exist")

func (r *RedisRepo) FindById(ctx context.Context, id uint64) (model.Order, error) {
	key := orderIdKey(id)

	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrNotExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("get order: %w", err)
	}

	var order model.Order
	err = json.Unmarshal([]byte(value), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to decode order json: %w", err)
	}

	return order, nil
}

func (r *RedisRepo) DeleteById(ctx context.Context, id uint64) error {
	key := orderIdKey(id)

	trx := r.Client.TxPipeline()

	err := trx.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		trx.Discard()
		return ErrNotExist
	} else if err != nil {
		trx.Discard()
		return fmt.Errorf("del order: %w", err)
	}

	err = trx.SRem(ctx, "orders", key).Err()
	if err != nil {
		trx.Discard()
		return fmt.Errorf("failed to remove from orders set: %w", err)
	}

	_, err = trx.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

func (r *RedisRepo) Update(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := orderIdKey(order.OrderId)

	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	return nil
}

type FindAllPage struct {
	Limit  uint64
	Offset uint64
}

type FindResult struct {
	Orders []model.Order
	Cursor uint64
}

func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	res := r.Client.SScan(ctx, "orders", page.Offset, "*", int64(page.Limit))

	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get order ids: %w", err)
	}

	if len(keys) == 0 {
		return FindResult{
			Orders: []model.Order{},
		}, nil
	}

	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get orders: %w", err)
	}

	orders := make([]model.Order, len(xs))
	for i, x := range xs {
		value := x.(string)

		var order model.Order
		err = json.Unmarshal([]byte(value), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to decode order json: %w", err)
		}
		orders[i] = order
	}
	return FindResult{
		Orders: orders,
		Cursor: cursor,
	}, nil
}
