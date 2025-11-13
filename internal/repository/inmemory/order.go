package inmemory

import (
	"context"
	"sync"

	"orderservice/internal/domain"

	"github.com/google/uuid"
)

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*domain.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[string]*domain.Order),
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.ID.String()]; ok {
		return domain.ErrOrderAlreadyExist
	}
	r.orders[order.ID.String()] = order

	return nil
}

func (r *OrderRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id.String()]
	if !ok {
		return nil, domain.ErrOrderNotFound
	}

	return order, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.ID.String()]; !ok {
		return domain.ErrOrderNotFound
	}
	r.orders[order.ID.String()] = order

	return nil
}

func (r *OrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[id.String()]; !ok {
		return domain.ErrOrderNotFound
	}

	delete(r.orders, id.String())

	return nil
}

func (r *OrderRepository) List(ctx context.Context) ([]*domain.Order, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]*domain.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}

	return orders, nil
}
