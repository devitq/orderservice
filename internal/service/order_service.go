package service

import (
	"context"
	"errors"
	"sync"

	pb "orderservice/pkg/api/order"

	"github.com/google/uuid"
)

var ErrOrderNotFound = errors.New("order not found")

func generateOrderID() string {
	return uuid.NewString()
}

type OrderServiceServer struct {
	pb.UnimplementedOrderServiceServer

	mu     sync.RWMutex
	orders map[string]*pb.Order
}

func NewOrderServiceServer() *OrderServiceServer {
	return &OrderServiceServer{
		orders: make(map[string]*pb.Order),
	}
}

func (s *OrderServiceServer) CreateOrder(
	_ context.Context,
	req *pb.CreateOrderRequest,
) (*pb.CreateOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := generateOrderID()
	order := &pb.Order{
		Id:       id,
		Item:     req.GetItem(),
		Quantity: req.GetQuantity(),
	}
	s.orders[id] = order

	return &pb.CreateOrderResponse{Id: id}, nil
}

func (s *OrderServiceServer) GetOrder(_ context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[req.GetId()]
	if !ok {
		return nil, ErrOrderNotFound
	}

	return &pb.GetOrderResponse{Order: order}, nil
}

func (s *OrderServiceServer) UpdateOrder(
	_ context.Context,
	req *pb.UpdateOrderRequest,
) (*pb.UpdateOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[req.GetId()]
	if !ok {
		return nil, ErrOrderNotFound
	}

	order.Item = req.GetItem()
	order.Quantity = req.GetQuantity()

	return &pb.UpdateOrderResponse{Order: order}, nil
}

func (s *OrderServiceServer) DeleteOrder(
	_ context.Context,
	req *pb.DeleteOrderRequest,
) (*pb.DeleteOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.orders[req.GetId()]
	if !ok {
		return nil, ErrOrderNotFound
	}

	delete(s.orders, req.GetId())

	return &pb.DeleteOrderResponse{Success: true}, nil
}

func (s *OrderServiceServer) ListOrders(
	_ context.Context,
	_ *pb.ListOrdersRequest,
) (*pb.ListOrdersResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]*pb.Order, 0, len(s.orders))
	for _, o := range s.orders {
		orders = append(orders, o)
	}

	return &pb.ListOrdersResponse{Orders: orders}, nil
}
