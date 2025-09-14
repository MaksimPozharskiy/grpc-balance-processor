package transport

import (
	"errors"

	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapDomainError(err error) error {
	if errors.Is(err, domain.ErrNotFound) {
		return status.Error(codes.NotFound, "account not found")
	}
	if errors.Is(err, domain.ErrDuplicateTx) {
		return status.Error(codes.AlreadyExists, "transaction already exists")
	}
	if errors.Is(err, domain.ErrNegativeBalance) {
		return status.Error(codes.InvalidArgument, "insufficient balance")
	}
	return status.Error(codes.Internal, "internal server error")
}
