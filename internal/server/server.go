package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"orderservice/internal/config"
	"orderservice/internal/interceptor"

	orderGrpcHandler "orderservice/internal/handler/grpc"
	orderInMemory "orderservice/internal/repository/inmemory"
	"orderservice/internal/service"

	pb "orderservice/pkg/api/order"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	config     *config.Config
}

func New(cfg *config.Config) *Server {
	loggerInterceptor := interceptor.NewLoggerInterceptor()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggerInterceptor.Unary()),
		grpc.StreamInterceptor(loggerInterceptor.Stream()),
	)

	return &Server{
		grpcServer: grpcServer,
		config:     cfg,
	}
}

func runHTTPHandler(s *Server, grpcServerEndpoint *string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf(":%d", s.config.HTTPPort)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) RegisterServices() {
	repo := orderInMemory.NewOrderRepository()
	orderService := service.NewOrderService(repo)
	orderHandler := orderGrpcHandler.NewOrderHandler(orderService)

	pb.RegisterOrderServiceServer(s.grpcServer, orderHandler)

	if s.config.GRPCEnableReflection {
		reflection.Register(s.grpcServer)
		log.Println("gRPC server will start with reflection")
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.GRPCPort)
	lis, err := net.Listen("tcp", addr) //nolint:noctx // no need to use context here
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	if s.config.EnableHTTPHandler {
		go func() {
			log.Printf("Starting HTTP gateway on port %d", s.config.HTTPPort)
			if err := runHTTPHandler(s, &addr); err != nil {
				log.Printf("HTTP gateway failed: %v", err)
			}
		}()
	}

	log.Printf("Starting gRPC server on port %d", s.config.GRPCPort)

	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
	log.Println("gRPC server stopped gracefully")
}
