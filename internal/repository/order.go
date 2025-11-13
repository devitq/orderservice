package repository

import (
	"context"

	"orderservice/internal/domain"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	Get(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*domain.Order, error)
}
