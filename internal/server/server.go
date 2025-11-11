package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"orderservice/internal/config"
	"orderservice/internal/interceptor"

	grpcHandlers "orderservice/internal/handler/grpc"
	orderPostgresRepo "orderservice/internal/repository/postgres"
	"orderservice/internal/service"

	pb "orderservice/pkg/api/order"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	config     *config.Config
	db         *sqlx.DB
	redisDB    *redis.Client
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

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterOrderServiceHandlerFromEndpoint(ctx, gwmux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf(":%d", s.config.HTTPPort)
	return http.ListenAndServe(addr, gwmux)
}

func getDatabase(cfg config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.BuildDsn())
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	_, err = db.Exec(orderPostgresRepo.Schema)
	if err != nil {
		return nil, fmt.Errorf("run schema: %w", err)
	}

	return db, nil
}

func getRedis(cfg config.Config) (*redis.Client, error) {
	conn, err := redis.ParseURL(cfg.RedisURI)
	client := redis.NewClient(&redis.Options{
		Addr: conn.Addr,
	})
	if err != nil {
		return nil, fmt.Errorf("parse Redis URI: %w", err)
	}

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("connect to Redis server: %w", err)
	}

	return client, nil
}

func (s *Server) RegisterServices() {
	db, err := getDatabase(*s.config)
	if err != nil {
		log.Print(err)
	}
	s.db = db

	redisDB, err := getRedis(*s.config)
	if err != nil {
		log.Print(err)
	}
	s.redisDB = redisDB

	orderRepo := orderPostgresRepo.NewOrderRepository(db, redisDB, &orderPostgresRepo.Config{CacheEnable: true})
	orderService := service.NewOrderService(orderRepo)
	orderHandler := grpcHandlers.NewOrderHandler(orderService)

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
