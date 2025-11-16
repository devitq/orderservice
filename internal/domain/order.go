package domain

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var (
	ErrOrderAlreadyExist = errors.New("order already exist")
	ErrOrderNotFound     = errors.New("order not found")
	ErrInvalidOrderData  = errors.New("invalid order data")
)

type Order struct {
	ID       uuid.UUID `db:"id"       json:"id"       validate:"required"`
	Item     string    `db:"item"     json:"item"     validate:"required"`
	Quantity int32     `db:"quantity" json:"quantity" validate:"required,gt=0"`
}

func NewOrder(id uuid.UUID, item string, quantity int32) (*Order, error) {
	order := &Order{
		ID:       id,
		Item:     item,
		Quantity: quantity,
	}

	err := order.Validate()
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *Order) Validate() error {
	validate := validator.New()

	return fmt.Errorf("%w: %w", ErrInvalidOrderData, validate.Struct(o))
}
