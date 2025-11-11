package service

import (
	"context"

	"orderservice/internal/domain"
	"orderservice/internal/repository"

	"github.com/google/uuid"
)

type OrderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) Create(ctx context.Context, item string, quantity int32) (*domain.Order, error) {
	order := &domain.Order{
		ID:       uuid.New(),
		Item:     item,
		Quantity: quantity,
	}

	err := order.Validate()
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) Get(ctx context.Context, id string) (*domain.Order, error) {
	return s.repo.Get(ctx, id)
}

func (s *OrderService) Update(ctx context.Context, id string, item string, quantity int32) (*domain.Order, error) {
	order, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	order.Item = item
	order.Quantity = quantity

	err = order.Validate()
	if err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *OrderService) List(ctx context.Context) ([]*domain.Order, error) {
	return s.repo.List(ctx)
}
