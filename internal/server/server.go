package server

import (
	"fmt"
	"log"
	"net"

	"orderservice/internal/config"
	"orderservice/internal/interceptor"
	"orderservice/internal/service"

	pb "orderservice/pkg/api/order"

	"google.golang.org/grpc"
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

func (s *Server) RegisterServices() {
	orderService := service.NewOrderServiceServer()
	pb.RegisterOrderServiceServer(s.grpcServer, orderService)

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
