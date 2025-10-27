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
		log.Printf("Request: %+v", req)

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		if err != nil {
			if st, ok := status.FromError(err); ok {
				log.Printf("Error: %s, Code: %s, Duration: %v",
					st.Message(), st.Code(), duration)
			} else {
				log.Printf("Error: %v, Duration: %v", err, duration)
			}
		} else {
			log.Printf("Response: %+v", resp)
			log.Printf("Method %s completed in %v", info.FullMethod, duration)
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
			log.Printf("Stream method %s failed: %v, Duration: %v",
				info.FullMethod, err, duration)
		} else {
			log.Printf("Stream method %s completed in %v",
				info.FullMethod, duration)
		}

		return err
	}
}
