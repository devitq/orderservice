package handler

import (
	"context"

	"orderservice/internal/domain"
	"orderservice/internal/service"
	pb "orderservice/pkg/api/order"

	"github.com/google/uuid"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer

	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func mapDomainStructToHandler(order *domain.Order) *pb.Order {
	return &pb.Order{
		Id:       order.ID.String(),
		Item:     order.Item,
		Quantity: order.Quantity,
	}
}

func (h *OrderHandler) CreateOrder(
	ctx context.Context,
	req *pb.CreateOrderRequest,
) (*pb.CreateOrderResponse, error) {
	order, err := h.service.Create(ctx, req.GetItem(), req.GetQuantity())
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.CreateOrderResponse{Id: order.ID.String()}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	parsedID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	order, err := h.service.Get(ctx, parsedID)
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.GetOrderResponse{Order: mapDomainStructToHandler(order)}, nil
}

func (h *OrderHandler) UpdateOrder(
	ctx context.Context,
	req *pb.UpdateOrderRequest,
) (*pb.UpdateOrderResponse, error) {
	parsedID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, mapError(domain.ErrInvalidID)
	}

	order, err := h.service.Update(ctx, parsedID, req.GetItem(), req.GetQuantity())
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.UpdateOrderResponse{Order: mapDomainStructToHandler(order)}, nil
}

func (h *OrderHandler) DeleteOrder(
	ctx context.Context,
	req *pb.DeleteOrderRequest,
) (*pb.DeleteOrderResponse, error) {
	parsedID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	err = h.service.Delete(ctx, parsedID)
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.DeleteOrderResponse{Success: true}, nil
}

func (h *OrderHandler) ListOrders(
	ctx context.Context,
	_ *pb.ListOrdersRequest,
) (*pb.ListOrdersResponse, error) {
	domainOrders, err := h.service.List(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	orders := make([]*pb.Order, 0, len(domainOrders))
	for _, o := range domainOrders {
		orders = append(orders, mapDomainStructToHandler(o))
	}

	return &pb.ListOrdersResponse{Orders: orders}, nil
}
