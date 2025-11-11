package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"orderservice/internal/domain"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	"github.com/redis/go-redis/v9"
)

//go:embed schema.sql
var Schema string

const (
	orderCachePrefix = "order:"
	cacheTTL         = 5 * time.Minute
	maxCacheRetries  = 2
	cacheRetryDelay  = 100 * time.Millisecond
)

type OrderRepository struct {
	db          *sqlx.DB
	cache       *redis.Client
	cacheEnable bool
}

type Config struct {
	CacheEnable bool
}

func NewOrderRepository(db *sqlx.DB, redisClient *redis.Client, config *Config) *OrderRepository {
	if config == nil {
		config = &Config{
			CacheEnable: true,
		}
	}

	return &OrderRepository{
		db:          db,
		cache:       redisClient,
		cacheEnable: config.CacheEnable,
	}
}

func (r *OrderRepository) cacheKey(id string) string {
	return orderCachePrefix + id
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		insert into orders (id, item, quantity)
		values (:id, :item, :quantity)
	`

	if _, err := tx.NamedExecContext(ctx, query, order); err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	if r.cacheEnable {
		if err := r.setCacheWithRetry(ctx, order); err != nil {
			log.Printf("WARN: cache set error for order %s: %v", order.ID, err)
		}
	}

	return nil
}

func (r *OrderRepository) Get(ctx context.Context, id string) (*domain.Order, error) {
	if r.cacheEnable {
		if order, err := r.getFromCache(ctx, id); err == nil {
			return order, nil
		} else if !errors.Is(err, redis.Nil) {
			log.Printf("WARN: cache get error for order %s: %v", id, err)
		}
	}

	const query = `
		select id, item, quantity
		from orders
		where id = $1
	`

	var order domain.Order
	if err := r.db.GetContext(ctx, &order, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, fmt.Errorf("get order by id: %w", err)
	}

	if r.cacheEnable {
		if err := r.setCacheWithRetry(ctx, &order); err != nil {
			log.Printf("WARN: cache set error for order %s: %v", id, err)
		}
	}

	return &order, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		update orders 
		set item = :item, quantity = :quantity
		where id = :id
	`

	result, err := tx.NamedExecContext(ctx, query, order)
	if err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrOrderNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	if r.cacheEnable {
		if err := r.setCacheWithRetry(ctx, order); err != nil {
			log.Printf("WARN: cache set error for order %s: %v", order.ID, err)
		}
	}

	return nil
}

func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	const query = `
		delete from orders
		where id = $1
	`

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrOrderNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	r.invalidateCache(ctx, id)
	return nil
}

func (r *OrderRepository) List(ctx context.Context) ([]*domain.Order, error) {
	const query = `
		select id, item, quantity
		from orders
		order by id
	`

	var orders []*domain.Order
	if err := r.db.SelectContext(ctx, &orders, query); err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}

	return orders, nil
}

func (r *OrderRepository) getFromCache(ctx context.Context, id string) (*domain.Order, error) {
	data, err := r.cache.Get(ctx, r.cacheKey(id)).Bytes()
	if err != nil {
		return nil, err
	}

	var order domain.Order
	if err := json.Unmarshal(data, &order); err != nil {
		r.cache.Del(ctx, r.cacheKey(id))
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) setCacheWithRetry(ctx context.Context, order *domain.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	key := r.cacheKey(order.ID.String())

	for i := range maxCacheRetries {
		err = r.cache.Set(ctx, key, data, cacheTTL).Err()
		if err == nil {
			return nil
		}

		if i < maxCacheRetries-1 {
			time.Sleep(cacheRetryDelay)
		}
	}

	return err
}

func (r *OrderRepository) invalidateCache(_ context.Context, id string) {
	if !r.cacheEnable {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := r.cache.Del(ctx, r.cacheKey(id)).Err(); err != nil {
			log.Printf("WARN: cache invalidation failed for order %s: %v", id, err)
		}
	}()
}
