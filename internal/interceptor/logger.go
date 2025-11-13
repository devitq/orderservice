package interceptor

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type LoggerInterceptor struct{}

func NewLoggerInterceptor() *LoggerInterceptor {
	return &LoggerInterceptor{}
}

func (i *LoggerInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		log.Printf("gRPC method %s called", info.FullMethod)

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		if err != nil {
			if st, ok := status.FromError(err); ok {
				log.Printf("error: %s, code: %s, duration: %v",
					st.Message(), st.Code(), duration)
			} else {
				log.Printf("error: %v, duration: %v", err, duration)
			}
		} else {
			log.Printf("method %s completed in %v", info.FullMethod, duration)
		}

		return resp, err
	}
}

func (i *LoggerInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv any,
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		log.Printf("gRPC stream method %s started", info.FullMethod)

		err := handler(srv, stream)

		duration := time.Since(start)

		if err != nil {
			log.Printf("stream method %s failed: %v, duration: %v",
				info.FullMethod, err, duration)
		} else {
			log.Printf("stream method %s completed in %v",
				info.FullMethod, duration)
		}

		return err
	}
}
