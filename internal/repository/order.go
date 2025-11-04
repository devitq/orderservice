package repository

import (
	"context"

	"orderservice/internal/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	Get(ctx context.Context, id string) (*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Order, error)
}
