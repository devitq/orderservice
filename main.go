package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	pb "orderservice/pkg/api/test"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

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

func (s *OrderServiceServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	order := &pb.Order{
		Id:       id,
		Item:     req.Item,
		Quantity: req.Quantity,
	}
	s.orders[id] = order

	return &pb.CreateOrderResponse{Id: id}, nil
}

func (s *OrderServiceServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[req.Id]
	if !ok {
		return nil, ErrOrderNotFound
	}

	return &pb.GetOrderResponse{Order: order}, nil
}

func (s *OrderServiceServer) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[req.Id]
	if !ok {
		return nil, ErrOrderNotFound
	}

	order.Item = req.Item
	order.Quantity = req.Quantity

	return &pb.UpdateOrderResponse{Order: order}, nil
}

func (s *OrderServiceServer) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.orders[req.Id]
	if !ok {
		return nil, ErrOrderNotFound
	}

	delete(s.orders, req.Id)

	return &pb.DeleteOrderResponse{Success: true}, nil
}

func (s *OrderServiceServer) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]*pb.Order, 0, len(s.orders))
	for _, o := range s.orders {
		orders = append(orders, o)
	}

	return &pb.ListOrdersResponse{Orders: orders}, nil
}

func main() {
	port := flag.Int("port", 50051, "port to run grpc error on")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, NewOrderServiceServer())

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
