package domain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrOrderAlreadyExist = errors.New("order already exist")
	ErrOrderNotFound     = errors.New("order not found")
	ErrInvalidOrderData  = errors.New("invalid order data")
)

func NewOrder(id uuid.UUID, item string, quantity int32) (*Order, error) {
	order := &Order{
		ID: id,
		Item: item,
		Quantity: quantity,
	}

	err := order.Validate()
	if err != nil {
		return nil, err
	}
	
	return order, nil
}

type Order struct {
	ID       uuid.UUID `db:"id"       json:"id"`
	Item     string    `db:"item"     json:"item"`
	Quantity int32     `db:"quantity" json:"quantity"`
}

func (o *Order) Validate() error {
	if strings.TrimSpace(o.Item) == "" {
		return fmt.Errorf("%w: item cannot be empty", ErrInvalidOrderData)
	}
	if o.Quantity <= 0 {
		return fmt.Errorf("%w: quantity must be positive", ErrInvalidOrderData)
	}
	if o.ID.String() == "" {
		return fmt.Errorf("%w: ID cannot be empty", ErrInvalidOrderData)
	}
	return nil
}
