package handler

import (
	"errors"
	"log"

	"orderservice/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	if errors.Is(err, domain.ErrOrderNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, domain.ErrOrderAlreadyExist) {
		return status.Error(codes.AlreadyExists, err.Error())
	}
	if errors.Is(err, domain.ErrInvalidOrderData) || errors.Is(err, domain.ErrInvalidID) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	log.Printf("internal server error: %v", err)
	return status.Error(codes.Internal, "internal server error")
}
