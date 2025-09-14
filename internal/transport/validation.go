package transport

import (
	pb "github.com/MaksimPozharskiy/grpc-balance-processor/proto/balance"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateProcessRequest(req *pb.ProcessRequest) error {
	if req.AccountId == "" {
		return status.Error(codes.InvalidArgument, "account_id is required")
	}
	if req.TxId == "" {
		return status.Error(codes.InvalidArgument, "tx_id is required")
	}
	if req.Amount == "" {
		return status.Error(codes.InvalidArgument, "amount is required")
	}
	if req.Source == pb.Source_SOURCE_UNSPECIFIED {
		return status.Error(codes.InvalidArgument, "source is required")
	}
	if req.State == pb.State_STATE_UNSPECIFIED {
		return status.Error(codes.InvalidArgument, "state is required")
	}
	return nil
}

func validateAndParseAccountID(accountID string) (uuid.UUID, error) {
	if accountID == "" {
		return uuid.Nil, status.Error(codes.InvalidArgument, "account_id is required")
	}

	parsed, err := uuid.Parse(accountID)
	if err != nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "invalid account_id format: must be valid UUID")
	}

	return parsed, nil
}

func validateAndParseAmount(amount string) (decimal.Decimal, error) {
	if amount == "" {
		return decimal.Zero, status.Error(codes.InvalidArgument, "amount is required")
	}

	parsed, err := decimal.NewFromString(amount)
	if err != nil {
		return decimal.Zero, status.Error(codes.InvalidArgument, "invalid amount format: must be valid decimal")
	}

	if parsed.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, status.Error(codes.InvalidArgument, "amount must be positive")
	}

	if parsed.Exponent() < -2 {
		return decimal.Zero, status.Error(codes.InvalidArgument, "amount must have at most 2 decimal places")
	}

	return parsed, nil
}

func validateTxID(txID string) error {
	if txID == "" {
		return status.Error(codes.InvalidArgument, "tx_id is required")
	}

	if len(txID) > 128 {
		return status.Error(codes.InvalidArgument, "tx_id must be at most 128 characters")
	}

	return nil
}
